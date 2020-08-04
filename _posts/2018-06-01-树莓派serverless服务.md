---
title: "树莓派serverless服务"
date: 2018-06-01
author: "Hsz"
category: introduce
tags:
    - RaspberryPi
    - Docker
    - DockerSwarm
    - Serverless
header-img: "img/home-bg-o.jpg"
update: 2018-06-01
series:
    raspberrypi_experiment:
        index: 6
    cluster_experiment:
        index: 4
---
# Serverless

Serverless不代表再也不需要服务器了,而是说:开发者再也不用过多考虑服务器的问题,计算资源作为服务而不是服务器的概念出现.Serverless是一种构建和管理基于微服务架构的完整流程,允许你在服务部署级别而不是服务器部署级别来管理你的应用部署,你甚至可以管理某个具体功能或端口的部署,这就能让开发者快速迭代,更快速地开发软件.

Serverless有以下几个特点：

1. Serverless意味无维护

    Serverless不代表完全去除服务器,而是代表去除有关对服务器运行状态的关心和担心,它们是否在工作,应用是否跑起来正常运行等等.Serverless代表的是你不要关心运营维护问题.有了Serverless,可以几乎无需Devops了.

2. Serverless不代表某个具体技术

    有些人会给他们的语言框架取名为Serverless,Serverless其实去除维护的担心,如果你了解某个具体服务器技术当然有帮助,但不是必须的.

3. Serverless中的服务或功能代表的只是微功能或微服务

    Serverless是思维方式的转变,从过去`构建一个框架运行在一台服务器上，对多个事件进行响应`变为`构建或使用一个微服务或微功能来响应一个事件`，你可以使用`django`,`node.js`和`express`等实现.但是serverless本身超越这些框架概念.框架变得也不那么重要了.

树莓派同样可以利用`open Faas`构建serverless服务.`OpenFaaS`是`Docker`下的一个框架,可以让任何进程或容器成为`Serverless`函数,由于`Docker`和`Golang`的便携性,它在树莓派上运行得很好.

## 安装

在manager节点上执行如下指令:

```shell
git clone https://github.com/openfaas/faas && \
cd faas && \
git checkout 0.7.8 && \
./deploy_stack.armhf.sh
```

安装完后可以在主节点的`8080`端口访问到`open Faas`的UI界面,也称为API网关.这是您可以定义,测试和调用函数的地方.


## 配置open Faas开发环境

需要注意,由于树莓派是arm架构的,和x86/amd64架构无法兼容,因此open Faas的代码开发可以在本地,但编译部署必须在树莓派上才可以执行.

我们需要在本地(比如我的mac下)和树莓派上(以manager机为例.)都安装`faas-cli`工具.本地的用于使用脚手架快速创建样板服务.树莓派上的则用于编译和部署.

mac下可以使用`brew install faas-cli`安装,而树莓派上则可以使用指令`curl -sSL cli.openfaas.com | sudo sh`安装.

这样我们可以使用git工具,在mac下编写,而在树莓派上编译执行.

要查看某个网关上挂载的函数,可以使用如下命令

```shell
faas-cli list --gateway http://192.168.3.40:8080
```

## 一个例子

下面是一个例子,通过这个例子我们来看看serverless服务时如何开发的.

1. 我们在我们的gogs上新建一个项目`functest`
2. 将其clone至我们的pc(mac)本地
3. 在项目根目录使用`faas-cli`工具创建模板

    ```shell
    faas-cli new --lang python-armhf functest
    ```

    这时会在新增如下内容:
    ```shell
    functest-|
             |-handler.py 项目的实际代码
             |-requirements.txt 执行代码依赖的库
    functest.yml 项目的配置文件
    template # 各种支持的模板
    ```
    `functest.yml`是项目的配置文件,其模板如下:
    ```yml
    provider:
        name: faas   # 此处
        gateway: http://<ip>:8080 #此处是启动该服务的机器ip,本处就是192.168.3.40,即`piNodeMaster`机

    functions:
        functest:
            lang: python-armhf
            handler: ./functest
            image: functest
    ```
    + `provider` 服务提供者信息配置
        + `gateway` 指定远程网关地址
    + `functions` 这个块用于定义本项目的编译配置
        + `lang`指定使用的模板语言.
        + `handler`指定`handler.py`及其所用到其他模块文件所在的文件夹
        + `image` 编译后的image名

    其他我们还需要在`requirement.txt`中写上依赖`requests`
    handler.py则写成这样:
    ```python
    import requests
    import json

    def handle(req):
        try:
            result = {"found": False}
            json_req = json.loads(req)
            r = requests.get(json_req["url"])
            if r.status_code == 200:
                if json_req["term"] in r.text:
                    result = {"found": True}
            print(json.dumps(result))
            return json.dumps(result)
        except Exception as e:
            print("{typ},{msg}".format(typ=type(e),msg=str(e)))
    ```


4. 将项目上传至gogs
5. 在随便哪台树莓派上clone这个项目
6. 在树莓派上这个项目的根目录下执行

    ```shell
    faas-cli build -f ./functest.yml
    ```
    之后就会开始编译打包,执行完后就有对应的镜像了

7. 执行`faas-cli push -f ./functest.yml`命令将镜像上传至docker或者指定的仓库
8. 在树莓派上这个项目的根目录下执行

    ```shell
    faas-cli deploy -f ./hello-python.yml
    ```

## 如何调用这些函数

调用这些函数的方式有两种:

1. 通过`open Faas`的网关界面走gui调用.
2. 通过api接口调用:

    ```shell
    curl <ip>:8080/function/functest --data-binary '{
        "url": "https://blog.hszofficial.site/",
        "term": "docker"
    }'
    ```