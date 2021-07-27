---
layout: post
title: "使用Jenkins代替GithubActions自动化工作流"
date: 2020-12-02
author: "Hsz"
category: recommend
tags:
    - DevOps
header-img: "img/home-bg-o.jpg"
update: 2021-07-21
series:
    get_along_well_with_github:
        index: 10
---
# 使用Jenkins代替GithubActions自动化工作流

CI/CD是现代开发体系中提高工作效率的基础.我们知道一般开发行为中开发编码的时间往往只占30%,剩下的全是测试和部署.如果可以降低测试和部署的时间那就可以大大提高工作效率.
CI/CD就是这样的工具,它的作用就是利用脚本自动化测试和部署.

在`Github Actions`出现之前,Github上也是使用第三方CI/CD工具的,那个时候[Jenkins](https://github.com/jenkinsci/jenkins)就是主流之一(另一个是Travis CI).
<!--more-->

## Jenkins简介

`Jenkins`是完全开源免费的项目,它的使用也没有什么限制.而且它的部署不依赖代码仓库,用户系统也是和代码仓库不通用的,这也就意味着天生的开发与测试,部署隔离.同时它支持图形化的流程配置,可以一定程度上降低维护人员的学习成本.我们只需要忍受丑陋的ui即可.

本文推荐对想了解更多的读者可以看[jenkins官方文档](https://jenkins.io/zh/doc/book/pipeline/)作为补充

顺道一提我之前公司一直使用的是gitlab套件,虽然体验上一致性还是不错的,但迁移了几次,每次都得重新部署全套,这相当让人厌烦.不少公司的运维人员恐怕连编程都不会,更不要提写脚本部署了,每次迁移都是一个伤筋动骨的过程,要不等上将近一周的时间等运维一个一个组件的部署完,要么自己动手.而权限分配也必然是我这个pm的任务,每次迁移都会重新配置一次权限.如果CI/CD和代码仓库分离,那么我只需要重新配置对应项目的仓库地址,这虽然也是个体力活,但已经比全部重配好得多.程序设计中有个单一职责原则,在工具选择上个人认为也是同样适用的.

<!-- 本项目对应的代码在[hszofficial/test_jenkins](https://github.com/hszofficial/test_jenkins) -->

## docker上搭建Jenkins

不啰嗦直接上docker-compose.yml

```yml
jenkins-server:
  image: jenkinsci/blueocean
  volumes:
    - /volume2/docker_deploy/devtools/jenkins/data:/var/jenkins_home
    - /var/run/docker.sock:/var/run/docker.sock
  mem_limit: 2g
  restart: on-failure
  user: root
  privileged: true
  ports:
    - "8080:8080"
  logging:
    <<: *default-log

```

下面是部署的几个注意点:

1. 如果要使用https,我们需要设置如下环境变量,当然还要把证书和私钥挂到volumes上,同时可以闪电`8080`端口的映射.

    ```yaml
    environment: 
      JENKINS_OPTS: "--httpPort=-1 --httpsPort=8083 --httpsCertificate=/certs/x.pem --httpsPrivateKey=/certs/x.key"
    ```

2. 如果要考虑后续的扩展性,可以打开`50000端口`,这个端口可以用于后续挂载slaver节点.
3. 如果我们要用到docker(现在几乎不会有用不到docker的情况),那么需要加上`user: root`和`privileged: true`,否则很容易docker会无权限使用
4. 目前该镜像只支持amd64指令集

## 配置jenkins的各项功能

安装完自然要配置,jenkins主要的配置项有:

+ 节点配置(非必须)
    jenkins支持多节点,其好处是可以在多台机器上做编译,测试工作,以提高吞吐量,当然小型团队完全没有必要搞.
    一样我们可以使用镜像[jenkins/inbound-agent](https://hub.docker.com/r/jenkins/inbound-agent/),最好将他部署到swarm集群上,注意目前该镜像也只支持amd64指令集

    部署好后再jenkins中`系统管理->节点管理`中对节点进行配置和监控.

+ 安装插件管理
    在`系统管理->插件管理`中可以管理插件.安装插件在右上角搜索框中查找到后点它安装即可,插件安装完后需要重启服务,这个是自动的我们不用人为干预.

    ![插件管理界面]({{site.url}}/img/in-post/jenkins/jenkins-插件管理.PNG)

    我们会安装如下插件:
    + `Git Parameter`用于选择分支执行
    + `Docker`和`Docker Pipeline`用于使用镜像来作为沙盒执行CI/CD任务

    其他实用插件我会在后面专门的章节补充介绍

+ 用户和安全配置
    在`系统管理->管理用户`中可以对用户进行管理.

    接着在`系统管理-->全局安全配置-->授权策略`中选择`项目矩阵授权策略`,然后为你的运维组成员设置不同的全局权限.

    不同项目的权限可以在`主页`点击`项目名`进入后再`配置`项中勾选`启用项目安全`来激活.在其中添加你希望添加的用户,并给他服务相应权限.

    通常一家规模不大的公司,运维可能只有1,2个人,这种时候其实就没有太大必要弄得这么复杂,直接给与权限就好

+ 邮箱配置
    邮箱需要安装`Email Extension Plugin`插件这个插件应该是默认安装的因此不用再额外安装,接着在`系统管理->系统设置`中找到`E-mail Notification`然后写上你可以发送消息的邮箱登录信息即可.

    ![设置邮箱]({{site.url}}/img/in-post/jenkins/emailconf1.png)
    注意在这之前还需要配置下管理员邮箱,这个邮箱要与发件箱一致.
    ![设置管理员邮箱]({{site.url}}/img/in-post/jenkins/emailconf2.png)

    接着配置`Extended E-mail Notification`

    前面和默认邮箱配置一致,后面会有一些新的内容,主要是设置默认的收件人`Default Recipients`,可填写多个,中间用空格隔开.这个值可以在`$DEFAULT_RECIPIENTS`变量中取到

## 使用gitea作为代码仓库

在使用gitea作为代码仓库的情况下我们需要分别在jenkins和gitea上做出如下设置:

> gitea上的设置

1. [可选]如果你的jenkins使用了https协议且签名是自己发布的,那么需要在git的配置文件中加上下面配置否则webhook会无法发送成功

    ```conf
    [webhook]
    ; Allow insecure certification
    SKIP_TLS_VERIFY = true
    ```

2. 创建一个用户专门用于给jenkins操作,这个用户并不需要是管理员,但创建完后我们需要进去(`用户->应用`中)为它生成一个`access token`

    ![生成令牌]({{site.url}}/img/in-post/jenkins/gitea-生成令牌.PNG)

3. 将所有需要让jenkins管理cicd的代码库都添加上新建的用户为协作者,并给它管理员权限

    ![设置协作者]({{site.url}}/img/in-post/jenkins/gitea-协作者.PNG)

> jenkins上的设置

1. 安装插件`Gitea`,这个插件可以自动将gitea代码仓库和jenkins关联.
2. 在`系统配置->Manage Credentials(凭据)->Stores scoped to jenkins-> Jenkins(域为全局)->系统->全局凭据`下用`添加凭据`创建一条类型为`Gitea Personal Access Token`的凭据,`Token`字段填上刚才新建出来的`access token`
    ![添加凭据1]({{site.url}}/img/in-post/jenkins/jenkins-添加凭据1.PNG)
    ![添加凭据2]({{site.url}}/img/in-post/jenkins/jenkins-添加凭据2.PNG)
3. 在`系统管理->系统配置`中设置插件`Gitea`的配置项`Gitea Servers`,将自己的gitea仓库注册进去,并勾选`Manage hooks`,并在`Credentials`上使用刚才配置的凭证
    ![Gitea插件配置]({{site.url}}/img/in-post/jenkins/jenkins-Gitea插件配置.PNG)
4. 点击`新建任务`,选择`Gitea Organization`进入配置页面,需要注意的只有如下几个配置项:
    1. `Project->Credentials`选择上面创建的凭证
    2. `Project->Owner`上写上要关注的命名空间(用户名或者组织名)
        ![基本配置]({{site.url}}/img/in-post/jenkins/jenkins-owner.PNG)
    3. `Behaviours`中选择需要的行为测lure,主要是针对`pull request`,另外建议点击`Add`将`Discover tags`添加上
        ![行为]({{site.url}}/img/in-post/jenkins/jenkins-behaviour.PNG)
    4. `扫描 Gitea Organization 触发器`部分根据需要选择扫描的时间间隔,个人建议使用默认的1 day
    5. `孤儿项策略`中设置好就流水线的删除策略,建议`保留旧的流水线的天数`设置为3,`保留旧的流水线的最大数`设置为30
        ![其他设置]({{site.url}}/img/in-post/jenkins/jenkins-其他设置.PNG)

上面这些都设置好了以后只要我们的仓库中包含文件`Jenkinsfile`并且语法没问题就会在扫描过后被添加进来.我们也可以不等它定时扫描,通过项目上点击`立刻扫描 Gitea Organization`来主动扫描添加仓库.

## 为项目构建CI/CD流程

`Jenkinsfile`是描述CI/CD流程的声明文件,使用[Groovy语法](https://www.w3cschool.cn/groovy/groovy_basic_syntax.html),我们可以每个分支都有一个`Jenkinsfile`,也可以只在`master`分支有一个`Jenkinsfile`个人推荐在项目创建之初就将`Jenkinsfile`创建好,后面每次创建分支就会都带上这个文件,需要修改的话就单独修改.

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
+ `post`用于定义不同时点的输出.角色类似python中的`try...except...else`语句,通常通过发送邮件起到通知作用

每个步骤包括名字和实现两个部分,实现部分又分为

+ `when`触发条件

+ `agent`使用的执行代理(现在一版都是docker)

+ `steps`具体的执行步骤,支持的所有`steps`可以在[这里找到](https://www.jenkins.io/doc/pipeline/steps/),不过最常用的还是只有`echo`,`sh`,和`dir`

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
      sh 'python -m coverage report -m > report.txt'
    }
  }
}
```

每隔步骤中也可以定义`agent`和`environment`,他们会在当前定义的步骤范围内生效,全局的环境变量则有两种:

1. `在设置->全局元素中`点击环境变量添加.
2. 默认会有一些环境变量,可以在<https://wiki.jenkins.io/display/JENKINS/Building+a+software+project>中查到

这一块的具体配置可以看[pipeline定义文档](https://www.w3cschool.cn/jenkins/jenkins-jg9528pb.html)和[使用script定义文档的关键字文档](https://jenkins.io/doc/pipeline/steps/workflow-basic-steps/#code-readfile-code-read-file-from-workspace)

### 展示html报告

我们可以在项目的`pipeline syntax`中设置`publish HTML`,然后只要我们在pipeline中有对应的`publishHTML`被执行了,就可以在分支的pipeline中左侧找到对应的链接了.

![设置html报告位置]({{site.url}}/img/in-post/jenkins/html.png)

我们甚至可以用它来部署接口文档

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

## 定义api触发任务

一些任务我们可能需要外部触发执行,我们可以在`新建任务`中选择`流水线`,然后在`构建触发器`位置选择`触发远程构建 (例如,使用脚本)`(如果要构造定时任务在这里选`定时构建`).
然后在`身份验证令牌`位置填一个随机的字符串作为token(推荐使用uuid4或者拿时间戳和一个密码做个md5),这样就可以通过api来触发执行了.

![api触发器]({{site.url}}/img/in-post/jenkins/jenkins-api触发器.PNG)

api有两种:

+ 不带参数的任务用`<jenkin_url>/job/test-api/build?token=<TOKEN>`
+ 带参数的任务`<jenkin_url>/job/test-api/buildWithParameters?token=<TOKEN>&<param_name>=<param_value>&...`

### 定义外部参数

在任务配置的`General`位置我们可以勾选`参数化构建过程`,勾选以后我们就可以向其中添加不同类型的参数,最常用的是`字符串型`,`布尔型`和`选项型`,其中需要注意的是`选项型`中的选项使用回车分隔,且第一个选项为默认值

![声明外部参数]({{site.url}}/img/in-post/jenkins/jenkins-声明外部参数.PNG)
