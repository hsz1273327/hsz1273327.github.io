---
title: "使用GithubActions自动化工作流"
date: 2020-11-30
author: "Hsz"
category: introduce
tags:
    - Github
    - DevOps
header-img: "img/home-bg-o.jpg"
update: 2020-11-30
series:
    get_along_well_with_github:
        index: 8
---
# 使用GithubActions自动化工作流

Github在2019年底开放了内置的CI/CD工具`GithubActions`.这样使用Github托管的代码终于有了不借助外部服务自动化测试打包部署的能力.
同时由于后发优势,`GithubActions`几乎是目前最易用的CI/CD工具.
<!--more-->

`GithubActions`类似于传统的CI/CD工具,都是使用代码配置脚本,执行器执行脚本,页面管理执行过程的结构.

+ 在代码中配置脚本放在根目录的`.github/workflow`文件夹下,使用`yaml`格式描述配置.
+ Github默认给每个用户配置3个的执行器,我们也可以自己创建`self-host`执行器
+ 每个代码仓库的顶部标签页都有专门的`actions`按钮,进去就是当前仓库的执行过程管理页面.

对`GithubActions`的详细描述可以看[官方文档](https://docs.github.com/cn/free-pro-team@latest/actions/reference)本文只是介绍和划重点.

## 术语和约束

在介绍使用方式之前我们先来了解下`GithubActions`的术语,借此了解下一次执行过程的流程.

+ `workflow`即工作流,一次执行过程.每个workflow用一个配置文件维护.
+ `Job`: workflow的分解,可串行存在依赖;可并行
+ `Step`: job的分解,即步骤,比如一个step是要给代码做单元测试,那可能会有三个步骤:下载依赖->测试->上传结果
+ `action`workflow最小执行单元.即每个执行步骤中的具体执行任务,我们可以自己定义action,也可使用Github社区定义好的action
+ `Artifact`： workflow运行时产生的中间文件.包括日志,测试结果等
+ `Event`: 触发workflow的事件

Github Action对`workflow`设有如下使用限制:

+ 一个仓库可最多同时开20个workflows;超过20则排队等待

+ 一个workflow下的每个job最多运行6小时,超过直接结束

+ 所有分支下的job根据github级别不同有不同的并行度限制,超过并行度进入队列等待

+ 1小时内最多1000次执行请求,也就是1.5api/1m

需要注意,Github对Github Action服务有最终解释权,也就是说乱用可能会被Github限制账户.Github也会生成相关使用统计情况

## 配置CI/CD

配置CI/CD过程本质上是向runner描述如下内容:

+ 什么时候执行
+ 执行什么操作.

我们的workflow配置文件也一样是干这个的.

一个典型的workflow配置文件如下:

```yaml
name: Python package # 定义workflow的名字

# 描述何时执行
on: 
  push:
    branches: [master]
  pull_request:
    branches: [master]

# 描述workflow要做什么
jobs:
  build:
    runs-on: ubuntu-latest #描述执行的操作系统
    strategy:
      matrix: #参数矩阵,每一个元素都会被带入步骤执行
        python-version: [3.6, 3.7, 3.8, 3.9]

    #描述执行步骤
    steps:
      - uses: actions/checkout@v2
      - name: Set up Python \$\{\{ matrix.python-version \}\}
        uses: actions/setup-python@v2
        with:
          python-version: \$\{\{ matrix.python-version \}\}
      - name: Install devDependence
        run: |
          python -m pip install --upgrade pip
          pip install mypy pycodestyle coverage lxml
      - name: Install dependencies
        run: |
          if [ -f requirements.txt ]; then pip install -r requirements.txt; fi
      - name: Lint with pep8
        run: |
          pycodestyle --max-line-length=140 --ignore=E501 --first --statistics schema_entry

      - name: Type Hint Check
        run: |
          mypy --ignore-missing-imports --show-column-numbers --follow-imports=silent --check-untyped-defs --disallow-untyped-defs --no-implicit-optional --warn-unused-ignores schema_entry
      - name: Unit Test
        run: |
          python -m coverage run --source=schema_entry -m unittest discover -v -s . -p *test*.py
          python -m coverage report
```

### workflow的触发

每个workflow的配置文件都需要定义`on`字段,它用来描述在何种情况(`Event`)下触发执行.我们可以定义`on`多种事件,这样**只要满足其中一个就会被触发**

我们可以将`Event`分为3类:

+ 定时事件:由定时任务触发的事件

+ 手动触发事件: 在actions页面中手动触发的事件

+ Webhook事件:由github网站的钩子行为触发的事件,通常Git操作都有钩子可以用于触发

#### 定时事件

最简单的事件就是定时事件其定义方式如下:

```yaml
on:
  schedule:
    # * is a special character in YAML so you have to quote this string
    - cron:  '*/15 * * * *'
```

上面定义了一个每隔15分钟执行依次的任务.Github Avtion目前只支持[crontab语法定义定时任务](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/crontab.html#tag_20_25_07)

这个事件只会拉取默认分支(一般是`master`或者`main`分支,可以在`仓库的settings->branches->Default branch`下修改)的最近一次提交进行执行.

#### 手动触发事件

手动触发事件分为两种:

+ `workflow_dispatch` 让用在`Actions`界面中手动触发workflow
    当在workflow中定义了`workflow_dispatch`后管理页面就会允许指定这个`workflow`被手动执行,执行时默认需要指定分支,如果我们在配置中定义了参数,则手动执行时也会需要填参数.

    ![手动触发事件][2]

    一个典型例子如下:

    ```yaml
    on:
        workflow_dispatch:
            inputs:
                name:
                    description: 'Person to greet'
                    required: true
                    default: 'Mona the Octocat'
                home:
                    description: 'location'
                    required: false
                    default: 'The Octoverse'

    jobs:
        say_hello:
            runs-on: ubuntu-latest
            steps:
            - run: |
                echo "Hello \$\{\{ github.event.inputs.name \}\}!"
                echo "- in \$\{\{ github.event.inputs.home \}\}!"
    ```

    上面在`workflow_dispatch`下通过定义`inputs`设定参数.在`jobs`中我们则可以在`github.event.inputs`中取到对应的参数.
    **注意**如果不定义手动触发事件那么就无法手动触发.

+ `repository_dispatch`让用户通过API批量手动执行

    这个event的主要作用是让其他的程序通过api调用,通过自定义事件类型来驱动执行.这个event对应的workflow必须在默认分支下定义.

    比如我们定义:

    ```yaml
    on:
        repository_dispatch:
            types: [opened, deleted]
    ```

    然后执行http请求:

    ```bash
    curl \
    -X POST \
    -H "Accept: application/vnd.github.v3+json" \
    https://api.github.com/repos/{namespace}/{repo_name}/dispatches \
    -d '{"event_type":"opened"}'
    ```

    那么就可以被执行了.其中的`opened`, `deleted`是用户自定义的事件.

#### Webhook事件

Webhook事件是借由Github的webhook事件触发的事件,具体有哪些可以看[官方文档](https://docs.github.com/en/free-pro-team@latest/actions/reference/events-that-trigger-workflows#webhook-events),本文将只介绍几个常用的和git操作相关的事件.

+ `create`

    分支,tag创建时触发

+ `delete`

    分支,tag删除时触发

+ `gollum`

    仓库的wiki创建或者更新时触发

+ `push`/`pull_request`

    `push`是当有对仓库的push操作时触发;`pull_request`则是在执行`pull request`中触发

    这两个事件可以额外限制:
    + `branches: [...]`指定符合条件的分支触发
    + `branches-ignore:[...]`指定除符合条件的分之外都触发
    + `tags:[...]`指定符合条件的tag触发
    + `tags-ignore:[...]`指定除符合条件的tag外都触发
    + `paths:[...]`代码中有符合条件的路径就触发(至少有一个存在)
    + `paths-ignore:[...]`代码中不存在指定的路径则都触发(至少有一个不存在)

    上面的限制都允许使用通配符做匹配,支持的通配符包括:
    + `*`: 表示匹配0个或多个非`/`字符
    + `**`: 表示匹配0个或多个字符.
    + `?`: 表示匹配0个或者一个字符
    + `+`: 表示匹配至少一个字符
    + `[]`: 表示匹配一个范围内的字符,比如`[0-9a-f]`表示数字和a到f间的字符可以匹配
    + `!`: 在匹配字符串的开头表示否,其他位置没有特殊含义

    `pull_request`默认的行为是在`merge`完成后处理`merge`后的那次提交中的代码.
    我们还可以通过`types: [...]`字段指定细分事件类型,包括:
    + `assigned`被分派到某个issue时触发
    + `unassigned`删除分派时触发
    + `labeled`打标签时触发
    + `unlabeled`取消标签时触发
    + `opened`创建pull request时触发
    + `edited`编辑pull request时触发
    + `closed`关闭pull request时触发
    + `reopened`重新打开pull request时触发
    + `synchronize`同步pull request代码时触发
    + `ready_for_review`,pull request处于ready_for_review状态时触发
    + `locked`,锁定时触发
    + `unlocked`,解锁时触发
    + `review_requested`,code review结束时触发
    + `review_request_removed`code review请求被删除时触发

+ `release`

    当执行`Github release`时触发,使用的代码时release时打tag的代码,类似于`pull_request`,也可以通过`types:[...]`来指定细分事件.
    + `published`公开后执行
    + `unpublished`取消公开后执行
    + `created` 创建后执行
    + `edited` 编辑后执行
    + `deleted` 删除后执行
    + `prereleased`预发布后执行
    + `released`发布后执行

### 模板语法

我们可以看到上面例子中会有`\$\{\{ ... \}\}`这样的文字,这是`Github Action`定义的模板语法,其中`...`的部分可以是常数,上下文变量,运算符或者预定义的函数调用.

#### 常数

模板语法支持所有json支持的简单数据类型,也就是`null`,`boolean`,`number`,`string`.

#### 上下文变量

每次workflow执行都会带上几个上下文变量用于描述自己和传递参数.具体的可以看[官方文档](https://docs.github.com/en/free-pro-team@latest/actions/reference/context-and-expression-syntax-for-github-actions#contexts).这边只介绍几个常用的

+ `matrix`,执行策略中定义的变量,每次执行每个key只会有一个取值
+ `env`,workflow中`env`定义的变量
+ `github`,通常用于获取仓库和分支的信息,比较值得关注的有:
    + `github.repository` 执行的仓库名,也就是`{namespace}/{repo_name}`,如果只要repo_name,可以使用`${GITHUB_REPOSITORY#*/}`
    + `github.ref`工作流的分支或tag,分支为`refs/heads/<branch_name>`格式,tag是`refs/tags/<tag_name>`格式,如果只要tag名可以使用`${GITHUB_REF/refs\/tags\//}`
    + `${GITHUB_SHA::8}`可以用于获得前8位的commit的id值
    + `github.event.inputs`由手动事件触发传入的参数

+ `secrets`,项目或命名空间定义的账号密码信息,可以在`项目的Settings->Secrets`中设置,一般用于上传package或者docker镜像.

#### 运算符

`Github Action`支持如下运算符

| 运算符 | 描述         |
| ------ | ------------ |
| `()`   | 逻辑分组     |
| `[ ]`  | 索引         |
| `.`    | 属性解除参考 |
| `!`    | 非           |
| `<`    | 小于         |
| `<=`   | 小于或等于   |
| `>`    | 大于         |
| `>=`   | 大于或等于   |
| `==`   | 等于         |
| `!=`   | 不等于       |
| `&&`   | 和           |
| `\|\|` | 或           |

可以看到这些运算符解百纳都是用于做谓词的.因此同擦汗给你都与`if`字段配合使用

```yml
steps:
  ...
  - name: The job has failed
    if: ${{ github.event.inputs.a >0 }}
```

`GitHub Action`进行的是宽松的等式比较,其原理是将不同类型的数据转换为数字进行比较:
| 类型     | 结果                                         |
| -------- | -------------------------------------------- |
| `null`   | 0                                            |
| `true`   | 返回 1                                       |
| `false`  | 返回 0                                       |
| `字符串` | 空字符串为0,符合数字格式的为对应数,否则为NaN |
| `Array`  | NaN,在为同一实例时才视为相等                 |
| `Object` | NaN,在为同一实例时才视为相等                 |

注意,类似SQL中的NULL,一个 NaN 与另一个 NaN 的比较不会产生 true.

#### 函数

Github Action支持一些内置函数,比较有用的有:

+ `contains( search, item )`,用于查看序列中是否存在元素
+ `startsWith( searchString, searchValue)`/`endsWith( searchString, searchValue)`,用于查看字符串中是否已特定字符串开头或者结尾
+ `format('Hello {0} {1} {2}', 'Mona', 'the', 'Octocat')`,类似python中的`string.format()`,使用模板字符串拼接字符串结果
+ `join( array, optionalSeparator )`,类似python中的join,用于拼接数组内容为字符串.
+ 作业状态检查函数`success()/always()/cancelled()/failure()`,这类函数返回的是bool型数据,因此一般作为谓词与`if`联合使用

### 执行策略

执行策略在一级关键字`strategy`中定义.它用于规定执行器执行workflow的行为.主要包括

+ `matrix`,定义执行矩阵,执行器会遍历矩阵执行作业,`matrix`中定义的值在执行时可以从上下文`matrix`中获取
+ `max-parallel`(int)最大并行度
+ `fail-fast`(bool,true)快速失败,任何`matrix`作业失败,GitHub将取消所有进行中的作业

上面的例子中我们定义了`python-version: [3.6, 3.7, 3.8, 3.9]`,这也就意味着执行器会以`matrix.python-version`为`3.6, 3.7, 3.8, 3.9`分别执行一次.

### 使用社区定义好的action

可以将action理解为执行过程的封装,使用的人只需要知道它的用法而不需要知道它具体怎么实现的,我们可以自己定义action也可以使用外面定义好的action就像我们编程调用函数一样.社区的actions可以在[marketplace](https://github.com/marketplace?type=actions)找到

上面的例子中我们就使用了一个外部定义好的action:`actions/setup-python@v2`

使用`action`用关键字`uses`来声明,如果action需要参数可以使用`with`来传入参数

```yaml
  - name: Set up Python \$\{\{ matrix.python-version \}\}
    uses: actions/setup-python@v2
    with:
        python-version: \$\{\{ matrix.python-version \}\}
```

比较常用的action有:

+ [actions/setup-python@v2](https://github.com/actions/setup-python),自动设置python环境
+ [actions/setup-node@v1](https://github.com/marketplace/actions/setup-node-js-environment),设置node环境
+ [actions/setup-go@v2](https://github.com/marketplace/actions/setup-go-environment),设置golang环境
+ [docker/build-push-action@v1](https://github.com/marketplace/actions/docker-build-push-action),登录docker 镜像仓库
+ [actions/upload-artifact@v2](https://github.com/marketplace/actions/upload-a-build-artifact),将`Artifact`发送到workflow的管理界面用于下载
+ [getsentry/action-release@v1](https://github.com/marketplace/actions/sentry-release),发送消息到sentry

## jobs间的依赖关系

当我们单纯定义job时这些job会并行执行,而如果希望明确其中的依赖关系,则可以使用关键字`needs`.`needs`后的值可以是字符串也可以是字符串为元素的列表

```yaml
jobs:
  build_and_pub_to_pypi:
    ...
  docker-build:
    needs: build_and_pub_to_pypi
```

## workflow执行器

github默认给每个用户提供了3个的执行器(两核7g内存16g硬盘),这三个执行器的配置是不可变的,但我们可以在`仓库的settings->Actions`中配置使用安全策略.

默认的允许所有行为自然是不安全的,但其实一般用问题也不大.当项目是私有项目时我们就需要对action进行限制了.

![管理执行器][3]

### self-host

我们也可以配置自己的执行器,点击`add runner`,进入页面后选好操作系统和平台,然后按指示的配置自己的机器就可以了.

selfhost机器的优势是可以提供更加丰富的配置方式.比如我们做深度学习项目,要用gpu,那可以配一台self-host的gpu机器,然后指定它执行任务;比如我们的部署服务器在内网环境并没有开外网,那么可以让执行器监听外网,然后部署到内网环境,这样也就相对安全了.

使用self-host,只要在配置中一级关键字`runs-on`上指定即可

```bash
runs-on: [self-hosted, linux, ARM64]
```

## 管理和查看workflow

在仓库中我们可以在顶部`Actions`标签中管理.
![actions][1]

进入后我们可以创建新的workflow,或者查看之前执行过的workflow.
![管理仓库的actions][4]

我们可以点击进入某一个workflow中查看详情,在详情页可以重跑任务.如果有上传Artifact也可以在其中下载到
![详情][5]

[1]: {{site.url}}/img/in-post/githubaction/actionplace.PNG
[2]: {{site.url}}/img/in-post/githubaction/actions-manually-run-workflow.png
[3]: {{site.url}}/img/in-post/githubaction/action_runner.PNG
[4]: {{site.url}}/img/in-post/githubaction/actionsmanager.PNG
[5]: {{site.url}}/img/in-post/githubaction/actions-detial.png