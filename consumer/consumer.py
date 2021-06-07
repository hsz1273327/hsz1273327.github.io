import asyncio
from uuid import uuid4
from typing import Any
from aredis import StrictRedis

REDIS_URL = "redis://localhost"
redis = StrictRedis.from_url(REDIS_URL, decode_responses=True)
maxlen = 100
sourceQ = "sourceQ"
resultQ = "resultQ"
exitCh = "exitCh"
group = "testconsumer"
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


async def process() -> None:
    """监听数据流获得数据后计算平方放入结果流."""
    print("process start")
    while True:
        entries = await redis.xreadgroup(group, client_id, count=1, **dict(sourceQ=">"))
        data = entries.get(sourceQ)
        if data:
            for event in data:
                eid = event[0]
                msg = event[1]
                print(f"received {msg} from {sourceQ}")
                result = {"source": int(msg.get("source"))**2}
                await redis.xadd(resultQ, result, max_len=100)
                print(f"send {result} @ {resultQ}")
                await redis.xack(sourceQ, group, eid)
        else:
            await asyncio.sleep(1)


async def main() -> None:
    try:
        await redis.xinfo_consumers(sourceQ, group)
    except Exception as e:
        print(f"get err {e}")
        await xgroup_create(sourceQ, group, stream_id='$', auto_create=True)

    try:
        await redis.xinfo_consumers(resultQ, group)
    except Exception as e:
        print(f"get err {e}")
        await xgroup_create(resultQ, group, stream_id='$', auto_create=True)

    task = asyncio.create_task(process())
    print("watch exitCh start")
    await xgroup_create(exitCh, client_id, stream_id='$', auto_create=True)
    try:
        while True:
            entries = await redis.xreadgroup(client_id, client_id, count=1, **{"exitCh": ">"})
            signals = entries.get("exitCh")
            if signals:
                print("get exit signal from producer")
                for signal in signals:
                    eid = signal[0]
                    msg = signal[1]
                    print(f"received signal {msg}")
                    await redis.xack("exitCh", client_id, eid)
                break
            else:
                await asyncio.sleep(0.1)
    finally:
        await redis.xgroup_destroy(exitCh, client_id)
        task.cancel()
        print("exit")


if __name__ == "__main__":
    asyncio.run(main())
