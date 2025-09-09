---
layout: post
title: "属于macos用户的Ubuntu配置指南"
series:
    ubuntu_experiment:
        index: 1
date: 2024-12-04
author: "Hsz"
category: recommend
tags:
    - Linux
    - 美化
header-img: "img/home-bg-o.jpg"
update: 2025-09-09
---
# 属于MacOs用户的Ubuntu配置指南

说来惭愧,作为一个整天折腾Linux服务器程序员,自己却没有真正折腾过Linux桌面系统,这主要是被MacOs和docker给惯坏了.这次搞aipc顺便就折腾下.

作为一个老MacOs用户,我对这个桌面系统的要求是

+ 精简流畅.
+ 使用兼容性和支持尽量好的发行版.
+ 界面和操作尽量接近MacOs.
+ 能充分利用所有硬件,廉价的在本地折腾ai相关工具
+ 能打游戏

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

## 设置grub

ubuntu安装好后,默认的grub菜单是隐藏的.这对于一般使用来说没什么问题,但如果我们需要切换内核,需要做内存测试,那就不方便了.我们可以通过设置让它显示出来.

1. 打开终端,编辑`/etc/default/grub`文件

    ```bash
    sudo nano /etc/default/grub
    ```

2. 修改 GRUB 选项

    找到以下几行,并根据需要进行修改:

    + `GRUB_TIMEOUT_STYLE`,这个值用于设置grub的样式行为.可选值为`hidden`和`menu`.默认为`hidden`,改为`menu`,这样就会显示菜单
    + `GRUB_TIMEOUT`,这个值表示grub菜单显示的时间,单位是秒.如果你想让它一直显示直到你选择,可以设置为`-1`,如果你想让它显示5秒钟,可以设置为`5`.默认是`0`,表示不显示菜单直接进入系统
    + `GRUB_DEFAULT`,这个值表示默认启动的菜单项,默认是`0`,表示第一个菜单项.如果你想让它默认启动第二个菜单项,可以设置为`1`.如果你想让它默认启动最后一个菜单项,可以设置为`-1`.如果你想让它默认启动某个特定的菜单项,可以设置为`saved`,然后使用`sudo grub-set-default "菜单项名称"`来设置.

    我个人习惯设置如下:

    ```txt
    GRUB_TIMEOUT_STYLE=menu # 将 hidden 改为 menu
    GRUB_TIMEOUT=5          # 设置等待时间为 5 秒,可以根据需要调整
    GRUB_DEFAULT=0          # 默认启动第一个菜单项,可以根据需要调整
    ```

3. 保存并退出

    在 nano 中,按 `Ctrl + O` 保存文件,然后按 `Ctrl + X` 退出编辑器.

4. 更新 GRUB 配置

    ```bash
    sudo update-grub
    ```

5. 重启系统后生效

    ```bash
    sudo reboot
    ```

### 更换linux内核[危险操作不推荐]

并不是说内核版本越高越好,根据你需要的软件,有的时候需要根据需求换内核版本.

最新的ubuntu 24.04lts使用的是6.14内核,这个内核比较新,一些老旧设备可能不支持,比如最近比较火的MI50显卡,它最后一个支持的rocm版本是6.3.3,而这个版本的rocm只支持在ubuntu 24.04lts中6.8和6.11版本的内核上运行,因此我们就需要降级内核使用.当然如果你并不需要rocm,只要驱动,最新驱动一般也能支持最新的内核,那就没必要降级内核了.

一般会需要注意linux内核版本的情况主要就是显卡相关软件和虚拟机相关软件.我们需要根据自己的需求选择合适的内核版本.

#### 更新linux内核

如果你希望用上发行版支持的最新的linux内核,可以像下面这样执行

```bash
sudo apt update # 更新软件包的索引或包列表
sudo apt-get upgrade linux-image-generic #更新内核
sudo reboot
```

#### 安装指定版本的linux内核

1. 首先更新软件包列表：

    ```bash
    sudo apt update
    ```

2. 然后,列出所有可用的内核映像包:

    ```bash
    sudo apt-cache search linux-image
    ```

    这样你会找到一大堆可用的内核版本. 通常你会看到几个不同的包名它们可能带有不同的后缀,例如:

    + `linux-image-6.11.0-xx-generic`: 这是标准的通用内核,适用于大多数桌面和服务器.

    + `linux-image-6.11.0-xx-lowlatency`: 这是一个低延迟内核,专为音频制作,实时控制或其他需要极低延迟的应用设计

    + `linux-image-unsigned-6.11.0-xx-generic`: 这是一个未签名的内核,通常用于需要自定义内核模块或在特定硬件上调试时.对于大多数用户来说你不需要安装这个版本

    其中`xx`表示具体的版本号,比如`1019`,`1020`等,在大版本号不变的情况下,版本号越大表示内核更新的补丁越多.

3. 选择你需要的版本,比如说你想安装`linux-image-6.8.0-1019-generic`,你可以使用如下命令安装:

    ```bash
    sudo apt install linux-image-6.8.0-1019-generic
    ```

4. 安装完成后,你需要更新引导加载器以确保新内核被正确识别:

    ```bash
    sudo update-grub
    ```

    在大多数情况下,当你使用apt命令安装新内核时apt的安装脚本会自动为你运行update-grub,但为了确保万无一失,手动执行一次这个命令总是一个好习惯.

5. 最后,重启你的计算机以使用新内核:

    ```bash
    sudo reboot
    ```

6. 重启后,在grub菜单中选择`advance`,找到你刚刚安装的内核版本,选择用它启动系统即可.

7. 重新进入ubuntu后你可以使用`uname -r`命令来验证当前正在运行的内核版本:

    ```bash
    uname -r
    ```

    这应该显示你刚刚安装的内核版本.

#### 后遗症

最大的后遗症就是很多设备的驱动可能无法使用,比如鼠标,显卡等,然后就是虚拟机可能无法启动等.这些问题一般都是由于内核模块和驱动不匹配导致的.当然最简单的办法就是直接在安装好操作系统后就更换好内核,然后再安装显卡驱动和虚拟机等软件.

另外就是一些本来直接二进制分发的软件可能会因为内核版本不匹配而只能编译安装.

因此,如果你不是特别需要,我并不推荐更换内核.

## 修复依赖问题

下面是总结的常见依赖错误

### --configure出错

apt安装的包都是系统级别,最怕不知为啥忽然报冲突,一个典型的的报错像是

```bash
...
pkg: 处理软件包 xxxx (--configure)时出错：
子进程 已安装 post-installation 脚本 返回错误状态 1

...

dpkg: 处理软件包 xxxx (--configure)时出错：
依赖关系问题 - 仍未被配置
...
```

等等,碰到这种报错我们需要做一些手动操作来修复,思路是:

1. 先备份`dpkg`下的`info`
2. 然后自己创建一个新的空的`info`
3. 然后更新索引,让apt解决冲突
4. 再把新的info里的内容拷贝到旧的里面
5. 旧的再改成info,删除自己

```bash
sudo mv /var/lib/dpkg/info /var/lib/dpkg/info_old 
sudo mkdir /var/lib/dpkg/info
sudo apt-get update
sudo apt-get -f install  
sudo mv /var/lib/dpkg/info/* /var/lib/dpkg/info_old   
sudo rm -rf /var/lib/dpkg/info  
sudo mv /var/lib/dpkg/info_old /var/lib/dpkg/info
```

## 显卡驱动和计算库[2025-09-03更新]

要说和macos最大的区别恐怕就是linux下可以用独显了.这既是Linux系统相对macos的优势,也是麻烦的来源.因为显卡驱动和计算库的安装和配置往往是最麻烦的.

当然如果你只是要用来打游戏,装驱动就行,但如果你要用来跑AI相关的程序,那就需要安装计算库了,比如N卡的cuda和A卡的rocm.

注意:

1. 驱动和计算库是分开的,驱动负责显卡的基本功能,计算库负责提供GPU计算能力.但往往计算库会捆绑驱动,因此安装计算库时会需要安装对应版本的驱动.
2. 如果你只是要打游戏,那只需要安装最新驱动就行了,一般最新驱动都能支持老硬件,但如果你要跑AI相关程序,那就需要根据程序要求的计算库版本来安装对应版本的驱动和计算库了.
3. 单独安装最新驱动和安装计算库是两码事.
4. 除了硬件绑定的计算库(N卡的cuda和A卡的Rocm)外,还有一些通用的计算库,比如`opencl`,`vulkan`等,这些计算库不依赖于具体硬件,但性能可能不如硬件绑定的计算库.

### MOC

`MOC (Machine Owner Key)`是一种用于在 Linux 系统上签署内核模块的机制.它允许用户创建和管理自己的密钥对,并使用这些密钥对内核模块进行签名.这样,即使是第三方或自定义的内核模块,也可以被系统信任并加载,而不会被安全机制阻止.

无论是amd显卡还是N卡,第一次安装显卡驱动过程中会让你填一个密码,密码要求6~8位.这个密码即`MOK`,它是用来给驱动签名的,如果你不填这个密码,驱动是无法加载的.在设置好密码后,重启时会进入一个蓝色的界面让你选,注意,**默认第一个选项是让你直接进系统,不能选**,选第二个`Enroll MOK`,一路跟着引导最后会让你输入这个密码,输入后驱动就可以正常加载了.

这个moc的验证过程只会在第一次安装驱动时进行,后续更新驱动不会再让你输入密码了.

### 仅安装最新驱动

这里介绍的是最简单的安装最新驱动的方式,如果你只是要打游戏,那就足够了.

#### N卡驱动安装

ubuntu 24.04lts中安装n卡驱动已经相当简单,只需要在安装系统时选择`第三方驱动安装`即可,这样会自动安装n卡的社区驱动.这个驱动已经可以让你跑大部分的程序了.在关机后你可以将显示器换到n卡上,也可以取bios中关闭核显.

这是最不折腾的N卡驱动安装方式,你可以在`软件和更新(Software & Updates)`应用中切换到`附加驱动(Additional Drivers)`选项卡进行版本确认,也可以在terminal中使用`nvidia-smi`命令查看驱动版本.

要更新驱动也可以在`软件和更新(Software & Updates)`应用中进行更新.

#### A卡驱动安装

a卡的驱动需要在amd官网下载最新版本的安装器,然后运行安装器进行安装.

如果你是amd的消费级显卡或者专业卡可以在[这个页面](https://www.amd.com/zh-cn/support/download/linux-drivers.html)下载安装器.如果你是用的amd的服务器卡则需要去[这个网页](https://www.amd.com/zh-cn/support/download/drivers.html)在搜索框中搜索你的显卡型号,然后下载对应的linux驱动.注意我们需要选择`Ubuntu x86 64-Bit`这个板块下的东西,而且需要注意选择ubuntu的版本,比如你是ubuntu 24.04lts就需要选择最上面的和你系统匹配的比如`Radeon™ Software for Linux® version 25.10.2 for Ubuntu 22.04.5 HWE`这个版本的安装器

我们下载好后会得到一个`.deb`的安装包(比如`amdgpu-install_6.4.60402-1_all.deb`),然后我们可以使用如下命令进行安装

```bash
# 清理之前的残留(如果有的话)
amdgpu-install --uninstall
sudo reboot
# 彻底清理之前的amdgpu-install
sudo apt-get purge amdgpu-install
sudo reboot
# 更新系统确保驱动可以安装
sudo apt-get update
sudo apt-get dist-upgrade
sudo reboot

# 安装依赖
sudo apt install "linux-headers-$(uname -r)" "linux-modules-extra-$(uname -r)"
sudo apt install python3-setuptools python3-wheel
# 安装amdgpu-install
cd ~/Downloads # 进入下载目录
sudo apt-get install ./amdgpu-install_6.4.60402-1_all.deb
sudo apt-get update
# 安装驱动,这里以安装最新的`amdgpu-pro`驱动为例,如果你只是要安装开源的`amdgpu`驱动,可以去掉`--pro`参数
sudo usermod -a -G render,video $LOGNAME # 添加当前用户到渲染和视频分组
amdgpu-install --usecase=graphics
sudo reboot # 重启后生效
```

amd的最新驱动安装程序即便是已经不支持老旧硬件也是可以安装的,对应的rocm和hip也可以安装只是有很多功能无法使用.

```bash
amdgpu-install --usecase=graphics,rocm,hip
```

如果只是用来跑跑ollama什么的其实也够了.

### 安装计算库

这里介绍的是安装计算库的方式,如果你要跑AI相关程序,那就需要安装计算库了.安装计算库追求的是稳定,因此一般装了就很少会更新.而且很多性价比设备支持计算库已经停在了某个版本,因此我们需要根据自己的硬件和程序要求来安装对应版本的计算库.

#### N卡cuda安装

对于最新的ubuntu 24.04tls来说安装n卡的设置已经相当简单.但还是有如下几个注意点:

1. 硬件连接和bios设置--你需要cpu带核显,并且核显在bios中是开启的,在第一次安装ubuntu时显示器也需要连接在核显上(即主板的hdmi/dp接口上)
2. 安装ubuntu时选择第三方驱动安装,这样会自动安装n卡的社区驱动.这个驱动已经可以让你跑大部分的程序了.在关机后你可以将显示器换到n卡上,也可以取bios中关闭核显.
3. 如果你需要cuda支持,那社区的驱动就不行了,在较早的ubuntu版本中我们需要先卸载社区驱动,然后再安装官方驱动,但在24.04中我们可以直接安装官方驱动,它会自动覆盖社区驱动.安装命令如下

```bash
# 添加nvidia的官方源
wget https://developer.download.nvidia.com/compute/cuda/repos/ubuntu2404/x86_64/cuda-ubuntu2404.pin
sudo mv cuda-ubuntu2404.pin /etc/apt/preferences.d/cuda-repository-pin-600
# 下载cuda的本地安装包(以cuda 12.9为例)
wget https://developer.download.nvidia.com/compute/cuda/12.9.0/local_installers/cuda-repo-ubuntu2404-12-9-local_12.9.0-575.51.03-1_amd64.deb
sudo dpkg -i cuda-repo-ubuntu2404-12-9-local_12.9.0-575.51.03-1_amd64.deb
sudo cp /var/cuda-repo-ubuntu2404-12-9-local/cuda-*-keyring.gpg /usr/share/keyrings/
# 安装cuda
sudo apt-get update
sudo apt-get -y install cuda-toolkit-12-9
# 安装驱动
sudo apt-get install -y cuda-drivers
sudo reboot
```

需要注意:

1. 安装驱动可以安装有开源部分的`nvidia-open`也可以安装完全闭源的`cuda-drivers`.为了兼容性和稳定性考虑,我推荐安装完全闭源的`cuda-drivers`,它包含了所有的功能,而且通常稳定性也更好.
2. 安装好后可以使用`nvidia-smi`命令查看驱动和cuda是否安装成功.

关于cuda的版本选择,我们可以总结如下:

1. 截止到2025-09-03,pytorch最新稳定版本为2.8支持的最新的cuda版本是12.9,这个版本支持的显卡是从`GeForce GTX 1650`开始的,如果你的显卡比这个老,那就只能安装老版本的cuda了.它支持最新的ubunut 24.04.03lts使用的是6.14的内核,最低支持的n卡驱动版本是575版本

2. 如果你需要兼容一些老显卡,那你就必须给linux内核降级,比如你需要给MI50显卡安装rocm,那你可以把内核降级到6.11,然后安装n卡驱动570版本,再安装cuda 12.8版本.

    ```bash
    ...
    wget https://developer.download.nvidia.com/compute/cuda/12.8.0/local_installers/cuda-repo-ubuntu2404-12-8-local_12.8.0-570.86.10-1_amd64.deb
    sudo dpkg -i cuda-repo-ubuntu2404-12-8-local_12.8.0-570.86.10-1_amd64.deb
    sudo cp /var/cuda-repo-ubuntu2404-12-8-local/cuda-*-keyring.gpg /usr/share/keyrings/
    sudo apt-get update
    sudo apt-get -y install cuda-toolkit-12-8
    ...
    ```

#### A卡rocm安装

对于amd的显卡来说安装rocm还是比较简单的,你如果是比较新的amd官方支持rocm的显卡,那直接在安装好最新安装器后执行如下命令即可

```bash
amdgpu-install --usecase=rocm,graphics,hip
sudo reboot
```

至于如何确定自己的卡是否官方支持rocm,可以查看[这个网页](https://rocm.docs.amd.com/projects/install-on-linux/en/latest/reference/system-requirements.html)

如果你是比较老的amd显卡,那就需要指定版本安装了,比如我现在用的MI50显卡,它最后一个支持的rocm版本是6.3.3,而这个版本的rocm只支持在ubuntu 24.04lts中6.8和6.11版本的内核上运行,因此我们就需要降级内核,然后下载6.3.3的安装器.

```bash
sudo apt update
sudo apt install "linux-headers-$(uname -r)" "linux-modules-extra-$(uname -r)"
sudo apt install python3-setuptools python3-wheel
sudo usermod -a -G render,video $LOGNAME # Add the current user to the render and video groups
wget https://repo.radeon.com/amdgpu-install/6.3.3/ubuntu/noble/amdgpu-install_6.3.60303-1_all.deb
sudo apt install ./amdgpu-install_6.3.60303-1_all.deb
sudo apt update
# 安装驱动,这里以安装最新的`amdgpu-pro`驱动为例,如果你只是要安装开源的`amdgpu`驱动,可以去掉`--pro`参数
sudo usermod -a -G render,video $LOGNAME # 添加当前用户到渲染和视频分组
amdgpu-install --usecase=graphics,rocm,hip
sudo reboot # 重启后生效
```

#### vulkan[2025-09-09更新]

现在的通用计算库中vulkan是最推荐的,它既支持n卡也支持a卡,i卡也支持,而且性能也不错,生态算是比较完善的.linux上跑游戏基本就是靠的vulkan.而且通常各家显卡的驱动中都会捆绑vulkan的支持.也就是说如果你只是要调用vulkan(所谓runtime运行时)那并不需要单独安装vulkan.但如果要查看vulkan的支持情况我们可以安装`vulkan-tools`,它提供了`vulkaninfo`命令可以查看vulkan的支持情况.

```bash
# 安装vulkan-tools
sudo apt-get install vulkan-tools

# 使用vulkaninfo查看vulkan支持情况
vulkaninfo
```

如果你想查看目前的vulkan环境是否支持amd显卡或NVIDIA显卡,可以使用如下命令

```bash
# 查看amd显卡支持情况
vulkaninfo | grep "AMD"
# 查看NVIDIA显卡支持情况
vulkaninfo | grep "NVIDIA" 
```

##### vulkan sdk

如果你要开发vulkan程序(所谓sdk开发包)那就需要单独安装vulkan sdk.

在ubuntu 24.04lts中安装vulkan的sdk推荐使用[vulkan官方提供的sdk压缩包](https://vulkan.lunarg.com/sdk/home#linux)
,下载后解压到某个目录,然后设置环境变量即可.

```bash
# 假设你下载并解压到了~/vulkan-sdk目录
export VULKAN_SDK=~/vulkan-sdk/x.x.x.x/x86_64
export PATH=$VULKAN_SDK/bin:$PATH
export LD_LIBRARY_PATH=$VULKAN_SDK/lib:$LD_LIBRARY_PATH
# export VK_ICD_FILENAMES=$VULKAN_SDK/etc/vulkan/icd.d:$VK_ICD_FILENAMES
# export VK_LAYER_PATH=$VULKAN_SDK/etc/vulkan/explicit_layer.d:$VK_LAYER_PATH
```

### 多显卡问题[2025-09-09更新]

我们装ubuntu需要核显,对于消费级主机,核显大致可以分为两类

+ 亮机核显,就是intel的核显和amd除了apu系列外的核显,这种核显性能很差,基本不能指望做计算任务,但足够带得起显示器.
+ 计算核显,就是amd的apu系列的核显,这种核显性能还不错,可以做一些轻量级的计算任务,但也不能指望它能做重型计算任务.

通常多显卡的情况主要是如下几种情况:

+ 亮机核显+a卡/n卡,最常见的情况,在linux下一般也是屏蔽核显使用的,所以可以忽略对这个配置的讨论.
+ apu+a卡/n卡,apu本身带一个还不错的核显,而且由于砍pcie通道所以一般最多配一张显卡而且基本是70ti级别以下的独显,这种配置往往并不是追求性能的方案而是追求性价比的方案,往往是渐进式升级的结果,比如先买个apu,然后再加个独显.也算一种挺常见的配置.我们会在下面讨论单独讨论如何针对这种配置进行设置.
+ 亮机核显+双n卡,一般是将消费级平台当小型工作站的用户,通常是用来跑ai的,这种配置在linux下也比较常见,我们也会在下面讨论如何针对这种配置进行设置.
+ 亮机核显+n卡+a卡,比较奇葩的配置,但其实是一个万金油配置,也是可玩性比较高的一个配置,我们下面会有针对性的进行讨论.
+ 亮机核显+双a卡,很少见,一般是用来跑ai推理的,这里就不讨论了.

对于多显卡用户,我们需要注意如下几点:

1. 如果是amd的带核显的cpu,它的核显也会受amd驱动的影响.因此如果安装amd驱动失败,核显也不会工作.因此建议配张n卡亮机卡,并且在安装系统时勾选`第三方驱动安装`,这样就不会因为amd驱动问题导致屏幕不亮无法进入系统了(可以安装amd驱动后屏幕接n卡以确保屏幕能亮).
2. 打游戏基本依赖vulkan,因此多显卡用户打游戏时需要确保vulkan能正常工作.

### apu+a卡/n卡

对于apu+a卡/n卡的配置,我们主要看独显是amd的还是nvidia的

> amd独显

apu的核显也会受amd驱动的影响,因此实际上最好别配amd的独显,如果只有amd的独显,安驱动出问题很可能完全亮不起来.

> nvidia独显

对于独显为N卡的情况,我建议先安装amd的驱动,确保核显能用,然后再安装n卡的驱动.
  
#### 充分利用apu的核显

为了充分利用apu,我建议用独显做主力卡,但apu也不要屏蔽核显,更加推荐显示器接核显使用混合模式.这并不会有什么坏处,反而在一些场景下可能还是更好的选择.原因有:

+ 显示器接核显可以减少独显的负担,尤其是对于一些4k显示器来说,独显不需要处理桌面显示的任务,可以专心做计算任务.
+ 可以利用比如小黄鸭(Lossless Scaling)插帧,插帧对显卡算力要求不高,核显完全可以胜任,独显可以专心做计算任务;同时核显比走pcie的独显离cpu更近,输出延迟可以得到更好的控制.这在游戏场景很常用(steam部分会做进一步介绍).
+ 在重型任务时我们主要用独显,空出来的核显可以用来跑一些轻量级任务,比如显卡跑渲染的时候拿核显打打游戏,或者显卡跑训练的时候拿核显跑点推理什么的,这样可以更充分利用硬件资源.
+ 核显使用内存,可以用来跑像qwen3 30BA3这样的吃显存但激活层少的moe模型,消费级显卡很少有足够的显存可以跑得动这种模型,但核显可以,还不影响独显的使用.

### 亮机核显+双n卡

这种搭配一般是用来跑ai的,通常是用来跑训练的,因此我们主要关注的是如何让两张卡都能被程序识别并使用.主要是围绕独显的驱动和计算库来设置的.核显建议不要屏蔽,显示器接核显使用混合模式,理由和上面apu的理由类似.这个核显只要不打游戏,看看视频足够了,

至于amd的核显要不要装驱动,我个人建议不要装,amd的亮机核显性能太差,装了驱动反而可能会带来不必要的麻烦,而且也不会带来什么性能提升.

### 亮机核显+n卡+a卡

这种配置比较奇葩,但其实是一个万金油配置,也是可玩性比较高的一个配置.这种配置可以让你同时利用n卡和a卡的优势,

+ n卡有cuda,对于ai相关生态来说非常好用,尤其是训练方面
+ a卡有rocm,hip,天生对opencl有更好的支持,很多工业软件对opencl的支持也不错,而且a卡的显存一般都比较大,对于一些大模型推理来说更有优势.
+ 核显可以用来跑一些轻量级任务,带带显示器看看视频什么的.
+ 想在打游戏时充分利用两张卡的性能可以像apu插帧一样借助小黄鸭(Lossless Scaling)来实现,具体就是

    1. 显示器接a卡
    2. 用小黄鸭(Lossless Scaling)将游戏窗口放大到全屏,并开启插帧功能
    3. 渲染用n卡跑,插帧用a卡跑

    这样就可以充分利用两张卡的性能了.

但这种配置坑是最多的,主要是驱动和计算库的冲突问题.因此我建议像下面这样安装:

+ 插满卡后核显接显示器安装ubuntu,并在安装过程中勾选`第三方驱动安装`.
+ 安装好系统后先安装amd的驱动和计算库,然后关机显示器改接n卡开机.确保驱动安装好核显能用.
+ 然后关机,显示器再接回核显,再安装n卡的驱动和cuda.注意安装n卡驱动时选择完全闭源的`nvidia-driver`,不要选择带开源部分的`nvidia-open`,因为`nvidia-open`会缺失vulkan的支持,会影响游戏等程序使用n卡运行.
+ 安装好后重启,然后使用`nvidia-smi`和`rocminfo`命令查看两张卡是否都能被识别.


#### 多显卡的功耗限制

多显卡情况下功耗往往是一个问题,一般来说apu功耗不会太高,120w基本就跑满了,但独显现在一个个都是耗电大户,因此如果是双独显方案我们需要先查清楚独显的功耗墙并准备一个足够冗余的大电源,建议1200w以上. 

除了电源外,我们还可以通过软件来限制显卡的功耗,以确保电源不会超载.

> AMD显卡功耗限制

对于amd显卡来说可以使用`rocm-smi`命令来限制功耗,这个命令需要安装`rocm`.比如我有一张MI50,它的功耗墙是220w,我可以将它限制在150w以内,性能大致损失在5%以内,但可以大幅降低功耗和发热.

```bash
rocm-smi --setpoweroverdrive 150 -d [显卡ID]
```

这个`显卡ID`可以通过`rocm-smi`命令查看,一般从0开始编号.先根据信息判断下要设置卡的编号,然后使用`rocm-smi --id [GPU_ID]`进行确认.没问题了再设置功耗.

> nvidia显卡功耗限制

对于n卡,限制功耗的功能在驱动中已经集成了,我们可以使用`nvidia-smi`命令来限制功耗.比如我有一张3080,它的功耗墙是320w,我可以将它限制在200w以内,性能大致损失在5%以内,但可以大幅降低功耗和发热.

```bash
nvidia-smi -i 0 -pl 200
```

这个`-i 0`表示对编号为0的显卡进行设置,如果有多张卡可以通过`nvidia-smi`命令查看编号,然后对每张卡分别设置即可.

> 持久化显卡功耗限制

上面介绍的设置都是临时设置,在机器重启后会失效,如果要持久化功耗限制,可以利用`systemd`来实现.

>> n卡限制功耗的设置:

1. 创建服务文件，例如 `/etc/systemd/system/nvidia-powerlimit.service`:

    ```conf
    [Unit]
    Description=Set NVIDIA GPU power limit on startup
    After=multi-user.target

    [Service]
    Type=oneshot
    ExecStart=/usr/bin/nvidia-smi -i 0 -pl 200  # 将参数改为你想要的功耗和显卡ID

    [Install]
    WantedBy=multi-user.target
    ```

2. 启用并启动服务

    ```bash
    sudo systemctl enable nvidia-powerlimit.service
    sudo systemctl start nvidia-powerlimit.service
    ```

>> a卡限制功耗的设置:

1. 创建服务文件，例如 `/etc/systemd/system/amd-powerlimit.service`:

    ```conf
    [Unit]
    Description=Set AMD GPU power limit on startup
    After=multi-user.target

    [Service]
    Type=oneshot
    ExecStart=/usr/bin/rocm-smi --setpoweroverdrive 150 -d 0  # 将参数改为你想要的功耗和显卡ID

    [Install]
    WantedBy=multi-user.target
    ```

2. 启用并启动服务

    ```bash
    sudo systemctl enable rocm-powerlimit.service
    sudo systemctl start rocm-powerlimit.service
    ```

## 基础设置

在开始其他设置之前,我们还需要做如下设置

### git的安装设置

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

### 配置代理

作为一个写代码的,代理几乎波不可少.作为一台主力办公机,没道理不按一个clash客户端.[clashverge](https://github.com/clash-verge-rev/clash-verge-rev/releases/tag/v1.7.7)自然是ubuntu下的首选,安装也只是下载下来,先安装依赖再安装本体即可,具体可以看[我的这一篇博文](https://blog.hszofficial.site/introduce/2024/11/29/clash%E9%83%A8%E7%BD%B2/#linuxdebian%E7%B3%BB%E4%B8%8B%E7%9A%84%E5%AE%89%E8%A3%85).

和在macos上一样,我们可以在`.zshrc`或其他shell的配置文件中像这样配置代理

```bash
#========================================================================= proxy
# 设置使用本机代理
alias setproxy="export https_proxy=http://127.0.0.1:7897 http_proxy=http://127.0.0.1:7897 all_proxy=socks5://127.0.0.1:7897"
# 设置使用本地局域网代理
alias setlocalproxy="export https_proxy=http://192.168.50.177:7890 http_proxy=http://192.168.50.177:7890 all_proxy=socks5://192.168.50.177:7890"
# 清空代理配置
alias unsetproxy="unset https_proxy;unset http_proxy;unset all_proxy"
```

### ssh客户端设置

+ `.ssh/config`

    ```txt
    ControlPersist 4h
    ServerAliveInterval 30
    ControlMaster auto
    ControlPath /tmp/ssh_mux_%h_%p_%r
    ForwardX11Trusted yes

    HostKeyAlgorithms +ssh-rsa
    PubkeyAcceptedKeyTypes +ssh-rsa

    ```

### 安装gnome-tweaks

在开始设置之前，我们需要先安装`gnome-tweaks`

```bash
sudo apt install gnome-tweaks
```

这是gnome环境的`优化`工具,安装好后会放在`显示应用`->`工具`->`优化`.

## 美化terminal

作为linux最重要的当然是terminal.ubuntu默认使用的是`Gonme termial`,美化terminal其实也分成两个部分

### terminal本体的美化

依然是使用`Solarized Dark`.在`Default->颜色->内置方案`中找到并设置即可.剩下的就是根据个人习惯做细微修改,
比如我会将文本颜色和光标改为更醒目的白色,将粗体字改为紫色,高亮字改为橙色,将光标改为竖线,并稍微设置点透明度.
ubuntu 默认使用的是,他只支持统一的透明度,比较可惜.

### zsh的美化

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

## 美化字体

ubuntu中字体分为`界面字体`,`文档字体`,`等宽字体`.所谓`界面字体`就是gui界面上展示的字体,比如应用title,菜单栏的字体;所谓`文档字体`就是打字时的字体;`等宽字体`则是每个字符固定宽度的字体,在一些特殊位置会用到.

在中文环境中,我们对字体的要求自然是对中问支持良好,而作为编程写文档为主用途来说,字体优雅辨识度高就是重中之重. 下面是我推荐的几种字体

+ [思源宋体/黑体](https://github.com/notofonts/noto-cjk),思源宋体/黑体是由Adobe和Google两家公司合作开发的中文字体,具有良好的字形结构和视觉效果.该字体支持的字符编码达到了GB18030字符集的完整覆盖,在Linux系统中运行稳定且无需进行额外的配置,舒适自然,适用于流畅的中英文混排场景.这个字体适合用作`界面字体`或者`文档字体`

    ```bash
    # `noto sans cjk sc`为黑体,`noto serif cjk sc`为宋体
    sudo apt-get install fonts-noto-cjk
    ```

+ [nerd-fonts](https://github.com/ryanoasis/nerd-fonts),是一套专门为编程设计的字体,适合给terminal用

    ```bash
    cd workspace/init_source
    git clone https://github.com/ryanoasis/nerd-fonts.git --depth 1
    ```

### 设置字体

设置全局字体在`显示应用`->`工具`->`优化`.选中`字体`这个位置,我个人习惯设置成如下:

+ `界面文本`: `Noto Sans CJK SC Meduim`
+ `文档文本`: `Noto Serif CJK SC Meduim`
+ `等宽文本`: `Ubuntu San Font`

之后我们可以进terminal,在`首选项`->`default`->`自定义字体`中选择`Ubuntu Normal Nerd Font`.

## 系统优化

ubuntu默认状态下是很原始的,我们需要做如下操作才能让它用起来舒服些

### 增加空格键预览功能

mac下一个经典操作就是选中目标后按空格键

```bash
sudo apt-get install gnome-sushi unoconv
```

### 常用软件预加载

[preload](https://www.fosslinux.com/135006/how-to-use-preload-to-speedup-app-launches-in-ubuntu.htm)是一个空间换时间的应用预加载工具.
它会根据用户行为预先加载常用软件进内存.对于拥有大内存的机器这个工具可以大幅提高应用的打开速度.

我们可以使用如下命令在terminal中安装

```bash
sudo apt-get install preload
```

`preload`安装好后会被`systemd`统一管理.

### 全面的编解码器支持

ubuntu受限于版权并不能直接提供全部的编解码支持,但我们可以安装`ubuntu-restricted-extras`来或者这一能力

```bash
sudo apt install ubuntu-restricted-extras
```

安装好后市面上大部分的音频视频格式我们就都可以使用了

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
3. 在`/etc/environment`中设置环境变量

    ```bash
    LANG=zh_CN.UTF-8
    LANGUAGE=zh_CN:en_US

    XMODIFIERS=@im=fcitx
    QT_IM_MODULE=fcitx
    GTK_IM_MODULE=fcitx
    ```

    之后重启

4. 在`显示应用->工具->优化`中将`Fcitx 5`添加自启动项,重启即可

5. 安装Gnome插件`Input Method Panel`
6. 着全局代理执行如下命令

    ```bash
    curl -sSL https://www.debuggerx.com/fcitx5_customizer/fcitx5_customizer.sh | bash -s -- recommend
    ```

7. 之后除了去掉`搜狗词库`外都按推荐的来即可

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

这么改有个缺陷就是只能使用`Control + Alt + C`在terminal中中断程序了

## 美化系统

ubuntu系统界面大致可以分为如下几个部分

![系统界面][1]

美化系统我们大致可以分为如下几个步骤

1. 美化桌面
2. 美化登录页面
3. 添加实用插件

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

一个是主题美化,这是一个很复杂的问题,我们知道linux下生态是很'去中心化'的,因此美化的效果很可能一些生态下支持一些就不支持.主题美化不光是自带软件的美化,还得让各种gui应用都可以获得相应的美化,这就很难了,这就会造成一个体验上的不统一.这个真没办法,我们只能尽量覆盖

对于我们这种要求不高的其实就很简单.直接执行`WhiteSur-gtk-theme`项目的`./install.sh`脚本即可

```bash
cd workspace/init_source/WhiteSur-gtk-theme
./install.sh
./install.sh -t all  # 安装指定颜色主题,如果要全部颜色可以使用`-t all`,要指定颜色则是类似`-t [purple/pink/red/orange/yellow/green/grey]`
./install.sh -N mojave # 改变文件管理器分栏样式,可选为默认,`mojave`和`glassy`
./install.sh -l  # 安装对`libadwaita`软件的适配,目前并不完美
sudo flatpak override --filesystem=xdg-config/gtk-3.0 && sudo flatpak override --filesystem=xdg-config/gtk-4.0 # 适配flatpak应用
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

安装好的壁纸我们需要重启后去`设置->外观->背景`中选择使用.

#### 美化图标

然后是美化图标.实话讲ubuntu的图标确实丑.我们可以使用[vinceliuice/WhiteSur-icon-theme](https://github.com/vinceliuice/WhiteSur-icon-theme)项目提供的图标来美化它,这个图标库就很还原macos了.

```bash
cd workspace/init_source
git clone https://github.com/vinceliuice/WhiteSur-icon-theme.git
cd WhiteSur-icon-theme
./install.sh
./install.sh -a # 安装macos风格的替换图标
./install.sh -b # 安装右上角下拉菜单的图标
```

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

但即便为了Gnome插件,firefox也是值得优化下体验的,因为我们还要靠它管理本地的gnome插件.这个优化主要包括3个方面

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

浏览器要作为管理本地插件的工具显然光靠一个浏览器插件是不够的.我们需要先安装一个专用连接器

```bash
sudo apt-get install gnome-browser-connector
```

之后在firefox中安装插件[GNOME Shell integration](https://addons.mozilla.org/zh-CN/firefox/addon/gnome-shell-integration/?utm_source=addons.mozilla.org&utm_medium=referral&utm_content=search).安装好后firefox的右上角插件栏中就会有一个脚印一样的图标,它就是gnome的浏览器插件,点击它就可以进入插件搜索页面.

#### gnome插件安装

安装gnome插件很简单,用firefox的gnome浏览器插件进入到gnome插件页面后点击插件名后面的开关到开的状态即可.

当安装好后我们可以在`扩展`应用中对插件进行开关和设置,而已经安装了哪些插件可以在[installed extentions页面中查看](https://extensions.gnome.org/local/).

需要注意我们的系统`ubuntu 2024.04`使用的是`Gnome 46`,插件需要支持这个版本才能安装.

下面是我认为比较有必要的gnome插件汇总

| 插件                                                                                                          | 推荐等级 | 用途                               | 补充说明                                                                                                                             |
| ------------------------------------------------------------------------------------------------------------- | -------- | ---------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| [user-themes](https://extensions.gnome.org/extension/19/user-themes/)                                         | 高       | 管理用户主题                       | ---                                                                                                                                  |
| [Dash to Dock](https://extensions.gnome.org/extension/307/dash-to-dock/)                                      | 高       | 一个对主题更友好的dash             | 关闭`Ubuntu Docker`,功能重复了                                                                                                       |
| [Blur my Shell](https://extensions.gnome.org/extension/3193/blur-my-shell/)                                   | 高       | 一个提供桌面模糊的插件             | 建议修改`Dash to Dock`中的拐角半径                                                                                                   |
| [Clipboard Indicator](https://extensions.gnome.org/extension/779/clipboard-indicator/)                        | 高       | 剪切板功能,可以保存近期的复制内容  | ---                                                                                                                                  |
| [Lock Keys](https://extensions.gnome.org/extension/1532/lock-keys/)                                           | 低       | 大小写锁定提示                     | ---                                                                                                                                  |
| [Removable Drive Menu](https://extensions.gnome.org/extension/7/removable-drive-menu/)                        | 高       | 顶栏的移动存储操作工具             | ---                                                                                                                                  |
| [Screenshort-cut](https://extensions.gnome.org/extension/6868/screenshort-cut/)                               | 低       | 顶栏截图工具                       | 需要额外安装`sudo apt install gir1.2-gtop-2.0 lm-sensors`获取硬盘信息                                                                |
| [Vitals](https://extensions.gnome.org/extension/1460/vitals/)                                                 | 中       | 顶栏系统监控                       | ---                                                                                                                                  |
| [GSConnect](https://extensions.gnome.org/extension/1319/gsconnect/)                                           | 高       | 快速连接移动端设备                 | 需要配合app`kde connect`                                                                                                             |
| [Todoit](https://extensions.gnome.org/extension/7538/todo-list/)                                              | 低       | 顶部todolist                       | ---                                                                                                                                  |
| [Lunar Calendar 农历](https://extensions.gnome.org/extension/675/lunar-calendar/)                             | 中       | 日历改为农历                       | 需要先额外安装[Nei/ChineseCalendar](https://gitlab.gnome.org/Nei/ChineseCalendar/-/archive/20240107/ChineseCalendar-20240107.tar.gz) |
| [Compiz alike magic lamp effect](https://extensions.gnome.org/extension/3740/compiz-alike-magic-lamp-effect/) | 中       | 仿macos的最小化动画                | ---                                                                                                                                  |
| [Forge](https://extensions.gnome.org/extension/4481/forge/)                                                   | 低       | 多应用划分窗口                     | ---                                                                                                                                  |
| [Input Method Panel](https://extensions.gnome.org/extension/261/kimpanel/)                                    | 高       | 输入法相关                         | ---                                                                                                                                  |
| [Click to close overview](https://extensions.gnome.org/extension/3826/click-to-close-overview/)               | 高       | 点击空白处关闭预览                 | ---                                                                                                                                  |
| [Hide Top Bar](https://extensions.gnome.org/extension/545/hide-top-bar/)                                      | 低       | 自动隐藏顶部工具栏                 | ---                                                                                                                                  |
| [desktop-lyric](https://extensions.gnome.org/extension/4006/desktop-lyric/)                                   | 中       | 桌面歌词                           |
| [applications-menu](https://extensions.gnome.org/extension/6/applications-menu/)                              | 低       | 顶部提供应用的归类入口             |
| [weather-or-not](https://extensions.gnome.org/extension/5660/weather-or-not/)                                 | 低       | 顶部天气插件,需要有`gnome weahter` |
| [GNOME Fuzzy App Search](https://extensions.gnome.org/extension/3956/gnome-fuzzy-app-search/)                 | 中       | 模糊搜索工具                       |
| [No overview at start-up](https://extensions.gnome.org/extension/4099/no-overview/)                           | 高       | 开机后自动进入第一个桌面           |

除此之外,我个人推荐对系统默认插件做如下处理

+ 如果安装了`dash to dock`我们最好把它关掉,因为自带插件`Ubuntu Dock`会在安装了`dash to dock`的情况下直接调用它,开了就相当于启动了两份`dash to dock`会造成混乱.

+ 如果安装有`applications-menu`则在全部配置完成后`dash to dock`插件中关闭`显示应用程序`

+ `weather-or-not`插件依赖`gnome weahter`,我们需要先安装`gnome weahter`之后打开,输入自己所在的城市(用中文搜),重启后就可以激活,推荐将插件的位置设置为`clock left center`

+ [可选]禁用Desktop Icons,这个插件会让桌面有图标(默认会有你的home目录文件夹).如果无法关闭(我就碰到了这种情况)我们可以进去它的设置把`在桌面显示个人文件夹`关掉

## 安装docker环境

可以这么说,原生的docker环境是linux系统最大的竞争力之一,docker自然是要装的,而且必须装原生docker!

+ 添加docker源

    ```bash
    # Add Docker's official GPG key:
    sudo apt-get update
    sudo apt-get install ca-certificates curl
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc

    # Add the repository to Apt sources:
    echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
    $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
    sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    sudo apt-get update
    ```

+ 安装docker

    ```bash
    sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    ```

+ 设置非root权限使用docker

    ```bash
    sudo groupadd docker
    sudo usermod -aG docker $USER
    newgrp docker
    ```

+ 修改docker配置增加镜像站

    ```bash
    sudo nano /etc/docker/daemon.json
    ```

    ```json
    {
        "registry-mirrors":[
            "https://hub.geekery.cn",
            "https://hub.littlediary.cn",
            "https://registry.dockermirror.com"
        ]
    }
    ```

+ 安装轻量级docker监控工具[whaler](https://flathub.org/apps/com.github.sdv43.whaler),这个工具可以监控系统中的容器和compose-stack.已经相当够用.

+ 重启后生效,然后就可以使用了

## 安装常用环境

ubuntu的apt是最常用的环境安装工具,但apt工具安装的包和应用都是系统级的,一旦出问题那就问题大了.因此遵循macos上的习惯,我们装[homebrew](https://docs.brew.sh/Installation).

*注意*:

1. `arm`和`x86`架构无法使用`homebrew`,只有`x86_amd64`才可以
2. linux下无法使用`Homebrew Cask`

```bash
# 安装依赖
sudo apt-get install build-essential procps curl file git

# 安装homebrew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
# 配置环境
echo >> /home/hsz/.zshrc
echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"' >> /home/hsz/.zshrc
eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"

# 测试安装没问题
brew install hello
hello
```

`homebrew`会被安装在`linuxbrew`用户目录下(根目录为`/home/linuxbrew`),而安装的包都会被放在`/home/linuxbrew/.linuxbrew/Cellar`目录下

安装好`homebrew`后我们自然也就需要安装`cmake,protobuf,grpc,go,node,python3.11,micromamba`等等这些基本环境了.

修改我们的`.zshrc`以固定各个包或工具的配置

```bash
#========================================================================= golang
# 主要是将go的mod缓存和gopath改到`Libraries`目录
export GO111MODULE=on
export GOPROXY=https://goproxy.cn
export GOPATH=/home/hsz/Libraries/go
export GOMODCACHE=/home/hsz/Libraries/go/pkg/mod
export PATH="$PATH:/home/hsz/Libraries/go/bin"
```

## 安装和管理Linux应用

linux上目前应用的分发是相当碎片化的,为了可以安装管理各种渠道的应用我们就得做对应的优化.不考虑虚拟机中使用应用,我们可以将应用分为原生应用,沙盒化应用两种.

### 沙盒化应用

应对应用的安全性问题,盒化方案在厂商看来更有优势,所以目前的三个主流沙盒化方案中有两个由厂商主导,剩下的一个也有厂商背书,而国产的统信系统使用的应用管理工具`玲珑`也是这个路线.

沙盒化应用大多可以理解为更轻量的docker.最基础的就是通过Linux的容器技术(`namespace`)隔离应用与操作系统,再单独为应用提供一套统一的lib作为runtime.不同的方案可能还会隔离应用的静态文件,做额外的权限控制等等.

沙盒化由于都多少使用了容器技术,计算性能多少会受一些影响(io性能基本没啥影响),但优点在于环境统一,这样分发效率就高太多了.
所有的沙盒化应用方案都是相似的结构: 服务端提供应用仓库;客户端提供命令行或gui的应用管理工具;再提供一套打包工具给开发者就齐活了.看看docker是不是似曾相识?就知道为啥厂商都爱这套方案了.

#### Flatpak

`Flatpak`是`Gnome`社区推的一套沙盒化应用方案,背后是老牌开源解决方案供应商`redhat`,主打的就是一个开源.也因为开源,linux社区普遍对它风评还不错.

`Flatpak`没有`snap`那么严格,包的没那么多,但性能损失也和`snap`差不太多,在10%~15%的程度,只是相同的应用通常比`snap`少损失个1%左右.至于制成品包的大小,`Flatpak`算的是除去各种依赖层各种runtime后的大小,因此看起来会比`snap`小的多,但熟悉docker的都知道这只相当于docker镜像的最外面一层.因此一般装`Flatpak`应用它还会下一堆依赖runtime什么的,总体看时间差不多.

> 安装和配置

在ubuntu中我们要安装`Flatpak`应用需要先安装如下内容:

+ `Flatpak`客户端工具

    ```bash
    sudo apt install flatpak
    ```

+ `Gnome`环境下的`Flatpak`支持插件

    ```bash
    sudo apt install gnome-software-plugin-flatpak
    sudo reboot
    ```

> 设置`flatpak`仓库

虽然`flatpak`号称去中心化完全开源,但也还是有一个社区维护的主仓库的,这就是[flathub](https://flathub.org).
`flathub`基本就是一个网页版的`flatpak`应用中心,不想用命令行安装`flatpak`应用,我们可以直接上去搜索找到要装的应用,点击`Install`按钮下载安装文件,然后双击就可以安装了.

如果要用命令行工具`flatpak`安装,则还需要做如下设置将`flathub`添加到仓库索引:

```bash
flatpak remote-add --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo
sudo reboot
```

`flathub`明显是一个境外网站,因此网速往往尴尬,尤其它很多运行时源用的是github,因此很多时候

> 管理`Flatpak应用`

为了更好的管理flatpak应用,我们最好安装[Warehouse](https://flathub.org/apps/io.github.flattool.Warehouse)来统一管理.

安装方式也很简单,在网页上点击`Install`按钮下载安装文件,然后双击就可以安装了.不过我更喜欢用命令来安装

```bash
flatpak install flathub io.github.flattool.Warehouse
```

#### snap

`Ubuntu`的开发公司`Canonical`搞出来的一套沙盒化应用方案,也是`Ubuntu`现在主推的标准,目前`Ubuntu`的应用中心也只能管理`snap`软件.

`snap`有最严格的权限管理,而且它的包仓库是`Canonical`自己维护的并不开源,也就是说你也没办法自己部署一个`snap`的私有仓库,这在开源精神泛滥的linux社区中那是触了逆鳞了,这也是`snap`最不受人待见的一点.

然后就是性能损失.由于`snap`包的最多,性能损失也是最严重的,相比原生版本一般有10%~15%的性能(计算性能)损失,而它的打包方式也类似docker,是包含完整runtime的,因此制成品包也是最大的.

虽然说了这么多缺点,但它作为`Ubuntu`默认的应用分发模式体验其实还算不错--开箱即用,下载速度也还不错,个人认为没啥必要黑它.

### 原生应用

原生应用直接运行在裸机上,提供了最大程度灵活性和最大程度资源利用,但问题就在于安全性.由于直接接触系统,原生应用之间可能会相互影响,造成依赖冲突等问题,一个不小心甚至可能将系统整崩溃.

#### deb

ubnutu是从debian发展而来的,继承了其deb应用的生态.我们可以使用`apt`工具很轻易的安装库或者应用,而这些都是使用`deb`分发的.

由于apt工具会处理依赖,`deb`包往往大量相互依赖而本体非常小.但请注意,如果不是万不得已,用户级别的应用我是不推荐用deb安装的,因为一旦依赖出了问题那就很难修复,想要重回干净的系统基本只能重装系统了,而python用户应该都清楚依赖冲突问题是多么容易发生.

管理deb软件包我们可以安装`synaptic`.

```bash
sudo apt install synaptic ### 安装资源监控工具
```

#### PORTABLE LINUX APPS

`PORTABLE LINUX APPS`这种"独立自主"路线固然是一个方案,但显然并不是唯一方案,沙盒化应用就是另一个方案.显然沙
所谓`PORTABLE LINUX APPS`是在GNU/Linux上可以独立运行(理论上)的应用程序,这种应用程序可以在任何地方运行甚至在U盘上都可以.

这些应用程序可以是`AppImage应用`,也可以是独立归档应用(比如`Firefox`,`Blender`,`Thunderbird`).

> 独立归档应用

`独立归档应用`本质上是一个文件夹,它里面包含了应用需要的所有资源,从静态资源到依赖包,而其中的可执行文件可以不依赖系统本身的lib就可以执行.

这种应用可以在所有linux发行版中直接使用的.

一般`独立归档应用`并不多见不多,只有重型应用像上面例子里的firefox,blender才会支持这样分发.

一般如果我们用的是开源软件,可以在github上找找,看看它的`release`中是否有给linux使用的非source压缩包,如果有一般就是`独立归档应用`,下载下来解压出来根据文档设置就可以正常使用了.

> AppImage应用

[AppImage应用](https://appimage.org/)是独立归档应用更进一步的产物,它提供了一套规范将一个完整的包含所有依赖的应用打包成了一个文件.这个文件除了可以提取出应用所需的文件外还可以直接赋予可执行权限后直接执行.

这个分发方式就像windows上的`绿色软件`和macos中应用的结合,可以说相当创新.但很可惜似乎从2019年开始发展就慢下来了.不过我个人依然推荐如果要开发linux下的应用,可以将这种方式作为一个分发选项.

一般如果我们用的是开源软件,可以在github上找找,看看它的`release`中是否有以`.AppImage`为后缀的文件.有的话就是了,下载下来就可以直接使用了.

注意`AppImage`实际有两个版本,版本1基本已经废弃,只会在一些已经停止维护的古早软件上有机会看到,这里讨论的都是版本2的`AppImage应用`

>> 直接执行`AppImage应用`

直接执行`AppImage`文件的原理很简单,`AppImage`文件里面有个[SquashFS](https://docs.kernel.org/filesystems/squashfs.html)的文件系统,运行的时候会借助[libfuse](https://github.com/libfuse/libfuse)挂在到一个临时的endpoint,里面的内容类似`独立归档应用`,有所有只读的文件,包括可执行文件,库文件,静态数据等.除此之外都和普通的软件运行环境是一样.也就是说直接运行`AppImage`文件相当于挂载了一块只读盘执行其中的可执行文件.

我们要直接执行`AppImage`文件需要确保系统中有`libfuse`.很遗憾ubuntu中是没有的,我们可以使用apt进行安装.以我所用的`Ubuntu 24.04`为例,使用如下语句

```bash
sudo add-apt-repository universe
sudo apt install libfuse2t64
```

在满足这个条件后,安装`AppImage应用`就很简单了

1. 下载`.AppImage`为后缀的文件
2. `chmod +x 文件名.AppImage`给`AppImage应用`服务可执行权限
3. 双击或者进terminal输入文件名就可以使用了

需要注意现在很多应用是`electron`开发的,他们底层是`Chromium`,对于这类应用我们需要执行是增加参数`--no-sandbox`,否则会报错,这是`Chromium`的限制

>> 提取`AppImage应用`后当做`独立归档应用`处理

`AppImage应用`自带的运行时支持使用flage`--appimage-extract`将其中的内容解包成一个文件夹([AppDir](https://docs.appimage.org/reference/appdir.html))放到路径`squashfs-root`下.
提取出来以后我们完全可以将它当做`独立归档应用`处理

> `AppImage应用`的兼容性问题

没错,`AppImage`实际上并没有完全解决兼容性问题,因为不同的linux发行版甚至同意linux发行版的不同版本,他们的`glibc`都不一定版本一样,而很多应用即便封装成了`AppImage`也并没有将`glibc`封装进去.这就造成了一些`AppImage应用`有兼容性问题.毕竟`AppImage`这套方案并没有任何审核,也没有任何强制手段,怎么封装全靠开发者自觉,这也是开源去中心化的代价.

如果碰到用不了的`AppImage应用`,我们基本只能放弃,看看有没有其他方案的可以用了.

> 系统集成

光是能执行并不够,我们希望在操作系统中可以方便的找到这些`PORTABLE LINUX APPS`.这就需要手动将他们集成到桌面系统已方便打开,怎么注册呢?通常我们用如下方式:

1. 下载软件压缩包,并解药到`~/Applications`目录
2. 进入`~/.local/share/applications`目录,创建`xxx.desktop`并编辑

    ```toml
    [Desktop Entry]
    Version=<应用版本>
    Type=Application
    Encoding=UTF-8
    Name=<应用名>
    Comment=<应用说明>
    Icon=<应用图标路径>
    Exec=<应用可执行文件路径> %U
    Terminal=false
    Categories=Application;
    ```

    其中
    + `Version`: 应用版本
    + `Type`: 指定快捷方式类型
    + `Encoding`: 指定编码格式
    + `Name`: 快捷方式的名称
    + `Comment`: 快捷方式的简短描述
    + `Icon`: 快捷方式的图标路径,一般解压后都有
    + `Exec`: 要执行的命令. 是一个占位符如果快捷方式是通过文件管理器启动的,则使用`%U`这个占位符指定文件将被作为参数传递给命令.
    + `Terminal`: 如果程序需要在终端中运行,则设置为true;否则设置为false
    + `Categories`: 用于帮助应用程序菜单对快捷方式进行分类

    其他字段可以参考[desktop-entry-spec相关规定](https://specifications.freedesktop.org/desktop-entry-spec/latest/recognized-keys.html)

3. 为这个`.desktop`文件赋予可执行权限

    ```bash
    chmod +x ~/Desktop/xxx.desktop
    ```

> `PORTABLE LINUX APPS`管理工具

`PORTABLE LINUX APPS`可以说是真正的去中心化,社区驱动,以至于时至今日已经有点驱动不下去了的感觉.

但不管怎样`PORTABLE LINUX APPS`也是一类重要的解决方案.这类应用的维护相对麻烦些,工具简陋,维护相对不上心,我们能做的就是多一些手动介入让他们在系统中相对更简洁好维护而已.

为了做到这一点,我个人会给`PORTABLE LINUX APPS`的部署做如下限制

+ 限制`PORTABLE LINUX APPS`仅用户级安装
+ 限制`PORTABLE LINUX APPS`的应用全部安装在`~/Applications`目录下,且一个应用一个文件夹或一个`.appimage`文件不做嵌套

说回到管理工具,`PORTABLE LINUX APPS`的管理工具其实也是有的,但一样也是相当的"去中心化"而且完成度相对低.我体验的相对靠谱的是[ivan-hc/AM](https://github.com/ivan-hc/AM)这套工具.
这套的缺点是只有命令行工具没有带gui的管理应用.但相对的index整理的算比较全的,而且提供了一个轻量级沙盒功能可以进一步隔离应用和操作系统.

我更加推荐安装它的非root版本[ivan-hc/AppMan](https://github.com/ivan-hc/AppMan),他们功能一样,只是appman是用户级而非系统级.

`am`和`appman`共用一份[应用索引](https://portable-linux-apps.github.io/apps.html)作为应用仓库,

>> 安装`appman`

安装`appman`可以通过如下命令实现

```bash
wget -q https://raw.githubusercontent.com/ivan-hc/AM/main/AM-INSTALLER && chmod a+x ./AM-INSTALLER && ./AM-INSTALLER
```

注意由于数据在github上,最好先挂上代理再执行上面的命令.

执行好后会让你选安装`am`还是`appman`,选`2`就是`appman`.

安装完成后会有如下几个位置有相关文件

+ 执行命令的目录下会有安装脚本`AM-INSTALLER`,可以删掉
+ `~/.local/bin`下会有可执行文件`appman`
+ `~/Applications`下会有文件夹`appman`

要卸载`appman`时先将`appman`安装的应用都卸载了,然后直接将这些文件,文件夹删除即可

>> 使用`appman`

+ 搜索am数据库中索引的应用

    ```bash
    appman -q <关键字>
    ```

+ 安装am数据库中索引的应用

    ```bash
    appman -i <应用名>
    ```

    安装即可.如果希望安装`AppImage`而不是`独立归档应用`,可以添加flag`-a`

    ```bash
    appman -ia <应用名>
    ```

    如果希望安装的同时进行沙盒化,可以加个flag`-s`

    ```bash
    appman -is <应用名>
    ```

    `am`工具的沙盒化是使用的[Aisap](https://github.com/mgord9518/aisap).

+ 安装来自github而am数据库中索引未收录的应用

    如果我们希望让`apman`安装并管理未被收录的github上的`appimage`,则可以使用`-e`命令替代`-i`和`-ia`

    ```bash
    appman -e <github地址> <应用名> [关键字]
    ```

    如果这个地址中有不止一个`appimage`文件,可选的`关键字`可以用来指定实际需要的是哪个.而`<应用名>`则是这个安装好的`appimage`文件在系统内的命名.

+ 将本地已经下载好的`appimage`文件注册到`appman`并集成进系统

    ```bash
    appman --launcher /path/to/File.AppImage
    ```

+ 查看本地由`appman`维护的应用

    ```bash
    appman -f
    ```

+ 更新应用

    ```bash
    appman -u <应用名>
    ```

+ 卸载应用

    ```bash
    appman -R <应用名>
    ```

### 选择应用安装方式的总结

在ubuntu中可以顺利使用的应用分发方式我做了如下总结

| 对比项目       | `独立归档应用`                                                                      | `deb`                                       | `AppImage`                                                                                  | `Flatpak`                       | `snap`         |
| -------------- | ----------------------------------------------------------------------------------- | ------------------------------------------- | ------------------------------------------------------------------------------------------- | ------------------------------- | -------------- |
| 发行版支持     | 跨发行版                                                                            | 仅debian系                                  | 跨发行版                                                                                    | 跨发行版                        | 跨发行版       |
| 是否独立执行   | 是                                                                                  | 否                                          | 是                                                                                          | 否                              | 否             |
| 是否沙盒运行   | 否但可以                                                                            | 否                                          | 否但可以                                                                                    | 是                              | 是             |
| 性能损失       | 无损失                                                                              | 无损失                                      | 启动有略微损失,运行无损失                                                                   | 损失较小                        | 损失小         |
| 空间占用       | 大                                                                                  | 最小                                        | 大                                                                                          | 本体较小,算上依赖大             | 最大           |
| 资源平台       | [portable-linux-apps](https://portable-linux-apps.github.io/apps.html)</br>应用官网 | `apt`</br>`github release`</br>各种官网下载 | [portable-linux-apps](https://portable-linux-apps.github.io/apps.html)</br>`github release` | [flathub](https://flathub.org/) | ubuntu应用商店 |
| 平台资源数量   | 较少                                                                                | 多                                          | 较少                                                                                        | 较多                            | 较少           |
| 平台版本活跃度 | 高                                                                                  | 取决于渠道                                  | 取决于渠道                                                                                  | 高                              | 高             |
| 维护工具       | `am`/`appman`                                                                       | `apt`/`synaptic`                            | `am`/`appman`                                                                               | `Warehouse`                     | 系统应用商店   |
| 安装难度       | 大                                                                                  | 小                                          | 极小                                                                                        | 小                              | 小             |
| 维护难度       | 大                                                                                  | 小                                          | 大                                                                                          | 小                              | 小             |

可以看到如果这些平台都各有利弊.我们应该针对不同需求使用不同的平台.

+ 对于比较大的公司或知名大项目提供的产品,比如谷歌微软blender什么的,尽量用官方提供的安装方案
+ 对性能要求比较高的尤其是吃计算资源的应用
    + 优先尝试使用`独立归档应用`
    + 如果有`Appimage`也可以尝试用`Appimage`
    + 如果没有就找官网有没有提供`deb`包
+ 对于扩展性要求较高的比如要自己装插件改文件的应用
    + 优先尝试使用`独立归档应用`
    + 如果有`Appimage`,可以解压出来当做独立归档应用来尝试使用
+ 其他应用无脑`Flatpak`即可.

### 常用软件

linux下我会尽量推荐开源工具.下面是一些常用软件的安装信息

#### 系统工具

| 软件                                                             | 渠道                                                               | 说明                                                             |
| ---------------------------------------------------------------- | ------------------------------------------------------------------ | ---------------------------------------------------------------- |
| Disk Usage Analyzer                                              | [flathub](https://flathub.org/apps/org.gnome.baobab)               | 查看硬盘使用情况的工具                                           |
| Remmina                                                          | [flathub](https://flathub.org/apps/org.remmina.Remmina)            | 远程桌面的连接客户端                                             |
| timeshift                                                        | `sudo apt install timeshift`                                       | 系统快照,我们需要自备一块</br>空的U盘专门做快照,每半年做一份快照 |
| [missioncenter](https://missioncenter.io/)                       | [flathub](https://flathub.org/apps/io.missioncenter.MissionCenter) | 有着window上资源管理器风格的综合性资源监控软件,可以监控GP        |
| [CPU-X](https://thetumultuousunicornofdarkness.github.io/CPU-X/) | `sudo apt install cpu-x`                                           | windows上cpu-z在linux上的平替                                    |
| AMD GPU TOP                                                      | [官网下载deb](https://github.com/Umio-Yasuno/amdgpu_top)           | amdgpu的运行详细信息监控                                         |
| bleachbit                                                        | `sudo apt install bleachbit`                                       | 系统清理工具                                                     |
| gufw                                                             | `sudo apt install gufw`                                            | 防火墙工具,使用`sudo ufw enable`启动                             |
| chrome                                                           | [官网下载deb](https://www.google.cn/intl/zh-CN/chrome/)            | google的浏览器                                                   |
| Warehouse                                                        | [flathub](https://flathub.org/apps/io.github.flattool.Warehouse)   | flatpak应用管理                                                  |
| synaptic                                                         | `sudo apt install synaptic`                                        | 新立得软件包管理器,管理deb应用和包                               |
| appman                                                           | [官网脚本下载](https://github.com/ivan-hc/AppMan)                  | 管理`PORTABLE LINUX APPS`                                        |
| whaler                                                           | [flathub](https://flathub.org/apps/com.github.sdv43.whaler)        | 轻量级docker容器监控工具                                         |

##### 补充设置

> chrome

硬件加速:

1. 地址栏输入`chrome://settings/`
2. 左侧选中系统,右边激活`使用图形加速功能(如果可用)`

gpu加速:

在地址栏输入`chrome://flags/#disable-accelerated-video-decode`,找到其中的

+ `Hardware-accelerated video decode`硬件解码设置,确保置为`已启用`

#### 生产力工具

 | 软件                           | 渠道                                                                   | 说明                                                               |
 | ------------------------------ | ---------------------------------------------------------------------- | ------------------------------------------------------------------ |
 | Calculator                     | [flathub](https://flathub.org/apps/org.gnome.Calculator)               | 基本的计算器工具                                                   |
 | gimp                           | [flathub](https://flathub.org/apps/org.gimp.GIMP)                      | 开源的图像编辑软件,ps平替                                          |
 | Inkscape                       | [flathub](https://flathub.org/apps/org.inkscape.Inkscape)              | 矢量图编辑工具                                                     |
 | freecad                        | `appman -i freecad`                                                    | 开源的工程制图,autocad平替                                         |
 | blender                        | `appman -i blender`                                                    | 开源的3d建模渲染工具,maya平替                                      |
 | godot                          | `appman -i godot`                                                      | 开源的轻量级游戏引擎                                               |
 | unrealengine5                  | [官网下载](https://www.unrealengine.com/zh-CN/download)                | 大名鼎鼎的虚幻引擎,                                                |
 | shotcut                        | `appman -i shotcut`                                                    | 轻量级的开源视频剪辑工具                                           |
 | DaVinci Resolve                | [官网下载](http://www.blackmagicdesign.com/cn/products/davinciresolve) | 大名鼎鼎的生产级视频剪辑工具达芬奇,有免费的社区版                  |
 | KiCad                          | [flathub](https://flathub.org/apps/org.kicad.KiCad)                    | 知名的开源eda工具                                                  |
 | balenaEtcher                   | [官网下载](https://etcher.balena.io/)                                  | 镜像写入工具                                                       |
 | 飞书                           | [官网下载deb](https://www.feishu.cn/download)                          | 知名的办公协作工具,flathub版本无法后台挂载因此用官网版本           |
 | 微信                           | [官网下载deb](https://linux.weixin.qq.com/)                            | 知名的聊天工具,由于flathub版本的后台功能有缺陷因此使用AppImage版本 |
 | wps                            | [官网下载deb](https://www.wps.cn/product/wpslinux)                     | 知名的office套件,flatpak版本过低                                   |
 | [obs](https://obsproject.com/) | [flathub](https://flathub.org/apps/com.obsproject.Studio)              | 知名的开源直播录屏工具                                             |
 | vscode                         | [官网下载deb](https://code.visualstudio.com/Download)                  | 文本编辑器                                                         |
 | github desktop                 | [github下载deb](https://github.com/shiftkey/desktop)                   | github desktop的第三方linux fork                                   |
 | Xmind                          | [flathub](https://flathub.org/apps/net.xmind.XMind)                    | 知名的脑图工具                                                     |
 | minder                         | [flathub](https://flathub.org/apps/com.github.phase1geo.minder)        | xmind的开源替代                                                    |

##### 补充设置

> vscode额外设置: vscode默认会将标题栏和工具栏分开,非常的丑也非常的不紧凑.我们可以进入`文件->首选项->设置`,在其中搜索`window.titleBarStyle`,将其设置为`custom`.这样标题栏就会和工具栏合并,好看很多.

#### 娱乐工具

 | 软件                     | 渠道                                                                          | 说明                                                     |
 | ------------------------ | ----------------------------------------------------------------------------- | -------------------------------------------------------- |
 | vlc                      | [flathub](https://flathub.org/apps/org.videolan.VLC)                          | 知名的开源视频播放器                                     |
 | NetEase Cloud Music Gtk4 | [flathub](https://flathub.org/apps/com.github.gmg137.netease-cloud-music-gtk) | 网易云音乐的开源第三方客户端                             |
 | steam                    | [官网下载](https://store.steampowered.com/about/)                             | 知名的pc游戏平台                                         |
 | ProtonUp-Qt              | [flathub](https://flathub.org/apps/net.davidotek.pupgui2)                     | 为steam管理GE-Proton                                     |
 | protontricks             | [flathub](https://flathub.org/apps/com.github.Matoking.protontricks)          | 为Steam/Proton游戏以及其他常见Wine功能运行Winetricks命令 |

## Linux下神奇的Steam

steam大家都知道,知名的游戏平台嘛.但在linux下steam是一个神奇的存在,某种程度上可以说是steam盘活了桌面linux生态都不过分.

steam最基本的能力当然是让玩家可以方便的打游戏,阀门社为了这个目标煞费苦心,他们的努力直接给linux平台大范围扩展了生态,以至于steam即便不打游戏都是linux上的必装软件.

### N卡安装steam的额外依赖

steam为了可以兼容老游戏，会要安装很多32位的依赖，如果你是n卡，直接安装steam它会给你报错

```txt
You are missing the following 32-bit libraries, and Steam may not run:

libc.so.6
```

要解决这个问题，我们需要先额外安装两个依赖

1. 安装32位的`libc6`

    ```bash
    sudo apt install libc6:i386 
    ```

2. 安装`nvidia-driver`对应版本的32位驱动

    ```bash
    dpkg -l | grep nvidia-driver # 获取到驱动版本 比如575
    sudo apt install libnvidia-gl-575:i386 #安装对应32位依赖
    ```

### 添加非steam应用

steam除了可以管理它自己平台上的软件,也可以自己给它添加软件由他启动

![添加非steam应用-入口][2]

![添加非steam应用-选应用][3]

### steam runtime

如果游戏也算应用的话steam在linux上可能是最大的应用分发渠道了.steam通过[steam-linux-runtime](https://github.com/ValveSoftware/steam-runtime)支持原生linux应用.

而这个`steam-linux-runtime`实际上就是一个通用的轻量级沙盒,包含了一个统一的库环境并用namespace与操作系统隔离.一些游戏开发相关的开源应用比如`blender`,`godot`,原生支持linux的游戏比如`dota`也都可以使用这种方式运行.

由于这个环境也是用namespace做抽象的,所以性能损失也是有的,只是没那么多.

话虽如此,我个人还是不太推荐用`steam`管理除游戏以外的linux原生应用的,一方面这种应用太少,另一方面这类应用都往往是对计算性能要求比较高的类型,白白损失性能还是太可惜了.

### 转译应用

其实很久之前就有一个开源项目[wine](https://github.com/wine-mirror/wine)致力于让linux和macos可以执行windows程序.起工作原理可以理解为将windows程序调用的各种系统依赖在linux/macos上重新实现一遍,然后让windows程序调用这些被转译的依赖以执行.这一套方案原理上逻辑上没啥问题,但一直以来支持并不好,主要原因是图形接口的转译跟不上.

1. 一方面是厂商不开源,让转译工作难以进行,
2. 另一方面是开源的工具(`opengl`这类)性能拉胯.
3. 最后就是项目是社区驱动的,没钱没资源,维护自然也跟不上

而就在最近几年转译应用忽然就迎来了春天,主要取决于两点

1. 微软开始拥抱开源了,
2. 阀门社开始做steam deck了.

总而言之,阀门社搞了个开源的转译层项目[proton](https://github.com/ValveSoftware/Proton).这直接让linux下可以正常跑大部分windows平台的游戏,还顺便让其他windows软件也可以借助steam进行管理运行.

我们就以官网下载的原生版本的steam为基准,介绍转译应用.

#### 使用steam安装proton

首先我们要先安装steam,由于steam本身是一个启动器,后需可能会有很多需要手动操作的部分,因此最好是原生安装方便手动调试.然后进入steam,进入`库`,勾选`工具`,找到`proton`,一般装个最新版本就够用了.也可以多装一个`proton experimental`备用,然后视游戏安装相关的反作弊运行时即可.

#### 使用ProtonUp-Qt安装proton-ge

proton项目是开源的,除了steam官方的proton,自然也会有魔改proton.其中最知名的自然就是[proton-ge](https://github.com/GloriousEggroll/proton-ge-custom),这个第三方的魔改`proton`据说比官方版本效率更高些,而且兼容性会更好些,虽然也一样可能有兼容性问题.

`proton-ge`官方给的安装方式还是比较麻烦的,我们可以使用[ProtonUp-Qt](https://flathub.org/apps/net.davidotek.pupgui2)来进行自动安装,只要打开后选择最新版本安装即可.需要注意steam必须是原生安装的,否则安装位置会不对了.

![安装proton-ge][4]

#### 使用proton运行windows软件

即便是安装好`proton`后,windows平台的软件也是不能直接下载了用的,我们需要选中已经购买好的游戏,点击右侧`设置`(齿轮按钮)->`属性`->`兼容性`,选一个`proton`版本强制指定兼容层.这样游戏就可以下载运行了.

![设置proton-入口][5]

![设置proton-选择版本][6]

之于能不能玩,这就得群策群力了,[protondb](https://www.protondb.com/)是一个steam上游戏兼容性的数据网站,主要靠社区玩家提交报告来确定特定游戏的兼容性.我们可以先在其中查看

### 游戏帧数检测

我们看到很多评测视频里会有帧数(FPS)等指示用来观察游戏的流畅程度.在linux下也有,需要使用[flightlessmango/MangoHud](https://github.com/flightlessmango/MangoHud)这个软件.

```bash
sudo apt install mangohud
```

安装完成后我们可以在steam中设置`启动选项修改`为`mangohud %command%`

![启动选项修改][7]

默认会在游戏界面的左上角给出cpu,gpu占用和帧数数据.我们可以在`~/.config/MangoHud/MangoHud.conf`文件中对展示的内容进行自定义,具体有哪些参数可以看[这张表](https://github.com/flightlessmango/MangoHud?tab=readme-ov-file#environment-variables).建议像我下面这样设置

```txt
cpu_temp
gpu_temp
ram
vram
position=top-left
```

<!-- https://www.bilibili.com/video/BV1zD4y1b7Jj?vd_source=08b668b29d50d7b81093d4adee9dfde0&spm_id_from=333.788.videopod.sections -->

<!-- 
https://www.mapeditor.org/
 https://itch.io/game-assets/free/tag-tilemap -->

### 小黄鸭(Lossless Scaling)伪双卡交火[2025-09-09更新]

`Lossless Scaling`是steam上的一个小软件,它可以让我们在低分辨率下运行游戏,然后将画面放大到全屏显示,还可以进行插帧.如果我们刚好有两块显卡,那么就可以用这个软件让一块显卡专门负责游戏的渲染,另一块显卡专门负责桌面的插帧显示,这样就相当于双卡交火了.遗憾的是官方版本是windows限定,不过有个开源插件[lsfg-vk](https://github.com/PancakeTAS/lsfg-vk),用它就可以在linux下使用这项伪双卡交火技术了.注意这个软件用到了`Vulkan`,你得确保你的显卡驱动支持`Vulkan`.



### 用steam串流

说到串流可能大家想到的更多的是[sunshine](https://github.com/LizardByte/Sunshine)-[moonlight](https://github.com/moonlight-stream/moonlight-android)的串流组合.但实际情况是在用N卡的情况下sunshine强无敌,而在用a卡尤其是apu的核显的情况下sunshine并不好用,尤其在linux下根本用不了.反而是steam link效果还不错.

用法也不复杂

1. 接收端的机器安装`steam link`,苹果直接在app store里下,安卓则去[官网下](https://help.steampowered.com/zh-cn/faqs/view/7112-CD02-7B57-59F8)即可,之后打开它
2. 发送端机器打开steam,在`steam->设置`中找到`远程畅玩`,激活后`配对流式应用`即可.如果你已经有配对好了的就需要先`取消设备配对`然后再做配对
3. 这时你接收端的机器会有个pin码,在发送端机器中将这个码填上就配对完成了.

![steamlink设置入口-1][8]
![steamlink设置入口-2][9]

steam串流非常吃网络,发送端一定要用有线的方式接入网络,接收端也尽量在网络覆盖范围内.网络稳定串流才能稳定

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


android: microsoft Remote Desktop
ios: microsoft windows app mobile
linux: Remmina


之后要远程使用的时候就点击这个远程应用即可 -->

## 其他补充

### aria2下载器

当然了我们也可以用docker安装aria2下载器,但那种方式比较适合家用服务器,因为它会固定下载到一个指定的文件夹下,而作为终端pc,更合适的用法是直接安装配置,在需要的时候命令行执行下载,这样灵活的多.

1. 安装

    ```bash
    sudo apt install aria2
    ```

2. 配置

    我们需要先创建配置相关文件

    ```bash
    sudo mkdir /etc/aria2 # aria2的配置文件夹
    sudo touch /etc/aria2/aria2.session    #新建session文件
    sudo chmod 777 /etc/aria2/aria2.session    #设置aria2.session可写
    sudo nano /etc/aria2/aria2.conf
    ```

    之后编辑`aria2.conf`文件

    ```txt
    ## 文件保存相关 ##

    # 文件的保存路径(可使用绝对路径或相对路径), 默认: 当前启动位置
    #dir=~/downloads
    # 启用磁盘缓存, 0为禁用缓存, 需1.16以上版本, 默认:16M
    disk-cache=128M
    # 文件预分配方式, 能有效降低磁盘碎片, 默认:prealloc
    # 预分配所需时间: none < falloc ? trunc < prealloc
    # falloc和trunc则需要文件系统和内核支持
    # NTFS建议使用falloc, EXT3/4建议trunc, MAC 下需要注释此项
    file-allocation=trunc 
    # 断点续传
    continue=true

    ## 下载连接相关 ##

    # 最大同时下载任务数, 运行时可修改, 默认:5
    #max-concurrent-downloads=5
    # 同一服务器连接数, 添加时可指定, 默认:1
    #max-connection-per-server=5
    # 最小文件分片大小, 添加时可指定, 取值范围1M -1024M, 默认:20M
    # 假定size=10M, 文件为20MiB 则使用两个来源下载; 文件为15MiB 则使用一个来源下载
    #min-split-size=10M
    # 单个任务最大线程数, 添加时可指定, 默认:5
    #split=5
    # 整体下载速度限制, 运行时可修改, 默认:0
    max-overall-download-limit=0
    # 单个任务下载速度限制, 默认:0
    max-download-limit=0
    # 整体上传速度限制, 运行时可修改, 默认:0
    max-overall-upload-limit=0
    # 单个任务上传速度限制, 默认:0
    max-upload-limit=0
    # 禁用IPv6, 默认:false
    #disable-ipv6=true
    # 连接超时时间, 默认:60
    #timeout=60
    # 最大重试次数, 设置为0表示不限制重试次数, 默认:5
    #max-tries=5
    # 设置重试等待的秒数, 默认:0
    #retry-wait=0

    ## 进度保存相关 ##

    # 从会话文件中读取下载任务
    #input-file=/etc/aria2/aria2.session
    # 在Aria2退出时保存`错误/未完成`的下载任务到会话文件
    #save-session=/etc/aria2/aria2.session
    # 定时保存会话, 0为退出时才保存, 需1.16.1以上版本, 默认:0
    #save-session-interval=60

    ## RPC相关设置 ##

    # 启用RPC, 默认:false
    enable-rpc=true
    # 允许所有来源, 默认:false
    rpc-allow-origin-all=true
    # 允许非外部访问, 默认:false
    rpc-listen-all=true
    # 事件轮询方式, 取值:[epoll, kqueue, port, poll, select], 不同系统默认值不同
    #event-poll=select
    # RPC监听端口, 端口被占用时可以修改, 默认:6800
    rpc-listen-port=6860
    # 设置的RPC授权令牌, v1.18.4新增功能, 取代 --rpc-user 和 --rpc-passwd 选项
    #rpc-secret=<TOKEN>
    # 设置的RPC访问用户名, 此选项新版已废弃, 建议改用 --rpc-secret 选项
    #rpc-user=<USER>
    # 设置的RPC访问密码, 此选项新版已废弃, 建议改用 --rpc-secret 选项
    #rpc-passwd=<PASSWD>
    # 是否启用 RPC 服务的 SSL/TLS 加密,
    # 启用加密后 RPC 服务需要使用 https 或者 wss 协议连接
    #rpc-secure=true
    # 在 RPC 服务中启用 SSL/TLS 加密时的证书文件,
    # 使用 PEM 格式时，您必须通过 --rpc-private-key 指定私钥，这里可以是CRT或者其他格式，标准内都支持。
    #rpc-certificate=/path/to/certificate.pem
    # 在 RPC 服务中启用 SSL/TLS 加密时的私钥文件
    #rpc-private-key=/path/to/certificate.key

    ## BT/PT下载相关 ##

    # 当下载的是一个种子(以.torrent结尾)时, 自动开始BT任务, 默认:true
    #follow-torrent=true
    # BT监听端口, 当端口被屏蔽时使用, 默认:6881-6999
    #listen-port=51413
    # 单个种子最大连接数, 默认:55
    #bt-max-peers=55
    # 打开DHT功能, PT需要禁用, 默认:true
    enable-dht=false
    # 打开IPv6 DHT功能, PT需要禁用
    enable-dht6=false
    # DHT网络监听端口, 默认:6881-6999
    #dht-listen-port=6881-6999
    # 本地节点查找, PT需要禁用, 默认:false
    bt-enable-lpd=false
    # 种子交换, PT需要禁用, 默认:true
    enable-peer-exchange=false
    # 每个种子限速, 对少种的PT很有用, 默认:50K
    #bt-request-peer-speed-limit=50K
    # 客户端伪装, PT需要
    peer-id-prefix=-TR2770-
    user-agent=Transmission/2.77
    peer-agent=Transmission/2.77
    # 当种子的分享率达到这个数时, 自动停止做种, 0为一直做种, 默认:1.0
    seed-ratio=2.0
    # 强制保存会话, 即使任务已经完成, 默认:false
    # 较新的版本开启后会在任务完成后依然保留.aria2文件
    #force-save=false
    # BT校验相关, 默认:true
    #bt-hash-check-seed=true
    # 继续之前的BT任务时, 无需再次校验, 默认:false
    bt-seed-unverified=true
    # 保存磁力链接元数据为种子文件(.torrent文件), 默认:false
    bt-save-metadata=true
    ```

3. 检查

    执行`sudo aria2c --conf-path=/etc/aria2/aria2.conf`,能正常跑不报错就说明没问题了

安装完成后重启,`aria2c`的服务就会被`systemd`接管了.由于上面的配置中我们启动了rpc,其他支持的软件也就可以借助这个rpc来通过`aria2`下载东西了.

之后每次要下载只要执行命令`aria2c [-d 保存目录 [-o 文件名] | -Z] <URI | 磁力链 | torrent文件 | METALINK文件>...`

`-d`指定下载文件保存的目录.当只有一个要下载的链接时可以用`-o`为下载文件的结果改名;如果有多个链接,可以使用`-Z`将并行下载各个链接改为串行下载.

个人认为作为个人终端的pc并不需要有界面管理`aria2c`的下载任务,就像`git`一样命令行足够了

### mycard

[mycard](https://mycard.world/)很神奇的自带linux支持,这个平台主要是用来玩ygopro和东方的,

不想花钱买卡又想打游戏王的牌佬就可以在linux下很轻松的玩到ygopro2了.在linux下首页下载提供的是一个`appimage`,下载下来后给个`+x`权限人工集成下就能用了.

有如下注意事项:

+ mycard下载游戏依赖`aria2c`,需要先安装`aria2c`再安装`mycard`
+ mycard解压游戏依赖`zstd`,需要先安装`zstd`(`sudo apt install zstd`)再安装`mycard`
+ 游戏王只能使用`ygopro2`,因为`ygopro`的依赖编译有问题,会缺少`irrKlang`,而`ygopro2`是unity开发,没有依赖问题
+ 东方游戏都是window游戏,mycard仅能用作下载,运行需要借助steam的proton




<!-- ## 串流 -->

[1]: {{site.url}}/img/in-post/ubuntu/gnome桌面结构.jpg
[2]: {{site.url}}/img/in-post/ubuntu/add_app.png
[3]: {{site.url}}/img/in-post/ubuntu/add_app_2.png
[4]: {{site.url}}/img/in-post/ubuntu/install_proton_ge.png
[5]: {{site.url}}/img/in-post/ubuntu/proton_set_1.png
[6]: {{site.url}}/img/in-post/ubuntu/proton_set_2.png
[7]: {{site.url}}/img/in-post/ubuntu/framepannel.png
[8]: {{site.url}}/img/in-post/ubuntu/steamlink-1.png
[9]: {{site.url}}/img/in-post/ubuntu/steamlink-2.png