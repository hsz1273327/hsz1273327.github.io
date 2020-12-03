# tls-example

演示https场景下的服务端和客户端配置.

服务端执行参数

+ `--redis_url`绑定链接的redis的url
+ `--authclient`强制双向认证

客户端执行参数:

+ `--authclient`强制双向认证

## 使用:

1. 启动服务端`python app.py [--authclient]`

2. 尝试使用客户端连接`python cli.py [--authclient]`

3. 将`ca/ca.crl`加入本地证书,并尝试使用浏览器访问非双向认证的服务

4. 将`clientkey/client-cert.pfx`加入本地证书,并尝试使用浏览器访问双向认证的服务