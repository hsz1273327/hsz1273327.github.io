import sys
import random
import asyncio
from aiokafka import AIOKafkaProducer
from aiokafka import AIOKafkaConsumer


loop = asyncio.get_event_loop()

sourceQ = "sourceQ"
resultQ = "resultQ"
exitCh = "exitCh"

kfkproducer = AIOKafkaProducer(loop=loop, bootstrap_servers='localhost:9092')
consumer = AIOKafkaConsumer(resultQ,loop=loop, bootstrap_servers='localhost:9092',group_id=resultQ)
async def push(Q, value):
    await kfkproducer.send_and_wait(Q,str(value).encode())
    print(f"send {value} to {Q}")


async def producer():
    while True:
        data = random.randint(1, 400)
        await push(sourceQ, data)
        await asyncio.sleep(1)


async def collector():
    sum = 0
    async for msg in consumer:
        print(f"received {msg}")
        sum += int(msg.value.decode("utf-8"))
        print(f"get sum {sum}")


async def main():
    task = asyncio.ensure_future(collector())
    print("start")
    await producer()


if __name__ == "__main__":
    loop.run_until_complete(kfkproducer.start())
    loop.run_until_complete(consumer.start())
    try:
        loop.run_until_complete(main())
    except KeyboardInterrupt:
        loop.run_until_complete(push(exitCh, "Exit"))
    except Exception as e:
        raise e
    finally:
        loop.run_until_complete(kfkproducer.stop())
        loop.run_until_complete(consumer.stop())
        loop.run_until_complete(loop.shutdown_asyncgens())
        loop.close()
