import uuid
import asyncio
import aio_pika


AMQP_URL = 'amqp://guest:guest@localhost:5672/'


async def consumer(sourceQ, channel) -> None:
    async with sourceQ.iterator() as queue_iter:
        # Cancel consuming after __aexit__
        async for message in queue_iter:
            async with message.process():
                print(f"received {message.body}")
                result = int(message.body.decode("utf-8"))**2
            await channel.default_exchange.publish(
                aio_pika.Message(body=str(result).encode()),
                routing_key="resultQ"
            )
            print(f"send {result} to resultQ")


async def main() -> None:
    connection = await aio_pika.connect_robust(AMQP_URL)
    async with connection:
        # Creating channel
        channel = await connection.channel(1)    # type: aio_pika.Channel
        sourceQ = await channel.declare_queue("sourceQ")
        await channel.declare_queue("resultQ")
        exitCh = await connection.channel(2)
        extE = await exitCh.declare_exchange('exitCh', type=aio_pika.exchange.ExchangeType.FANOUT)
        q_id = str(uuid.uuid4())
        extQ = await exitCh.declare_queue(q_id, auto_delete=True)
        await extQ.bind(extE)
        task = asyncio.ensure_future(consumer(sourceQ, channel))
        async with extQ.iterator() as queue_iter:
            async for message in queue_iter:
                async with message.process():
                    print(f"received {message.body}")
                    if message.body == b"Exit":
                        task.cancel()
                        return


if __name__ == "__main__":
    print("consumer start")
    asyncio.run(main())
