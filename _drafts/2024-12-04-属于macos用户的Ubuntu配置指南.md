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

对应的工具链叫[rocm](https://rocm.docs.amd.com/en/latest/),需要注意目前rocm在同时存在amd独显和amd核显的情况下会有错误.因此如果你用的是amd独显需要在bios中禁用核显(amd的这个操作真的很神奇,因此一般推荐au配n卡).

可以使用如下步骤安装:

1. 安装安装器

    ```bash
    sudo apt update # 更新软件包的索引或包列表
    sudo apt install "linux-headers-$(uname -r)" "linux-modules-extra-$(uname -r)" # 根据linux内核来安装对应的linux-headers和linux-modules-extra
    sudo usermod -a -G render,video $LOGNAME # 添加当前用户到渲染和视频分组
    wget https://repo.radeon.com/amdgpu-install/6.3/ubuntu/noble/amdgpu-install_6.3.60300-1_all.deb # 下载amdgpu安装工具,这里以6.3.6030为例
    sudo apt install ./amdgpu-install_6.3.60300-1_all.deb #安装rocm安装工具
    # sudo apt update # 更新软件包的索引或包列表
    # sudo apt install amdgpu-dkms rocm # 安装amdgpu-dkms(驱动) rocm(rocm包)
    sudo reboot #重启后生效
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

    正常情况下使用`amdgpu-install --usecase=rocm,graphics` 安装即可

3. 设置系统连接

    也就是设置相关工具的查找位置

    ```bash
    sudo tee --append /etc/ld.so.conf.d/rocm.conf <<EOF
    /opt/rocm/lib
    /opt/rocm/lib64
    EOF
    sudo ldconfig
    ```

4. 使用`update-alternatives`更新配置ROCm二进制文件的路径.

    ```bash
    update-alternatives --list rocm
    ```

5. 设置环境变量

    rocm安装好后会被放在`/opt/rocm-<ver>`目录:
    + rocm的可执行文件会放在`/opt/rocm-<ver>/bin`目录.
        如果无法使用rocm工具,可以将它的`bin`目录加入到PATH中

        ```bash
        export PATH=$PATH:/opt/rocm-6.3.0/bin
        ```

    + rocm的动态链接库会放在`/opt/rocm-<ver>/lib`目录.
        如果要用到这些动态链接库,可以将它加入到`LD_LIBRARY_PATH`

        ```bash
        export LD_LIBRARY_PATH=/opt/rocm-6.3.0/lib
        ```

    + rocm的模块则会被放在`/opt/rocm-<ver>/lib/rocmmod`目录.

6. 检查驱动是否正常

    ```bash
    dkms status
    ```

    这个命令会打印出显卡的状态

7. 检查rocm是否正常安装

    ```bash
    rocminfo # 检查rocm状态
    clinfo # 检查opencl状态
    ```

8. 检查包是否安装正常

    ```bash
    apt list --installed
    ```

##### rocm版本更新

更新版本我们需要完全卸载已有的rocm,驱动,和rocm安装器

```bash
sudo amdgpu-install --uninstall # 卸载驱动和库
sudo apt purge amdgpu-install # 卸载安装器
sudo apt autoremove # 卸载对应依赖
sudo reboot # 重启后生效
```

之后在下载新版本的安装器重新安装配置一次即可

<!-- #### Nvidia显卡

对应的工具链叫[cuda](),安装 -->

## 美化系统

美化系统我们大致可以分为如下几个步骤

1. 美化桌面
2. 美化登录页面
3. 添加实用插件
4. 美化terminal
5. 优化快捷键

一般我也会按这个次序进行设置

macos风格的Gnome桌面美化一般使用的是[vinceliuice/WhiteSur-gtk-theme](https://github.com/vinceliuice/WhiteSur-gtk-theme.git)这个项目.
这个项目其实已经可以包办大部分的美化任务了.

我们可以先找个地方(比如`~/workspace/beautify`)来安装它

```bash
sudo apt install git # 安装git
mkdir -p workspace/beautify # 构造目录
cd workspace/beautify
git clone https://github.com/vinceliuice/WhiteSur-gtk-theme.git --depth=1
```

### 美化桌面

对于美化桌面,其实也可以认为是3个任务

#### 美化主题

一个是主题美化,这是一个很复杂的问题,我们知道ubuntu最开始使用的是源自debian的`deb`软件分发,但最近几个版本他们搞了个`snap`.而且还存在一种`flatpak`应用,也就是说有三种gui应用.

+ `deb`软件
+ `snap`软件
+ `flatpak`软件

主题美化不光是自带软件的美化,还得让各种gui应用都可以获得相应的美化.这就很难了.

但对于我们这种要求不高的其实就很简单.直接执行`WhiteSur-gtk-theme`项目的`./install.sh`脚本即可

```bash
cd workspace/beautify/WhiteSur-gtk-theme
./install.sh
./install.sh -t all  # 安装指定颜色主题,如果要全部颜色可以使用`-t all`,要指定颜色则是类似`-t [purple/pink/red/orange/yellow/green/grey]`
./install.sh -N mojave # 改变文件管理器分栏样式,可选为默认,`mojave`和`glassy`
./install.sh -l  # 安装对`libadwaita`软件的适配,目前并不完美
sudo flatpak override --filesystem=xdg-config/gtk-3.0 && sudo flatpak override --filesystem=xdg-config/gtk-4.0 # 适配非snap的flatpak应用
```

#### 美化壁纸

另一个桌面的美化点就是壁纸,我们可以使用[vinceliuice/WhiteSur-wallpapers](https://github.com/vinceliuice/WhiteSur-wallpapers)项目提供的macos风格的壁纸.安装好然后去设置中替换即可.

```bash
cd workspace/beautify
git clone https://github.com/vinceliuice/WhiteSur-wallpapers.git
cd WhiteSur-wallpapers
# 安装会根据时间变化的桌面壁纸,
#可以使用`-t [whitesur|monterey|ventura]`指定壁纸,默认全装;
#可以使用`-s [1080p|2k|4k]`指定分辨率,默认4k;
sudo ./install-gnome-backgrounds.sh
# 安装静态壁纸
#可以使用`-t [whitesur|monterey|ventura]`指定壁纸,默认全装;
#可以使用`-s [1080p|2k|4k]`指定分辨率,默认4k;
#可以使用`-c [night|light|dark]`指定颜色风格,默认全装;
#可以使用`-n [whitesur|monterey|ventura]`安装灰化壁纸,默认不灰化;
./install-wallpapers.sh
```

#### 美化图标

最后是美化图标.实话讲ubuntu的图标确实丑.我们可以使用[vinceliuice/WhiteSur-icon-theme](https://github.com/vinceliuice/WhiteSur-icon-theme)项目提供的图标来美化它,这个图标库就很还原macos了.

```bash
cd workspace/beautify
git clone https://github.com/vinceliuice/WhiteSur-icon-theme.git
cd WhiteSur-icon-theme
./install.sh
./install.sh -a # 安装macos风格的替换图标
./install.sh -b # 安装右上角下拉菜单的图标
```

需要注意这个库对snap应用并不原生支持

### 美化登录页面

对于登录界面我们还是使用`vinceliuice/WhiteSur-gtk-theme`这个项目

```bash
cd workspace/beautify/WhiteSur-gtk-theme
sudo ./tweaks.sh -g # 我们可以增加`-nd`(不将背景变暗)或`-nb`(不将背景变模糊)或`-b default`(默认,背景变暗变模糊)来设置效果.
```

这个登录页除了ubuntu字样外就完全果里果气了.

### 添加实用插件

Gnome支持插件.插件可以增加功能也可以增加动画效果等.而gnome的插件搜索和安装在ubuntu下我们一般依赖于firefox浏览器.那不妨我们就先优化下firefox浏览器的体验.

#### firefox浏览器优化

就像ie/edge之于windows,safari之于macos,firefox浏览器是ubuntu自带的默认浏览器.讲道理它和chrome一样很好用,甚至在chrome之前它是最好用的浏览器.但由于我个人chrome用的太久,要迁移太麻烦,所以我还是会将chrome作为主力浏览器.

但即便为了Gnome插件,firefox也是值得优化下体验的.这个优化主要包括2个方面

> 美化

依然借助`vinceliuice/WhiteSur-gtk-theme`项目,这个项目提供了对firefox的专门优化

```bash
cd workspace/beautify/WhiteSur-gtk-theme
./tweaks.sh -f monterey # 可选flat和monterey,monterey比较紧凑
```

> 墙

处理墙的问题我们只要安装[zeroomega](https://github.com/zero-peak/ZeroOmega)即可.进去网页<https://addons.mozilla.org/en-US/firefox/addon/zeroomega/>点击安装就可以了,剩下的就是设置ip和端口.

#### 安装gnome浏览器插件

虽然是个浏览器插件,但我们得用apt安装.

```bash
sudo apt-get install gnome-browser-connector
```

安装好后firefox的右上角插件栏中就会有一个脚印一样的图标,它就是gnome的浏览器插件,点击它就可以进入插件搜索页面.

#### gnome插件安装

安装gnome插件很简单,用firefox的gnome浏览器插件进入到gnome插件页面后点击插件名后面的开关到开的状态即可.

当安装好后我们可以在`扩展`应用中对插件进行开关和设置,而已经安装了哪些插件可以在[installed extentions页面中查看](https://extensions.gnome.org/local/)

下面是我认为比较有必要的gnome插件汇总

| 插件                                                                                                          | 推荐等级    | 用途                              | 补充说明                           |
| ------------------------------------------------------------------------------------------------------------- | ----------- | --------------------------------- | ---------------------------------- |
| [user-themes](https://extensions.gnome.org/extension/19/user-themes/)                                         | 高          | 管理用户主题                      | ---                                |
| [Dash to Dock](https://extensions.gnome.org/extension/307/dash-to-dock/)                                      | 高          | 一个对主题更友好的dash            | 关闭`Ubuntu Docker`,功能重复了     |
| [Blur my Shell](https://extensions.gnome.org/extension/3193/blur-my-shell/)                                   | 高          | 一个提供桌面模糊的插件            | 建议修改`Dash to Dock`中的拐角半径 |
| [Clipboard Indicator](https://extensions.gnome.org/extension/779/clipboard-indicator/)                        | 高          | 剪切板功能,可以保存近期的复制内容 | ---                                |
| [Compiz alike magic lamp effect](https://extensions.gnome.org/extension/3740/compiz-alike-magic-lamp-effect/) | 中          | 仿macos的最小化动画               | ---                                |
| [Lock Keys](https://extensions.gnome.org/extension/1532/lock-keys/)                                           | 高          | 大小写锁定提示                    | ---                                |
| [Removable Drive Menu](https://extensions.gnome.org/extension/7/removable-drive-menu/)                        | 高          | 顶栏的移动存储操作工具            | ---                                |
| [Bluetooth Quick Connect](https://extensions.gnome.org/extension/1401/bluetooth-quick-connect/)               | 高          | 右侧顶部下拉菜单快速连接蓝牙      | ---                                |
| [Screenshort-cut](https://extensions.gnome.org/extension/6868/screenshort-cut/)                               | 中          | 顶栏截图工具                      | ---                                |
| [Audio output selector](https://extensions.gnome.org/extension/1400/audio-output-selector/)                   | 高 (待验证) | 右侧顶部下拉菜单音频输出设备选择  | ---                                |
| [No overview at start-up](https://extensions.gnome.org/extension/4099/no-overview/)                           | 高          | 取消开机时自动进入overview        | ---                                |---                                |                    
| [Vitals](https://extensions.gnome.org/extension/1460/vitals/)                                                 | 中          | 顶栏系统监控                      | ---                                |



除此之外,我个人推荐对系统默认插件做如下处理

+ 禁用Desktop Icons,这个插件会让桌面有图标(默认会有你的home目录文件夹)

### 美化terminal


### 优化快捷键

我们可以安装[](https://github.com/rbreaves/kinto)这个项目来获得不同风格且统一的快捷键布局


## 安装监控工具


## 安装常用软件

mpv

## 安装docker
## 安装常用开发环境
## 安装steam