---
layout: post
title: "使用Supervisor做服务监控和管理"
date: 2015-02-19
author: "Hsz"
category: introduce
tags:
    - DevOps
header-img: "img/post-bg-2015.jpg"
update: 2019-03-15
---
# 使用Supervisor做服务监控和管理

我们的服务往往启动后如果关闭终端就会跟着关闭,这时候我们就需要把进程设置成守护进程,而[Supervisor](http://supervisord.org/)就是一个可以管理守护进程的工具.

Supervisor虽然现在在某种程度上已经被容器技术所取代,但在一些不用容器的情况下它依然有用.

安装可以直接pip安装.

## 配置

实际上supervisord是一个c/s结构的软件,我们要先来指定配置文档,默认的配置文档一般放在`/etc/supervisor/supervisord.conf`,而细化的各个服务的配置
/etc/supervisor/conf.d/gogs.conf 

## 任务配置

,在配置文档中至少要有一个可运行的命令:

```conf
[program:foo]
command=<foo>.py
```

一般的样例:

```conf
[program:bar]
command=python <bar>.py;需要执行的命令wd
user =root  ;默认是当前用户,如果需要可以改到root用户
autostart=true  ;是否在supervisord启动时自启动(default: true)
autorestart=true  ;可以设置是否自动重启,何时自动重启(default: unexpected)
startsecs=3  ;最大启动时间设置( def . 1)
stderr_logfile=/tmp/bandwidth_err.log  ;标准错误输出位置定向
stdout_logfile=/tmp/bandwidth.log  ;log输出位置定向
```

我们可以用`echo_supervisord_conf > supervisord.conf`来生成一个模板,这样修改模板就可以了.

### 关于虚拟环境的处理方式

我们如果要使用python的虚拟环境来启动项目的话,可以添加上`environment=PATH="/env/path/bin:%(ENV_PATH)s"`

### 配置文件套接

我们的默认配置文件中有这样一句:

```conf
[include]
files = /etc/supervisor/conf.d/*.conf
```

也就是说我们可以将任务一个一个放在`/etc/supervisor/conf.d/`文件夹下,只要后缀为`.conf`就可以被识别引入.习惯上我们以项目名作为对应项目的配置文件名.

files后面的参数可以不止一条,他们中间用空格分隔

### 程序分组

我们可以在单独的配置文件中设置`[group:x]`来将几个程序分组管理,我们可以自己在`conf.d`文件夹下新建一个文件夹`groups`,再在其中设置分组


```conf
[group:foo]
programs=bar,baz
priority=999
```

### 配置http管理工具

将`supervisord.conf`中`inet_http_server`部分做相应配置，在`supervisorctl`中`reload`即可启动web管理界面,如果要使外部可以访问,需要将port后面的具体host改为`*`

```conf
[inet_http_server]
port = 127.0.0.1:9001
username = user
password = 123
```

### 其他配置

+ `supervisord`: supervisor的服务器端部分,用于supervisor启动,他的设置可以有:

```conf
[supervisord]
logfile = /tmp/supervisord.log
logfile_maxbytes = 50MB
logfile_backups=10
loglevel = info
pidfile = /tmp/supervisord.pid
nodaemon = false
minfds = 1024
minprocs = 200
umask = 022
user = chrism
identifier = supervisor
directory = /tmp
nocleanup = true
childlogdir = /tmp
strip_ansi = false
environment = KEY1="value1",KEY2="value2"
```

+ `supervisorctl`:启动supervisor的命令行窗口,在该命令行中可执行start、stop、status、reload等操作.他的设置可以有:

```conf
[supervisorctl]
serverurl = unix:///tmp/supervisor.sock
username = chris
password = 123
prompt = mysupervisor
```

## 开启服务

`supervisord -c supervisord.conf`可以指定配置文档并启动.

## 命令行控制项目启停

supervisorctl就是supervisord的命令行客户端.其操作有:

操作|说明
---|---
`supervisorctl update`|更新新的配置到supervisord,之后常用`supervisorctl reload`来重载配置
`supervisorctl reload`|重载配置
`supervisorctl shutdown`|关闭supervisord
`supervisorctl status`|查看项目状态
`supervisorctl start`|启动项目
`supervisorctl restart`|重启项目
`supervisorctl stop`|停止项目

## 后台监控

[cesi](https://github.com/Gamegos/cesi)是官方推荐的后台监控工具.它本质上是一个supervisor的插件.

cesi可以监控你注册的节点服务器,利用supervisor来监控管理其上的项目.

### 配置node

node代表服务器节点,在cesi的配置文件`cesi.conf`上这样注册:

```conf
[node:gulsah]
username = gulsah
password = ***
host = gulsah.xyz.com
port = 9001
```
注意这里填的信息是上文中配置http管理工具中的内容

### 配置environment

接着可以将上面的节点分组,用的是`environment`设置

```conf
[environment:market]
members = gulsah, kaan
```

### 基本配置

最后是设置cesi自己的一些信息,主要是设置log信息存放位置和指定sqlite3数据库

```conf
[cesi]
database = /x/y/userinfo.db
activity_log = /x/y/cesi_activity.log
```

cesi有完善的权限管理功能,可以为不同的管理员设置不同的权限

## 其他插件

有许多基于supervisor的[插件](http://supervisord.org/plugins.html#dashboards-and-tools-for-multiple-supervisor-instances)可以在官网找到.有兴趣的可以取拿来试试.