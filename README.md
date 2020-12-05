# nginx配置跨域

这个例子用于演示如何在nginx上配置CORS.

## 依赖

+ [docker环境](https://www.docker.com/get-started)
+ [docker-compose](https://docs.docker.com/compose/install/)

## 使用

+ 静态网页托管的配置在`config/conf.d/static.d/static.conf`
+ 静态页面放在项目的`static/www`文件夹下
+ api服务在`server`文件夹下
+ 执行容器可以在`该项目根目录下`打开`terminal`使用`docker-compose up -d`
+ 浏览器中打开页面`localhost:8000/test`可以看到api的返回,打开`localhost:8001`可以看到一个包含一个按钮的页面,点击这个按钮可以看到`Hello World`,证明解决了跨域问题.