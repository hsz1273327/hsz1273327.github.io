---
title: "树莓派配置spark on yarn"
date: 2016-11-26
author: "Hsz"
category: introduce
tags:
    - Spark
    - Hadoop
    - RaspberryPi
header-img: "img/home-bg-o.jpg"
update: 2016-11-26
series:
    raspberrypi_experiment:
        index: 8
    cluster_experiment:
        index: 5

---
# 基于树莓派的集群实验(一)--spark on yarn

终于有时间尝试集群安装了，可惜没有多余的电脑，手头有树莓派就拿来凑活用了。
本文将具体讲解如何再树莓派上安装hadoop和spark，并与自己本地的mac尝试构建分布式系统集群。
由于树莓派本身性能有限，这篇文章更多的是尝试而非可以用于大规模运算。

sudo apt-get install -y maven build-essential autoconf automake libtool cmake zlib1g-dev pkg-config libssl-dev libfuse-dev libsnappy-dev libsnappy-java libbz2-dev oracle-java8-jdk subversion

# 单机安装

## 树莓派集群网络设置

树莓派单机配置coder系统查看我的[这篇博文](!).

单机配置好后集群还需要做一些设定

主机名|ip|账户名|用途
---|---|---|---
coderMaster|192.168.1.40|pi|作为主机并允许外网ssh访问
coderSlaver01|192.168.1.41|pi|作为从机
coderSlaver02|192.168.1.42|pi|作为从机
coderSlaver03|192.168.1.43|pi|作为从机
coderSta|192.168.1.30|pi|作为反向代理服务器和通信入口,

因此,修改各机的/etc/hosts


修改权限:

sudo chmod ugo+w /etc/hosts

之后在vim中修改:

再在其中添加

    192.168.1.40 coderMaster
    192.168.1.41 coderSlaver01
    192.168.1.42 coderSlaver02
    192.168.1.43 coderSlaver03
    192.168.1.33 coderSta01

当然也可以用echo

    echo "192.168.1.40 coderMaster" >> /etc/hosts
    echo "192.168.1.41 coderSlaver01" >> /etc/hosts
    echo "192.168.1.42 coderSlaver02" >> /etc/hosts
    echo "192.168.1.50 coderSta01" >> /etc/hosts

不过一旦重启修改就会失效,怎么修改呢?

sd卡拔下来修改其中的hosts.txt文件即可

**ps:** 需要将主机的`127.0.1.1       coder`一行删除


## ssh配置


四台机器通信需要使用ssh无密码通讯

### 生成本机ssh密匙

首先为各自机器生成一对公钥，私钥

    ssh-keygen -t dsa -P "" //用dsa因为它比rsa速度快
    cat ~/.ssh/id_dsa.pub >> ~/.ssh/authorized_keys
    ssh localhost



第一次ssh localhost会要你设置下，以后就可以无密码登陆了。每台机器都这么设置下，要求都能
`ssh localhost`无密码登陆。

### master与slave建立联系

将三台机器上.ssh/文件夹下下载主机的公钥

    ssh-copy-id -i ~/.ssh/id_dsa.pub pi@coderSlaver01
    ssh-copy-id -i ~/.ssh/id_dsa.pub pi@coderSlaver02
    ssh-copy-id -i ~/.ssh/id_dsa.pub pi@coderSlaver03
    ssh-copy-id -i ~/.ssh/id_dsa.pub pi@coderSta

这样主机就可以无密码连到从机了.

同样的,将各个从机的公钥传给主机

    ssh-copy-id -i ~/.ssh/id_dsa.pub pi@coderMaster

这就算配置好了

## 编译环境配置

+ 首先是更新apt-get版本和已安装软件的版本(每台都要)

    sudo apt-get update
    sudo apt-get upgrade -y

+ 之后是安装编译所需库文件(每台最好装下)

    sudo apt-get install -y maven build-essential autoconf automake libtool cmake zlib1g-dev pkg-config libssl-dev libfuse-dev libsnappy-dev libsnappy-java libbz2-dev oracle-java8-jdk subversion

*期间因为有些库依赖java1.6所以安装了1.6版本的jvm,需要用*

    sudo update-alternatives --config java

 *将其选回java1.8*

+ 安装protocol buffers v2.5 (必须2.5)

    wget https://protobuf.googlecode.com/files/protobuf-2.5.0.tar.gz
    sudo tar xzvf protobuf-2.5.0.tar.gz
    sudo chown -R pi protobuf-2.5.0/
    cd protobuf-2.5.0
    ./configure --prefix=/usr
    make
    make check
    sudo make install

需要注意的是编译安装很慢....要有耐心,不要以为是死机了

+ 安装scala语言(每台都要)


    cd ~
    wget http://www.scala-lang.org/files/archive/scala-2.10.4.tgz
    sudo mkdir lib/scala
    sudo tar vxzf scala-2.10.4.tgz -C lib/scala


vim 修改`/etc/profile`

    #----------------------------------------------------------------------scala

    export SCALA_HOME=/home/pi/lib/scala/scala-2.10.4
    export PATH=PATH:SCALA_HOME/bin


    #---------------------------------------------------------------------------

## hadoop

### 安装hadoop

因为是编译安装所以很慢很慢,要有耐心

    cd ~
    wget http://apache.mirrors.spacedump.net/hadoop/core/hadoop-2.6.0/hadoop-2.6.0-src.tar.gz
    tar xzvf hadoop-2.6.0-src.tar.gz
    sudo chown -R pi hadoop-2.6.0-src/
    cd hadoop-2.6.0-src/

打开其中的` pom.xml `文件.将

    <additionalparam>-Xdoclint:none</additionalparam>

添加进`<properties></properties>`标签中然后保存退出编辑

因为是arm处理器,所以要打个补丁

    cd hadoop-common-project/hadoop-common/src
    wget https://issues.apache.org/jira/secure/attachment/12570212/HADOOP-9320.patch
    patch < HADOOP-9320.patch

这样安装的准备工作就完成了,开始使用maven编译安装

    sudo mvn package -Pdist,native -DskipTests -Dtar

这个要花近两个小时才能编译完...

编译后的hadoop在`hadoop-2.6.0-src/hadoop-dist/target/hadoop2.6.0`这个位置

同时会有一个你的本地编译程序.tar包在hadoop-2.6.0-src/hadoop-dist/target/hadoop2.6.0.tar.gz


### 配置hadoop

首先是环境变量配置,在/etc/profile中添加如下内容(先修改为可读`sudo chmod ugo+w /etc/profile `)

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

创建用于存储的文件夹

    sudo mkdir -p /data/hadoop/{pids,storage}
    sudo mkdir -p /data/hadoop/storage/{hdfs,tmp}
    sudo mkdir -p /data/hadoop/storage/hdfs/{name,data}
    sudo chown -R pi /data
    sudo chmod -R ugo+w /data


> 配置 core-site.xml

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

> 配置 hdfs-site.xml

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


> 配置 mapred-site.xml

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

> 配置 yarn-site.xml

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

> 配置 hadoop-env.sh

修改JAVA_HOME

    JAVA_HOME=/usr/lib/jvm/jdk-8-oracle-arm-vfp-hflt

> 配置slaves

    coderSlaver01

    coderSlaver02

    coderSlaver03

> 将 hadoop复制给子节点

    sudo scp -r ~/lib/hadoop/hadoop-2.6.0/ pi@coderSlaver01:~/lib/hadoop
    sudo scp -r ~/lib/hadoop/hadoop-2.6.0/ pi@coderSlaver02:~/lib/hadoop
    sudo scp -r ~/lib/hadoop/hadoop-2.6.0/ pi@coderSlaver03:~/lib/hadoop



### 尝试运行

    hdfs namenode -format #初始化节点
    start-dfs.sh          #开启文件系统
    start-yarn.sh         #开启yarn
    jps                   #查看开启的java进程


如果看到主节点有

    ResourceManager
    NameNode
    SecondaryNameNode

从节点有

    NodeManager
    DataNode

就说明成功了

## Spark


### 安装spark


    cd ~
    wget http://apache.crihan.fr/dist/spark/spark-1.3.1/spark-1.3.1-bin-hadoop2.6.tgz
    sudo mkdir ~/lib/spark/
    sudo tar vxzf spark-1.3.1-bin-hadoop2.6.tgz -C ~/workspace/lib/spark/
    echo "#-----------------------------------spark" >> /etc/profile
    echo "export SPARK_HOME=/home/pi/lib/spark/spark-1.3.1-bin-hadoop2.6" >> /etc/profile
    echo "#----------------------------------------" >> /etc/profile

之后修改`$SPARK_HOME/conf/slaves`中的值为(没有就新建)

    coderMaster
    coderSlaver01
    coderSlaver02
    coderSlaver03

`$SPARK_HOME/conf/spark-env.sh`中的值为(没有就新建)

    export JAVA_HOME=/usr/lib/jvm/jdk-8-oracle-arm-vfp-hflt
    export SCALA_HOME=/home/pi/lib/scala/scala-2.10.4
    export HADOOP_HOME=/home/pi/lib/hadoop/hadoop-2.6.0
    export HADOOP_CONF_DIR=$HADOOP_HOME/etc/hadoop


    export SPARK_WORKER_WEBUI_PORT=2023
    export SPARK_MASTER_WEBUI_PORT=2022
    export SPARK_MASTER_IP=192.168.1.40

    export LD_LIBRARY_PATH=/home/pi/lib/hadoop/hadoop-2.6.0/lib/native


将spark复制去各个节点,并如上设置环境变量

    sudo scp -r /home/pi/lib/spark/spark-1.3.1-bin-hadoop2.6 pi@coderSlaver01:~/lib/spark
    sudo scp -r /home/pi/lib/spark/spark-1.3.1-bin-hadoop2.6 pi@coderSlaver02:~/lib/spark
    sudo scp -r /home/pi/lib/spark/spark-1.3.1-bin-hadoop2.6 pi@coderSlaver03:~/lib/spark

### 测试

运行
```shell
$SPARK_HOME/sbin/start-all.sh
```
依然用jps查看进程,如果主节点多了`Master`,从节点多了`Worker`

打开网页`http:你的Masterip:2022`如果显示有4个worker则说明集群运行正常,如果发现有问题,可以去运行脚本时显示的对应log(.out)
文件中查看哪里出错.


### spark on yarn 运行

事实上有两种运行方式,一种是集群模式,集群模式是用来运行application的,还有一种是客户端模式,是在spark-shell中采用集群的方式
就用第二种写个简单的word count演示吧

+ 第一步 将要分析的文本送入hdfs


    hdfs dfs -mkdir /user/pi/sources
    hdfs dfs -put README.md  /user/pi/sources


+ 第二步 进入spark-shell
```shell
$SPARK_HOME/bin/spark-shell  --num-executors 3 --master yarn --executor-memory 200m
```
+ 第三步 运行一个word count


    scala> var text_file = spark.textFile("hdfs://localhost:9000/user/pi/source/README.md")
    scala> var counts = text_file.flatMap(lambda line: line.split(" ")).map(lambda word: (word, 1)).reduceByKey(lambda a, b: a + b)
    scala> counts.saveAsTextFile("hdfs://coderMaster:9000/user/pi/result/wc")

+ 第四步 另起一个shell 查看结果

    hdfs dfs -copyToLocal hdfs://coderMaster:9000/user/pi/result/wc ~/
    cd ~/wc
    cat part-00000 part-00001 >result
    cat result

##后日谈:

其实这东西娱乐都不能算,小pi内存太小了,加上节点少,性能差的叫人无语...只能说这是个集群实验了.没钱真可怜.
