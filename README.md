# nginx使用https

这个例子托管的前端页面使用的是我的[javascript攻略](https://tutorialforjavascript.github.io/)中[前端概览](https://tutorialforjavascript.github.io/web%E5%89%8D%E7%AB%AF%E6%8A%80%E6%9C%AF/%E5%89%8D%E7%AB%AF%E6%A6%82%E8%A7%88/)一篇的[helloworld项目](https://github.com/TutorialForJavascript/frontend-basic/tree/master/code/C0)中的成品

## 依赖

+ [docker环境](https://www.docker.com/get-started)
+ [docker-compose](https://docs.docker.com/compose/install/)

## 使用

+ 这个静态网页托管的配置在`config/conf.d/static.d/static.conf`
+ 静态页面放在项目的`bbs`文件夹和`www`文件夹下
+ 执行容器可以在`该项目根目录下`打开`terminal`使用`docker-compose up -d`
+ 浏览器中打开页面`https://localhost:8443`可以看到helloworld项目
