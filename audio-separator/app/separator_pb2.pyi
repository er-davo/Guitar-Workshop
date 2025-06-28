from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class SeparateRequest(_message.Message):
    __slots__ = ("audio_data",)
    AUDIO_DATA_FIELD_NUMBER: _ClassVar[int]
    audio_data: AudioFileData
    def __init__(self, audio_data: _Optional[_Union[AudioFileData, _Mapping]] = ...) -> None: ...

class SeparateResponse(_message.Message):
    __slots__ = ("stems",)
    class StemsEntry(_message.Message):
        __slots__ = ("key", "value")
        KEY_FIELD_NUMBER: _ClassVar[int]
        VALUE_FIELD_NUMBER: _ClassVar[int]
        key: str
        value: AudioFileData
        def __init__(self, key: _Optional[str] = ..., value: _Optional[_Union[AudioFileData, _Mapping]] = ...) -> None: ...
    STEMS_FIELD_NUMBER: _ClassVar[int]
    stems: _containers.MessageMap[str, AudioFileData]
    def __init__(self, stems: _Optional[_Mapping[str, AudioFileData]] = ...) -> None: ...

class AudioFileData(_message.Message):
    __slots__ = ("file_name", "audio_bytes")
    FILE_NAME_FIELD_NUMBER: _ClassVar[int]
    AUDIO_BYTES_FIELD_NUMBER: _ClassVar[int]
    file_name: str
    audio_bytes: bytes
    def __init__(self, file_name: _Optional[str] = ..., audio_bytes: _Optional[bytes] = ...) -> None: ...
