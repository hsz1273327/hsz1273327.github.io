import ssl
import argparse
from sanic import Sanic
from sanic.request import Request
from sanic.response import HTTPResponse, json
from aredis import StrictRedis


async def getfoo(request: Request) -> HTTPResponse:
    value = await request.app.redis.get('foo')
    return json({"result": value})


async def ping(_: Request) -> HTTPResponse:
    print("inside")
    return json({"result": "pong"})


async def setfoo(request: Request) -> HTTPResponse:
    value = request.args.get("value", "")
    await client.set('foo', value)
    return json({"result": "ok"})

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='执行https服务')
    parser.add_argument('--redis_url', type=str, default="redis://host.docker.internal?db=0", help='是否双向认证')
    parser.add_argument('--authclient', action='store_true', default=False, help='是否双向认证')
    args = parser.parse_args()
    if args.authclient:
        # 双向校验
        print("双向校验")
        context = ssl.SSLContext(ssl.PROTOCOL_SSLv23)
        # 设置服务证书和私钥
        context.load_cert_chain("./serverkey/server-cert.pem", keyfile="./serverkey/server-key.pem")
        # 设置根证书(双向校验有必要控制吊销的证书,文件夹下ca.pem为根证书,ca.crl为吊销证书的列表)
        context.load_verify_locations(capath='./ca')
        # 强制双向认证
        context.verify_mode = ssl.CERT_REQUIRED
        # 允许校验吊销的CRL列表
        context.verify_flags = ssl.VERIFY_CRL_CHECK_LEAF

    else:
        # 单向校验
        #context = ssl.create_default_context(purpose=ssl.Purpose.SERVER_AUTH)
        context = ssl.create_default_context(purpose=ssl.Purpose.CLIENT_AUTH)
        # 设置服务令牌和私钥
        context.load_cert_chain("./serverkey/server-cert.pem", keyfile="./serverkey/server-key.pem")
    app = Sanic("hello_example")
    client = StrictRedis.from_url(args.redis_url, decode_responses=True)
    app.redis = client
    app.add_route(getfoo, '/foo', methods=['GET'])
    app.add_route(ping, '/ping', methods=['GET'])
    app.add_route(setfoo, '/set_foo', methods=['GET'])
    app.run(host="0.0.0.0", port=5000, ssl=context)
