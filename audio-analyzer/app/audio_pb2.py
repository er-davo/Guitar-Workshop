# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: audio.proto
# Protobuf Python Version: 5.29.0
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    29,
    0,
    '',
    'audio.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0b\x61udio.proto\x12\x05\x61udio\"D\n\x0c\x41udioRequest\x12\x12\n\naudio_path\x18\x01 \x01(\t\x12 \n\x04type\x18\x02 \x01(\x0e\x32\x12.audio.RequestType\"9\n\rAudioResponse\x12(\n\rnote_features\x18\x01 \x03(\x0b\x32\x11.audio.AudioEvent\"b\n\nAudioEvent\x12\x0c\n\x04time\x18\x01 \x01(\x02\x12\r\n\x05pitch\x18\x02 \x01(\x02\x12\x11\n\tmain_note\x18\x03 \x01(\t\x12\x0e\n\x06octave\x18\x04 \x01(\x05\x12\x14\n\x0c\x63hroma_notes\x18\x05 \x03(\t*$\n\x0bRequestType\x12\x08\n\x04\x46ILE\x10\x00\x12\x0b\n\x07YOUTUBE\x10\x01\x32J\n\rAudioAnalyzer\x12\x39\n\x0cProcessAudio\x12\x13.audio.AudioRequest\x1a\x14.audio.AudioResponseB\x15Z\x13internal/audioprotob\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'audio_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\023internal/audioproto'
  _globals['_REQUESTTYPE']._serialized_start=251
  _globals['_REQUESTTYPE']._serialized_end=287
  _globals['_AUDIOREQUEST']._serialized_start=22
  _globals['_AUDIOREQUEST']._serialized_end=90
  _globals['_AUDIORESPONSE']._serialized_start=92
  _globals['_AUDIORESPONSE']._serialized_end=149
  _globals['_AUDIOEVENT']._serialized_start=151
  _globals['_AUDIOEVENT']._serialized_end=249
  _globals['_AUDIOANALYZER']._serialized_start=289
  _globals['_AUDIOANALYZER']._serialized_end=363
# @@protoc_insertion_point(module_scope)
