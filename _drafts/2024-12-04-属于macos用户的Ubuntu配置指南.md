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

## git的安装设置

git工具我们需要好好设置下,毕竟ubuntu下很多东西,尤其是大文件的项目都需要git下载

1. 安装`git`,ubuntu并不会默认安装git工具,因此我们需要先安装git
        
    ```bash
    sudo apt install git # 安装git
    sudo apt-get install git-lfs
    ```

2.设置以支持大文件传输

+ `.gitconfig`

    ```txt
    [filter "lfs"]
        clean = git-lfs clean -- %f
        smudge = git-lfs smudge -- %f
        process = git-lfs filter-process
        required = true
    [user]
        name = HUANG SIZHE
        email = hsz1273327@gmail.com
    [http]
        postBuffer = 157286400
        version = HTTP/1.1

    [core] 
        packedGitLimit = 512m 
        packedGitWindowSize = 512m 
    [pack] 
        deltaCacheSize = 2047m 
        packSizeLimit = 2047m 
        windowMemory = 2047m
    ```

## 配置代理

和在macos上一样,我们可以在`.zshrc`或其他shell的配置文件中像这样配置代理

```bash
# 设置使用本机代理
alias setproxy="export https_proxy=http://127.0.0.1:7897 http_proxy=http://127.0.0.1:7897 all_proxy=socks5://127.0.0.1:7897"
# 设置使用本地局域网代理
alias setlocalproxy="export https_proxy=http://192.168.50.177:7890 http_proxy=http://192.168.50.177:7890 all_proxy=socks5://192.168.50.177:7890"
# 清空代理配置
alias unsetproxy="unset https_proxy;unset http_proxy;unset all_proxy"
```

## 美化系统

美化系统我们大致可以分为如下几个步骤

1. 美化桌面
2. 美化登录页面
3. 添加实用插件
4. 美化terminal
5. 美化字体

一般我也会按这个次序进行设置

macos风格的Gnome桌面美化一般使用的是[vinceliuice/WhiteSur-gtk-theme](https://github.com/vinceliuice/WhiteSur-gtk-theme.git)这个项目.
这个项目其实已经可以包办大部分的美化任务了.

我们可以先找个地方(比如`~/workspace/init_source`)来安装它

```bash

mkdir -p workspace/init_source # 构造目录
cd workspace/init_source
git clone https://github.com/vinceliuice/WhiteSur-gtk-theme.git --depth=1
```

### 美化桌面

对于美化桌面,其实也可以认为是3个任务

#### 美化主题

一个是主题美化,这是一个很复杂的问题,我们知道ubuntu最开始使用的是源自debian的`deb`软件分发,但最近几个版本他们搞了个`snap`.而且还存在一种`flatpak`应用,也就是说有三种gui应用.

+ `deb`软件,debian的软件包格式,最轻量,也最多.`apt`工具安装的就是这种
+ `flatpak`软件,可以在[flathub](https://flathub.org/)中下载并安装.一种完全开源的,面向跨发行版分发的软件包格式.我们的美化主题需要额外设置以适配
+ `snap`软件,ubuntu主推的软件包格式,也是为了跨发行版分发而推出的,但由于其服务端并不开源,不太受开源社区待见.

主题美化不光是自带软件的美化,还得让各种gui应用都可以获得相应的美化.这就很难了.

但对于我们这种要求不高的其实就很简单.直接执行`WhiteSur-gtk-theme`项目的`./install.sh`脚本即可

```bash
cd workspace/init_source/WhiteSur-gtk-theme
./install.sh
./install.sh -t all  # 安装指定颜色主题,如果要全部颜色可以使用`-t all`,要指定颜色则是类似`-t [purple/pink/red/orange/yellow/green/grey]`
./install.sh -N mojave # 改变文件管理器分栏样式,可选为默认,`mojave`和`glassy`
./install.sh -l  # 安装对`libadwaita`软件的适配,目前并不完美
sudo apt install flatpak
sudo flatpak override --filesystem=xdg-config/gtk-3.0 && sudo flatpak override --filesystem=xdg-config/gtk-4.0 # 适配非snap的flatpak应用
```

#### 美化壁纸

另一个桌面的美化点就是壁纸,我们可以使用[vinceliuice/WhiteSur-wallpapers](https://github.com/vinceliuice/WhiteSur-wallpapers)项目提供的macos风格的壁纸.安装好然后去设置中替换即可.

```bash
cd workspace/init_source
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
cd workspace/init_source
git clone https://github.com/vinceliuice/WhiteSur-icon-theme.git
cd WhiteSur-icon-theme
./install.sh
./install.sh -a # 安装macos风格的替换图标
./install.sh -b # 安装右上角下拉菜单的图标
```

需要注意这个库对snap应用并不原生支持

在安装好后我们还需要设置激活

在`显示应用->工具->优化`中选中`外观`然后设置.我个人习惯用如下设置

+ `图标`: `WhiteSur-light`
+ `shell`:  `WhiteSur-light-solid`
+ `过时应用程序`: `WhiteSur-light-solid`

#### 美化鼠标

这也属于不优化也没什么大不了的项目,但如果能优化下确实体验是能更好些.

美化鼠标可以使用项目[ful1e5/apple_cursor](https://github.com/ful1e5/apple_cursor)这个项目,去它最新的release中下载`macOS.tar.xz`
,然后解压,可以获得两个文件夹`macOS`和`macOS-White`,将他们都移动到`/usr/share/icons`目录后重启就安装完成了.

```bash
cd macOS
sudo mv macOS* /usr/share/icons/ 
sudo reboot
```

之后在`显示应用->工具->优化`中选中`外观`然后设置`光标`为`macOS`即可.

### 美化登录页面

对于登录界面我们还是使用`vinceliuice/WhiteSur-gtk-theme`这个项目

```bash
cd workspace/init_source/WhiteSur-gtk-theme
sudo ./tweaks.sh -g # 我们可以增加`-nd`(不将背景变暗)或`-nb`(不将背景变模糊)或`-b default`(默认,背景变暗变模糊)来设置效果.
```

这个登录页除了ubuntu字样外就完全果里果气了.

### 添加实用插件

Gnome支持插件.插件可以增加功能也可以增加动画效果等.而gnome的插件搜索和安装在ubuntu下我们一般依赖于firefox浏览器.那不妨我们就先优化下firefox浏览器的体验.

#### firefox浏览器优化

就像ie/edge之于windows,safari之于macos,firefox浏览器是ubuntu自带的默认浏览器.讲道理它和chrome一样很好用,甚至在chrome之前它是最好用的浏览器.但由于我个人chrome用的太久,要迁移太麻烦,所以我还是会将chrome作为主力浏览器.

但即便为了Gnome插件,firefox也是值得优化下体验的.这个优化主要包括3个方面

> 美化

依然借助`vinceliuice/WhiteSur-gtk-theme`项目,这个项目提供了对firefox的专门优化

```bash
cd workspace/init_source/WhiteSur-gtk-theme
./tweaks.sh -f monterey # 可选flat和monterey,monterey比较紧凑
```

> 墙

处理墙的问题我们只要安装[zeroomega](https://github.com/zero-peak/ZeroOmega)即可.进去网页<https://addons.mozilla.org/en-US/firefox/addon/zeroomega/>点击安装就可以了,剩下的就是设置ip和端口.

> 内建翻译

firefox有内建翻译,但不支持中文,我们可以安装插件[划词翻译](https://addons.mozilla.org/zh-CN/firefox/addon/hcfy/?utm_source=addons.mozilla.org&utm_medium=referral&utm_content=search)借助谷歌百度等的翻译来实现和chrome中的相同效果

#### 安装gnome浏览器插件

虽然是个浏览器插件,但我们得用apt安装.

```bash
sudo apt-get install gnome-browser-connector
```

安装好后firefox的右上角插件栏中就会有一个脚印一样的图标,它就是gnome的浏览器插件,点击它就可以进入插件搜索页面.

#### gnome插件安装

安装gnome插件很简单,用firefox的gnome浏览器插件进入到gnome插件页面后点击插件名后面的开关到开的状态即可.

当安装好后我们可以在`扩展`应用中对插件进行开关和设置,而已经安装了哪些插件可以在[installed extentions页面中查看](https://extensions.gnome.org/local/).

需要注意我们的系统`ubuntu 2024.04`使用的是`Gnome 46`,插件需要支持这个版本才能安装.

下面是我认为比较有必要的gnome插件汇总

| 插件                                                                                                          | 推荐等级 | 用途                              | 补充说明                                                                                                                             |
| ------------------------------------------------------------------------------------------------------------- | -------- | --------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| [user-themes](https://extensions.gnome.org/extension/19/user-themes/)                                         | 高       | 管理用户主题                      | ---                                                                                                                                  |
| [Dash to Dock](https://extensions.gnome.org/extension/307/dash-to-dock/)                                      | 高       | 一个对主题更友好的dash            | 关闭`Ubuntu Docker`,功能重复了                                                                                                       |
| [Blur my Shell](https://extensions.gnome.org/extension/3193/blur-my-shell/)                                   | 高       | 一个提供桌面模糊的插件            | 建议修改`Dash to Dock`中的拐角半径                                                                                                   |
| [Clipboard Indicator](https://extensions.gnome.org/extension/779/clipboard-indicator/)                        | 高       | 剪切板功能,可以保存近期的复制内容 | ---                                                                                                                                  |
| [Lock Keys](https://extensions.gnome.org/extension/1532/lock-keys/)                                           | 高       | 大小写锁定提示                    | ---                                                                                                                                  |
| [Removable Drive Menu](https://extensions.gnome.org/extension/7/removable-drive-menu/)                        | 高       | 顶栏的移动存储操作工具            | ---                                                                                                                                  |
| [Screenshort-cut](https://extensions.gnome.org/extension/6868/screenshort-cut/)                               | 中       | 顶栏截图工具                      | ---                                                                                                                                  |
| [Vitals](https://extensions.gnome.org/extension/1460/vitals/)                                                 | 中       | 顶栏系统监控                      | ---                                                                                                                                  |
| [GSConnect](https://extensions.gnome.org/extension/1319/gsconnect/)                                           | 高       | 快速连接移动端设备                | 需要配合app`kde connect`                                                                                                             |
| [Todoit](https://extensions.gnome.org/extension/7538/todo-list/)                                              | 低       | 顶部todolist                      | ---                                                                                                                                  |
| [Lunar Calendar 农历](https://extensions.gnome.org/extension/675/lunar-calendar/)                             | 中       | 日历改为农历                      | 需要先额外安装[Nei/ChineseCalendar](https://gitlab.gnome.org/Nei/ChineseCalendar/-/archive/20240107/ChineseCalendar-20240107.tar.gz) |
| [Compiz alike magic lamp effect](https://extensions.gnome.org/extension/3740/compiz-alike-magic-lamp-effect/) | 中       | 仿macos的最小化动画               | ---                                                                                                                                  |
| [No overview at start-up](https://extensions.gnome.org/extension/4099/no-overview/)                           | 高       | 取消开机时自动进入overview        | ---                                                                                                                                  |
| [gTile](https://extensions.gnome.org/extension/28/gtile/)                                                     | 高       | 多应用划分窗口                    | ---                                                                                                                                  |
| [Input Method Panel](https://extensions.gnome.org/extension/261/kimpanel/)                                    | 高       | 输入法相关                        | ---                                                                                                                                  |
| [Click to close overview](https://extensions.gnome.org/extension/3826/click-to-close-overview/)               | 低       | 点击空白处关闭预览                | ---                                                                                                                                  |
| [Hide Top Bar](https://extensions.gnome.org/extension/545/hide-top-bar/)                                      | 低       | 自动隐藏顶部工具栏                | ---                                                                                                                                  |

除此之外,我个人推荐对系统默认插件做如下处理

+ 禁用Desktop Icons,这个插件会让桌面有图标(默认会有你的home目录文件夹)

### 增加空格键预览功能

mac下一个经典操作就是选中目标后按空格键

```bash
sudo apt-get install gnome-sushi unoconv
```

### 美化terminal

ubuntu默认就是zsh,美化terminal其实也分成两个部分

#### terminal本体的美化

依然是使用`Solarized Dark`.在`Default->颜色->内置方案`中找到并设置即可.剩下的就是根据个人习惯做细微修改,
比如我会将文本颜色和光标改为更醒目的白色,将粗体字改为紫色,高亮字改为橙色,将光标改为竖线,并稍微设置点透明度.
ubuntu 默认使用的是`Gonme termial`,他只支持统一的透明度,比较可惜.

#### zsh的美化

我们还是把bash环境替换为zsh,毕竟确实更好用.

```bash
# 安装 zsh
sudo apt install zsh

# 更改默认shell为zsh
chsh -s /bin/zsh
```

> 安装`oh-my-zsh`

```bash
sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
```

> zsh的经典插件安装

+ autojump

    ```bash
    sudo apt install autojump
    ```

+ zsh-autosuggestions

    ```bash
    git clone https://github.com/zsh-users/zsh-autosuggestions $ZSH_CUSTOM/plugins/zsh-autosuggestions
    ```

+ `zsh-syntax-highlighting`

    ```bash
    git clone https://github.com/zsh-users/zsh-syntax-highlighting $ZSH_CUSTOM/plugins/zsh-syntax-highlighting
    ```

在`~/.zshrc`中配置启用这些插件:

```txt
...
plugins=(git autojump zsh-autosuggestions zsh-syntax-highlighting)
...
```

> 安装主题[可选]

我们可以使用`powerlevel10k`这个主题

```bash
cd workspace/init_source
git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ${ZSH_CUSTOM:-$HOME/.oh-my-zsh/custom}/themes/powerlevel10k
```

在`~/.zshrc`中配置启用新的主题

```txt
ZSH_THEME="powerlevel10k/powerlevel10k"
```

关闭后重启终端,会有几个问题:

```bash
# 这看起来像菱形吗？
Does this look like a diamond (rotated square)?

# 这看起来像一个锁吗？
Does this look like a lock?

# 这看起来像一个向上的箭头吗？
Does this look like an upwards arrow?

# 这些图标都正常显示放在 x 之间，没有重叠吗？
Do all these icons fit between the crosses?
```

回答好后就选择终端前面的提示样式.下面是推荐的样式

| 项目                    | 推荐选择       | 说明                 |
| ----------------------- | -------------- | -------------------- |
| Prompt Style            | Rainbow        | 选择提示样式         |
| Character Set           | Unicode        | 选择字符集           |
| Show current time?      | 24-hour format | 选择提示时间样式     |
| Prompt Separators       | Round          | 提示分隔符的样式     |
| Prompt Heads            | Blurred        | 选择提示箭头的样式   |
| Prompt Tails            | Flat           | 选择提示尾样式       |
| Prompt Height           | Two line       | 提示是否独立一行显示 |
| Prompt Connection       | Disconnected   | 头尾连接样式         |
| Prompt Frame            | Left           | 两行间的联系提示样式 |
| Frame Color             | Light          | 两行间的联系提示颜色 |
| Prompt Spacing          | Compact        | 行间距离             |
| Icons                   | Many icons     | 是否开启图标         |
| Prompt Flow             | Concise        | 提示细节             |
| Enable Transient Prompt | Yes            | 是否启用瞬时提示     |

就行了,要重置可以执行`p10k configure`再来一遍

### 美化字体

ubuntu中的默认字体其实还行,但我更喜欢用[nerd-fonts](https://github.com/ryanoasis/nerd-fonts)

```bash
cd workspace/init_source
git clone https://github.com/ryanoasis/nerd-fonts.git --depth 1
```

之后进入`显示应用`->`工具`->`优化`.选中`字体`,我个人习惯设置成如下:

+ `界面文本`: `AnonymicePro Nerd Font`
+ `文档文本`: `UbuntuMono Nerd Font`
+ `等宽文本`: `Ubuntu San Font`

## 系统优化

ubuntu默认状态下是很原始的,我们需要做如下操作才能让它用起来舒服些

### 安装监控工具

对于cpu,gpu这类常规设备来说,监控就很简单,用系统自带的`系统监视器`即可.但对于gpu,系统监视器就无能为力了.

### 常用软件预加载

[preload](https://www.fosslinux.com/135006/how-to-use-preload-to-speedup-app-launches-in-ubuntu.htm)是一个空间换时间的应用预加载工具.
它会根据用户行为预先加载常用软件进内存.对于拥有大内存的机器这个工具可以大幅提高应用的打开速度.

我们可以使用如下命令在terminal中安装

```bash
sudo apt-get install preload
systemctl status preload
```

### 系统快照

类似macos的time machine,我们也可以通过系统快照锁定系统当前状态的切片以随时回退.这个操作通常使用[timeshift-gtk](https://github.com/linuxmint/timeshift)工具

安装只需要在terminal中使用如下命令

```bash
sudo apt install timeshift
```

这个工具做快照目前只支持使用本地的linux文件系统,这也就意味着我们需要自备一块空的U盘专门做系统快照.个人建议这块u盘可以就一直插在机器上面,第一次完整装完系统后做一份快照,之后设置定时每半年做一份快照.

### 浏览器中激活显卡渲染

> chrome:

硬件加速:

1. 地址栏输入`chrome://settings/`
2. 左侧选中系统,右边激活`使用图形加速功能(如果可用)`

gpu加速:
 
在地址栏输入`chrome://flags/#disable-accelerated-video-decode`,找到其中的

+ `Hardware-accelerated video decode`硬件解码设置,确保置为`已启用`

### 全面的编解码器支持

ubuntu受限于版权并不能直接提供全部的编解码支持,但我们可以安装`ubuntu-restricted-extras`来或者这一能力

```bash
sudo apt install ubuntu-restricted-extras
```

安装好后市面上大部分的音频视频格式我们就都可以使用了

<!-- ### amd显卡的监控

amdGPU状态的监控使用[AMD GPU TOP](https://github.com/Umio-Yasuno/amdgpu_top).我们可以直接在项目的release中下载`.deb`文件,双击安装即可

这个工具同样可以监控核显

<!-- ### Nvidia显卡的监控

nvidia-smi --> -->

### 安装资源监控工具

https://missioncenter.io/

https://thetumultuousunicornofdarkness.github.io/CPU-X/

### 安装防火墙

linux只是用的人少,并不是就没有安全隐患.我们还是应该装个防火墙来保护下

```bash
sudo apt install gufw
sudo ufw enable
```

### 应用管理器优化

Ubuntu现在在主推snap应用,应用中心也只能管理snap软件,这显然是不够的,毕竟ubnutu更多的还是deb应用.管理deb软件包我们可以安装`synaptic`.

```bash
sudo apt install synaptic
```

### 安装系统清理工具

linux上也是会产生垃圾文件的,所以一样的我们也应该定期清理,清理系统可以使用`bleachbit`

```bash
sudo apt install bleachbit
```

### 优化蓝牙连接

ubuntu中蓝牙设备在机器长期不用后会自动断开连接.这对于一般的设备来说挺好,毕竟还省电了.但对于鼠标那就尴尬了.

为了让蓝牙设备不丢失,可以这样设置

1. 从terminal中进入蓝牙的设置文件`cd /etc/bluetooth`
2. 找到其中的`main.conf`,在`[LE]`块下修改`Autoconnecttimeout=0`

### 优化输入法

在linux桌面环境下输入法似乎是一个比较麻烦的问题.在Linux系统上常见的输入法框架(Keyboard input method system)有三种:

+ IBus(Intelligent Input Bus),也是ubuntu的默认输入法
+ Fcitx(FlexibleInput Method Framework),谷歌搜狗等用的框架,
+ XIM(X Input Method).一般没啥用

上面的3种输入法框架Ubuntu 24.04都自带,但相对而言Fcitx可能是更智能的选择,因为字库比较多.

我们可以用如下步骤设置出一个相对好用的中文输入环境(参考自[大佬DebuggerX的这篇博客](https://www.debuggerx.com/2023/09/20/fcitx5-customizer/))

1. 在`设置->系统->区域与语言`中设置中文环境.
2. 在`语言支持`中选择默认的输入法框架为`Fcitx 5`
3. 挂着全局代理执行如下命令

    ```bash
    curl -sSL https://www.debuggerx.com/fcitx5_customizer/fcitx5_customizer.sh | bash -s -- recommend
    ```

4. 之后除了去掉`搜狗词库`外都按推荐的来即可

5. 在`.zshrc`(如果希望是系统级修改则在`/etc/profile`)中设置环境变量

    ```bash
    export XMODIFIERS=@im=fcitx
    export QT_IM_MODULE=fcitx
    export GTK_IM_MODULE=fcitx
    ```

6. 安装Gnome插件`Input Method Panel`

7. 在`显示应用->工具->优化`中将`Fcitx 5`添加自启动项,重启即可

<!-- sudo apt install gnome-tweaks -->

### 优化快捷键

快捷键优化我们需要分为在桌面环境和在terminal中,在ubuntu中主要的快捷键逻辑还是跟随的windows,这就有点尴尬了.我个人会对快捷键做出如下设置

1. 重新设置一些系统级快捷键

    我们可以进入`设置->键盘->键盘快捷键->查看及自定义快捷键`

    | 分类   | 操作             | 按键组合              |
    | ------ | ---------------- | --------------------- |
    | 启动器 | 启动终端         | `Control + Alt + T`   |
    | 启动器 | 全局搜索         | `Super+ Space`        |
    | 导航   | 切换应用         | `Control+ Tab`        |
    | 截图   | 全屏截图         | `Shift + Control + 2` |
    | 截图   | 窗口截图         | `Shift + Control + 3` |
    | 截图   | 截图             | `Shift + Control + 4` |
    | 截图   | 录屏             | `Shift + Control + 5` |
    | 窗口   | 全屏模式切换     | `Super + Control + F` |
    | 窗口   | 最大化模式切换   | `Alt + Control + F`   |
    | 窗口   | 隐藏窗口(最小化) | `Control + M`         |
    | 窗口   | 关闭应用         | `Control + Q`         |
    | 系统   | 锁定屏幕         | `Control + Super + Q` |


2. 重新设置一些terminal快捷键

    右键terminal,点击`首选项`,在`快捷键`中进行如下修改

    | 操作       | 按键组合      |
    | ---------- | ------------- |
    | 复制       | `Control + C` |
    | 黏贴       | `Control + V` |
    | 编辑器全选 | `Control + A` |
    | 关闭窗口   | `Control + Q` |

3. 设置vscode.

    进入`文件->首选项->键盘快捷方式`通过搜索关键字替换快捷键

    | 操作       | 按键组合          | 冲突                                                |
    | ---------- | ----------------- | --------------------------------------------------- |
    | 格式化文档 | `Shift + Alt + F` | 会有两个冲突,分别设置为`Alt + F` 和`Shift + F` 即可 |

4. 左侧`Control`和`Win`键(Linux下`Super`键,Mac下`Command`键)互换映射.

    点击`显示应用`,打开`工具->优化`,选择`键盘->布局->其他布局`,勾选`Ctrl位置->交换左win左ctrl`

    这样操作起来就和mac中一样了.如果你的键盘是机械键盘,你也可以将`Contrl`和`Win`键互换位置

这样我们就可以将整机的快捷键统一为如下,需要注意现在你的键盘`Win`键和`Control`已经互换了位置:

| 操作           | 按键组合              |
| -------------- | --------------------- |
| 复制           | `Control + C`         |
| 黏贴           | `Control + V`         |
| 切换输入法     | `Control + Space`     |
| 全局搜索       | `Super+ Space`        |
| 全屏截图       | `Shift + Control + 3` |
| 截图           | `Shift + Control + 4` |
| 录屏           | `Shift + Control + 5` |
| 全屏模式切换   | `Super + Control + F` |
| 最大化模式切换 | `Alt + Control + F`   |
| 应用最小化     | `Control + M`         |
| 关闭应用       | `Control + Q`         |
| 切换应用       | `Control+ Tab`        |
| 锁屏           | `Control + Super + Q` |
| 编辑器格式化   | `Shift + Alt + F`     |
| 编辑器搜索     | `Control + F`         |
| 存储文件       | `Control + S`         |
| 编辑器全选     | `Control + A`         |
| 编辑器撤销     | `Control + Z`         |
| 编辑器注释     | `Control + /`         |

ubuntu特有操作

| 操作                         | 按键组合                   |
| ---------------------------- | -------------------------- |
| 快速启动terminal             | `Control + Alt + T`        |
| workspace管理                | `Super`                    |
| 切换窗口                     | `Alt + Tab`                |
| 移动到左边的工作空间         | `Super + PageUp`           |
| 移动到右边的工作空间         | `Super + PageDown`         |
| 将窗口向左移动一个工作空间   | `Shift + Super + PageUp`   |
| 将窗口向右移动一个工作空间   | `Shift + Super + PageDown` |
| 切换到第一个工作空间         | `Super + Home`             |
| 切换到最后一个工作空间       | `Super + End`              |
| 将窗口移动到第一个工作空间   | `Shift + Super + Home`     |
| 将窗口移动到最后一个工作空间 | `Shift + Super + End`      |

由于键盘`Win`键和`Control`互换了位置,因此即便有一些没有覆盖到的快捷键也会和mac上的基本一致不用额外修改

如果你不在乎wayland,可以接受桌面环境运行在x11上,那我们也可以安装[kinto](https://github.com/rbreaves/kinto)这个项目来获得不同风格且统一的快捷键布局.

<!-- ## 远程桌面

远程桌面分为本地开放远程桌面给远端连接和本地连接远端远程桌面.

ubuntu自带远程桌面服务,我们打开来就可以了

1. 打开`设置`,找到``,点击`远程桌面`
2. 启用`远程桌面`,不要钩选vnc
3. 启用`远程控制`
4. 在`如何连接`中将设备名设置为你的机器名,地址设置为`ms-rd://<设备名>.local`
5. 在登录验证中设置用户名和密码

这样在本地的设置就好了,之后我们在ddnsto中进行设置

1. 进入`我的设备`,选择`远程应用`,点`+`
2. 选择`远程rdp`
3. 设置应用名(随便写);ip为上面开启远程桌面的机器在内网中的ip;端口不变;用户名和密码就是你再上面设置的对应值,NLA认证设置为True即可.

之后要远程使用的时候就点击这个远程应用即可 -->

## 安装常用开发环境

我个人比较习惯借助vscode的`Dev Containers`插件,在完全独立的容器中做开发.但即便是这样,本地的开发环境还是必不可少的,毕竟vscode还要用对应语言的生态做编程辅助.而且也难免会需要本地的轻量级开发

### C/C++环境


### Golang环境


### android开发环境


### protobuffer编译环境


### python环境


### node.js环境

### latex环境

## 安装常用软件

linux下我会尽量推荐开源工具.开源配开源嘛,他好我也好.

下面是一些常用软件的安装信息

| 分类       | 软件                            | 渠道                                                                                       | 说明                                                         |
| ---------- | ------------------------------- | ------------------------------------------------------------------------------------------ | ------------------------------------------------------------ |
| 平台       | steam                           | 应用商店                                                                                   | 必装,很多windows下的软件可以靠它在linux下运行,并不仅仅是游戏 |
| 平台       | docker                          | [官网指导](https://docs.docker.com/engine/install/ubuntu/)                                 | 必装,linux下开发神器                                         |
| 生产力工具 | gimp                            | 应用商店                                                                                   | 开源的图像编辑软件,ps平替                                    |
| 生产力工具 | freecad                         | 应用商店                                                                                   | 开源的工程制图,autocad平替                                   |
| 生产力工具 | blender                         | steam                                                                                      | 开源的3d建模渲染工具,maya平替                                |
| 生产力工具 | godot                           | steam                                                                                      | 开源的轻量级游戏引擎                                         |
| 生产力工具 | unrealengine5                   | [官网下载](https://www.unrealengine.com/zh-CN/download)                                    | 大名鼎鼎的虚幻引擎,                                          |
| 生产力工具 | [shotcut](https://shotcut.org/) | 应用商店                                                                                   | 轻量级的开源视频剪辑工具                                     |
| 生产力工具 | DaVinci Resolve                 | [官网下载](http://www.blackmagicdesign.com/cn/products/davinciresolve)                     | 大名鼎鼎的生产级视频剪辑工具达芬奇,有免费的社区版            |
| 生产力工具 | 飞书                            | [官网下载](https://www.feishu.cn/download)                                                 | 知名的办公协作工具                                           |
| 生产力工具 | 微信                            | [官网下载](https://linux.weixin.qq.com/)                                                   | 知名的聊天工具                                               |
| 生产力工具 | wps                             | [官网下载](https://www.wps.cn/product/wpslinux)                                            | 知名的office套件                                             |
| 生产力工具 | obs                             | [官网下载](https://obsproject.com/)                                                        | 知名的开源直播录屏工具                                       |
| 生产力工具 | vscode                          | 应用商店                                                                                   | 文本编辑器                                                   |
| 生产力工具 | github desktop                  | [fork版本下载](https://github.com/shiftkey/desktop)                                        | github desktop的第三方linux fork                             |
| 生产力工具 | clashverge                      | [github下载](https://github.com/clash-verge-rev/clash-verge-rev/releases/tag/v1.7.7)       | clash的客户端,方便离开有公共代理的换进时使用                 |
| 生产力工具 | balenaEtcher                    | [官网下载](https://etcher.balena.io/)                                                      | 镜像写入工具                                                 |
| 娱乐工具   | mpv                             | [官网指导](https://mpv.io/installation/)                                                   | 知名的开源视频播放器                                         |
| 娱乐工具   | YesPlayMusic                    | [官网下载](https://github.com/qier222/YesPlayMusic)                                        | 网易云音乐的开源第三方客户端                                 |


https://www.mapeditor.org/
 https://itch.io/game-assets/free/tag-tilemap

### steam环境补充

steam在其他操作系统中只是一个游戏平台,但在linux下它是必备软件,因为它提供了转译层[proton](https://github.com/ValveSoftware/Proton).这太伟大了,直接让linux下可以正常跑大部分windows平台的游戏,还顺便让其他windows软件也可以借助steam进行管理运行.

https://github.com/flightlessmango/MangoHud

https://www.bilibili.com/video/BV1zD4y1b7Jj?vd_source=08b668b29d50d7b81093d4adee9dfde0&spm_id_from=333.788.videopod.sections

### docker环境补充

如果希望有个图形界面用于管理docker,可以安装[docker desktop](https://docs.docker.com/desktop/setup/install/linux/)代替原生docker.但需要注意,docker desktop以及其管理的所有容器都运行在虚拟机中.这带来的的好处是

1. 有和windows,mac os上基本一致的体验
2. 更好的安全性

但同时也带来了如下缺陷

1. 性能损失
2. 失去调用显卡的能力

那是用原生的docker engine还是用docker desktop,这个取舍就需要根据需求自己来做了.
<!-- 
### Waydroid环境的补充

安装waydroid我们可以简单的使用如下命令

```bash
# sudo免密码
sudo su
exit
# 导入 repo 源头、
curl https://repo.waydro.id | sudo bash

# 安装 waydroid
sudo apt install waydroid -y
```

在安装好后我们可以适当优化下界面

```bash
sudo waydroid init
waydroid prop set persist.waydroid.width 480
waydroid prop set persist.waydroid.height 900

# waydroid session stop
```

这个模拟器虽然很丝滑,但是默认情况下是没法跑ARM的APK.而国内很少有原生的x86 APP,所以还是有必要安装一下ARM相关的转译依赖.

```bash
sudo apt install python3.12-venv
cd workspace/init_source
git clone https://github.com/casualsnek/waydroid_script
cd waydroid_script
python -m venv env # 给项目构造虚拟环境并执行设置脚本
sudo su # 需要root用户
source env/bin/activate
# 必须挂代理
export https_proxy=http://127.0.0.1:7897 http_proxy=http://127.0.0.1:7897 all_proxy=socks5://127.0.0.1:7897
python -m pip install --upgrade pip -i https://mirrors.tuna.tsinghua.edu.cn/pypi/web/simple 
python -m pip install -r requirements.txt -i https://mirrors.tuna.tsinghua.edu.cn/pypi/web/simple
python main.py
```

进入到命令行图形化界面后我们选择

+ 安卓版本: android 11
+ action: Install
+ app:
    + microg
    + libndk
    + magisk
    + fdroidpriv
    + libhoudini
    + widevine
+ MicroG:
    + standard

安装完毕后重启

```bash
# 当前会话关机
waydroid session stop
# 重启 waydroid 服务
sudo systemctl restart waydroid-container.service
```

这样后续app就可以正常安装android软件了

```bash
waydroid app install <app>.apk
```

我们也可以给waydroid模拟器设置多窗口模式

```bash
waydroid prop set persist.waydroid.multi_windows true
``` -->