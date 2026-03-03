from dataclasses import dataclass
from datetime import datetime
from typing import Optional


@dataclass(slots=True)
class AudioSepTask:
    id: str
    status: str
    input_audio_name: str
    separated_dir_name: Optional[str]
    error: Optional[str]
    created_at: datetime

    def audio_object_name(self) -> str:
        return f"{self.id}/{self.input_audio_name}"
