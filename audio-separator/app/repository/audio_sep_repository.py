from entities import AudioSepTask
from typing import Optional
import asyncpg


class AudioSepTaskRepository:
    def __init__(self, pool: asyncpg.Pool):
        self.pool = pool

    async def get(self, task_id: str) -> Optional[AudioSepTask]:
        async with self.pool.acquire() as conn:
            row = await conn.fetchrow("""
                SELECT id,
                       status,
                       input_audio_name,
                       separated_dir_name,
                       error,
                       created_at
                FROM audio_sep_tasks
                WHERE id = $1
            """, task_id)

        if row is None:
            return None

        return AudioSepTask(
            id=row["id"],
            status=row["status"],
            input_audio_name=row["input_audio_name"],
            separated_dir_name=row["separated_dir_name"],
            error=row["error"],
            created_at=row["created_at"],
        )

    async def try_set_processing(self, task_id: str) -> bool:
        """
        Атомарно переводит задачу из PENDING -> PROCESSING.
        Возвращает True если удалось захватить задачу.
        """

        async with self.pool.acquire() as conn:
            result = await conn.execute("""
                UPDATE audio_sep_tasks
                SET status = 'PROCESSING'
                WHERE id = $1
                  AND status = 'PENDING'
            """, task_id)

        return result.endswith("1")

    async def mark_done(self, task_id: str, separated_dir: str):
        async with self.pool.acquire() as conn:
            await conn.execute("""
                UPDATE audio_sep_tasks
                SET status = 'DONE',
                    separated_dir_name = $1
                WHERE id = $2
            """, separated_dir, task_id)

    async def mark_error(self, task_id: str, error: str):
        async with self.pool.acquire() as conn:
            await conn.execute("""
                UPDATE audio_sep_tasks
                SET status = 'FAILED',
                    error = $1
                WHERE id = $2
            """, error, task_id)