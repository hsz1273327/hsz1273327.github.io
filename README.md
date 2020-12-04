# 反向代理示例

这个例子代理[使用Javascript构建RESTful接口服务](https://tutorialforjavascript.github.io/%E4%BD%BF%E7%94%A8Javascript%E6%90%AD%E5%BB%BA%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1/RESTful%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1.html)文中的[C0](https://github.com/TutorialForJavascript/js-server/tree/master/code/RESTful%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1/C0)和[C2](https://github.com/TutorialForJavascript/js-server/tree/master/code/RESTful%E6%8E%A5%E5%8F%A3%E6%9C%8D%E5%8A%A1/C2)两个例子,这两个例子都已经打包好了上传在我的dockerhub下.

## 依赖

+ [docker环境](https://www.docker.com/get-started)
+ [docker-compose](https://docs.docker.com/compose/install/)

## 使用

+ 执行容器可以在`该项目根目录下`打开`terminal`使用`docker-compose up -d`
+ 浏览器中打开页面`https://localhost:5000`可以看到`/`下是C0的接口,`/api/`下是C2的接口
