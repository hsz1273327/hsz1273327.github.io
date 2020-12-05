# 反向代理的负载均衡示例

这个例子这个例子使用了3个服务做反向代理的负载均衡.代理的服务是[使用Javascript构建RESTful接口服务](https://tutorialforjavascript.github.io/%E4%BD%BF%E7%94%A8Javascript%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/RESTful%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1.html)文中的[C0](https://github.com/TutorialForJavascript/js-server/tree/master/code/RESTful%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1/C0).例子已经打包好了上传在我的dockerhub下.有兴趣的可以去看下具体的实现.

## 依赖

+ [docker环境](https://www.docker.com/get-started)
+ [docker-compose](https://docs.docker.com/compose/install/)

## 使用

+ 执行容器可以在`该项目根目录下`打开`terminal`使用`docker-compose up -d`
+ 浏览器中打开页面`https://localhost:5000`可以看到`/`下代理了接口
