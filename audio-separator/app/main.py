import asyncio
import logging
import signal

import asyncpg

from config import settings
from repository import AudioSepTaskRepository
from service import AudioSepService
from broker.consumer import AudioSepStartConsumer, ConsumerConfig
from broker.producer import AudioSepCompletedProducer, ProducerConfig
from separator import SpleeterSeparator
from storage import Storage


logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


async def create_postgres_pool():
    return await asyncpg.create_pool(
        dsn=settings.database.dsn,
        min_size=1,
        max_size=10,
    )


async def main():
    pg_pool = await create_postgres_pool()

    async with Storage(
        endpoint=settings.s3.endpoint,
        access_key=settings.s3.access_key,
        secret_key=settings.s3.secret_key,
        region_name=settings.s3.region,
    ) as storage:

        repo = AudioSepTaskRepository(pg_pool)

        separator = SpleeterSeparator()

        service = AudioSepService(
            separator=separator,
            repo=repo,
            storage=storage,
            input_bucket=settings.app.audio_bucket,
            output_bucket=settings.app.audio_bucket,
            concurrency=settings.service.concurrency,
        )

        producer_config = ProducerConfig(
            bootstrap_servers=settings.kafka.bootstrap_servers,
            topic = settings.kafka.audio_sep_completed,
        )

        completed_producer = AudioSepCompletedProducer(producer_config)

        consumer_config = ConsumerConfig(
            bootstrap_servers=settings.kafka.bootstrap_servers,
            group_id=settings.kafka.group_id,
            topic=settings.kafka.audio_sep_start,
            poll_timeout=settings.kafka.poll_timeout,
        )

        consumer = AudioSepStartConsumer(
            config=consumer_config,
            service=service,
            completed_producer=completed_producer,
            max_concurrency=settings.consumer.concurrency,
        )

        stop_event = asyncio.Event()

        def shutdown():
            logger.info("Shutdown signal received")
            stop_event.set()

        loop = asyncio.get_running_loop()
        loop.add_signal_handler(signal.SIGINT, shutdown)
        loop.add_signal_handler(signal.SIGTERM, shutdown)

        consumer_task = asyncio.create_task(consumer.start())
        consumer_task.add_done_callback(lambda t: logger.error(t.exception()) if t.exception() else None)

        await stop_event.wait()

        await consumer.stop()
        consumer_task.cancel()

    await pg_pool.close()


if __name__ == "__main__":
    asyncio.run(main())