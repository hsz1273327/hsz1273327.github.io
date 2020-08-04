---
layout: post
title: "树莓派linux系统安装和配置"
date: 2018-02-19
author: "Hsz"
category: blog
tags:
    - linux
    - opencl
    - 树莓派
header-img: "img/post-bg-js-module.jpg"
---

### opencl支持

树莓派也可以使用opencl支持了,这个项目叫[VC4CL](https://github.com/doe300/VC4CL),它是一个OpenCL 1.2的实例,虽然还没有完全支持,但更新很积极,本文写到这里的时候上一次更新是在1月21日.

安装它需要安装几个依赖:

+ cmake

    使用`sudo apt-get install cmake`安装

+ llvm
+ clang

+ Khronos ICD加载器
    
    使用`sudo apt-get install ocl-icd-opencl-dev ocl-icd-dev`安装.安装好后需要创建一个只有一行内容的文件`/etc/OpenCL/vendors/VC4CL.icd`并在其中第一行写上`VC4CL`库的绝对地址,当然这个可以安装好后再写上
+ OpenCL 头文件,使用`sudo apt-get install opencl-headers`

