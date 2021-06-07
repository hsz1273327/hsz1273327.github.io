# 消息中间件演示(Rabbitmq)

本例子演示使用redis实现分发随机数求平方和的功能

+ 生产消费模式: 生产者向`sourceQ`队列发送数据,消费者从`sourceQ`队列取数据,消费者计算完成平方后将结果放入队列`resultQ`,生产者接收`resultQ`队列中的结果更新累加结果并打印在标准输出中.
+ 广播模式: 生产者在收到`KeyboardInterrupt`错误时向频道`exitCh`发出消息,消费者订阅频道`exitCh`,当收到消息时退出.

这个实现使用`aio-pika`配合`asyncio`实现.

## 使用

1. `docker-compose up -d`启动redis
2. `pip install -r requirements.txt`安装依赖,建议使用虚拟环境
2. 执行`producer.py`和`consumer.py`
