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



## 作为备份库

