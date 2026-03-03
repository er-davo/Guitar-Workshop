from dataclasses import dataclass
import json
import logging

from aiokafka import AIOKafkaProducer


logger = logging.getLogger(__name__)


@dataclass
class ProducerConfig:
    bootstrap_servers: str
    topic: str



class AudioSepCompletedProducer:
    def __init__(self, config: ProducerConfig):
        self.topic = config.topic
        self.bootstrap_servers = config.bootstrap_servers
        self._producer: AIOKafkaProducer | None = None

    async def start(self):
        self._producer = AIOKafkaProducer(
            bootstrap_servers=self.bootstrap_servers,
            acks="all",
            linger_ms=5,
            retry_backoff_ms=100,
        )
        await self._producer.start()

    async def stop(self):
        if self._producer:
            await self._producer.stop()

    async def publish(self, payload: dict, key: str | None = None):
        if not self._producer:
            raise RuntimeError("Producer not started")

        value = json.dumps(payload).encode()
        key_bytes = key.encode() if key else None

        metadata = await self._producer.send_and_wait(
            self.topic,
            value=value,
            key=key_bytes,
        )

        logger.info(
            "Delivered topic=%s partition=%s offset=%s",
            metadata.topic,
            metadata.partition,
            metadata.offset,
        )