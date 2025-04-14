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
    __slots__ = ("notes",)
    NOTES_FIELD_NUMBER: _ClassVar[int]
    notes: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, notes: _Optional[_Iterable[str]] = ...) -> None: ...
