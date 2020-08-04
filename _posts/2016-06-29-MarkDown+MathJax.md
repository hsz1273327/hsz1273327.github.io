---
title: "MarkDown+MathJax"
date: 2016-06-29
author: "Hsz"
category: recommend
tags:
    - DocumentTool
header-img: "img/home-bg-o.jpg"
update: 2019-03-13
series:
    get_along_well_with_github:
        index: 4
---
# MarkDown+MathJax

## MarkDown

简明强大的markdown文本标记语言是做笔记的好帮手,相比起word文档,你不用双手离开键盘来影响你的速记效率,
相比起tex你不用费心码一堆代码来维护格式.
当然了,markdown是为速记而生的,它的最大好处是内容与形式分离,因此不要对齐排版有过多期待,
它的特点便是可以用最简单最统一的形式清楚的表达.废话不多说,开始介绍语法

### 基本语法

markdown中换行是两个回车

如果一个回车,会接在后面继续写.

你可以是试试,
上面这段的代码如下:

```markdown
markdown中换行是两个回车

如果一个回车,会接在后面继续写.

你可以是试试,
上面这段的代码如下:
```

### html支持

markdown是html的子集,所以对于一些html标签和样式可以直接插入,但js无效.

```html
<table border="1">
<tr>
<th>姓名</th>
<td>Bill Gates</td>
</tr>
<tr>
<th rowspan="2">电话</th>
<td>555 77 854</td>
</tr>
<tr>
<td>555 77 855</td>
</tr>
</table>
```

<table  class = "zebra">
<tr>
<th>姓名</th><td>Bill Gates</td>
</tr>
<tr>
<th rowspan="2">电话</th>
<td>555 77 854</td>
</tr>
<tr>
<td >555 77 855</td>
</tr>
</table>

### 文章结构相关语法

#### 分割

就像你下面看到的一样

---
```markdown
---

---
```
---

#### 标题

markdown的标题用复数个连续`#`表示,`#`越多标题层级越低,连续`#`后空一格再写标题内容即可

---
```markdown
# 第一级
## 第二级
```

#第一级

## 第二级

---

#### 区块引用 Blockquotes

区块引用可以嵌套表现层级

```markdown
> This is the first level of quoting.
>
> > This is nested blockquote.
>
> Back to the first level.
>
> ## 这是一个标题。
>
> 1.   这是第一行列表项。
> 2.   这是第二行列表项。
>
> 给出一些例子代码：
>
>     return shell_exec("echo $input | $markdown_script");
```


> This is the first level of quoting.
>
> > This is nested blockquote.
>
> Back to the first level.
>
> ## 这是一个标题。
>
> 1.   这是第一行列表项。
> 2.   这是第二行列表项。
>
> 给出一些例子代码：
>
>     return shell_exec("echo $input | $markdown_script");

#### 列表

Markdown 支持有序列表和无序列表。

---
```markdown

+ Red

    蓝色

+ Green
+ Blue


1. Bird

    鸟

2. McHale

    >hi

3. Parish

        print("hello")
```


+ Red

    蓝色

+ Green
+ Blue


1. Bird

    鸟

2. McHale

    >hi

3. Parish

        print("hello")

---

### 字形相关

---
```markdown
**这是斜体**

*这是斜体*

***这是粗斜体***
```

**这是斜体**

*这是斜体*

***这是粗斜体***



---

### 页面元素

### 链接


---
链接可以使用快捷,内连和脚注两种方式引用

#### 快捷方式

```markdown
<http://google.com/>
```

<http://google.com/>

#### 内连方式

```markdown
[Google](http://google.com/ "Google search")
[Google](http://google.com/)
[利用github建立blog](./利用github建立blog.md)
```

[Google](http://google.com/ "Google search")

[Google](http://google.com/)

[利用github建立blog](./利用github建立blog.md)

#### 脚注方式

```markdown
[Google][1]
[1]: http://google.com/
```

[Google][1]

[1]: http://google.com/


---

### 表格

```markdown
名字|电话
---|---
小马哥|0987
小黑|9870
```

其效果:

名字|电话
---|---
小马哥|0987
小黑|9870



### 图片

图片可以有内连和脚注两种引用方式

#### 内连方式:

```markdown
![Alt google](https://lh4.googleusercontent.com/-v0soe-ievYE/AAAAAAAAAAI/AAAAAAAAAAA/OixOH_h84Po/photo.jpg)
```

![Alt google](https://lh4.googleusercontent.com/-v0soe-ievYE/AAAAAAAAAAI/AAAAAAAAAAA/OixOH_h84Po/photo.jpg)


#### 脚注方式

```markdown
![Alt google][2]

[2]: https://lh4.googleusercontent.com/-v0soe-ievYE/AAAAAAAAAAI/AAAAAAAAAAA/OixOH_h84Po/photo.jpg
```
![Alt google][2]

[2]: https://lh4.googleusercontent.com/-v0soe-ievYE/AAAAAAAAAAI/AAAAAAAAAAA/OixOH_h84Po/photo.jpg


### 插入代码

插入代码也很简单,有3种方式

+ 使用\`符号包裹,一般这种事嵌入在一段文字中

    你好\`s=hello\`

    对应的结果是:

    你好`s=hello`



+ 通过tab键缩进(不建议)

tab $ git clone xxxxx

    $ git clone xxxxx

+ 使用三个\`符号包裹,首行的```后面接上代码的类型可以高亮关键字

        ```python

        def func(a:int)->int:
            b:int = 10
            return a+b

        ```
    对应的结果:

    ```python

    def func(a:int)->int:
        b:int = 10
        return a+b

    ```

### 视屏

可以通过html标签实现(github不可以)

```html
<video width="640" height="360" preload="auto" poster="http://media.w3.org/2010/05/sintel/poster.png" controls >
  <source src="http://media.w3.org/2010/05/sintel/trailer.mp4">
</video>
```

<video width="640" height="360" preload="auto" poster="http://media.w3.org/2010/05/sintel/poster.png" controls >
  <source src="http://media.w3.org/2010/05/sintel/trailer.mp4">
</video>

### 音频

可以通过html标签实现(github不可以)
```html
<audio controls>
  <source src="{{site.url}}/source/{{page.title}}/蔡志展 - 战斗音乐.mp3">
</audio>
```
<audio controls>
  <source src="{{site.url}}/source/{{page.title}}/蔡志展 - 战斗音乐.mp3">
</audio>


#### 数学公式的支持

本来markdown是不支持数学公式的,
但如果和mathjax结合使用便可以利用"latex"的公式编辑语言在markdown中编写公式了.内嵌公式使用`$a \ne b$`表示.
效果就像$a \ne b$这样.

公式段落使用像这样:

```
$$
f(x)=\sum (x_i)^2
$$
```

效果如下:

$$
f(x)=\sum (x_i)^2
$$

latex符号可以查看[这篇总结](http://blog.163.com/goldman2000@126/blog/static/167296895201221242646561/)和[这篇总结](http://mlworks.cn/posts/introduction-to-mathjax-and-latex-expression/)


#### 流程图

流程图不是markdown默认支持的语法,需要通过js模块提供支持,有的环境有支持但大多数没有.不过我的博客是支持的.

通常我们使用[flowchart.js](http://flowchart.js.org/)语法

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

效果如下:

```flow
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

另一种是使用[sequence](https://bramp.github.io/js-sequence-diagrams/)

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

效果如下:

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
