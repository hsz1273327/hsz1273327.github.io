---
title: "使用docker封装和部署应用"
date: 2017-05-15
author: "Hsz"
category: experiment
tags:
    - Docker
    - DevOps
    - MicroService
header-img: "img/home-bg-o.jpg"
update: 2019-03-16
---
# 使用docker封装和部署应用

各种编程语言的社区都或多或少的为其应用的部署提供解决方案.

比如python中从开发调试测试到部署,有一系列简单好用的工具:
+ 开发使用标准库`venv`
+ 依赖管理使用`pip`
+ 批量部署使用[fabric]

在js中默认就是使用虚拟环境,只需要通过git在目标机器上拉取仓库,并执行`npm install`就可以完美部署好执行环境;

而如果是Go那就更加简单了,本地可以交叉编译目标平台的可执行文件,编译完成后直接放上去就可以执行.

从运维的角度来说docker的所谓一次编辑随处部署在多数情况下是多余的.那为什么还要用docker呢?

+ docker比所有的虚拟环境都底层,它只依赖于系统内核,因此更加不容易出问题
+ docker的使用方式是配置化的,这意味着可以将开发和运维分离,运维不需要知道要部署和维护的软件是如何开发的,运行需要什么参数,run就行了,借助一些可视化工具运维甚至都不用懂linux命令行操作(前提是它只是用docker而非维护和部署docker环境).
+ docker由于本质上是一种轻量级的虚拟机技术,它可以配置每个容器的使用资源,这样可以更好的控制整个系统的资源分配.
+ docker天生为集群服务,当你的应用受众变大后它可以更容易的横向扩展


## 从一个简单的项目出发

亲身经历docker的学习过程最好借助于例子,一个例子走一遍基本什么都知道怎么回事了.我们先做一个最简单的TODOList项目,用它来介绍如何

## 依赖

项目依赖`flask`和`pymongo`

docker构建依赖`docker`

## 封装python应用

我们以一个简单的flask项目作为例子,看怎么封装的,项目在[我的github](https://github.com/hsz1273327/python_docker_example/)上.

该项目由一个`requirements.txt`文件和一个`app.py`文件构成,它需要外部的mongodb来存取数据.
项目的共能也非常简单--为一个collection提供增,查功能.代码如下

```python
from flask import Flask, request, json
from flask.views import MethodView
from pymongo import MongoClient
import argparse


def args_parse():
    parser = argparse.ArgumentParser()
    parser.add_argument("-H", "--host", type=str,
                        help="主机", default="localhost")
    parser.add_argument("-p", "--port", type=int, help="端口", default=27017)
    parser.add_argument("-d", "--database", type=str,
                        help="数据库", default="default")
    args = parser.parse_args()
    return args


args = args_parse()
client = MongoClient(args.host, args.port)
db = client[args.database]
user_collection = db.user


class UserAPI(MethodView):

    def get(self):
        if not request.args:
            result = [i for i in user_collection.find()]
        else:
            if request.args["user_name"]:
                result = [i for i in user_collection.find(
                    {"user_name": request.args["user_name"]})]

            else:
                result = [i for i in user_collection.find()]

        return json.dumps(str(result))

    def post(self):
        obj = request.get_json()
        result = user_collection.save(obj)
        return json.dumps({'result': str(result), 'ok': True})


app = Flask(__name__)
app.add_url_rule('/users/', view_func=UserAPI.as_view("users"))


def main():
    app.run(debug=True, host="0.0.0.0")


if __name__ == '__main__':
    main()

```

app默认链接localhost的27017的mongodb.可以通过命令行修改链接的mongo主机和端口
默认启动在本地5000端口,并且可以内网访问

## Dockerfile

Dockerfile

```dockerfile
FROM python:3.6
ADD . /code
WORKDIR /code
RUN pip install -r requirements.txt
```

在项目目录下使用`docker build -t example_app .`创建镜像

跑项目只要:

`docker run -p 5000:5000  example_app python /code/app.py -H <你的mongo主机> -p <你的端口> -d <你的数据库名>`

## 将你的镜像上传至仓库


## 部署镜像为容器

## 限制资源

这边是重点了,之所以要使用docker甚至不惜牺牲一些性能主要也就是为了可以控制资源使用.

### Cpu限制

要限制cpu使用某几个cpu

+ 默认的libcontainer引擎

    只要`run`后面加上参数`--cpuset=0,3`来指定.这样就是只是用了1号和4号cpu了.

    如果要设置相对的cpu权重比,可以使用`--cpu-shares=n`n必须为一个非负数

+ 如果使用lxc引擎，

    可以指定 --lxc-conf="lxc.cgroup.cpuset.cpus = 0,1"
    如果要设置相对的cpu权重比,可以使用`--lxc-conf="lxc.cgroup.cpu.shares = 1234"`

### 内存限制

`-m`参数可以直接设定内存的最大使用量

`-m 100m`这样就是只允许使用100m大小的内存,多了就会自己崩溃

本文主要是针对python玩家,专业的运维还需要了解`docker-compose`这类容器编排工具.这就不在本文的范围内了.

### 自动重启

## 将所有依赖编排到一起


## 将项目部署在集群上




