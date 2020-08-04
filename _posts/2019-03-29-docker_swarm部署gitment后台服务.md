---
layout: post
title: "docker swarm 部署gitment后台服务"
date: 2019-03-29
author: "Hsz"
category: recommend
tags:
    - Github
    - Docker
    - DockerSwarm
header-img: "img/home-bg-o.jpg"
update: 2019-03-29
---
# docker swarm 部署gitment后台服务

这个博客最早的时候使用的是多说,但后来我把账号给忘了...于是想想要是能没有依赖就可以有评论功能该有多好呀,然后就发现了一个只需要github账号就可以登录发表评论的工具[gitment](https://github.com/imsun/gitment).这边先感谢下这位作者大大,无论如何比多说这种确实方便很多--借用github的账号系统,评论放项目工单,怎么都比把评论交给一个不知哪里来的公司强多了.

但有一段时间了gitment失效了...为什么呢?经过查找发现这个js会向一个服务器发送请求...然后这个服务器是作者自己搭建的,而且已经停了....

那要怎么修复就很简单了.
1. 自己建立这个服务器
2. browser.js中引用这个网址的地方，改为自己的网址
3. 主题中引用的browser.js，不能是原来gitment的js，而要是自己修改过的js

但不是光这样就好刚好借这个机会来演示下如何在`docker swarm`集群上构建这个服务(算是最优实践吧)

## 介绍下背景

我之前在腾讯云买了两个服务器,用这两台服务器搭了一个`docker swarm`的集群.这个集群使用`portainer`进行管理,对外有一个nginx用于做反向代理和虚拟主机.要对外的容器使用`net-output`这个网络相连

## 构建镜像

这个服务的代码在github上<https://github.com/imsun/gh-oauth-server>我们可以fork下来为它添加dockerfile,docker-compose.yml这些东西用于构建镜像

这个项目在[我的这个仓库](https://github.com/hszofficial/gh-oauth-server)里.

+ dockerfile

    ```dockerfile
    FROM node:latest
    MAINTAINER hsz
    ADD server.js /app/server.js
    ADD package.json /app/package.json
    WORKDIR /app
    #安装依赖
    RUN npm install
    #对外暴露的端口
    EXPOSE 3000
    CMD [ "npm","start"]
    ```

+ build.sh

    ```bash
    docker build -t hsz1273327/gh-oauth-server:0.0.1 -t hsz1273327/gh-oauth-server:latest .
    dcoker push hsz1273327/gh-oauth-server
    ```

执行好后这个镜像就在[我的docker hub](https://cloud.docker.com/u/hsz1273327/repository/docker/hsz1273327/gh-oauth-server)中了.

## 部署服务

### 服务配置

下面是docker-compose.yml文件的内容

```yml
version: "3.6"
services:
  gh-oauth-server:
    image: hsz1273327/gh-oauth-server:latest
    networks:
      - net-output
    deploy:
      restart_policy: #设置重启配置
        condition: on-failure
        delay: 5s
        max_attempts: 3
        window: 120s
      replicas: 4
      update_config:
        parallelism: 2
        delay: 10s
      placement:
        constraints: [node.labels.codehub == codehub]

networks:
  net-output:
    external: true
```

这里我们使用的是已经存在的网络,所以`networks`中使用`external: true`这个属性表示是挂载到一个一直网络

顺道一提如果要自己建一个网络让别的服务挂载可以在manager节点的命令行中输入`docker network create --attachable --driver overlay --scope swarm xxxxx`创建


`restart_policy`用于设置重启的选项,一旦程序报错退出,docker会自动按这个配置重启容器,上面这个配置的意思是当有错误时延迟5s重启,最大重试3次,每次间隔120s


`replicas`表示我们会起多少个实例


`update_config`则指定更新service时的行为,上面配置的意思是更新时先更新两个,隔10s再更新剩下的两个

### 部署服务

有两种方式:
1. 进去[我的这个项目仓库](https://github.com/hsz1273327/gh-oauth-server),吧`cocker-compose.yml`的内容复制下来,之后打开`portainer`,按路径`stacks=>addstack`进入创建stack的页面,然后选择`web editor`,把内容贴上然后点`deploy the stack`.
2. 进去[我的这个项目仓库](https://github.com/hsz1273327/gh-oauth-server),fork下,之后打开`portainer`,按路径`stacks=>addstack`进入创建stack的页面,然后选择`git Repository`,把自己的登录信息和项目仓库地址填上,并在开始的地方为这个stack取个名字.然后点`deploy the stack`

### 注册一个二级域名

我的域名是在阿里云买的,上去给他设置个就叫`gitment`的二级域名就行

### 修改nginx为服务提供反向代理

nginx如下这么配置

```conf
upstream gh-oauth-server{
    server  gh-oauth-server:3000;
}
server {
    listen       80;
    server_name  gitment.hszofficial;

    location / {
        proxy_pass http://gh-oauth-server;
    }
}
```

改好后build更新镜像然后更新下nginx的stack

## 修改gitment中的代码

查找到`gitment.browser.js`中有`gh-oauth`的代码,然后用我们的网址替换掉代码中对应的网址,这样就修复好了