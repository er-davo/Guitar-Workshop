from broker.producer import AudioSepCompletedProducer
from service import AudioSepService

from dataclasses import dataclass
import logging
import json
import asyncio

from aiokafka import AIOKafkaConsumer, ConsumerRecord

logger = logging.getLogger(__name__)


@dataclass
class ConsumerConfig:
    bootstrap_servers: str
    group_id: str
    topic: str
    poll_timeout: float


class AudioSepStartConsumer:
    def __init__(
        self,
        config: ConsumerConfig,
        service: AudioSepService,
        completed_producer: AudioSepCompletedProducer,
        max_concurrency: int = 4,
    ):
        self.config = config
        self.service = service
        self.completed_producer = completed_producer

        self.consumer = AIOKafkaConsumer(
            config.topic,
            bootstrap_servers=config.bootstrap_servers,
            group_id=config.group_id,
            enable_auto_commit=False,
            auto_offset_reset="earliest",
            session_timeout_ms=10000,
            max_poll_interval_ms=900000,
        )

        self.semaphore = asyncio.Semaphore(max_concurrency)
        self._running = False

    async def start(self):
        logger.info("Before consumer.start()")
        await self.consumer.start()
        logger.info("After consumer.start()")
        await self.completed_producer.start()
        logger.info("After completed_producer.start()")

        self._running = True
        logger.info("AudioSepConsumer started")

        try:
            async for msg in self.consumer:
                logger.info("Received message: %s", msg)
                await self.semaphore.acquire()
                asyncio.create_task(self._handle_message(msg))
        finally:
            await self.consumer.stop()
            await self.completed_producer.stop()

    async def _handle_message(self, msg: ConsumerRecord):
        try:
            event = json.loads(msg.value.decode())
            task_id = event["id"]

            logger.info(
                "Processing task %s (partition=%s offset=%s)",
                task_id,
                msg.partition,
                msg.offset,
            )

            await self.service.handle(task_id)

            await self.completed_producer.publish(
                payload={"id": task_id},
                key=task_id,
            )

            await self.consumer.commit()

        except Exception as e:
            logger.error("Task failed: %s", e, exc_info=True)
            await self.consumer.commit()

        finally:
            self.semaphore.release()

    async def stop(self):
        self._running = False
        await self.consumer.stop()
        await self.completed_producer.stop()