from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Iterable as _Iterable, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class RequestType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    FILE: _ClassVar[RequestType]
    YOUTUBE: _ClassVar[RequestType]
FILE: RequestType
YOUTUBE: RequestType

class AudioRequest(_message.Message):
    __slots__ = ("audio_path", "type")
    AUDIO_PATH_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    audio_path: str
    type: RequestType
    def __init__(self, audio_path: _Optional[str] = ..., type: _Optional[_Union[RequestType, str]] = ...) -> None: ...

class AudioResponse(_message.Message):
    __slots__ = ("note_features",)
    NOTE_FEATURES_FIELD_NUMBER: _ClassVar[int]
    note_features: _containers.RepeatedCompositeFieldContainer[AudioEvent]
    def __init__(self, note_features: _Optional[_Iterable[_Union[AudioEvent, _Mapping]]] = ...) -> None: ...

class AudioEvent(_message.Message):
    __slots__ = ("time", "pitch", "main_note", "octave", "chroma_notes")
    TIME_FIELD_NUMBER: _ClassVar[int]
    PITCH_FIELD_NUMBER: _ClassVar[int]
    MAIN_NOTE_FIELD_NUMBER: _ClassVar[int]
    OCTAVE_FIELD_NUMBER: _ClassVar[int]
    CHROMA_NOTES_FIELD_NUMBER: _ClassVar[int]
    time: float
    pitch: float
    main_note: str
    octave: int
    chroma_notes: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, time: _Optional[float] = ..., pitch: _Optional[float] = ..., main_note: _Optional[str] = ..., octave: _Optional[int] = ..., chroma_notes: _Optional[_Iterable[str]] = ...) -> None: ...
