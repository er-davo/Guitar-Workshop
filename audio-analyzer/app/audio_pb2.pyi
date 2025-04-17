from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class AudioRequest(_message.Message):
    __slots__ = ("audio_path",)
    AUDIO_PATH_FIELD_NUMBER: _ClassVar[int]
    audio_path: str
    def __init__(self, audio_path: _Optional[str] = ...) -> None: ...

class AudioResponse(_message.Message):
    __slots__ = ("pitches", "times", "chromagram", "tempo", "sr", "hop_length")
    PITCHES_FIELD_NUMBER: _ClassVar[int]
    TIMES_FIELD_NUMBER: _ClassVar[int]
    CHROMAGRAM_FIELD_NUMBER: _ClassVar[int]
    TEMPO_FIELD_NUMBER: _ClassVar[int]
    SR_FIELD_NUMBER: _ClassVar[int]
    HOP_LENGTH_FIELD_NUMBER: _ClassVar[int]
    pitches: _containers.RepeatedScalarFieldContainer[float]
    times: _containers.RepeatedScalarFieldContainer[float]
    chromagram: bytes
    tempo: float
    sr: float
    hop_length: int
    def __init__(self, pitches: _Optional[_Iterable[float]] = ..., times: _Optional[_Iterable[float]] = ..., chromagram: _Optional[bytes] = ..., tempo: _Optional[float] = ..., sr: _Optional[float] = ..., hop_length: _Optional[int] = ...) -> None: ...
