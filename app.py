import ssl
import argparse
from sanic import Sanic
from sanic.request import Request
from sanic.response import HTTPResponse, json
from aredis import StrictRedis

app = Sanic("hello_example")
client = StrictRedis.from_url("redis://host.docker.internal?db=0", decode_responses=True)


@app.get("/foo")
async def getfoo(_: Request) -> HTTPResponse:
    value = await client.get('foo')
    return json({"result": value})


@app.get("/ping")
async def ping(_: Request) -> HTTPResponse:
    return json({"result": "pong"})


@app.get("/set_foo")
async def setfoo(request: Request) -> HTTPResponse:
    value = request.args.get("value", "")
    await client.set('foo', value)
    return json({"result": "ok"})

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='执行https服务')
    parser.add_argument('--authclient', action='store_true', default=False, help='是否双向认证')
    args = parser.parse_args()
    if args.authclient:
        # 双向校验
        context = ssl.create_default_context(purpose=ssl.Purpose.CLIENT_AUTH)
        # 设置服务证书和私钥
        context.load_cert_chain("/path/to/cert", keyfile="/path/to/keyfile")
        # 设置根证书
        context.load_verify_locations('./path/to/caroot.crt')

    else:
        # 单向校验
        context = ssl.create_default_context(purpose=ssl.Purpose.SERVICE_AUTH)
        # 设置服务令牌和私钥
        context.load_cert_chain("/path/to/cert", keyfile="/path/to/keyfile")

    app.run(host="0.0.0.0", port=5000, ssl=context)
