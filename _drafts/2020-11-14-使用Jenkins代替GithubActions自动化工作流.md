---
layout: post
title: "使用Jenkins做CICD"
date: 2019-04-10
author: "Hsz"
category: experiment
tags:
    - DevOps
    - Github
header-img: "img/home-bg-o.jpg"
update: 2019-04-10
series:
    get_along_well_with_github:
        index: 6
---
# 使用Jenkins做CI/CD

CI/CD是提高工作效率的基础.我们知道一般开发行为中开发编码的时间往往只占30%,剩下的全是测试和部署.如果可以降低测试和部署的时间那就可以大大提高工作效率.

CI/CD就是这样的工具,它的作用就是利用脚本自动化测试和部署.

那为啥使用[Jenkins](https://github.com/jenkinsci/jenkins)呢?毕竟这已经是一个比较旧的工具了,而且用于配置`Jenkinsfile`使用的是类似C语言的写法(本质是groovy),与当今流行的json或者yml完全不通用.因为它可以不依赖任何其他的工具.

当今最流行的CI/CD工具有3种:

+ `Travis CI` 项目必须使用github托管,必须由github账号
+ `Drone` 依赖docker集群
+ `Gitlab runner` 依赖gitlab套件

而`Jenkins`是完全开源免费的项目,它的使用也完全没有这种限制.而且它的部署不依赖代码仓库,用户系统也是和代码仓库不通用的,这也就意味着天生的开发与测试,部署隔离.同时它支持图形化的流程配置,可以一定程度上降低维护人员的学习成本.我们只需要忍受丑陋的ui即可.

本文推荐对想了解更多的读者可以看[jenkins官方文档](https://jenkins.io/zh/doc/book/pipeline/)作为补充

顺道一提我之前公司一直使用的是gitlab套件,虽然体验上一致性还是不错的,但迁移了几次,每次都得重新部署全套,这相当让人厌烦.不少公司的运维人员恐怕连编程都不会,更不要提写脚本部署了,每次迁移都是一个伤筋动骨的过程,要不等上将近一周的时间等运维一个一个组件的部署完,要么自己动手.而权限分配也必然是我这个pm的任务,每次迁移都会重新配置一次权限.如果CI/CD和代码仓库分离,那么我只需要重新配置对应项目的仓库地址,这虽然也是个体力活,但已经比全部重配好得多.程序设计中有个单一职责原则,在工具选择上个人认为也是同样适用的.

本项目对应的代码在[hszofficial/test_jenkins](https://github.com/hszofficial/test_jenkins)

## docker集群上搭建Jenkins

不啰嗦直接上docker-compose.yml

```yml
version: "3.6"
services:
  jenkins-server:
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

这边有两个预先创建的项:

+ `jenkins-data`,用于挂载运行过程中产生的数据
+ `net-output`,对外网络,让网关的nginx可以识别到,默认端口是8080,我个人是将它反向代理到一个二级域名上.

我是在swarm集群上部署的,可以使用protainer进行管理,新建一个stack把上面配置的贴上就可以部署了.

jenkins至今依然是一个活跃的开源项目,依然会有更新,其中的插件也会有更新,在protainer中更新的方式很简单,先拉取最新的镜像,之后进入stack使用`update`即可.

## 配置jenkins的各项功能

安装完自然要配置,jenkins主要的配置项有:

+ 节点配置(非必须)

jenkins支持多节点,其好处是可以在多台机器上做编译,测试工作,以提高吞吐量,当然小型团队完全没有必要搞.
一样我们可以使用镜像[jenkins/jnlp-slave.](https://hub.docker.com/r/jenkin/jnlp-slave),将它添加到上面的docker-compose.yml中,注意最好新建一个networks专门用于jenkins组件间通信,因为一个swarm网络最多只可以有255个服务.

部署好后再jenkins中`系统管理->节点管理`中对节点进行配置和监控.

+ 插件管理

在`系统管理->插件管理`中可以管理插件.安装插件在右上角搜索框中查找到后点它安装即可,插件安装完后需要重启服务,这个是自动的我们不用人为干预.

我们会安装`Git Parameter`.

+ 用户和安全配置

在`系统管理->管理用户`中可以对用户进行管理.

接着在`系统管理-->全局安全配置-->授权策略`中选择`项目矩阵授权策略`,然后为你的运维组成员设置不同的全局权限.

不同项目的权限可以在`主页`点击`项目名`进入后再`配置`项中勾选`启用项目安全`来激活.在其中添加你希望添加的用户,并给他服务相应权限.

通常一家规模不大的公司,运维可能只有1,2个人,这种时候其实就没有太大必要弄得这么复杂,直接给与权限就好

+ 邮箱配置

邮箱需要安装`Email Extension Plugin`插件,接着在`系统管理->系统设置`中找到`E-mail Notification`然后写上你可以发送消息的邮箱登录信息即可.

![设置邮箱]({{site.url}}/img/in-post/jenkins/emailconf1.png)
注意在这之前还需要配置下管理员邮箱,这个邮箱要与发件箱一致.
![设置管理员邮箱]({{site.url}}/img/in-post/jenkins/emailconf2.png)

接着配置`Extended E-mail Notification`

前面和默认邮箱配置一致,后面会有一些新的内容,主要是设置默认的收件人`Default Recipients`,可填写多个,中间用空格隔开.这个值可以在`$DEFAULT_RECIPIENTS`变量中取到

## 构建一个基于git的项目

基于git的项目我们一般使用`多分支流水线`,针对不同的分支和行为进行不同的管理.创建方法是:

+ 在`凭据->系统`下添加一个域比如叫`admin`,进入其中使用`添加凭据`创建一条凭据,凭据可以是用户名密码也可以是ssh的秘钥.
+ `New 任务->输入项目名->多分支流水线`进入多分支流水线设置页面
+ 设置`Build Configuration`选`Mode: by Jenkinsfile`,`Script Path:Jenkinsfile`这样我们就是可以使用项目中的`Jenkinsfile`来定义过程了
+ 如果用于编译的库不在dockerhub拉取而是在自己的镜像仓库拉取,name可以设置`Pipeline Model Definition`中的内容
    + `Docker Label`用于设置默认拉取的镜像标签(如果没显式的写出来),
    + `Docker registry URL`用于设置私有仓库地址
    + `Registry credentials`用于设置私有仓库的登录凭证
+ `Scan 多分支流水线 Triggers`勾选`Periodically if not otherwise run`,将周期设为`1h`一般就可以了.这个设置可以隔段时间扫描下仓库创建出合适的分支

这些都设置好了还不能保存,我们开始重点--针对`gogs`和`github`项目的配置,这个项目需要配置`Branch Sources`

### 配置gogs项目

我的个人私有仓库就是使用的gogs,这东西用go写的,轻量简单,有工单系统,有wiki系统,作为一个代码仓库已经很好用了.

> jenkins端的设置

1. 在`Branch Sources`中选择`Git`,复制上项目的仓库地址并填上登录凭证,
2. `Behaviours`中添加行为,包括
   1. Discover branches
   2. Discover tags

之后保存,保存好了后gogs会扫描项目创建pipeline

> gogs端的配置

pipeline配置gogs需要在gogs中进入要配置的项目

![gogs上配置]({{site.url}}/img/in-post/jenkins/gogs测试.png)

之后保存,然后再次进去这个webhook可以在gogs上测试线是否可以连通

### 配置github项目

> jenkins端的项目设置

1. 在`Branch Sources`中选择`Github`,复制上项目的仓库地址并填上登录凭证,
2. 填上拥有者,也就是组织或者用户的名字
3. 这时候就可以在下面的下拉菜单中找到要配置的项目了
4. 下面的行为一块都有翻译,就不介绍了

之后保存,保存好了后gogs会扫描项目创建pipeline

> github端的权限配置

在github上`个人页面->setting->Developer settings->Personal access tokens`生成一个token,注意选择下权限.

jenkins必须使用token才能进行访问

> jenkins端的访问设置

1. 在`凭据->系统`中添加一个新的密钥,选择`secret text`后填上刚拿到的token,这个就叫做`github token`,再添加一个token填上一段自己的文字叫`github access token`
2. 在`系统管理->系统配置->	GitHub`上`添加github服务器`,将`github token`选上,点击测试下能不能通
3. 选择`advance`后勾选`为 Github 指定另外一个 Hook URL`,这时我们可以获得一个url,并选上`github access token`.

![jenkins端的访问设置]({{site.url}}/img/in-post/jenkins/jenkins端的访问设置.png)

> github端的项目配置

1. 进入项目,在`setting->Webhooks->add Webhooks`填上获得一个url,并填上`github access token`对应的文本

![github上配置]({{site.url}}/img/in-post/jenkins/github测试.png)


## 为项目构建pipeline

我们可以每个分支都有一个`Jenkinsfile`,也可以只在`master`分支有一个`Jenkinsfile`个人推荐在项目创建之初就将`Jenkinsfile`创建好,后面每次创建分支就会都带上这个文件,需要修改的话就单独修改.

一个典型的`Jenkinsfile`如下:

```groovy
pipeline {
  agent none
  environment {
    sendmail = 'yes'
    version = '0.0.2'
  }
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
          echo 'install requirement'
          sh 'python -m pip install --user -r requirements.txt'
          echo 'start test'
          sh 'python -m coverage run --source=test_drone -m unittest discover -v -s .'
          echo 'send report'
          sh 'python -m coverage html -d report/coverage'
        }
      }
      post {
        success {
          publishHTML([
            allowMissing: true, 
            alwaysLinkToLastBuild: true, 
            keepAll: true, 
            reportDir: 'report/coverage', 
            reportFiles: 'index.html', 
            reportName: 'Coverage Report - Unit Test'
            ])
          emailext body: "${git_url}:${git_branch} test success",
            subject: "${git_url}:${git_branch} test success",
            to: "hsz1273327@gmail.com"
        }
        failure {
          emailext body: "${git_url}:${git_branch} test failure",
            subject: "${git_url}:${git_branch} test failure",
            to: "hsz1273327@gmail.com"
        }
      }
    }
    stage('Release') {
      when {
        branch "release-*"
      }
      agent any
      steps {
        withEnv(["HOME=${env.WORKSPACE}"]) {
          sh 'docker build -t hsz1273327/test_drone:latest -t hsz1273327/test_drone:'+version+' .'
          sh 'docker login -u hsz1273327 -p hsz881224'
          sh 'docker push hsz1273327/test_drone'
        }
      }
    }
  }
  post{
    success {
      script {
        if (sendmail == 'yes') {
          emailext body: '''pipelie succeed:
          构建名称:${JOB_NAME}
          构建结果:${BUILD_STATUS}
          构建编号：${BUILD_NUMBER}
          GIT 地址：${git_url}
          GIT 分支：${git_branch}
        ''',
          subject: 'Jenkins build ${PROJECT_NAME} succeed', 
          to: 'hsz1273327@gmail.com'
        }
      }
    }
    failure {
      script {
        if (sendmail == 'yes') {
          emailext body: '''pipelie failure:
            构建名称:${JOB_NAME}
            构建结果:${BUILD_STATUS}
            构建编号：${BUILD_NUMBER}
            GIT 地址：${git_url}
            GIT 分支：${git_branch}
            ${BUILD_LOG}''',
            subject: 'Jenkins build ${PROJECT_NAME} is ${currentBuild.result}: ${env.JOB_NAME} #${env.BUILD_NUMBER}',
            to: 'hsz1273327@gmail.com'
        }
      }
    }
  }
}
```

我们可以看到`pipeline`定义一条流程管道,其中

+ `agent`用于定义全局使用的执行代理,
+ `environment`用于定义全局的变量
+ `stages`则用于定义管道中的步骤.
+ `post`用于定义不同时点的输出.一般用于发送邮件

每个步骤包括名字和实现两个部分,实现部分又分为

+ `when`触发条件

+ `agent`使用的执行代理

+ `steps`具体的执行步骤

```groovy
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
      sh 'python -m coverage report -m >report.txt'
    }
  }
}
```

这一块的具体配置可以看[pipeline定义文档](https://www.w3cschool.cn/jenkins/jenkins-jg9528pb.html)和[使用script定义文档的关键字文档](https://jenkins.io/doc/pipeline/steps/workflow-basic-steps/#code-readfile-code-read-file-from-workspace)


### 展示html报告

我们可以在项目的`pipeline syntax`中设置`publish HTML`,然后只要我们在pipeline中有对应的`publishHTML`被执行了,就可以在分支的pipeline中左侧找到对应的链接了.

![设置邮箱]({{site.url}}/img/in-post/jenkins/html.png)

### 发送邮件

我们可以在post中定义发送邮件的逻辑

```text
emailext body: '''pipelie failure:
  构建名称:${JOB_NAME}
  构建结果:${BUILD_STATUS}
  构建编号：${BUILD_NUMBER}
  GIT 地址：${git_url}
  GIT 分支：${git_branch}
  ${BUILD_LOG}''',
  subject: 'Jenkins build ${PROJECT_NAME} is ${currentBuild.result}: ${env.JOB_NAME} #${env.BUILD_NUMBER}',
  to: 'hsz1273327@gmail.com'
```

body中可以定义html模板,subject是主题,to指定发送去的邮箱

## 分支管理

在实际项目中我的经验是主干分支策略会比较高效.

> 如果没有专职的运维,那么应该把代码和配置分开,配置打包进image中

```shell

{代码开发 v0-base(v0-base分支)  ==>   dev-0.0.0(dev分支)  ==>  0.0.0(master分支,同时打好latest标签)  ==>  release-0.0.0(release-0.0.0分支,留档) }
                                            |                          |
                                            V                          V
                                {代码调试 (debug分支)}        {配置文件部分 deploy-xxx(release-xxx分支,使用latest标签)}
                                                                        |
                                                                        V
                                                            {部署部分 test分支,使用标签为deploy-xxx的镜像}
                                                                        |
                                                                        V
                                                            {部署部分 production分支,使用标签为deploy-xxx的镜像}
```

> 如果有专职的运维,那么代码仓库中应该只管代码

```shell

{代码开发 v0-base(v0-base分支)  ==>   dev-0.0.0(dev分支)  ==>  0.0.0(master分支,同时为image打好latest标签)  ==>  v0.0.0(v0.0.0 tag,留档) }
                                     |           |
                                     V           V
                     {代码调试 (debug分支)}      {test分支,使用标签为test-xxx的镜像}
```
