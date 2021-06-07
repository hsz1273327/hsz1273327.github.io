import random
import asyncio
import aio_pika

AMQP_URL = 'amqp://guest:guest@localhost:5672/'


async def producer(channel):
    while True:
        data = random.randint(1, 400)
        await channel.default_exchange.publish(
            aio_pika.Message(body=str(data).encode()),
            routing_key="sourceQ"
        )
        print(f"send {data} to sourceQ")
        await asyncio.sleep(1)


async def collector(Q):
    sum = 0
    async with Q.iterator() as queue_iter:
            # Cancel consuming after __aexit__
        async for message in queue_iter:
            async with message.process():
                print(f"received {message.body}")
                sum += int(message.body)
                print(f"get sum {sum}")


async def pubExt():
    print("exit publishing")
    connection = await aio_pika.connect_robust(AMQP_URL)
    async with connection:
        exitCh = await connection.channel(2)
        ex = await exitCh.declare_exchange('exitCh', type=aio_pika.exchange.ExchangeType.FANOUT)
        await ex.publish(aio_pika.Message(body=b"Exit"), routing_key="")
    print("exit published")


async def main():
    connection = await aio_pika.connect_robust(AMQP_URL)
    async with connection:
        # Creating channel
        channel = await connection.channel(1)    # type: aio_pika.Channel
        await channel.declare_queue("sourceQ")
        resultQ = await channel.declare_queue("resultQ")
        prodcor = producer(channel)
        collcor = collector(resultQ)
        await asyncio.gather(prodcor, collcor)


if __name__ == "__main__":
    print("producer start")
    try:
        asyncio.get_event_loop().run_until_complete(main())
    except KeyboardInterrupt:
        asyncio.get_event_loop().run_until_complete(pubExt())
    finally:
        print("done")
