---
layout: post
title: "用deskmini和8700g组一台aipc"
series:
    aipc_experiment:
        index: 2
date: 2024-11-27
author: "Hsz"
category: experiment
tags:
    - linux
    - aipc
header-img: "img/home-bg-o.jpg"
update: 2024-11-27
---
# amd8700g折腾记录



我在2024年8月就入手了[华擎的deskmini x600准系统](https://www.asrock.com/nettop/AMD/DeskMini%20X600%20Series/index.cn.asp).它包含了

+ 一个stx规格的无芯片组主板,提供
    + 6+2相供电,一般认为最大可以带220W板载功率
    + 一个AMD的AM5接口,支持最大tcp65W的zen4,zen5带核显的cpu
    + 两个SO-DIMM DDR5插槽,支持最高6400+MHz,最大96GB容量的笔记本DDR5内存
    + 两个sata2.5寸硬盘接口
    + 一个pcie 5.0x4的m.2固态硬盘(正面)
    + 一个pcie 4.0x4的m.2固态硬盘(背面)
    + 一个 M.2 Key-E 无线网卡接口
    + 前置一个USB 3.2 Gen1的Type-C接口;一个USB 3.2 Gen1的Type-A接口;一个麦克风接口和一个带麦克风的耳机接口
    + 后置一个dp1.4接口(4K@120Hz);一个hdmi接口(4K@120Hz);一个D-Sub接口(1080p@60Hz);一个RJ-45 LAN网口(2.5G);两个USB 3.2 Gen1 Type-A接口
    + 板载1个USB 2.0 接针;两个4pinpwm风扇接口,1个2pin的清cmos接针
    + BIOS目前(截止到2024年11月11日)更新到4.03版本,支持了AMD AGESA 1.2.0.0
+ 一个与之对应的1.9L机箱,支援高度小于等于47mm的的下压式散热器,右侧板提供了安装7.5x7.5标准VESA支架的接口,背面io挡板面提供了3个8mm天线孔和一个db9孔位可以掰掉挡板使用,左侧板提供了两个usb type-A孔位可以掰掉挡板使用
+ 一个标称可以解热65w的卡扣式下压散热器
+ 一个标称120w的19v,2.5/5.0规格的电源适配器

我为这套除了配置了8700g外,还配置了

+ 英睿达48gx2 标称xmp5600MHZ的的笔记本内存
+ 闪迪至尊高速PCIE 3.0x4的2t固态硬盘,放在背面
+ 官方的Intel® Wi-Fi Kit,提供蓝牙和wifi6的支持

我为这套系统额外配置了如下设备

+ 利民apx90 x45纯铜版散热器用于替代原装散热器
+ 一对铝制内存散热装甲
+ 一条利民的2.5mm厚的散热硅胶垫用于背面ssd散热
+ 猫头鹰9214黑款风扇用于替换上面散热器的风扇
+ 第三方的DeskMini USB 2.0 排线,用于充分利用板载USB 2.0 接针以额外提供2个usb2.0接口
+ 自制clear cmos按钮(ds-402无锁自复位小微型按钮+母对母杜邦线,拆开一段后套入按钮针脚,按实并焊接加固,最后用热缩管保护加固)
+ 一个12cm的防尘网放在进风口
+ 官方提供的7.5x7.5标准VESA支架
+ 一块带7.5VESA接口的1080p@120hz的便携屏幕随时备用
+ 氮化镓200w 19v 2.5/5.0接口适配器替代原装电源适配器
