---
title: "使用Gitea代替Github托管代码"
date: 2020-11-12
author: "Hsz"
category: recommend
tags:
    - Git
    - DevOps
header-img: "img/home-bg-o.jpg"
update: 2020-11-12
series:
    get_along_well_with_github:
        index: 9
---
# 使用Gitea代替Github托管代码

Github虽好,但它不是开源的!它是代码托管的服务商.确实Github做的非常好用,大家用起来也非常方便,但不要忘记,它不是开源的!使用它你需要冒被限制账号甚至被封号的风险.毕竟伊朗人已经被封了.

没错这不符合开源精神,但是它本来就不是开源的!因此我们有必要给Github找个备胎,甚至拜托Github自己来.[Gitea](https://github.com/go-gitea/gitea)就是一个相当靠谱的选择.

<!--more-->

## Gitea简介

Gitea脱胎于gogs项目,是一个使用go编写的git代码托管工具.之所以推荐用它代替,有如下考虑:

+ Gitea几乎支持所有Github提供的服务.Gitea官方提供了一个[与其他各种代码托管工具横向对比的文档](https://docs.gitea.io/zh-cn/comparison/),从其中可以看出Gitea在功能上几乎和Github的社区版一致.截至至目前只有`Github Pages`,`Github Actions`和`Github Project`中的多仓库关联功能没有支持.`Github Actions`可以使用其他开源CI/CD工具替代,`Github Project`中的多仓库关联功能目前看是正在做,相信不久就可以有了,`Github Pages`则可以通过`CI/CD`工具配合Nginx自行部署.况且有wiki其实基本也够用.

+ Gitea部署极其方便,因为是go写的,可以直接二进制部署.而如果用docker部署则更加方便.这点对比Gitlab是极大的优势,Gitlab的部署非常麻烦.

+ Gitea自身部署成本极低,可以在x86-64/arm等各种机器上部署,只要有个树莓派的性能就可以跑的很流畅.这点同样对比Gitlab是极大的优势,Gitlab官方要求4核8g的机器.

## 安装与配置

个人更加推荐使用docker的单机方式部署.

+ 在你的宿主机上找个合适的地方创建一个文件夹用于映射gitea的`/data`目录,我们以`/volume2/docker/gitea/data`为例

+ 单机模式部署`docker-compose.yml`

    ```yml
    version: '2.4'
    services:
    gitea:
        image: gitea/gitea:latest
        logging:
        options:
            max-size: "10m"
            max-file: "3"
        ports:
        - 10022:22
        - 10080:3000
        volumes:
        - "/volume2/docker/gitea/data:/data"

    ```

+ 修改`/volume2/docker/gitea/data/gitea/conf/app.ini`,这个文件是配置文件.主要是要修改如下几项

    + `APP_NAME`,你的程序名
    + `server`,主要是设置
        + `PROTOCOL` http还是https
        + `CERT_FILE`(如果是https就需要设置,可以创建一个文件夹`/volume2/docker/gitea/data/keys`后将证书放在其中)
        + `KEY_FILE`(如果是https就需要设置,可以创建一个文件夹`/volume2/docker/gitea/data/keys`后将私钥放在其中)
        + `DOMAIN`服务器域名,如果没有就填ip地址
        + `SSH_DOMAIN`SSH的服务器域名,如果没有就填ip地址
        + `ROOT_URL`,服务对外的URL

    + `database`,设置元数据的存放,建议使用外部的数据库,推荐pg.
    + `mailer`,设置发送消息的发件箱

+ 重启`service`

个人推荐要么直接将Gitea部署在内网不要暴露到外网,然后外网通过vpn进入内网后使用;要么就暴露到外网,但一定要设置域名并设置为https方式访问.

第一次登陆后我们需要设置一些基本配置.之后就可以正常使用了

## 像Github一样使用Gitea

登陆进Gitea服务后界面如下:

![Gitea的个人页][1]

可以看到Gitea自带中文,并且基本功能的入口和Github的也很相似,我们可以一目了然的找到入口.

+ 在文档方面,Gitea可以使用`README`和`Wiki`管理项目文档;
+ 在项目管理方面,Gitea支持`pull request`,支持对`pull request`的内容做`code review`,也支持工单,标签,里程碑和kanban.
+ 在成果分发方面,Gitea也支持`release`.
+ 在项目组织方面,Gitea支持team也支持组织.
+ 在社交方面Gitea甚至还支持点赞...

## 作为备份库

Gitea的一大作用就是作为github的备份库,Gitea直接提供了对应的入口--左上角的`+`

![备份github][2]

填好提交就可以了

[1]: {{site.url}}/img/in-post/gitea/mainpage.PNG
[2]: {{site.url}}/img/in-post/gitea/qianyi.PNG