---
title: "树莓派配置docker及swarm"
date: 2018-05-31
author: "Hsz"
category: introduce
tags:
    - RaspberryPi
    - Docker
    - DockerSwarm
header-img: "img/home-bg-o.jpg"
update: 2018-05-31
series:
    raspberrypi_experiment:
        index: 5
    cluster_experiment:
        index: 4
---
# docker

## 安装

```shell
curl -sSL https://get.docker.com | sh
```

## 配置自启动

使用下面的命令将docker注册到开机自启动

```shell
sudo usermod -aG docker pi

sudo systemctl enable docker

sudo systemctl start docker
```

## 将pi用户添加进docker组

```shell
sudo gpasswd -a ${USER} docker
```

之后重启服务器或者重启docker服务即可

```shell
sudo service docker restart
```

## swarm配置节点

swarm集群中有两种角色

+ manager节点 
    
    manager角色用于管理集群中的节点,也可以运行容器并提供容灾.
    `sudo docker swarm init`


+ worker节点 

    工人,就是供容器执行的平台
    `sudo docker swarm join --token <token> <location>`
    先在一台机器上执行init后就会获得token和location


## docker运维可视化

swarm集群可以使用[portainer](https://github.com/portainer/portainer)进行可视化.功能很简单不做详细叙述.


### 安装和启动

在你希望启动该服务的节点上执行

```shell
docker volume create portainer_data
docker run -d -p 9000:9000 --name portainer --restart always -v /var/run/docker.sock:/var/run/docker.sock -v portainer_data:/data portainer/portainer:linux-arm
```

之后就可以访问9000端口管理集群了.当然了还是推荐将这个服务配置在manager节点上,这样如果manager节点奔溃了直接就没法访问这个服务,比较好辨识

### 集群配置

在manager节点上执行如下命令,这样每个节点就都会被放入portainer的管理范围中:

```shell
curl -L https://portainer.io/download/portainer-agent-stack.yml -o portainer-agent-stack.yml
docker stack deploy --compose-file=portainer-agent-stack.yml portainer
```

## 注意点

由于树莓派是arm架构,因此一般的镜像无法使用,需要使用支持arm的镜像才行.

有两种方式:
1. 去dockhub查找想要的镜像有没有`armv7/raspberrypi/linux-arm`字样的标签
2. 去找官方支持armv7的镜像作为基镜像自己编译.这些镜像在<https://hub.docker.com/u/arm32v7/>可以找到.

## 基于portainer的swarm集群使用

portainer是一个足够好的运维工具,在网络足够好的情况下足以应付小规模swarm集群的一般运维.

以下是其常用操作:

### 创建一个网络

左侧`Networks`可以用于定义网络,比较常用的是如下设置:

+ Driver--overlay
+ scope--swarm
+ attachable--true

这种方式可以让swarm中各个节点都可以访问

### 创建一个硬盘挂载

左侧`Volumes`可以用于定义挂载的卷,在集群环境下卷只有宿主机上的容器才能用,

没啥可以选的,创建的时候driver选择local,然后为其指定一个mountpoint即可即可.

如果需要数据共享,可以结合nfs实现.此处不多阐述,以后有机会再说.使用方法是创建的时候

形式如下:
```
+ driver: local
+ driver_opts:
    + type: "nfs"
    + o: "addr=fs-0214f623.efs.ap-northeast-1.amazonaws.com,nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,noresvport"
    + device: "fs-0214f623.efs.ap-northeast-1.amazonaws.com:/consul/certs"
```

### 创建一个配置

点击左侧`configs`就可以进入配置的管理界面,点击`Add config`来新建一个配置,只要写上配置名和配置内容即可.

Configs 是安装在容器的文件系统中的而不是使用 RAM 磁盘.可以随时添加或删除,服务可以共享一个配置,与 Environments 或 Labels 结合使用可以获得最大的灵活性.

config可以在compose中申明,也可以在外部定义

在使用时当做文件访问,例子如下:

```yml
version: "3.3"
services:
  redis:
    image: redis:latest
    deploy:
      replicas: 1
    configs:
      - my_config
      - my_other_config
configs:
  my_config:
    file: ./my_config.txt
  my_other_config:
    external: true
```

### 创建一个秘钥

点击左侧`Secrets`就可以进入秘钥管理界面,点击`Add secret`来新建一个秘钥.在其中填入秘钥名和对应值就好,要访问时就当做文件访问,比方说秘钥名为`hszofc_crt`地址就是`/run/secrets/hszofc_crt`

要使用需要在compose文件中显式的声明.


### 创建一个stack

创建stack一样在左侧,可以通过上传文件,直接复制内容,或者给出git仓库和分支实现.

下面是一个简单的stack示例

```yml
version: '3.6'
services:
  nginx-output:
    image: registry.hszofficial.site/lib/cus_nginx:latest
    ports:
      - 80:80
      - 443:443
    networks:
      - net-output
    logging:
      options:
          max-size: "10m"
          max-file: "3"

    secrets:
      - hszofc_crt
      - hszofc_key
      - monitor_crt 
      - monitor_key
      - code_crt
      - code_key
      - registry_crt
      - registry_key
    volumes:
      - static_pages:/usr/local/static/
    deploy:
      restart_policy:
        condition: on-failure
      replicas: 3
      placement:
        constraints: [node.labels.type == worker]
secrets:
  hszofc_crt:
    external: true
  hszofc_key:
    external: true
  monitor_crt:
    external: true
  monitor_key:
    external: true
  code_crt:
    external: true
  code_key:
    external: true
  registry_crt:
    external: true
  registry_key:
    external: true
volumes:
  static_pages:
    external: true
networks:
  net-output:
    external: true
```