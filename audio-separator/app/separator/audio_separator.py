from abc import ABC, abstractmethod
from pathlib import Path


class AudioSeparator(ABC):
    @abstractmethod
    async def separate(self, input_path: Path, output_dir: Path, stems: int = 4) -> None:
        pass