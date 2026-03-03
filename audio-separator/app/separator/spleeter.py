import asyncio
from pathlib import Path

from .audio_separator import AudioSeparator

class SpleeterSeparator(AudioSeparator):
    def __init__(self):
        pass

    async def separate(self, input_path: Path, output_dir: Path, stems: int = 4) -> None:
        output_dir.mkdir(parents=True, exist_ok=True)

        process = await asyncio.create_subprocess_exec(
            "spleeter",
            "separate",
            "-p", f"spleeter:{stems}stems",
            "-o", str(output_dir),
            str(input_path),
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
        )

        stdout, stderr = await process.communicate()

        if process.returncode != 0:
            raise RuntimeError(
                "Spleeter failed:\n"
                f"STDOUT:\n{stdout.decode()}\n"
                f"STDERR:\n{stderr.decode()}"
            )