from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Style(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    UNSPECIFIED: _ClassVar[Style]
    ROCK: _ClassVar[Style]
    BLUES: _ClassVar[Style]
    METAL: _ClassVar[Style]
    JAZZ: _ClassVar[Style]
    FUNK: _ClassVar[Style]
UNSPECIFIED: Style
ROCK: Style
BLUES: Style
METAL: Style
JAZZ: Style
FUNK: Style

class RiffRequest(_message.Message):
    __slots__ = ("tone", "style")
    TONE_FIELD_NUMBER: _ClassVar[int]
    STYLE_FIELD_NUMBER: _ClassVar[int]
    tone: str
    style: Style
    def __init__(self, tone: _Optional[str] = ..., style: _Optional[_Union[Style, str]] = ...) -> None: ...

class RiffResponse(_message.Message):
    __slots__ = ("riff",)
    RIFF_FIELD_NUMBER: _ClassVar[int]
    riff: str
    def __init__(self, riff: _Optional[str] = ...) -> None: ...
