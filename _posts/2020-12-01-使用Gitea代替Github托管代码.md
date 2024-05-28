---
title: "使用Gitea代替Github托管代码"
date: 2020-12-01
author: "Hsz"
category: recommend
tags:
    - Git
    - DevOps
header-img: "img/home-bg-o.jpg"
update: 2024-05-27
series:
    get_along_well_with_github:
        index: 9
---
# 使用Gitea代替Github托管代码

Github虽好,但它不是开源的!它是代码托管的服务商.确实Github做的非常好用,大家用起来也非常方便,但不要忘记,它不是开源的!使用它你需要冒被限制账号甚至被封号的风险.毕竟伊朗人已经被封了.

没错这不符合开源精神,但是它本来就不是开源的!因此我们有必要给Github找个备胎,甚至摆脱Github自己来.[Gitea](https://github.com/go-gitea/gitea)就是一个相当靠谱的选择.

<!--more-->

## Gitea简介

Gitea脱胎于gogs项目,是一个使用go编写的git代码托管工具.之所以推荐用它代替,有如下考虑:

+ Gitea几乎支持所有Github提供的服务.Gitea官方提供了一个[与其他各种代码托管工具横向对比的文档](https://docs.gitea.io/zh-cn/comparison/),从其中可以看出Gitea在功能上几乎和Github的社区版一致.截至至目前只有`Github Pages`功能没有支持.`Github Pages`则可以通过`CI/CD`工具配合Nginx自行部署.况且有wiki其实基本也够用.

+ Gitea部署极其方便,因为是go写的,可以直接二进制部署.而如果用docker部署则更加方便.这点对比Gitlab是极大的优势,Gitlab的部署非常麻烦.

+ Gitea自身部署成本极低,可以在x86-64/arm等各种机器上部署,只要有个树莓派的性能就可以跑的很流畅.这点同样对比Gitlab是极大的优势,Gitlab官方要求4核8g的机器.

+ Gitea支持分布式部署,有相当的扩展性,可以渐进式的扩展以满足各个阶段的需求

同时gitea有一个还算挺活跃的社区,[awesome-gitea](https://gitea.com/gitea/awesome-gitea)项目也总结了一些相关的有用资源.对于小型团队来说其实挺够用的

## 运维部分

gitea需要自己部署,自然也就有了一份运维的工作.好在由于gitea是go开发的应用,跨平台部署非常方便而且天然适应分布式部署非常容易扩展和维护.但在好做也还是要有维护的,这部分我们详细介绍gitea的运维操作.

### 规模与部署模式

下面是总结的不同用户规模下比较合适的部署模式.我们以

| 用户规模               | 部署方式          | 存储后端  | 元数据存储 |
| ---------------------- | ----------------- | --------- | ---------- |
| 3人以下(家用级别)      | docker部署        | local     | sqlite     |
| 10人以下(小工作室级别) | docker部署        | nas local | pg         |
| 50人以下(小微公司级别) | docker部署        | minio/nfs | pg         |
| 50人以上(集团级别)     | Kubernetes HA部署 | minio ha  | pg ha      |

通常用的到git的都是技术人员.

技术人员数量在50人以下都可以直接使用单节点部署,不同点仅为仓库存储和元数据存储方案的选择.如果是自己在家里开个夫妻老婆店,那拿个闲置的笔记本搭个gitea就够用了,使用默认的本地方式所有数据都以文件的形式放在机器上,而机器上的访问又十分方便,应付这种需求是足够的.如果是小工作室或小微公司,建议整个纯ssd的nas(sata接口的就够)把gitea挂在上面,可以两块ssd做个简单raid1(读1.8写0.9,容量0.5),然后把pg装上面.如果大文件比较多或发布比较频繁的,比如是做游戏开发的机器学习模型开发的,剪视频的,也可以考虑用minio或nfs作为存储的后端,然后给这个后端一个足够大的硬盘空间.需要注意ssd有读写次数限制,因此好做好备份,个人更加推荐使用minio.需要注意代码仓库无论如何都是保存在硬盘上的无法保存到minio存储后端.但

超过50人基本都是很大的公司了.大公司一般多base,每个地区部署一份gitea,数据则通过异地多活技术共享即可.

### 安装

我们用最简单的docker单机方式部署.gitea官方提供两个版本的镜像

+ 原始版本,比如`gitea/gitea:latest`
+ rootless版本,比如`gitea/gitea:latest-rootless`

这两者之间使用上基本没什么区别,仅仅是设置上有如下区别,一般我们需要根据实际情况选版本.

| 项目                  | 原始版本                   | rootless版本                        |
| --------------------- | -------------------------- | ----------------------------------- |
| 适合的部署用户        | root                       | 非root                              |
| 执行用户              | `1000`                     | 随意                                |
| 执行用户组            | `1000`                     | 随意                                |
| ssh端口               | `22`                       | `2222`                              |
| 默认data目录          | `/data/`                   | `/var/lib/gitea`                    |
| 默认服务目录          | `/data/gitea`              | `/var/lib/gitea/custom`             |
| 默认配置位置          | `/data/gitea/conf/app.ini` | `/etc/gitea/app.ini`                |
| 默认的git仓库存放位置 | `/data/git/repositories`   | `/var/lib/gitea/gitea-repositories` |

rootless在挂在nfs时会相对比较方便,可以避免用户权限设置造成的可访问性问题,如果不用nfs且可以用root用户部署,建议还是原始版本部署,因为挂载volumes比较方便只要挂个`/data`即可.

本文也将以原始版本为例介绍部署方式,如果需要其他的部署方式可以看[官方文档](https://docs.gitea.com/category/installation)

部署步骤基本可以总结为

1. 在你的宿主机上找个合适的地方创建一个文件夹用于映射gitea的`/data`目录,我们以`/volume2/docker/gitea/data`为例

2. 单机模式部署`docker-compose.yml`

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

个人推荐要么直接将Gitea部署在内网不要暴露到外网,然后外网通过vpn进入内网后使用

### 首次登录的基本设置

第一次登陆后我们需要设置一些基本配置.之后就可以正常使用了

首次登陆进Gitea服务后界面如下:

<!-- todo -->

我们需要进行初始化配置,其中主要是:

+ 是设置第一个用户(管理员用户).当设置完成提交后就可以正常使用了.
+ 元数据使用的数据库.建议pg.
+ email设置,我们需要配置一个可以发送email的邮箱用于执行注册确认,密码找回等操作
+ 设置访问域名,一般我们内网用的话需要软路由(AdGuardHome)或nas套件进行相应设置,在内网中构造域名解析,然后将需要的域名设置进去

### 进一步设置

gitea的配置项有哪些可以查看[官方文档](https://docs.gitea.com/zh-cn/administration/config-cheat-sheet),需要注意随着版本的迁移
修改配置文件,主要是要修改如下几项:

+ `server`,主要是设置
    + `PROTOCOL` http还是https
    + `CERT_FILE`(如果是https就需要设置,可以创建一个文件夹`/volume2/docker/gitea/data/keys`后将证书放在其中)
    + `KEY_FILE`(如果是https就需要设置,可以创建一个文件夹`/volume2/docker/gitea/data/keys`后将私钥放在其中)
    + `DOMAIN`服务器域名,如果没有就填ip地址
    + `SSH_DOMAIN`SSH的服务器域名,如果没有就填ip地址
    + `ROOT_URL`,服务对外的URL
    + `LFS_START_SERVER`设置是否可以传输大文件

+ `database`,设置元数据的存放
+ `mailer`,设置发送消息的发件箱
+ `lfs`,设置大文件的传输存储
+ `indexer`,设置索引,gitea默认使用内部的`bleve`索引,也支持`elasticsearch`,可以用`ISSUE_INDEXER_TYPE`和`REPO_INDEXER_TYPE`分别设置工单和代码仓库的索引类型
+ `queue`,设置队列,gitea中的任务都是在队列中顺序执行的,默认使用`level`(即`LevelDB`）,也可以使用`redis`以增加扩展性,可以通过`TYPE`子项进行设置
+ `cache`,设置缓存,默认`memory`(即用内存作为缓存),也支持`redis`和`redis-cluster`,可以通过`ADAPTER`子项进行设置
+ `session`,设置会话缓存,默认`memory`,也支持`redis`和`redis-cluster`,可以通过`PROVIDER`子项进行设置
+ `packages`,设置包注册表配置
+ `storage`,设置存储配置,所谓的存储实际是静态存储,基本是放冷数据或温数据的.默认使用`local`方式存储即直接保存在硬盘上,也可以使用`minio`作为存储后端,可以通过`STORAGE_TYPE`子项进行设置.默认情况下如果使用`minio`会被保存在`MINIO_BUCKET`指定的桶中.每个存储项都有其默认的基本路径:

    | 存储项                                             | minio中的路径                                  | local中的路径                    |
    | -------------------------------------------------- | ---------------------------------------------- | -------------------------------- |
    | `attachments` (附件,一般是issue中带的)             | `attachments/`                                 | `/data/gitea/attachments/`       |
    | `lfs`  (大文件)                                    | `lfs/`                                         | `/data/git/lfs/`                 |
    | `avatars`                                          | `avatars/`                                     | `/data/gitea/avatars/`           |
    | `repo-avatars`                                     | `repo-avatars/`                                | `/data/gitea/repo-avatars/`      |
    | `repo-archive` (仓库release等操作保存的切片压缩包) | `repo-archive/`                                | `/data/gitea/repo-archive/`      |
    | `packages`                                         | `packages/`   (仓库发布的制品包)               | `/data/gitea/packages/`          |
    | `actions_log`                                      | `actions_log/`                                 | `/data/gitea/actions_log/`       |
    | `actions_artifacts`                                | `actions_artifacts/`(仓库action操作生成的制品) | `/data/gitea/actions_artifacts/` |

+ `actions`,设置actions的行为

我们也可以直接在部署的compose文件中用环境变量指定设置,环境变量的命名规则就是`GITEA__{分类项目}__{子项}`

```yml
version: '2.4'
services:
gitea:
    ...
    environment:
      - GITEA__database__DB_TYPE=postgres
      ...
    ...
```

### 升级版本

用docker部署的一大优势就是升版本非常简单,直接拉取最新的镜像重启即可.有时候会有兼容性问题,但一般都不是破坏性的,影响不大.

### 归档

命令行工具[gitea dump](https://docs.gitea.com/zh-cn/administration/command-line#dump)可以用于归档当前gitea,它会将所有存储,仓库,数据库中的元数据等的当前状态全部打包到一个命名为`gitea-dump-时间戳.zip`的zip中.其中包含如下内容：

+ `app.ini`配置文件
+ `custom/`所有保存在`custom/`目录下的配置和自定义的文件。
+ `data/`,数据目录,该目录下包括`attachments`,`avatars`,`lfs`,`indexers`,如果使用`SQLite`作为数据库则也包括`SQLite`文件
+ `repos/`,仓库目录的完整副本
+ `gitea-db.sql`,数据库dump出来的SQL,
+ `log/`,Logs文件,如果用作迁移不是必须的

需要注意由于我们是docker部署,执行命令行会麻烦些,可以用`docker exec`命令实现

#### 还原

归档当然是为了能还原了,对于数据库之外的部分,还原参照上面的两个路径表即可.

对于数据库,`SQLite`拷到设置中的位置即可,而mysql和pg则需要将`gitea-db.sql`中的数据转存到对应数据库中

```bash
# mysql
mysqldump -u$USER -p$PASS --database $DATABASE > gitea-db.sql
# postgres
pg_dump -U $USER $DATABASE > gitea-db.sql
```

### 迁移数据

另一种需求是水平扩展时迁移数据.对于数据库的部分我们可以利用归档还原功能,而对于其他部分,如果我们要将存储转到minio,我们可以利用命令行工具`gitea migrate-storage`

其具体用法是

```bash
gitea migrate-storage \
    -t <migrate_type(`attachments', 'lfs', 'avatars', 'repo-avatars', 'repo-archivers', 'packages', 'actions-log', 'actions-artifacts')> \ # 指定要迁移的项目
    --minio-base-path <minio_bucket_path> \ # 比如lfs就是lfs/
    --storage minio \
    --minio-endpoint <minio_url> \
    --minio-access-key-id <minio_key_id> \
    --minio-secret-access-key <minio_key> \
    --minio-bucket <minio_bucket> \
    --minio-use-ssl true \
    --minio-insecure-skip-verify true
```

## 使用部分

gitea一般是作为内部库或外部库的镜像库而存在的.在用法上也没什么特别之处.

### 作为备份库

Gitea的一大作用就是作为github的备份库,Gitea直接提供了对应的入口--左上角的`+`

![备份github][3]

填好提交就可以了

### 作为私有代码仓库

登陆进Gitea服务后界面如下:

![Gitea的个人页][2]

单纯作为代码仓库来说我们像用github一样用它即可.接下来的部分更多的是介绍它和github对标的各种实用功能

### 项目管理工具

gitea提供了和github几乎完全一致的项目管理体验

#### 成员管理

用户一旦创建了他就拥有了一个私人的仓库空间,在这个私人仓库空间中他是其中仓库的唯一拥有者,其他人则都仅仅有读权限.如果你想找帮手,你可以邀请他们称为`协作者`,而每个协作者的权限则是需要你设置的,可以设置读,写,管理三种权限.而有写权限的协作者也可以使用`分支保护`功能控制他对不同分支的写入权限.

而更多的时候我们会组成`组织`以和一些人一起完成一些项目.对于`组织`来说,成员需要被划分为`团队`,`团队`的作用类似linux中的`group`用户组,一般是按功能划分的.同一个团队中的成员有相同的基础权限.

成员则有三种角色



+ 所有者团队
创建组织时将自动创建所有者团队，创建者将成为所有者团队的第一名成员。所有者团队不可删除，且至少有一名成员。

管理员团队
创建团队时，有两种类型的团队。一种是管理员团队，另一种是普通团队。可以创建一个管理员团队来管理某些版本库，其成员可以对这些版本库做任何事情。只有所有者或管理员团队的成员才能创建新团队。

普通团队

我们可以添加协作者.


gitea支持用组织的形式组织仓库,组织中可以包含固定成员,这些成员可以设定对其中仓库的角色



#### 
+ 在人员组织方面,Gitea支持team也支持组织.
+ 在事项交流方面提供了工单功能,工单还支持点赞,支持上传附件
+ 

+ 在项目管理方面,Gitea支持`pull request`,支持对`pull request`的内容做`code review`,也支持工单,标签,里程碑和kanban.

+ 在交流方面Gitea也提供了工单功能,甚至还支持点赞...

### CI/CD

+ 事件钩子
+ actions

### 分发仓库

+ 在成果分发方面,Gitea也支持`release`.
+ package


### 文档工具

在文档方面,Gitea可以使用`README`和`Wiki`管理项目文档;







[2]: {{site.url}}/img/in-post/gitea/mainpage.PNG
[3]: {{site.url}}/img/in-post/gitea/qianyi.PNG