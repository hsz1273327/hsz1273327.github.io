import sys
import uuid

import asyncio
from aiokafka import AIOKafkaProducer
from aiokafka import AIOKafkaConsumer

loop = asyncio.get_event_loop()

kfkproducer = AIOKafkaProducer(loop=loop, bootstrap_servers='localhost:9092')

sourceQ = "sourceQ"
resultQ = "resultQ"
exitCh = "exitCh"


async def push(Q, value):
    await kfkproducer.send_and_wait(Q, str(value).encode())
    print(f"send {value} to {Q}")


async def get_result(producer):
    PARTITION_0 = 0
    consumer = AIOKafkaConsumer(sourceQ, loop=loop, group_id=sourceQ, bootstrap_servers='localhost:9092')
    await consumer.start()
    try:
        async for msg in consumer:
            print(f"received {msg}")
            data = int(msg.value.decode("utf-8"))
            await push(resultQ, int(data)**2)
    finally:
        await consumer.stop()


async def main():
    await kfkproducer.start()
    try:
        asyncio.ensure_future(get_result(kfkproducer))
        group_id = str(uuid.uuid4())
        print(f"group_id:{group_id}")
        consumer = AIOKafkaConsumer(exitCh, loop=loop, group_id=group_id, bootstrap_servers='localhost:9092')
        await consumer.start()
        try:
            async for msg in consumer:
                print(msg)
                if msg.value == b"Exit":
                    sys.exit(0)
        finally:
            await consumer.stop()
    finally:
        await kfkproducer.stop()


if __name__ == "__main__":
    try:
        loop.run_until_complete(main())
    except Exception as e:
        raise e
    finally:
        loop.run_until_complete(loop.shutdown_asyncgens())
        loop.close()
