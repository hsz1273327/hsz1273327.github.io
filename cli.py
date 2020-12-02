import argparse
import requests as rq

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='执行https服务')
    parser.add_argument('--authclient', action='store_true', default=False, help='是否双向认证')
    args = parser.parse_args()
    if args.authclient:
        res = rq.get('https://github.com', verify='/path/to/certfile',cert=('/path/client.cert', '/path/client.key'))
    else:
        res = rq.get('https://github.com', verify='/path/to/certfile')
    print(res.json())