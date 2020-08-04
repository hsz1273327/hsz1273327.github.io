---
layout: post
title: "树莓派上的集群实验--Spark on Yarn"
date: 2018-02-25
author: "Hsz"
category: blog
tags:
    - Spark
    - Hadoop
    - Linux
    - 集群
    - 树莓派
header-img: "img/home-bg-o.jpg"
---
# spark

spark是前几年非常火的一个分布式计算工具套件,它提供了一套很完备的分布式计算工具,并可以与hadoop生态链很好的结合.本文的spark搭建在yarn上,这需要先安装hadoop,hadoop的安装在[]()一文中已经写过这边不再复述.

## 单机安装

spark单机就可以自己运转,spark是jvm上的软件,只要安装过hadoop的环境就已经可以直接下载使用spark了.
注意看好自己的hadoop版本选择对应的spark二进制包

```shell
wget https://www.apache.org/dyn/closer.lua/spark/spark-2.2.1/spark-2.2.1-bin-hadoop2.7.tgz
sudo mkdir ~/lib/spark/
sudo tar vxzf spark-2.2.1-bin-hadoop2.7.tgz -C ~/workspace/lib/spark/
echo "#-----------------------------------spark" >> /etc/profile
echo "export SPARK_HOME=/home/pi/lib/spark/spark-2.2.1-bin-hadoop2.7.tgz" >> /etc/profile
echo "#----------------------------------------" >> /etc/profile
```

<!-- ## 编译安装

+ 安装依赖

    spark是用scala写的,要编译安装需要先有scala
    ```shell
    wget http://www.scala-lang.org/files/archive/scala-2.10.4.tgz
    sudo mkdir lib/scala
    sudo tar vxzf scala-2.10.4.tgz -C lib/scala
    ```
    然后添加环境变量
    ```shell
    #---------------------------------------------------scala
    export SCALA_HOME=/home/pi/lib/scala/scala-2.10.4
    export PATH=$PATH:$SCALA_HOME/bin
    #---------------------------------------------------- -->
    ```
### 配置

+ 修改`$SPARK_HOME/conf/slaves`中的值为(没有就新建)

    ```shell
    piNodeMaster
    piNode01
    piNode02
    piNode03
    ```

+ 修改`$SPARK_HOME/conf/spark-env.sh`中的值为(没有就新建)

    ```shell
    export JAVA_HOME=/usr/lib/jvm/jdk-8-oracle-arm-vfp-hflt
    export SCALA_HOME=/home/pi/lib/scala/scala-2.10.4
    export HADOOP_HOME=/home/pi/lib/hadoop/hadoop-2.6.0
    export HADOOP_CONF_DIR=$HADOOP_HOME/etc/hadoop
    export SPARK_WORKER_WEBUI_PORT=2023
    export SPARK_MASTER_WEBUI_PORT=2022
    export SPARK_MASTER_IP=192.168.1.40
    export LD_LIBRARY_PATH=/home/pi/lib/hadoop/hadoop-2.6.0/lib/native
    ```

## 集群安装配置

将spark复制去各个节点,并如上设置环境变量

```shell
sudo scp -r /home/pi/lib/spark/spark-1.3.1-bin-hadoop2.6 pi@coderSlaver01:~/lib/spark
sudo scp -r /home/pi/lib/spark/spark-1.3.1-bin-hadoop2.6 pi@coderSlaver02:~/lib/spark
sudo scp -r /home/pi/lib/spark/spark-1.3.1-bin-hadoop2.6 pi@coderSlaver03:~/lib/spark
```

## 测试启动成功

命令行使用`$SPARK_HOME/sbin/start-all.sh`运行spark

用jps查看进程,如果主节点多了`Master`,从节点多了`Worker`

打开网页`http:你的Masterip:2022`如果显示有4个worker则说明集群运行正常,如果发现有问题,可以去运行脚本时显示的对应log(.out)
文件中查看哪里出错.

## spark运行测试

spark有两种运行方式,

+ 集群模式

    集群模式是用来运行application的,通过`$SPARK_HOME/bin/spark-submit`工具执行.
    spark-submit的参数有:

    + `--class <str>`应用程序的入口点,比如`org.apache.spark.examples.SparkPi`
    + `--master <master url>` 指定主机的url,合法的主机url包括:
        
    Master URL|含义
    ---|---
    `local`|使用一个线程本地运行Spark
    `local[K]`|使用K个worker线程本地运行Spark(理想情况下,设置这个值的数量应为机器的 core 数量)
    `local[K,F]`|使用K个worker线程本地运行Spark并允许最多失败F次
    `local[*]`|使用更多的worker线程作为逻辑的core在机器上来本地的运行Spark
    `local[*,F]`|使用更多的worker线程作为逻辑的core在您的机器上来本地的运行Spark并允许最多失败F次
    `spark://HOST:PORT`|连接至给定的spark独立模式下的集群.该端口必须有一个作为`master`配置来使用,默认是7077.
    `spark://HOST1:PORT1,HOST2:PORT2`|连接至给定的spark独立模式下的集群,并使用`Zookeeper`.该列表必须包含由`zookeeper`设置的高可用集群中的所有`master`主机,该端口必须有一个作为`master`配置来使用,默认是7077.
    `mesos://HOST:PORT`|连接至给定的Mesos集群.该端口必须有一个作为您的配置来使用,默认是5050.或者,对于使用了`ZooKeeper`的`Mesos cluster`来说,使用 `mesos://zk://....`也行.此时应使用`--deploy-mode cluster`来提交，该`HOST:PORT`应该被配置以连接到`MesosClusterDispatcher`.
    `yarn`|连接至一个YARN集群.该集群的位置将根据环境变量`HADOOP_CONF_DIR`或者`YARN_CONF_DIR` 来找到.

    + `--deploy-mode <str>`驱动部署方式,默认为`client`,包括

        + `cluster`在worker节点上部署,单机版本和spark的独立模式没有这个模式,一般来说这种模式效率更高些,因为没有驱动和执行器之间的延迟
        + `client`在本地作为外部客户端部署,客户端模式实际上就是这种部署方式

    + `--executor-memory <int + M/G>`
    + `--num-executors <int>` 指定执行器的个数
    + `--total-executor-cores <int>`指定执行器总共使用多少个核
    + `--supervise`用来确保driver自动重启

    这些参数后则是指定app的地址和参数

+ 客户端模式

    在`spark-shell`中使用交互模式操作spark,本质上就是`--deploy-mode`为`client`的application

## 客户端模式写个简单的word count演示

集群模式就不演示了,毕竟java技术栈几乎为0,此处演示的是利用`spark-shell`在交互模式下执行简单word count的操作.

+ 第一步 将要分析的文本送入hdfs

    ```shell
    hadoop fs -mkdir /user/pi/sources
    hadoop fs -put README.md  /user/pi/sources
    ```

+ 第二步 进入spark-shell

    ```shell
    $SPARK_HOME/bin/spark-shell  --num-executors 3 --master yarn --executor-memory 200m
    ```

+ 第三步 运行一个word count

    ```shell
    scala> var text_file = spark.textFile("hdfs://localhost:9000/user/pi/source/README.md")
    scala> var counts = text_file.flatMap(lambda line: line.split(" ")).map(lambda word: (word, 1)).reduceByKey(lambda a, b: a + b)
    scala> counts.saveAsTextFile("hdfs://coderMaster:9000/user/pi/result/wc")
    ```

+ 第四步 另起一个shell 查看结果

    ```shell
    hadoop fs -copyToLocal hdfs://coderMaster:9000/user/pi/result/wc ~/
    cd ~/wc
    cat part-00000 part-00001 >result
    cat result
    ```

## 使用pyspark

spark很早就支持python,pyspark是spark官方的python接口库.原本是和spark一起安装的,现在也可以使用pip安装.pyspark本质上就是使用py4j来将python代码转译到jvm上运行而已.

集群模式中使用pyspark就是写好python程序,然后通过`$SPARK_HOME/bin/spark-submit`工具执行即可,
而客户端模式也就是在python的交互模式下使用pyspark即可.通常,为了方便我们会将pyspark与jupyter notebook结合使用,jupyter notebook的相关介绍和经验我在[Jupyter攻略](http://blog.hszofficial.site/TutorialForJupyter/)系列文章中已经做过介绍,此处不再复述.

### 集群模式使用pyspark演示word count


### 客户端模式使用pyspark演示word count


## spark 结合 hive


## spark 结合 hbase