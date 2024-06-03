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

### action runnner

如果我们打算让用户使用gitea的actions功能做ci/cd,还需要额外安装配置action runnner.我们可以使用镜像[gitea/act_runner:latest](https://hub.docker.com/r/gitea/act_runner).我们可以遵循如下操作来实现部署

> 在gitea中

1. 决定runner级别,runner分为3个级别

    + 实例级别：Runner将为实例中的所有存储库运行Job
    + 组织级别：Runner将为组织中的所有存储库运行Job
    + 存储库级别：Runner将为其所属的存储库运行Job

    请注意即使存储库具有自己的存储库级别Runner,它仍然可以使用实例级别或组织级别Runner.

2. 获取令牌,注册令牌的格式是一个形如`D0gvfu2iHfUjNqCYVljVyRV14fISpJxxxxxxxxxx`的随机字符串.runner级别决定了从哪里获取注册令牌:

    + 实例级别: 管理员设置页面,例如`<your_gitea.com>/admin/actions/runners`
    + 组织级别: 组织设置页面,例如`<your_gitea.com>/<org>/settings/actions/runners`
    + 存储库级别: 存储库设置页面,例如`<your_gitea.com>/<owner>/<repo>/settings/actions/runners`
    如果无法看到设置页面,请确保使用的用户具有正确的权限并且gitea已启用Actions.

> 在本机

1. 先拉下来镜像(最好指定版本)

    ```bash
    docker pull gitea/act_runner:latest
    ```

2. 生成默认配置文件

    ```bash
    docker run --entrypoint="" --rm -it gitea/act_runner:latest act_runner generate-config > config.yaml
    ```

    拿到的默认配置大致长这样

    ```yml
    # Example configuration file, it's safe to copy this as the default config file without any modification.

    # You don't have to copy this file to your instance,
    # just run `./act_runner generate-config > config.yaml` to generate a config file.

    log:
    # The level of logging, can be trace, debug, info, warn, error, fatal
    level: info

    runner:
    # Where to store the registration result.
    file: .runner
    # Execute how many tasks concurrently at the same time.
    capacity: 1
    # Extra environment variables to run jobs.
    envs:
        A_TEST_ENV_NAME_1: a_test_env_value_1
        A_TEST_ENV_NAME_2: a_test_env_value_2
    # Extra environment variables to run jobs from a file.
    # It will be ignored if it's empty or the file doesn't exist.
    env_file: .env
    # The timeout for a job to be finished.
    # Please note that the Gitea instance also has a timeout (3h by default) for the job.
    # So the job could be stopped by the Gitea instance if it's timeout is shorter than this.
    timeout: 3h
    # Whether skip verifying the TLS certificate of the Gitea instance.
    insecure: false
    # The timeout for fetching the job from the Gitea instance.
    fetch_timeout: 5s
    # The interval for fetching the job from the Gitea instance.
    fetch_interval: 2s
    # The labels of a runner are used to determine which jobs the runner can run, and how to run them.
    # Like: "macos-arm64:host" or "ubuntu-latest:docker://gitea/runner-images:ubuntu-latest"
    # Find more images provided by Gitea at https://gitea.com/gitea/runner-images .
    # If it's empty when registering, it will ask for inputting labels.
    # If it's empty when execute `daemon`, will use labels in `.runner` file.
    labels:
        - "ubuntu-latest:docker://gitea/runner-images:ubuntu-latest"
        - "ubuntu-22.04:docker://gitea/runner-images:ubuntu-22.04"
        - "ubuntu-20.04:docker://gitea/runner-images:ubuntu-20.04"

    cache:
    # Enable cache server to use actions/cache.
    enabled: true
    # The directory to store the cache data.
    # If it's empty, the cache data will be stored in $HOME/.cache/actcache.
    dir: ""
    # The host of the cache server.
    # It's not for the address to listen, but the address to connect from job containers.
    # So 0.0.0.0 is a bad choice, leave it empty to detect automatically.
    host: ""
    # The port of the cache server.
    # 0 means to use a random available port.
    port: 0
    # The external cache server URL. Valid only when enable is true.
    # If it's specified, act_runner will use this URL as the ACTIONS_CACHE_URL rather than start a server by itself.
    # The URL should generally end with "/".
    external_server: ""

    container:
    # Specifies the network to which the container will connect.
    # Could be host, bridge or the name of a custom network.
    # If it's empty, act_runner will create a network automatically.
    network: ""
    # Whether to use privileged mode or not when launching task containers (privileged mode is required for Docker-in-Docker).
    privileged: false
    # And other options to be used when the container is started (eg, --add-host=my.gitea.url:host-gateway).
    options:
    # The parent directory of a job's working directory.
    # NOTE: There is no need to add the first '/' of the path as act_runner will add it automatically. 
    # If the path starts with '/', the '/' will be trimmed.
    # For example, if the parent directory is /path/to/my/dir, workdir_parent should be path/to/my/dir
    # If it's empty, /workspace will be used.
    workdir_parent:
    # Volumes (including bind mounts) can be mounted to containers. Glob syntax is supported, see https://github.com/gobwas/glob
    # You can specify multiple volumes. If the sequence is empty, no volumes can be mounted.
    # For example, if you only allow containers to mount the `data` volume and all the json files in `/src`, you should change the config to:
    # valid_volumes:
    #   - data
    #   - /src/*.json
    # If you want to allow any volume, please use the following configuration:
    # valid_volumes:
    #   - '**'
    valid_volumes: []
    # overrides the docker client host with the specified one.
    # If it's empty, act_runner will find an available docker host automatically.
    # If it's "-", act_runner will find an available docker host automatically, but the docker host won't be mounted to the job containers and service containers.
    # If it's not empty or "-", the specified docker host will be used. An error will be returned if it doesn't work.
    docker_host: ""
    # Pull docker image(s) even if already present
    force_pull: true
    # Rebuild docker image(s) even if already present
    force_rebuild: false

    host:
    # The parent directory of a job's working directory.
    # If it's empty, $HOME/.cache/act/ will be used.
    workdir_parent:

    ```

    需要注意`container.network`需要填上你`gitea`实例所在的网络名

3. 修改配置文件

    我们使用的是docker方式部署,这样就需要额外配置`cache`

    ```yml
    cache:
        # Enable cache server to use actions/cache.
        enabled: true
        # The directory to store the cache data.
        # If it's empty, the cache data will be stored in $HOME/.cache/actcache.
        dir: ""
        # The host of the cache server.
        # It's not for the address to listen, but the address to connect from job containers.
        # So 0.0.0.0 is a bad choice, leave it empty to detect automatically.
        host: "192.168.50.92"
        # The port of the cache server.
        # 0 means to use a random available port.
        port: 8088
        # The external cache server URL. Valid only when enable is true.
        # If it's specified, act_runner will use this URL as the ACTIONS_CACHE_URL rather than start a server by itself.
        # The URL should generally end with "/".
        external_server: ""
    ```

> 在目标机器

1. 先拉下来镜像(最好指定版本)

    ```bash
    docker pull gitea/act_runner:latest
    ```

2. 将本机写好的配置文件放到合适的位置

    ```bash
    scp config.yaml user@remote:/你的/配置/位置/config.yaml
    ```

3. 编辑启动的compose文件

    ```yaml

    services:
        ...
        runner:
            image: gitea/act_runner:nightly
            environment:
                CONFIG_FILE: /config.yaml
                GITEA_INSTANCE_URL: "${INSTANCE_URL}"
                GITEA_RUNNER_REGISTRATION_TOKEN: "${REGISTRATION_TOKEN}"
                GITEA_RUNNER_NAME: "${RUNNER_NAME}"
                GITEA_RUNNER_LABELS: "${RUNNER_LABELS}"
            ports:
                - "8088:8088"
            volumes:
                - ./config.yaml:/config.yaml
                - ./data:/data
                - /var/run/docker.sock:/var/run/docker.sock
    ```

注意如果你的gitea有dns解析有自己的域名,`GITEA_INSTANCE_URL`中应该使用该域名.

以组织级别为例,最后配好后就是这样

![runner状态][1]

要让用户可以比较安心的使用action功能,我们必须清晰的了解仓库的运行环境.

1. 如果我们的仓库有境外网站访问环境

    我们就什么都不需要额外设置直接使用即可.gitea的默认配置会从dockerhub下载镜像,从github下载actions使用,在外网使用没什么问题的情况下默认配置就可以了

2. 如果我们的仓库没有境外网站访问环境或希望尽量少的访问境外网站

    我们可以
    1. 在gitea设置项中设置

        ```conf
        ...
        [actions]
        ENABLED: true
        DEFAULT_ACTIONS_URL: self 
        ```

        其中`DEFAULT_ACTIONS_URL`用于设置action的默认查找位置,默认是`github`,这里改为`self`则会在本地gitea中查找

    2. 在gitea中新建一个组织`actions`,里面用*镜像仓库*方式将自己常用的action放到本地.为了可以拉取镜像,我们也可以在gitea实例的容器中设置代理,以docker compose方式为例

        ```yml
        version: '2.4'
        services:
        gitea:
            ...
            environment:
                https_proxy: "http://你的代理host:你的代理port"
                http_proxy: "http://你的代理host:你的代理port"
            ...
        ```

    3. 通过`docker save`和`docker load`命令将用到的镜像都放到runnner所在机器上,目前可以使用的镜像可以查看[这个列表](https://github.com/nektos/act/blob/master/IMAGES.md)

    这样action就不会去连github和dockerhub了.

## 使用部分

gitea一般是作为内部库或外部库的镜像库而存在的.在用法上也没什么特别之处.

### 作为备份库

Gitea的一大作用就是作为github的备份库,Gitea直接提供了对应的入口--左上角的`+`

![备份github][2]

填好提交就可以了

### 作为私有代码仓库

登陆进Gitea服务后界面如下:

![Gitea的个人页][3]

单纯作为代码仓库来说我们像用github一样用它即可.接下来的部分更多的是介绍它和github对标的各种实用功能

### 文档工具

在文档方面,Gitea默认可以使用`README`和`Wiki`管理项目文档.但目前并没有官方的`page`服务.很遗憾,但作为内部仓库,其实这样也够用了.

通常我们应该确保`README.md`中描述清楚了仓库的功能和目标以及一个最简单的安装和使用示例,而`Wiki`中我们应该尽量详尽的包含如下内容

1. 一个目录
2. 详细安装说明,包括安装包和源码编译安装,外部依赖等,确保按步骤走可以跑起来
3. 启动的配置项说明,包括配置的设置方式以及各个配置项的含义和取值范围等
4. 各功能特性的详细使用说明
5. 命令行说明
6. 网络API说明,如果项目已经有了openapi的描述json或yaml文件,可以通过[openapi-to-md](https://www.npmjs.com/package/openapi-to-md)直接工具生成,然后略微修改即可.如果项目是grpc的可以把protobuf文件贴上来
7. 二次开发的帮助文档,包括代码架构和一些实现上的生命周示意图等

### CICD

gitea的CICD有两种实现方式,一种是通过`Web Hook`触发外部CI/CD工具,这个可以看我的[<使用Jenkins代替GithubActions自动化工作流>这篇文章](https://blog.hszofficial.site/recommend/2020/12/02/%E4%BD%BF%E7%94%A8Jenkins%E4%BB%A3%E6%9B%BFGithubActions%E8%87%AA%E5%8A%A8%E5%8C%96%E5%B7%A5%E4%BD%9C%E6%B5%81/).另一种就是使用[gitea actions](https://docs.gitea.com/zh-cn/usage/actions/overview).这个功能是gitea社区跟进github actions做出来的,因此语法和github actions一样,使用上不同之处就是我们得自己部署runner.如何部署runner可以查看上面运维部分的相关内容.

在使用上gitea兼容[github action的配置语法](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idruns-on),得益于这一点,我们可以几乎无缝的从github action切过来.当然了也不会是完全没有不同之处.不同点可以看[这个文档总结](https://docs.gitea.com/zh-cn/usage/actions/comparison).比较重要的点包括:

1. `gitea action`的定义文件在`.gitea/workflows/`目录下而不是`.github/workflows/`下
2. `gitea`中`uses`可以直接指定仓库.这一特性是作为`actions.DEFAULT_ACTIONS_URL`配置项的补充,例如你可以使用`uses: https://gitea.com/actions/checkout@v4`的写法来代替`uses: actions/checkout@v4`的写法,从而使用`gitea.com`中的action.
3. 不支持手动触发(即`workflow_dispatch`字段定义),也没有ui操作手动触发

### 项目管理工具

gitea提供了和github几乎完全一致的项目管理体验

#### 成员管理

成员管理方面,gitea使用的是用户权限+分支保护的方案.

用户对于仓库的权限只有三种

+ `读`,只能访问
+ `写`,读权限之外还可以改变仓库内容
+ `管理员`,写权限外还可以管理人员等

可以设置权限的功能项包括

+ 代码
+ 工单
+ 合并请求
+ 版本发布
+ 百科
+ 访问外部百科
+ 访问外部工单
+ 项目
+ 软件包
+ Actions

而分支保护则是在分支层面控制不同用户修改仓库的更精细的设置.

##### 私人仓库成员管理

用户一旦创建了他就拥有了一个私人的仓库空间,在这个私人仓库空间中他是其中仓库的唯一拥有者,其他人则都仅仅有读权限.如果你想找帮手,你可以在`设置->协作者`中邀请他们成为`协作者`,每个协作者的权限则是需要你设置的.

![Gitea的个人页协作者][4]

而有写权限的协作者也可以使用`分支保护`功能控制他对不同分支的写入权限.

+ 分支管理页
![Gitea的个人页分支保护][5]

+ 分支保护设置

![分支保护设置][6]

通常来说,我们会希望主干分支无法被协作者修改,协作者只能通过`pull request`的方式提交特性到主干.

不过一般来说私人仓库都是个人项目,很少会需要协作者.

##### 组织成员管理

而更多的时候我们会组成`组织`以和一些人一起完成一些项目.对于`组织`来说,成员需要被划分为`团队`,`团队`的作用类似linux中的`group`用户组,一般是按功能划分的.同一个`团队`中的成员有相同的基础权限.而添加新的成员也必须从`团队`中进行添加

+ 团队管理界面

![组织团队][7]

+ 创建团队

![创建团队][8]

在组织中团队有三种类型

+ `所有者团队`,创建组织时将自动创建的团队,单例,就叫`所有者团队`,创建者将成为所有者团队的第一名成员.所有者团队不可删除且至少有一名成员.`所有者团队`中的成员除了拥有对组织中所有仓库的读写管理员操作权限外还对组织本身具有管理员权限.
+ `管理员团队`,需要手动创建,一般用于管理一些组织中的指定仓库,其成员可以对这些指定的仓库做任何事情.我们可以在创建管理员团队时选择是否允许创建新的仓库,也可以直接赋予对组织中所有仓库的管理员权限.如果需要手动指定可以操作的仓库,则需要进入具体的仓库中在`设置->协作者->团队`中添加.
+ `普通团队`,普通团队一般是有相同权限设置的一组用户,和上面`管理团队`一样也可以设置是否对全部仓库还是指定仓库有设置的权限,我们需要在创建时就将权限设置好,比如给实习生设置为全部只读权限且需要指定仓库,给普通开发者设置为只有`代码`,`工单`,`合并请求`为写权限其他都是读权限等.

三种团队的成员中只有所有者或管理员团队的成员才能创建新的团队.一个组织只可以维护自己内部的团队,如果你希望引入组织成员外部的人员来参与个别项目,也可以像在私人仓库中一样通过引入外部的`协作者`来实现.而对于分支的精细控制,我们一样是使用`分支保护`功能控制,只是其中增加可以用团队来控制的条件.

#### 项目模版管理

项目模版可以用来快速构建同一模式的代码仓库,比如用go写命令行工具怎么都会有个`go.mod`,怎么都会依赖一个命令行环境变量或从文件中读取参数的包.那就可以把这些共性的东西固定好作为项目模版管理起来,这样每次要写个新的go命令行工具时就可以快速开始.可以理解为是一个代码仓库级别的脚手架工具.

<!-- https://docs.gitea.com/zh-cn/usage/template-repositories -->

我们可以在一个公共可访问的组织中统一管理项目模版(此处我们假设该组织为`share`).创建项目模版本质上也是创建一个项目,只是在创建时需要勾选最下方的`模版 设置仓库为模版仓库`即可.

![创建模版仓库][9]

之后在其他项目中要使用该模版时,在创建页面中选择需要的模版即可

![使用模版仓库][10]

模版仓库可以看做是最佳实践的一种固化,一般由资深开发人员来进行.

模版仓库与一般仓库不同点在于:

1. 内容文件中可以使用`${VAR}`(或`$VAR`)的形式导入变量占位使用或使用`${VAR_转换器}`(或`$VAR_转换器`)的形式将变量经过转换器处理后使用,目前支持的变量包括

    | 变量                   | 含义                           | 是否可转换 |
    | ---------------------- | ------------------------------ | ---------- |
    | `REPO_NAME`            | 生成的仓库名称                 | ✓          |
    | `REPO_DESCRIPTION`     | 生成的仓库描述                 | ✘          |
    | `REPO_OWNER`           | 生成的仓库所有者               | ✓          |
    | `REPO_LINK`            | 生成的仓库链接(不包含hostname) | ✘          |
    | `REPO_HTTPS_URL`       | 生成的仓库的HTTP(S)克隆链接    | ✘          |
    | `REPO_SSH_URL`         | 生成的仓库的SSH克隆链接        | ✘          |
    | `TEMPLATE_NAME`        | 模板仓库名称                   | ✓          |
    | `TEMPLATE_DESCRIPTION` | 模板仓库描述                   | ✘          |
    | `TEMPLATE_OWNER`       | 模板仓库所有者                 | ✓          |
    | `TEMPLATE_LINK`        | 模板仓库链接(不包含hostname)   | ✘          |
    | `TEMPLATE_HTTPS_URL`   | 模板仓库的HTTP(S)克隆链接      | ✘          |
    | `TEMPLATE_SSH_URL`     | 模板仓库的SSH克隆链接          | ✘          |

    支持的转换器包括如下,我们以`go-sdk`作为输入的例子下面是转换器的效果

    | 转换器   | 目标形式 |
    | -------- | -------- |
    | `SNAKE`  | `go_sdk` |
    | `KEBAB`  | `go-sdk` |
    | `CAMEL`  | `goSdk`  |
    | `PASCAL` | `GoSdk`  |
    | `LOWER`  | `go-sdk` |
    | `UPPER`  | `GO-SDK` |
    | `TITLE`  | `Go-Sdk` |

    如果希望保留原始字面量不进行扩展,则使用`$${xxx}`或`$$xxx`的形式

2. 根目录下需要有一个`.gitea/template`用于标识需要进行变量插入和转换的文件,其规则和`.gitignore`一致.

    ```txt
    README.md
    pyproject.toml
    **_SNAKE
    **.py
    ```

    需要注意
    1. 文件名文件夹名中有变量需要扩展的也需要在`.gitea/template`中匹配到的才能执行.
    2. `REPO_LINK`和`TEMPLATE_LINK`仅为`path`部分,而`REPO_HTTPS_URL`和`TEMPLATE_HTTPS_URL`为完整url,但末尾会有`.git`

#### 工单

在项目管理方面gitea和github一样都使用的是基于工单的项目管理方案.工单本身可以看做是事项的最小单位,是各种项目管理工具的轴.

![工单列表页][11]

工单功能本身就功能丰富,可以打标签,可以上传附件,可以点赞,可以指派处理的成员,还可以关联其他工单,`pull request`,milestone,以及project

![工单创建页][12]

工单的创建者可以随时自己编辑工单信息,也可以关闭工单

![工单管理页][13]

##### 工单标签

gitea中的工单也可以用标签分类.我们可以在标签管理页中对标签进行管理

![工单标签管理页][14]

##### 工单模版

仓库的工单模版的表现形式有两种:

+ markdown模版.就是最基础的模版形式,模本文件写什么用它做模版的工单就写什么
+ yaml模版.表单问卷形式的模版.当用户使用该模版构造工单时,模版会根据yaml中描述的项目渲染一个表单,用户填好表单提交即完成了工单的提交工作

仓库的工单模版在仓库根目录下的`.gitea`文件夹下维护,分为两种:

+ 默认工单模版.使用`.gitea/issue_template.md`或`.gitea/issue_template.yaml`文件描述.
+ 其他工单模版.放在`.gitea/issue_template/`目录下,同样可以是`.md`文件也可以是`.yaml`文件.

我们也可以借助模版仓库将工单模版固定下来,这样相同类型的任务就可以有相同的工单形式,方便管理.

#### pull request模版
<!-- todo https://docs.gitea.com/zh-cn/usage/issue-pull-request-templates-->

<!-- https://axolo.co/blog/p/part-3-github-pull-request-template -->

#### 质量管理

`pull request`则是代码仓库的质量管理工具,它支持评论,`code review`.

<!-- todo https://docs.gitea.com/zh-cn/usage/issue-pull-request-templates-->

充分使用`pull request`配合CICD工具以及[状态检查功能](https://github.com/go-gitea/gitea/pull/7481)可以严格控制代码质量避免上线出问题.

通常我们会视项目的性质,组织的规模来设置代码质量控制的方式


##### 状态检查功能
<!-- todo https://docs.gitea.com/zh-cn/api/1.22/#tag/repository/operation/repoListStatuses-->

##### 质量管理的严格程度

理论上我们所有的项目在任何时候都应该进行严格的质量管理,但实际生产中很多时候是做不到的.难点在于

+ 工期不允许,就好像做卷子,同样的卷子1个小时做出来和2个小时做出来当然不一样,而生产中往往会在功能需求不变甚至增加的情况下压缩工期,那相当于卷子要优先都填上而不是优先作对,自然就会放宽质量管理
+ 人力不够,质量管理必然要求交叉验证,这就需要一个质量管理委员会来把控和组织code review,而code review本身也占用工作时间耗费人力,自然会人力不够

下面是我总结的不同情况下质量管理的严格程度和对应方案

| 人力规模               | 情况                           | 严格程度 | `code review`人数 | 其他限制                                            |
| ---------------------- | ------------------------------ | -------- | ----------------- | --------------------------------------------------- |
| 3人以下(家用级别)      | 所有情况的仓库                 | 0档      | `0`               | 检查提交后CICD工具的静态检查和单元测试结果          |
| 10人以下(小工作室级别) | 涉及记账等难以回退的业务的仓库 | 1档      | `1`               | 检查提交后CICD工具的静态检查和单元测试结果          |
| 10人以下(小工作室级别) | 其他业务的仓库                 | 0档      | `0`               | 检查提交后CICD工具的静态检查和单元测试结果          |
| 50人以下(小微企业级别) | 涉及记账等难以回退的业务       | 2档      | `3`               | 检查提交后CICD工具的静态检查,单元测试,回归测试结果  |
| 50人以下(小微企业级别) | 涉及关键基础设施的仓库         | 1档      | `1`               | 检查提交后CICD工具的静态检查和单元测试结果          |
| 50人以下(小微企业级别) | 其他的仓库                     | 0档      | `0`               | 检查提交后CICD工具的静态检查和单元测试结果          |
| 50人以上(集团级别)     | 涉及记账等难以回退的业务       | 3档      | `5`               | 检查提交后CICD工具的静态检查和单元测试,回归测试结果 |
| 50人以上(集团级别)     | 涉及关键基础设施的仓库         | 2档      | `3`               | 检查提交后CICD工具的静态检查和单元测试,回归测试结果 |
| 50人以上(集团级别)     | 其他的仓库                     | 1档      | `1`               | 检查提交后CICD工具的静态检查和单元测试,回归测试结果 |

#### 进度管理

gitea和github一样提供了milestone和project两种进度管理工具.这两者虽然在实现上都是依托于工单系统的,但在设计和使用思路上有很大区别.

+ `里程碑方式`,是[Scrum方法论](https://baike.baidu.com/item/Scrum/1698901?fr=ge_ala)的一个工具.适用于时间限制内迭代交付的敏捷工作流
+ `项目方式`,是[看板方法论](https://baijiahao.baidu.com/s?id=1768363734089014454&wfr=spider&for=pc)的一个工具.`项目`有利于持续交付和稳定的工作流程

通常这两种进度管理工具对立统一--我们通常用`里程碑方式`管理特定仓库中代码的迭代,用`项目方式`管理产品层面功能的迭代以及优化工作流

##### Milestone

##### Project




### 分发仓库

+ 在成果分发方面,Gitea也支持`release`.
+ package




[1]: {{site.url}}/img/in-post/gitea/gitea_runner.png
[2]: {{site.url}}/img/in-post/gitea/qianyi.PNG
[3]: {{site.url}}/img/in-post/gitea/mainpage.PNG
[4]: {{site.url}}/img/in-post/gitea/user_collaborator.png
[5]: {{site.url}}/img/in-post/gitea/user_branch.png
[6]: {{site.url}}/img/in-post/gitea/branch_protect.png
[7]: {{site.url}}/img/in-post/gitea/org_group.png
[8]: {{site.url}}/img/in-post/gitea/org_group_create.png
[9]: {{site.url}}/img/in-post/gitea/create_template_repo.png
[10]: {{site.url}}/img/in-post/gitea/use_template_repo.png
[11]: {{site.url}}/img/in-post/gitea/issue_list.png
[12]: {{site.url}}/img/in-post/gitea/issue_create.png
[13]: {{site.url}}/img/in-post/gitea/issue_manager.png
[14]: {{site.url}}/img/in-post/gitea/issue_tag.png