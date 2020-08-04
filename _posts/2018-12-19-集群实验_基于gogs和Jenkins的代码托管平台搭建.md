---
title: "集群实验_基于gogs和Jenkins的代码托管平台搭建"
date: 2018-12-19
author: "Hsz"
category: introduce
tags:
    - Linux
    - Docker
    - DockerSwarm
    - CI/CD
    - Git
header-img: "img/home-bg-o.jpg"
update: 2018-12-19
series:
    cluster_experiment:
        index: 6
---
# 代码托管平台

现代软件工程是一个系统工程,早就不是一个大神干就能干下来的时代了.软件开发意味着合作,合作就意味着需要有套机制管理代码的版本,有一套系统管理代码的权限.

这种天然的需求就催生了代码托管平台的诞生.当然有的代码托管平台还附带上了社交属性,这就是另一个故事了.

目前最常见的代码托管平台都是依托于git工具.开源界最知名的应该就是[github](https://github.com/).大量优质的开源项目都托管在上面.无论是程序员,数据科学家,还是开源硬件的玩家都会在其上分享自己的代码.而目前企业应用最广的代码托管平台[gitlab](https://about.gitlab.com/)则依赖代码开源,生态完善的优势被广泛应用.

但这两位都不今天的主角,今天的主角是[gogs](https://gogs.io/),这是一个go写的代码托管应用,很小,但也功能单一,它的优势是资源消耗少,甚至可以在树莓派上安装,非常适合10人以下的小型团队作为代码托管工具.它有完整的工单系统和分组系统(虽然只有一层)同时支持邮件推送.文档也可以使用自带的wiki模板解决.如果需要ci/cd工具,则可以使用[Jenkins](https://jenkins.io/)

本文使用docker部署gogs和Jenkins,因此需要[相关的前置知识](http://blog.hszofficial.site/blog/2018/05/31/%E6%A0%91%E8%8E%93%E6%B4%BE%E7%9A%84%E9%9B%86%E7%BE%A4%E5%AE%9E%E9%AA%8C_%E9%85%8D%E7%BD%AEdocker%E5%8F%8Aswarm/)

## gogs部署

gogs在docker swarm环境可以直接安装部署:

```yml
version: "3.6"
services:
  gogs:
    image: gogs/gogs
    ports:
      - "10022:22"
    networks:
      - net-output
    volumes:
      - gogs-data:/data
    deploy:
      placement:
        constraints: [node.labels.codehub == codehub]

volumes:
  gogs-data:
    external: true
networks:
  net-output:
    external: true
```

需要注意,为了数据不丢失,我们需要预先创建一个volume叫`gogs-data`,然后是使用`external: true`让stack使用它就好,如果是nfs的那最好,但如果不是,也不是不行,因为本来也就只起一个实例.
因此只要指定好node和volume所在的宿主机一致即可.

初始化时我们会被要求进行一些设定,这些设定的具体意思可以看[官网说明](https://gogs.io/docs/advanced/configuration_cheat_sheet).这些设定大多都是可以通过`gogs-data`下的`gogs/conf/app.ini`.

## Jenkins

Jenkins是老牌的ci/cd工具,老实说有点过时了,但为啥还用它呢?我尝试了drone,满地都是坑所以放弃了,相对来说再没有gitlab-ci的情况下Jenkins是几乎唯一的免费选择.

Jenkins部署十分简单

```yml
version: "3.6"
services:
  drone-server:
    image: jenkinsci/blueocean
    volumes:
      - jenkins-data:/var/jenkins_home
      - /var/run/docker.sock:/var/run/docker.sock
    networks:
      - net-output
    deploy:
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
        window: 20s
      replicas: 1
      placement:
        constraints: [node.labels.codehub == codehub]

volumes:
  jenkins-data:
    external: true

networks:
  net-output:
    external: true

```

他会开方8080端口作为http的访问端口,注意再nginx中进行反向代理.

安装好后需要去我们设置的volume`jenkins-data`中找到启动用的秘钥来解锁这个服务.而之后呢就是一些基础设置了.全设置完后我们需要在系统管理中找到插件管理器,安装两个插件`Gogs`和`Generic Webhook Trigger Plugin`.安装完后重启以加载插件.


## Jenkins用法

jenkins并不是为每个项目单独设置一套流水线,而是预先定义一套流程名,之后项目选择自己使用哪个流程,这种思路上的不同确实会让人一开始不太适应,但仔细想想,在资源有限的情况下似乎Jenkins才是最合理的选择.
Jenkins的中文文档已经有翻译<https://jenkins.io/zh/doc/book/>可以对着学习下.
Jenkins的用法大的方向上分为3块:

+ 为Jenkins配置权限
+ 在gogs中配置Jenkins
+ 在项目代码仓库中定义流程

### 为Jenkins配置权限

我们的Jenkins要使用docker的话需要为其配置访问权限,一个简单粗暴的方法是直接把宿主机的`docker.sock`权限设为777

```shell
chmod 777 /var/run/docker.sock
```
但这过于简单粗暴,同时暴露安全隐患,另一种方式是修改镜像.

```dockerfile
FROM jenkins:1.651.1
MAINTAINER xxx
USER root
RUN apt-get update \
  && apt-get install -y apt-transport-https ca-certificates \
  && echo "deb https://apt.dockerproject.org/repo debian-jessie main" > /etc/apt/sources.list.d/docker.list \
  && apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D \
  && apt-get update -y \
  && apt-get install -y docker-engine \
  && apt-get install sudo \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
ADD jenkins.sh /usr/local/bin/jenkins.sh
RUN chmod a+x /usr/local/bin/jenkins.sh
ENTRYPOINT ["/usr/local/bin/jenkins.sh"]
```

其中的jenkins.sh为:

```shell
#!/bin/bash -x

# this only works if the docker group does not already exist

DOCKER_SOCKET=/var/run/docker.sock
DOCKER_GROUP=docker

if [ -S ${DOCKER_SOCKET} ]; then
    DOCKER_GID=$(stat -c '%g' ${DOCKER_SOCKET})
    groupadd -for -g ${DOCKER_GID} ${DOCKER_GROUP}
    usermod -aG ${DOCKER_GROUP} jenkins
fi
```

本方案来自官方的工单<https://github.com/jenkinsci/docker/issues/263>

之后我们需要为Jenkins配置一个可以访问git仓库的用户,将其在`Jenkins主页左侧->凭据页面`中添加到全局域中.

然后我们需要创建一个job,我们创建一个多分支流水线项目,设置的时候在分支源处添加`Project Repository`为需要做ci/cd的项目仓库地址,`Credentials`选上之前定义好的权限凭据,
写好`Display Name`和`Description`即可.

之后其他项目要用这个job,修改这个job的配置,把新的`Project Repository`和`Credentials`添加进去即可.

### 在gogs中配置Jenkins

gogs的`仓库设置`中点击`管理web钩子`,将推送地址写为`http://<hostname>/gogs-webhook/?job=<jonname>`,之后测试推送试一下.

### 在项目代码仓库中定义流程

在我们的项目根目录中新建一个名叫`Jenkinsfile`的文件,其中按groovy语法写上执行流程就可以.语法可以看<https://jenkins.io/zh/doc/book/pipeline/syntax/>


下面是一个python项目的示例,这个项目会在test,dev,master分支执行Test,在master分支执行后续的Build和Release.


```groovy
pipeline {
  agent none 
  stages {
    stage('Test') {
      when {
        anyOf {
          branch 'test'
          branch 'dev'
          branch 'master'
        }
      }
      agent {
          docker {
              image 'python:3.6' 
          }
      }
      steps {
        withEnv(["HOME=${env.WORKSPACE}"]) {
          sh 'python -m pip install --user -r requirements.txt'
          sh 'python -m coverage run --source=test_drone -m unittest discover -v -s .'
          sh 'python -m coverage report -m' 
        }
      }
    }
    stage('Build') {
      when {
        branch 'master'
      }
      agent any
      steps {
        withEnv(["HOME=${env.WORKSPACE}"]) {
          sh 'docker build -t hsz1273327/test_drone:0.0.1 .'
        }
      }
    }
    stage('Release') {
      when {
        branch 'master'
      }
      agent any
      steps {
        withEnv(["HOME=${env.WORKSPACE}"]) {
          sh 'docker login -u hsz1273327 -p hsz881224'
          sh 'docker push hsz1273327/test_drone:0.0.1'
        }
      }
    } 
  }
}
```

上面是单节点Jenkins的用法,在资源叫少的情况下,一个master节点就很够用了,但一旦规模扩大,我们也可以通过添加执行节点的方式来满足需求.不过这边就不多叙述了.