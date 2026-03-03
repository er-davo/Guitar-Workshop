import asyncio
from pathlib import Path

from .audio_separator import AudioSeparator


class DemucsSeparator(AudioSeparator):
    def __init__(self, device: str = "cpu", model: str = "mdx_extra_q"):
        self.device = device
        self.model = model

    async def separate(self, input_path: Path, output_dir: Path, stems: int) -> None:
        output_dir.mkdir(parents=True, exist_ok=True)

        process = await asyncio.create_subprocess_exec(
            "demucs",
            str(input_path),
            "--out", str(output_dir),
            "--device", self.device,
            "--repo", "/root/.cache/torch/hub/checkpoints",
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
        )

        stdout, stderr = await process.communicate()

        if process.returncode != 0:
            raise RuntimeError(
                "Demucs failed:\n"
                f"STDOUT:\n{stdout.decode()}\n"
                f"STDERR:\n{stderr.decode()}"
            )