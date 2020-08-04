---
layout: post
title: "树莓派上的集群实验_hdfs和yarn"
date: 2018-02-22
author: "Hsz"
category: blog
tags:
    - Hadoop
    - Linux
    - 集群
    - 树莓派
header-img: "img/home-bg-o.jpg"
---
# HDFS和YARN

hadoop是目前最常见的分布式基础框架之一,其核心部件有3块:

1. HDFS集群:负责海量数据的存储,集群中的角色主要有`NameNode` / `DataNode`/`SecondaryNameNode`
2. YARN集群:负责海量数据运算时的资源调度,集群中的角色主要有`ResourceManager`/`NodeManager`
3. MapReduce:它其实是一个应用程序开发包.

Hadoop集群具体来说包含两个集群:HDFS集群和YARN集群,两者逻辑上分离,但物理上常在一起.

HDFS即Hadoop分布式文件系统(Hadoop Distributed Filesystem),这是一种增量式的分布式文件系统,

其优点有:

1. 高容错性
    + 数据自动保存多个副本
    + 副本丢失后，自动恢复

2. 适合批处理
    + 移动计算而非移动数据
    + 数据位置暴露给计算框架

3. 适合大数据处理
    + GB、TB、甚至PB级数据
    + 百万规模以上的文件数量
    + 10K+ 节点数量

4. 可构件在廉价机器上
    + 通过多副本提高可靠性
    + 提供了容错和恢复机制

HDFS缺点有:

1. 无法应付毫秒级低延迟需求和高吞吐量需求

2. 不适合小文件存取,会占用节点大量内存,寻道时间会超过读取的时间

3. 无写入锁,并发写入会造成文件随机修改,
4. 只支持增量修改.

HDFS设计时就没考虑文件内容的改变,这种分布式大数据存储一般只适用于增量存储.对于那些按条存储的数据,由于有时间标签,如果需要修改内容,则直接放入新的数据.如果需要删除,则采用标记的方式,将要删除的记录标记为删除.因此数据会越来越多,使用时则以最后的时间标签为准.一般情况下,我们应定期对数据进行整理,将文件读出,将历史标签的记录删除,然后再写入HDFS.

YARN 即为hadoop的资源调度系统,它的职能就是将资源调度和任务调度分开,hadoop本身已经没什么人使用它的mapreduce框架编写应用了,但YARN的`ResourceManager`具备通用性,很多其他工具都可以调用它来管理资源.


## 单机安装

### 依赖

hadoop依赖java8环境,在[树莓派linux系统安装和配置](http://blog.hszofficial.site/blog/2018/02/19/%E6%A0%91%E8%8E%93%E6%B4%BElinux%E7%B3%BB%E7%BB%9F%E5%AE%89%E8%A3%85%E5%92%8C%E9%85%8D%E7%BD%AE/)一文中已经有描述了如何安装java环境,照方抓药即可.

<!-- 另一个一来是google的protocol buffers,我们可以安装[protocol buffers v2.5](https://protobuf.googlecode.com/files/protobuf-2.5.0.tar.gz).随便找个文件夹下载后执行下述操作,等待编译完成即可. -->

另一个依赖是`google`的`protocol buffers`,可以使用`sudo apt install libprotobuf-dev`安装,目前的版本是3.0

<!-- ```shell
wget https://protobuf.googlecode.com/files/protobuf-2.5.0.tar.gz
sudo tar xzvf protobuf-2.5.0.tar.gz
sudo chown -R pi protobuf-2.5.0/
cd protobuf-2.5.0
./configure --prefix=/usr
make
make check
sudo make install
``` -->


### 直接下载使用

hadoop是java程序,直接下载发行版软件就可以使用`wget http://mirrors.tuna.tsinghua.edu.cn/apache/hadoop/common/hadoop-3.0.0/hadoop-3.0.0.tar.gz`

然后解压`tar -xzvf hadoop-3.0.0.tar.gz`


### 编译安装

+ 编译依赖安装

    如果要编译安装,那么就还要再安装些其他依赖

    `sudo apt-get install -y maven build-essential autoconf automake libtool cmake zlib1g-dev pkg-config libssl-dev libfuse-dev libsnappy-dev libsnappy-java libbz2-dev subversion`

    *期间因为有些库依赖java1.6所以安装了1.6版本的jvm,需要用`sudo update-alternatives --config java`将其选回java1.8*

+ 下载hadoop源文件

    ```shell
    wget http://apache.mirrors.spacedump.net/hadoop/core/hadoop-2.6.0/hadoop-2.6.0-src.tar.gz
    tar xzvf hadoop-2.6.0-src.tar.gz
    sudo chown -R pi hadoop-2.6.0-src/
    cd hadoop-2.6.0-src/
    ```

+ 修改配置

    打开其中的`pom.xml `文件.将
    ```
    <additionalparam>-Xdoclint:none</additionalparam>
    ```
    添加进`<properties></properties>`标签中然后保存退出编辑

+ 打arm处理器补丁

    因为是arm处理器,所以要打个补丁
    ```shell
    cd hadoop-common-project/hadoop-common/src
    wget https://issues.apache.org/jira/secure/attachment/12570212/HADOOP-9320.patch
    patch < HADOOP-9320.patch
    ```

+ 使用`maven`编译安装

    执行如下指令进行编译
    ```shell
    sudo mvn package -Pdist,native -DskipTests -Dtar
    ```
    编译会花近两个小时才能完成.编译后的hadoop在`hadoop-2.6.0-src/hadoop-dist/target/hadoop2.6.0`这个位置,同时会有一个你的本地编译程序`.tar`包在`hadoop-2.6.0-src/hadoop-dist/target/hadoop2.6.0.tar.gz`

### 配置hadoop

+ 首先是环境变量配置

    在`/etc/profile`中添加如下内容(先修改为可读`sudo chmod ugo+w /etc/profile`)
    ```shell
    #-------------------------------------------------------------------java

    export JAVA_HOME=/usr/lib/jvm/jdk-8-oracle-arm-vfp-hflt
    export CLASS_PATH=$JAVA_HOME/lib:$JAVA_HOME/jre/lib  

    #-------------------------------------------------------------------hadoop
    export HADOOP_HOME=/home/pi/lib/hadoop/hadoop-2.6.0
    export HADOOP_PID_DIR=/data/hadoop/pids
    export HADOOP_COMMON_LIB_NATIVE_DIR=$HADOOP_HOME/lib/native
    export HADOOP_OPTS="$HADOOP_OPTS -Djava.library.path=$HADOOP_HOME/lib/native"

    export HADOOP_MAPRED_HOME=$HADOOP_HOME
    export HADOOP_COMMON_HOME=$HADOOP_HOME
    export HADOOP_HDFS_HOME=$HADOOP_HOME
    export YARN_HOME=$HADOOP_HOME

    export HADOOP_CONF_DIR=$HADOOP_HOME/etc/hadoop
    export HDFS_CONF_DIR=$HADOOP_HOME/etc/hadoop
    export YARN_CONF_DIR=$HADOOP_HOME/etc/hadoop

    export JAVA_LIBRARY_PATH=$HADOOP_HOME/lib/native
    export PATH=$PATH:$HADOOP_HOME/bin:$HADOOP_HOME/sbin
    ```

+ 创建用于存储的文件夹
    ```shell
    sudo mkdir -p /data/hadoop/{pids,storage}
    sudo mkdir -p /data/hadoop/storage/{hdfs,tmp}
    sudo mkdir -p /data/hadoop/storage/hdfs/{name,data}
    sudo chown -R pi /data
    sudo chmod -R ugo+w /data
    ```

+ 配置 `core-site.xml`
    ```xml
    <configuration>
        <property>
            <name>fs.defaultFS</name>
            <value>hdfs://coderMaster:9000</value>
        </property>
        <property>
            <name>io.file.buffer.size</name>
            <value>131072</value>
        </property>
        <property>
            <name>hadoop.tmp.dir</name>
            <value>file:/data/hadoop/storage/tmp</value>
        </property>
        <property>
            <name>hadoop.proxyuser.hadoop.hosts</name>
            <value>*</value>
        </property>
        <property>
            <name>hadoop.proxyuser.hadoop.groups</name>
            <value>*</value>
        </property>
        <property>
            <name>hadoop.native.lib</name>
            <value>true</value>
        </property>
    </configuration>
    ```

+ 配置`hdfs-site.xml`
    ```xml
    <configuration>
        <property>
            <name>dfs.namenode.name.dir</name>
            <value>file:/data/hadoop/storage/hdfs/name</value>
        </property>
        <property>
            <name>dfs.datanode.data.dir</name>
            <value>file:/data/hadoop/storage/hdfs/data</value>
        </property>
        <property>
            <name>dfs.replication</name>
            <value>3</value>
        </property>
        <property>
            <name>dfs.webhdfs.enabled</name>
            <value>true</value>
        </property>
    </configuration>
    ```

+ 配置`mapred-site.xml`

    ```xml
    <configuration>
        <property>
            <name>mapreduce.framework.name</name>
            <value>yarn</value>
        </property>
        <property>
            <name>mapreduce.jobhistory.address</name>
            <value>coderMaster:10020</value>
        </property>

        <property>
            <name>mapreduce.jobhistory.webapp.address</name>
            <value>coderMaster:19888</value>
        </property>
    </configuration>
    ```

+ 配置 yarn-site.xml

    ```xml
    <configuration>
        <property>
            <name>yarn.nodemanager.aux-services</name>
            <value>mapreduce_shuffle</value>
        </property>
        <property>
            <name>yarn.nodemanager.aux-services.mapreduce.shuffle.class</name>
            <value>org.apache.hadoop.mapred.ShuffleHandler</value>
        </property>
        <property>
            <name>yarn.resourcemanager.scheduler.address</name>
            <value>coderMaster:8030</value>
        </property>
        <property>
            <name>yarn.resourcemanager.resource-tracker.address</name>
            <value>coderMaster:8031</value>
        </property>
        <property>
            <name>yarn.resourcemanager.address</name>
            <value>coderMaster:8032</value>
        </property>
        <property>
            <name>yarn.resourcemanager.admin.address</name>
            <value>coderMaster:8033</value>
        </property>
        <property>
            <name>yarn.resourcemanager.webapp.address</name>
            <value>coderMaster:8088</value>
        </property>
    </configuration>
    ```

+ 配置`hadoop-env.sh`

    修改JAVA_HOME的值为`JAVA_HOME=/usr/lib/jvm/jdk-8-oracle-arm-vfp-hflt`

+ 配置`slaves`
    ```shell
    piNode01

    piNode02

    piNode03
    ```

## 集群安装和配置

将单机安装配置好的hadoop复制给子节点即可

```shell
sudo scp -r ~/lib/hadoop/hadoop-2.6.0/ pi@piNode01:~/lib/hadoop
sudo scp -r ~/lib/hadoop/hadoop-2.6.0/ pi@piNode02:~/lib/hadoop
sudo scp -r ~/lib/hadoop/hadoop-2.6.0/ pi@piNode03:~/lib/hadoop
```

## 尝试运行

```shell
hdfs namenode -format #初始化节点
start-dfs.sh          #开启文件系统
start-yarn.sh         #开启yarn
jps                   #查看开启的java进程
```

如果看到主节点有

```shell
ResourceManager
NameNode
SecondaryNameNode
```

从节点有

```shell
NodeManager
DataNode
```

就说明安装成功了.

## hdfs基本操作

在terminal中操作hdfs使用的命令是`hadoop fs xxx`,下表是支持的操作和说明

+ hdfs系统上操作

操作指令                                  | 使用格式                   | 含义
----------------------------------------- | -------------------------- | ----
`-ls`|`-ls<路径>`|查看指定目录的当前目录结构
`-lsr`|`-lsr<路径>`|递归查看指定路径的目录结构
`-du`|`-du<路径>`|统计目录下文件大小
`-dus`|`-dus<路径>`|汇总统计目录下文件(夹)大小
`-count`|`-count[-q]<路径>`|统计文件夹数量
`-mv`|`-mv<源路径><目的路径>`|移动
`-cp`|`-cp<源路径><目的路径>`|复制
`-rm`|`-rm[-skipTrash]<路径>`|删除文件或者空白文件夹
`-rmr`|`-rmr[-skipTrash]<路径>`|递归删除
`-cat`|`-cat<hdfs文件内容>`|查看文件内容
`-text`|`-text<hdfs文件内容>`|查看文件内容
`-mkdir`|`-mkdir<hdfs路径>`|创建空白文件夹
`-setrep`|`-setrep[-r][-w]<副本数><路径>`|修改副本文件
`-touchz`|`-touchz<文件路径>`|创建空白文件
`-stat`|`-stat[format]<路径>`|显示文件统计信息
`-tail`|`-tail[-f]<文件>`|查看文件尾部信息
`-chmod`|`-chmod[-R]<权限模式>[路径]`|修改权限
`-chown`|`-chown[-R][属主][:[属组]]路径`|修改属主
`-chgrp`|`-chgrp[-R] 属组名称 路径`|修改属组

+ hdfs与本地文件系统交互

操作指令|使用格式|含义
---|---|---
`-get`|`-get < hdfs file > < local file or dir>`|将hdfs上的文件下载到本地
`-copyToLocal`|`-copyToLocal[-ignoreCrc][-crc][hdfs源路径][linux目的路径]`|从本地复制
`-moveToLocal`|`-moveToLocal [-crc] <hdfs 源路径> <linux目的路径>`|从本地移动
`-put`|`-put[多个linux上的文件><hdfs路径>`|上传文件
`-copyFromLocal`|`-copyFromLocal<多个linux上的文件><hdfs路径>`|从本地复制
`-moveFromLocal`|`-moveFromLocal<多个linux上的文件><hdfs路径>`|从本地移动
`-getmerge`|`-getmerge<源路径><linux路径>`|合并到文件

### 使用python操作hdfs

python包[hdfs](https://hdfscli.readthedocs.io/)是一个可以使用python操作hdfs的工具,
可以使用`pip install hdfs`安装.大致的用法是:

+ 使用`hdfs.client.Client(url, root=None, proxy=None, timeout=None, session=None)`连接hdfs

+ 使用Client实例的`status(hdfs_path, strict=True)`方法查看路径状态

    + hdfs_path：就是hdfs路径
    + strict：设置为True时，如果hdfs_path路径不存在就会抛出异常，如果设置为False，如果路径为不存在，则返回None

+ 使用Client实例的`list(hdfs_path, status=False)`方法查看路径下对应的子目录情况

    + status：为True时，也返回子目录的状态信息，默认为Flase

+ 使用Client实例的`makedirs(hdfs_path, permission=None)`创建目录

    + permission：设置权限

+ 使用Client实例的`rename(hdfs_path, local_path)`重命名目标

+ 使用Client实例的`delete(hdfs_path, recursive=False)`删除目标

    + recursive：删除文件和其子目录，设置为False如果不存在，则会抛出异常，默认为False

+ 使用Client实例的`upload(hdfs_path, local_path, overwrite=False, n_threads=1, temp_dir=None,chunk_size=65536,progress=None, cleanup=True, **kwargs)`上传数据

    + overwrite:是否是覆盖性上传文件
    + n_threads:启动的线程数目
    + temp_dir:当overwrite=true时,远程文件一旦存在,则会在上传完之后进行交换
    + chunk_size:文件上传的大小区间
    + progress:回调函数来跟踪进度,为每一chunk_size字节.它将传递两个参数,文件上传的路径和传输的字节数.一旦完成,-1将作为第二个参数
    + cleanup:如果在上传任何文件时发生错误,则删除该文件

+ 使用Client实例的`download(hdfs_path, local_path, overwrite=False, n_threads=1, temp_dir=None, **kwargs)`方法下载文件到本地

+ 使用Client实例的`read(*args, **kwds)`读取文件

    + hdfs_path：hdfs路径
    + offset：设置开始的字节位置
    + length：读取的长度（字节为单位）
    + buffer_size：用于传输数据的字节的缓冲区的大小。默认值设置在HDFS配置。
    + encoding：制定编码
    + chunk_size：如果设置为正数，上下文管理器将返回一个发生器产生的每一chunk_size字节而不是一个类似文件的对象
    + delimiter：如果设置，上下文管理器将返回一个发生器产生每次遇到分隔符。此参数要求指定的编码。
    + progress：回调函数来跟踪进度，为每一chunk_size字节（不可用，如果块大小不是指定）。它将传递两个参数，文件上传的路径和传输的字节数。称为一次与- 1作为第二个参数。