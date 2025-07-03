from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class OAFRequest(_message.Message):
    __slots__ = ("audio_data",)
    AUDIO_DATA_FIELD_NUMBER: _ClassVar[int]
    audio_data: AudioFileData
    def __init__(self, audio_data: _Optional[_Union[AudioFileData, _Mapping]] = ...) -> None: ...

class OAFResponse(_message.Message):
    __slots__ = ("notes",)
    NOTES_FIELD_NUMBER: _ClassVar[int]
    notes: _containers.RepeatedCompositeFieldContainer[NoteEvent]
    def __init__(self, notes: _Optional[_Iterable[_Union[NoteEvent, _Mapping]]] = ...) -> None: ...

class AudioFileData(_message.Message):
    __slots__ = ("file_name", "audio_bytes")
    FILE_NAME_FIELD_NUMBER: _ClassVar[int]
    AUDIO_BYTES_FIELD_NUMBER: _ClassVar[int]
    file_name: str
    audio_bytes: bytes
    def __init__(self, file_name: _Optional[str] = ..., audio_bytes: _Optional[bytes] = ...) -> None: ...

class NoteEvent(_message.Message):
    __slots__ = ("start_seconds", "midi_pitch", "velocity", "duration_seconds")
    START_SECONDS_FIELD_NUMBER: _ClassVar[int]
    MIDI_PITCH_FIELD_NUMBER: _ClassVar[int]
    VELOCITY_FIELD_NUMBER: _ClassVar[int]
    DURATION_SECONDS_FIELD_NUMBER: _ClassVar[int]
    start_seconds: float
    midi_pitch: int
    velocity: float
    duration_seconds: float
    def __init__(self, start_seconds: _Optional[float] = ..., midi_pitch: _Optional[int] = ..., velocity: _Optional[float] = ..., duration_seconds: _Optional[float] = ...) -> None: ...
