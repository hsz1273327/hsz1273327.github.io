# nginx配置websocket反向代理

这个例子用于演示如何在nginx上配置websocket反向代理,代理的ws服务的例子来自于[js攻略中服务部分](https://tutorialforjavascript.github.io/%E4%BD%BF%E7%94%A8Javascript%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/Websocket%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1.html)的[helloworld](https://github.com/TutorialForJavascript/js-server/tree/master/code/Websocket%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1/C0)

本例的测试使用的是node的[ws](https://github.com/websockets/ws)模块作为客户端

## 依赖

+ [docker环境](https://www.docker.com/get-started)
+ [docker-compose](https://docs.docker.com/compose/install/)

## 使用

+ `npm install`安装测试用客户端的依赖
+ `docker-compose up -d`启动容器
+ `npm install`安装测试的依赖
+ `npm test`测试是否可以连通ws服务
+ `npm run test_nginx`测试是否可以连通代理
