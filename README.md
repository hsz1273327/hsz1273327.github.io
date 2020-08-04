# 博客模板

本模板ui来自[黄玄的博客](https://github.com/Huxpro/huxpro.github.io),功能上参考了[scateu的博客](https://github.com/scateu/scateu.github.io)中的一些设计.

主要改进点:

+ 将背景图高度固定为40vh,避免整个页面大半被图片遮去
+ 更细化的模块方便扩展
+ 配置文件只保留特性开关,特性配置移至`data`
+ 新增了对流程图的支持
+ 抛弃第三方评论系统(比较不喜欢外部依赖)使用gitment直接将评论放到项目issue
+ 新加入更新时间概念,会在侧边栏提示新更新的前5条博文
+ 新增幻灯片功能,幻灯片放在`_keynotes`文件夹下
+ 新增文章系列,可以将几个文章聚合到一个系列之中

# 使用指南

这部分介绍如何使用这个博客模板,用之前请先把我自己的配置部分替换谢谢~

## 如何配置环境

1. 先安装依赖
    + Ruby
    + RubyGems
    + Nodejs

2. 安装jekyll环境

```shell
gem install jekyll
```

1. 安装bundle

```shell
gem install bundle
```


1. 安装github pages相同的插件环境,在项目根目录下执行
```shell
bundle install
```

### 小工具的使用环境

小工具依赖node环境,只要在项目根目录下执行`npm install`即可.

## 如何配置博客

在项目的`_config.yml`中对项目的一些元数据和支持特性使用情况进行配置,使用[YAML语法](https://yaml.org/spec/1.2/spec.html)

具体的可以看该文件.在修改完`_config.yml`后修改各个特性在`_data`文件夹下对应的数据,这些数据对应下面的特使常量.

### 特殊常量

特殊常量|说明
---|---
`site.data.books`|本人在github上挂的书籍项目信息
`site.data.comment`|使用[gitment](https://github.com/imsun/gitment)基于github构建评论,评论实际在对应项目的issue中
`site.data.donate`|存放收款二维码
`site.data.featured_tags`|存放tags过滤的一些条件,大于`featured-condition-size`数值出现次数的tag会被展示
`site.data.friends`|存放好友博客列表
`site.data.projects`|存放个人开源项目信息
`site.data.short_about`|存放个人简介内容
`site.data.analysis`|存放百度或谷歌站长统计设置信息
`site.keynotes`|幻灯片集合
`site.series`|文章系列

## 文章分类

这个项目非强制的限定了如下几个文章分类:

分类|说明
---|---
`introduce`|介绍类文章,主要是一些工具的使用配置过程
`recommend`|推荐向文章,推荐工具介绍用法
`comment`|评论类文章,主要是评论一些外界的消息
`analysis`|分析类文章,针对一个主题做分析.
`demo`|测试类文章,用于测试一些jekyll模板功能和一些前端技术
`experiment`|实验记录
`reading_note`|读书笔记
`travel_note`|游记
`essay`|随笔也就是瞎写


## 如何写博客

1. 直接touch出一个你希望的草稿到`_drafts`文件夹

2. 编辑你的文章，注意开头需要设定例如这个：

  ```markdown
  ---
  title: "你的ｔｉｔｌｅ"
  date: 2016-11-27
  author: "Hsz"
  tags:
      - atom
      - editor
      - 编辑器
      - tool
      - 小工具
  header-img: "img/home-bg-o.jpg"
  update: 2019-02-02
  ---
  ```
3. 使用markdown语法写文章
4. 文章标题改为型如`2016-06-10-xxxxx.md`的固定格式放入`_post`文件夹
5. 将项目推到仓库即发布.

### 博文支持的特性

+ 插入latex数学公式,推荐使用`$ c = m^e \mod n $`写法作为内嵌公式,

  ```
  $$ c = m^e \mod n $$
  ```


+ 插入流程图[flowchart.js](http://flowchart.js.org/)

  ```
  ```flowchart
  st=>start: Start:>http://www.google.com[blank]
  e=>end:>http://www.google.com
  op1=>operation: My Operation
  sub1=>subroutine: My Subroutine
  cond=>condition: Yes
  or No?:>http://www.google.com
  io=>inputoutput: catch something...

  st->op1->cond
  cond(yes)->io->e
  cond(no)->sub1(right)->op1
  ```
  ```

+ 插入流程图[sequence](https://bramp.github.io/js-sequence-diagrams/)
  ```
  ```sequence
  participant Device
  participant Browser
  participant Server
  Browser->Server: username and password
  Note over Server: verify password
  Note over Server: generate challenge
  Server->Browser:  challenge
  Browser->Device: challenge
  Note over Device: user touches button
  Device-->Browser: response
  Browser->Server: response
  Note over Server: verify response
  ```
  ```

### 编辑系列文章

如果希望将多个相关的文章组合成一个系列,那么可以使用这个功能.这个功能的配置有如下几步:

1. 在`_series`文件夹下新建一个`.md`文件,其内容类似:


```markdown
---
title: "玩转github"
series_name: "get_along_well_with_github"
description: "github是全球最流行的开源项目托管平台,也是全球最大的IT领域同性交友平台(😁),本系列将介绍如何利用github参与开源项目和结交大佬."
date: 2019-03-14
author: "Hsz"
category: recommend
tags:
    - github
update: 2019-03-14
---

{% assign series_posts = "" | split: "" %}
{% for post in site.posts %}
    {% if post.series contains page.series_name %}
    {% assign series_posts=series_posts | push: post %}
    {% endif %}
{% endfor %}
{% assign indexed_posts = series_posts |sort: "series.get_along_well_with_github.index" %}

{% for post in indexed_posts %}
+ [{{ post.title }}]({{site.url}}{{ post.url }})
{% endfor %}
```

元数据中

字段|说明
---|---
`title`|系列名
`series_name`|系列的标识名,注意要英语且无空格
`description`|系列简介
`date`|创建日期
`author`|编辑者
`tags`|系列关键字
`update`|更新时间

下面的模板中需要将`sort: "series.get_along_well_with_github.index"`这边的`get_along_well_with_github`替换为你的系列的`series_name`

2. 修改你要组成系列的每篇博文,为其增加一个形如

```yml
series:
    raspberrypi_experiment:
        index: 5
    cluster_experiment:
        index: 4
```

的元数据.

其中`raspberrypi_experiment`和`cluster_experiment`是系列文章的元数据`series_name`的值,下面的`index`则是在这个系列中文章的顺序

这样一个系列合集就完成了.

### 如何使用小工具辅助写博客

1. 创建新的草稿

  使用命令`npm run draft <title>`创建一个名为`<title>`的草稿,注意草稿不用带上文件后缀,注意如果草稿中已经有重名的了那么会在命令行提示已存在,不会进行覆盖

2. 发布草稿到文章

  使用命令`npm run publish <path>`将草稿发布为文章,path为草稿的位置而非草稿名,这样就可以发布任何地方的草稿了,如果已经发布过,那么会覆盖写入.

## 编写幻灯片

编写幻灯片和编写博文大致上是一致的,都要先写元数据后写内容

### 定义元数据

元数据部分定义这个特定幻灯片的基础信息和使用了些什么特性.
一个典型的幻灯片元数据这样定义:

```markdown
title: Reveal.js与Jekyll结合演示(幻灯片)
description : Reveal.js与Jekyll结合演示(幻灯片)
date: 2019-03-13
update: 2019-03-13
author: "Hsz"
category: keynote
tags:
    - frontend
theme: night #使用的主题
transition: slide #使用的转场动画
diagram: true
mathjax: true
multiplex: # 暂时没搞明白
  id: 99702b17e6a745f9 #Secret: 14830192734523284320
  url: https://multiplex.scateu.me
```

支持的主题有:

+ `black`: Black background, white text, blue links (default theme)
+ `white`: White background, black text, blue links
+ `league`: Gray background, white text, blue links (default theme for reveal.js < 3.0.0)
+ `beige`: Beige background, dark text, brown links
+ `sky`: Blue background, thin dark text, blue links
+ `night`: Black background, thick white text, orange links
+ `serif`: Cappuccino background, gray text, brown links
+ `simple`: White background, black text, blue links
+ `solarized`: Cream-colored background, dark green text, blue links


支持的transition有: `none`, `fade`, `slide`, `convex`, `concave`, `zoom`.

由于幻灯片是一个很低频的功能,所以它支持的特性不在配置文件中开启而是在每个幻灯片中在元数据中单独开启.支持的特性有:

特性|说明
---|---
`diagram: bool`|是否支持流程图
`mathjax: bool`|是否支持latex数学公式

### 写ppt

ppt的内容写法和博文大体是一样的,只是需要多出一个分页来

我们使用段落`<section></section>`来定义每一页的范围 可以定义属性

+ `markdown="1"`来标明使用markdown语法写内容
+ `data-transition="slide-in fade-out"`来标明这个段落的进入退出动画
+ `data-transition-speed="fast"`来标明这个段落的动画速度
+ `data-background-xxx`来定义单页的背景,详细可以参考<https://github.com/hakimel/reveal.js#slide-backgrounds>

# 再开发指南

## 项目结构

文件/文件夹|用途
---|---
`_data`|存项目中的常量
`_drafts`|存草稿
`_include`|存前端模板组件
`_layout`|存前端页面模板
`_pages`|存导航指向的页面
`_posts`|存博文
`_keynotes`|存幻灯片
`_series`|文章系列目录文件
`_site`|渲染完成后的项目
`assets`|css,字体,js脚本等静态文件
`img`|用到的图片
`source`|各篇文章用到的资源(图片,音频,视频)
`tools`|辅助开发工具(js写的)
`_config.yml`|配置文件,主要是用于开关特性
`index.html`|主页
`package.json`|npm的执行脚本在其中

## 模板引擎

本项目使用[jekyll](http://jekyllcn.com/docs/home/)作为静态网页渲染工具

jekyll使用liquid模板.

+ 基本语法,参考<https://help.shopify.com/en/themes/liquid/tags>

+ 内置过滤器,参考<https://help.shopify.com/en/themes/liquid/filters>

+ 全局变量,参考<https://www.jianshu.com/p/5c6d68bcd836>

## 系列文章的实现方式

系列文章使用了`iframe`来实现,这样做其实略偷懒,还有很多值得优化的地方

## 本地调试

根目录执行`npm start`

## 想到的值得修改的地方

+ series功能的模板(现在的略丑)
+ 用回多说这类第三方评论系统
+ 替换代码高亮
+ 页面编辑和执行js代码
+ 修改post模板使其在iframe中的行为与在页面中不同,比如,在iframe中的上一篇和下一篇改为系列文章中的上一篇和下一篇而不是按时间顺序