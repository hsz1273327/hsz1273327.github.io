# nginx使用https

这个例子托管的前端页面使用的是我的[javascript攻略](https://tutorialforjavascript.github.io/)中[前端概览](https://tutorialforjavascript.github.io/web%E5%89%8D%E7%AB%AF%E6%8A%80%E6%9C%AF/%E5%89%8D%E7%AB%AF%E6%A6%82%E8%A7%88/)一篇的[helloworld项目](https://github.com/TutorialForJavascript/frontend-basic/tree/master/code/C0)中的成品

## 依赖

+ [docker环境](https://www.docker.com/get-started)
+ [docker-compose](https://docs.docker.com/compose/install/)

## 使用

+ 这个静态网页托管的配置在`config/conf.d/static.d/static.conf`
+ 静态页面放在项目的`bbs`文件夹和`www`文件夹下
+ build镜像可以在mac或linux下在`该项目根目录下`打开`terminal`直接执行`bash build.sh`,windows下的可以自己写个cmd脚本,我不会
+ 执行容器可以在`该项目根目录下`打开`terminal`使用`docker-compose up -d`
+ 浏览器中打开页面`https://localhost:8043`可以看到helloworld项目

## 说明

使用http协议访问`http://localhost:8080`也会显示,因为我们使用的路径`/usr/share/nginx/html`本来就是nginx的默认静态页面路径