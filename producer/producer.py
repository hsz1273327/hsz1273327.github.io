from typing import Any
from uuid import uuid4
import asyncio
import random
from aredis import StrictRedis

REDIS_URL = "redis://localhost"
redis = StrictRedis.from_url(REDIS_URL, decode_responses=True)
group = "testconsumer"

sourceQ = "sourceQ"
resultQ = "resultQ"
exitCh = "exitCh"
client_id = str(uuid4())


async def xgroup_create(name: str, group: str, *, stream_id: str = '$', auto_create: bool = False) -> Any:
    """创建绑定用户组到流.

    Args:
        name (str): 流名
        group (str): 用户组
        stream_id (str, optional): 流中的默认读取位置,$表示默认从最近开始读. Defaults to '$'.
        auto_create (bool, optional): 流不存在是否创建流. Defaults to False.

    """
    if not auto_create:
        res = await redis.execute_command('XGROUP CREATE', name, group, stream_id)
        return res
    if not await redis.exists(name):
        res = await redis.execute_command('XGROUP CREATE', name, group, stream_id, "MKSTREAM")
        return res
    res = await redis.execute_command('XGROUP CREATE', name, group, stream_id)
    return res


async def producer() -> None:
    try:
        while True:
            data = {"source": random.randint(1, 400)}
            await redis.xadd(sourceQ, data, max_len=100)
            print(f"send {data} @ {sourceQ}")
            await asyncio.sleep(1)
    except Exception as e:
        raise e


async def collector() -> None:
    """监听结果流获得数据后计算和后打印出来."""
    sum = 0
    while True:
        entries = await redis.xreadgroup(group, client_id, count=1, **dict(resultQ=">"))
        data = entries.get(resultQ)
        if data:
            for event in data:
                eid = event[0]
                msg = event[1]
                print(f"received {msg} from {resultQ}")
                sum += int(msg.get("source"))
                print(f"get sum {sum}")
                await redis.xack(resultQ, group, eid)
        else:
            await asyncio.sleep(1)


async def main() -> None:
    try:
        await redis.xinfo_consumers(sourceQ, "testconsumer")
    except Exception as e:
        print(f"get err {e}")
        await xgroup_create(sourceQ, "testconsumer", stream_id='$', auto_create=True)

    task = asyncio.ensure_future(collector())
    try:
        await producer()
    finally:
        task.cancel()

if __name__ == "__main__":
    loop = asyncio.get_event_loop()
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        try:
            loop.run_until_complete(redis.xadd(exitCh, {"event": 'Exit'}, max_len=10))
        except Exception as e:
            print("$$$$$$")
            print(e)
    except Exception as e:
        raise e
