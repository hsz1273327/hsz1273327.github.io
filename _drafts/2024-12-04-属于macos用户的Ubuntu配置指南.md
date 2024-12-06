---
layout: post
title: "属于macos用户的Ubuntu配置指南"
series:
    aipc_experiment:
        index: 3
date: 2024-12-04
author: "Hsz"
category: recommend
tags:
    - Linux
    - Ubuntu
    - MacOs
    - 美化
header-img: "img/home-bg-o.jpg"
update: 2024-12-04
---
# 属于MacOs用户的Ubuntu配置指南

说来惭愧,作为一个整天折腾Linux服务器程序员,自己却没有真正折腾过Linux桌面系统,这主要是被MacOs和docker给惯坏了.这次搞aipc顺便就折腾下.

作为一个老MacOs用户,我对这个桌面系统的要求是

+ 精简流畅.
+ 使用兼容性和支持尽量好的发行版.
+ 可以使界面和操作尽量接近MacOs.
+ 可以的话增加一些window下的优秀交互工具
+ 尽量少折腾


选哪个发行版答案就呼之欲出了--标题上的ubuntu.

这篇文章虽然来自于折腾aipc,但本身并不涉及具体的配置和软件,因此可以看做通用美化和设置教程来使用.

## 基本操作逻辑

Ubuntu毕竟是个独立的操作系统,它有自己的一套操作逻辑和对应工具.

Ubuntu首先是个Linux系统,这就意味着它运行在linux内核上而且所有操作都可以在terminal中实现.而桌面图形界面其实可以理解为是一个gui化的bash命令.

Ubuntu的桌面使用的是[GNOME桌面](https://www.gnome.org/),具体到我现在使用的`Ubuntu 24.04 LTS`使用的`GNOME 46`,而`GNOME`是GNU计划中的一部份,它基于纯C实现的[GTK](https://www.gtk.org/)开发. 在Linux桌面这个语境下GUI的实现方案一般就两个--GTK和QT.QT好用但会给用户带来法律上授权上的麻烦;GTK难用但完全开源不会有纠纷.所以作为GNU的一部分自然是选用的GTK.
那由于用了相对开发难度大的GTK,自然的GNOME桌面的逻辑会更简单,资源占用也会相对多些,可配置性也会差不少.与之对应的是[arch linux](https://github.com/archlinux),它使用的就是QT,可配置性和资源占用都会更好,但小bug多,适合自定义桌面比较多的场景,比如steamos就是arch linux变体.

GNOME桌面的结构如下图,

可管理的元素和windows/macos基本是一致的,都是`应用->窗口->工作区`这样的3级归属模式,但使用逻辑就不一样了.GNOME非常强调工作区的作用,它的使用逻辑是

+ Ubuntu鼓励一个应用开一个窗口
+ 一个窗口下可以多开应用,但应该是相关的应用
+ 一个工作区下可以多开窗口,但这些窗口应该是相互关联的.

比如你有一个工作区就专门用来编程,一个工作区就专门用来打游戏.拿来编程的工作区开一个vscode,开一个浏览器查资料,开一个pdf阅读器看文档,那在编程时你就只需要在这个工作区内切窗口就好了.而如果你想打游戏了,直接切工作区到打游戏的哪个工作区就行了.

当然了我们想像macos一样基于应用使用也没啥问题.这只是使用习惯问题.



## 更新系统

安装完系统后我建议先全面更新下系统,因为一般下载到系统镜像都不会是真正的最新版本.通常我们会比较关心系统发行版本和内核版本.

查看系统内核版本可以使用命令`sudo uname -r`实现,查看发行版本可以用`cat /etc/os-release`查看

我们可以使用如下命令更新系统

> 更新整个系统

```bash
sudo apt update # 更新软件包的索引或包列表
sudo apt full-upgrade #更新系统上所有过时的软件包升级到最新版本并解决依赖问题,如果不要解决依赖问题,可以使用`sudo apt upgrade`
sudo reboot
```

> 仅更新linux内核

```bash
sudo apt update # 更新软件包的索引或包列表
sudo apt-get upgrade linux-image-generic #更新内核
sudo reboot
```

## 安装驱动

通常我们并不需要手动安装驱动,系统自带的通用驱动都可以正常使用.ubuntu本身也提供了一个应用`驱动管理`来检测和安装驱动.只有一些特殊硬件我们需要在特定情况下手动安装驱动.

### 显卡驱动

最常见的需要手动安装驱动的硬件就是显卡.这种通常是为了通过官方工具调用显卡计算接口.主流的显卡其实就amd(ati)显卡和Nvidia显卡两家,最近几年Intel也开始做独显了,虽然基本没人买.

我手头只有一个apu和一张Nvidia显卡,所以就只介绍这俩了.

我这边仅以ubuntu 24.04 TLS系统为例

#### amd显卡/apu核显

对应的工具链叫[rocm](https://rocm.docs.amd.com/en/latest/),可以使用如下步骤安装:

1. 安装安装器

    ```bash
    sudo apt update # 更新软件包的索引或包列表
    sudo apt install "linux-headers-$(uname -r)" "linux-modules-extra-$(uname -r)" # 根据linux内核来安装对应的linux-headers和linux-modules-extra
    sudo usermod -a -G render,video $LOGNAME # 添加当前用户到渲染和视频分组
    wget https://repo.radeon.com/amdgpu-install/6.3/ubuntu/noble/amdgpu-install_6.3.60300-1_all.deb # 下载amdgpu安装工具,这里以6.3.6030为例
    sudo apt install ./amdgpu-install_6.3.60300-1_all.deb #安装rocm安装工具
    # sudo apt update # 更新软件包的索引或包列表
    # sudo apt install amdgpu-dkms rocm # 安装amdgpu-dkms(驱动) rocm(rocm包)
    sudo reboot
    ```

    上面的代码只是例子,我们安装的是rocm 6.3的安装器,具体版本可以查看[rocm发布页](https://rocm.docs.amd.com/en/latest/release/versions.html)

2. 根据使用场景安装需要组件

    上面的代码安装了`amdgpu-install`这个工具,它是一个amdgpu的管理工具,可以用于安装和更新AMDGPU的驱动, `rocm`,`rocm 组件`等amdgpu相关的工具.

    在重启后我们重新进入命令行,然后运行`amdgpu-install`来安装所必须得组件

    ```bash
    amdgpu-install --usecase=rocm,graphics
    sudo reboot
    ```

    支持的usecase可以通过命令`sudo amdgpu-install --list-usecase`查看.

    主要的usecase包括

    + `dkms`,仅安装驱动,其他的所有usecase都会安装驱动所以一般不用这个
    + `graphics`,图形界面相关工具,如果你使用ubuntu桌面系统你就得装,不装很多软件会因为显卡报错无法打开(比如各种electron封装)
    + `multimedia`,开源多媒体库相关工具
    + `multimediasdk`,开源多媒体开发,包含`multimedia`
    + `workstation`,工作站相关工具,包含`multimedia`同时包含闭源的OpenGL工具
    + `rocm`,显卡做异构计算工具,包括OpenCL运行时,HIP运行时,机器学习框架,和rocm相关的库和工具
    + `rocmdev`,rocm开发工具,包含`rocm`和相关的调试开发工具
    + `rocmdevtools`,仅包含`rocm`和相关的调试开发工具
    + `amf`,基于amf编解码器(闭源)的多媒体工具
    + `lrt`,rocm的编译器,运行时和设备库等工具
    + `opencl`,异构计算库opencl相关工具,库和运行时
    + `openclsdk`,包含`opencl`,同时包含opencl的相关开发工具和头文件等
    + `hip`,高性能计算库hip的运行时
    + `hiplibsdk`,包含`hip`,同时包含hip开发相关库和工具以及ROCm的数学库
    + `openmpsdk`,并行计算库openmp的运行时和相关库和工具
    + `mllib`,机器学习相关工具和库,包括MIOpen核心和相关库,以及Clang OpenCL
    + `mlsdk`,包含`mllib`,额外附带MIOpen和Clang OpenCL的开发库
    + `asan`,支持ASAN(内存检测工具)的ROCm工具

3. 设置系统连接

    也就是设置相关工具的查找位置

    ```bash
    sudo tee --append /etc/ld.so.conf.d/rocm.conf <<EOF
    /opt/rocm/lib
    /opt/rocm/lib64
    EOF
    sudo ldconfig
    ```

4. 使用更新选项或环境模块Linux实用程序配置ROCm二进制文件的路径。ROCm安装过程将ROCm可执行文件添加到这些系统中，前提是它们已安装在系统上

在安装好chongqi重启后

##### rocm版本和更新

更新

#### Nvidia显卡

对应的工具链叫[cuda](),安装


### 驱动的更新


## 美化系统

Linux桌面大致可以分为如下几个部分

+ 窗口管理器
+ 窗口页面
+ 登录管理器
+ 插件系统



## 美化terminal
## 安装监控工具
## 安装常用软件

mpv

## 安装docker
## 安装常用开发环境
## 安装steam