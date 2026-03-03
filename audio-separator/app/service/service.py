from repository import AudioSepTaskRepository
from storage import Storage
from entities import AudioSepTask
from separator import AudioSeparator

import asyncio
import logging
from pathlib import Path
import tempfile

logger = logging.getLogger(__name__)


class AudioSepService:
    def __init__(
        self,
        separator: AudioSeparator,
        repo: AudioSepTaskRepository,
        storage: Storage,
        input_bucket: str,
        output_bucket: str,
        concurrency: int = 3,
    ):
        self.separator = separator
        self.repo = repo
        self.storage = storage
        self.input_bucket = input_bucket
        self.output_bucket = output_bucket
        self.semaphore = asyncio.Semaphore(concurrency)

    async def handle(self, task_id: str):
        logger.info("Handle started task_id=%s", task_id)

        async with self.semaphore:
            task = await self.repo.get(task_id)

            if task is None:
                logger.error("Task not found task_id=%s", task_id)
                return

            locked = await self.repo.try_set_processing(task.id)
            if not locked:
                logger.info("Task already processing task_id=%s", task.id)
                return

            logger.info("Task locked task_id=%s", task.id)

            try:
                separated_dir = f"{task.id}/{task.separated_dir_name}"
                logger.info(
                    "Processing task_id=%s separated_dir=%s",
                    task.id,
                    separated_dir,
                )

                await self._process(task, separated_dir)

                await self.repo.mark_done(task.id, separated_dir)
                logger.info("Task completed task_id=%s", task.id)

            except Exception as e:
                logger.exception("Separation failed task_id=%s", task.id)
                await self.repo.mark_error(task.id, str(e))

    async def _process(self, task: AudioSepTask, separated_dir: str):
        logger.info(
            "Downloading source file task_id=%s bucket=%s object=%s",
            task.id,
            self.input_bucket,
            task.audio_object_name(),
        )

        obj = await self.storage.get_file(
            self.input_bucket,
            task.audio_object_name(),
        )

        with tempfile.TemporaryDirectory() as tmpdir:
            tmpdir_path = Path(tmpdir)

            input_path = tmpdir_path / "input.wav"
            output_dir = tmpdir_path / "separated"

            body = await obj["Body"].read()
            input_path.write_bytes(body)

            logger.info(
                "Source file saved task_id=%s path=%s size=%s bytes",
                task.id,
                input_path,
                len(body),
            )

            logger.info("Starting separation task_id=%s", task.id)
            await self.separator.separate(
                input_path,
                output_dir,
            )
            logger.info("Separation finished task_id=%s", task.id)

            # --- УБИРАЕМ ЛИШНИЙ УРОВЕНЬ (Spleeter создаёт папку с именем входного файла)
            subdirs = [p for p in output_dir.iterdir() if p.is_dir()]

            if len(subdirs) == 1:
                real_output_dir = subdirs[0]
            else:
                real_output_dir = output_dir
            # ---

            await self._upload_all_stems(real_output_dir, separated_dir)

    async def _upload_all_stems(self, output_dir: Path, separated_dir: str):
        logger.info(
            "Uploading stems from %s to bucket=%s prefix=%s",
            output_dir,
            self.output_bucket,
            separated_dir,
        )

        for file_path in output_dir.rglob("*"):
            if not file_path.is_file():
                continue

            relative_path = file_path.relative_to(output_dir)
            object_name = f"{separated_dir}/{relative_path.as_posix()}"
            size = file_path.stat().st_size

            logger.info(
                "Uploading stem object=%s size=%s bytes",
                object_name,
                size,
            )

            with file_path.open("rb") as f:
                await self.storage.upload_file(
                    self.output_bucket,
                    object_name,
                    f,
                    size,
                )

        logger.info("All stems uploaded prefix=%s", separated_dir)