// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: separator.proto

package separator

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SeparateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AudioData     *AudioFileData         `protobuf:"bytes,1,opt,name=audio_data,json=audioData,proto3" json:"audio_data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SeparateRequest) Reset() {
	*x = SeparateRequest{}
	mi := &file_separator_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SeparateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SeparateRequest) ProtoMessage() {}

func (x *SeparateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_separator_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SeparateRequest.ProtoReflect.Descriptor instead.
func (*SeparateRequest) Descriptor() ([]byte, []int) {
	return file_separator_proto_rawDescGZIP(), []int{0}
}

func (x *SeparateRequest) GetAudioData() *AudioFileData {
	if x != nil {
		return x.AudioData
	}
	return nil
}

type SeparateResponse struct {
	state         protoimpl.MessageState    `protogen:"open.v1"`
	Stems         map[string]*AudioFileData `protobuf:"bytes,1,rep,name=stems,proto3" json:"stems,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SeparateResponse) Reset() {
	*x = SeparateResponse{}
	mi := &file_separator_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SeparateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SeparateResponse) ProtoMessage() {}

func (x *SeparateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_separator_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SeparateResponse.ProtoReflect.Descriptor instead.
func (*SeparateResponse) Descriptor() ([]byte, []int) {
	return file_separator_proto_rawDescGZIP(), []int{1}
}

func (x *SeparateResponse) GetStems() map[string]*AudioFileData {
	if x != nil {
		return x.Stems
	}
	return nil
}

type AudioFileData struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	FileName      string                 `protobuf:"bytes,1,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	AudioBytes    []byte                 `protobuf:"bytes,2,opt,name=audio_bytes,json=audioBytes,proto3" json:"audio_bytes,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AudioFileData) Reset() {
	*x = AudioFileData{}
	mi := &file_separator_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AudioFileData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AudioFileData) ProtoMessage() {}

func (x *AudioFileData) ProtoReflect() protoreflect.Message {
	mi := &file_separator_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AudioFileData.ProtoReflect.Descriptor instead.
func (*AudioFileData) Descriptor() ([]byte, []int) {
	return file_separator_proto_rawDescGZIP(), []int{2}
}

func (x *AudioFileData) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *AudioFileData) GetAudioBytes() []byte {
	if x != nil {
		return x.AudioBytes
	}
	return nil
}

var File_separator_proto protoreflect.FileDescriptor

const file_separator_proto_rawDesc = "" +
	"\n" +
	"\x0fseparator.proto\x12\tseparator\"J\n" +
	"\x0fSeparateRequest\x127\n" +
	"\n" +
	"audio_data\x18\x01 \x01(\v2\x18.separator.AudioFileDataR\taudioData\"\xa4\x01\n" +
	"\x10SeparateResponse\x12<\n" +
	"\x05stems\x18\x01 \x03(\v2&.separator.SeparateResponse.StemsEntryR\x05stems\x1aR\n" +
	"\n" +
	"StemsEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12.\n" +
	"\x05value\x18\x02 \x01(\v2\x18.separator.AudioFileDataR\x05value:\x028\x01\"M\n" +
	"\rAudioFileData\x12\x1b\n" +
	"\tfile_name\x18\x01 \x01(\tR\bfileName\x12\x1f\n" +
	"\vaudio_bytes\x18\x02 \x01(\fR\n" +
	"audioBytes2Z\n" +
	"\x0eAudioSeparator\x12H\n" +
	"\rSeparateAudio\x12\x1a.separator.SeparateRequest\x1a\x1b.separator.SeparateResponseB\x1aZ\x18internal/proto/separatorb\x06proto3"

var (
	file_separator_proto_rawDescOnce sync.Once
	file_separator_proto_rawDescData []byte
)

func file_separator_proto_rawDescGZIP() []byte {
	file_separator_proto_rawDescOnce.Do(func() {
		file_separator_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_separator_proto_rawDesc), len(file_separator_proto_rawDesc)))
	})
	return file_separator_proto_rawDescData
}

var file_separator_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_separator_proto_goTypes = []any{
	(*SeparateRequest)(nil),  // 0: separator.SeparateRequest
	(*SeparateResponse)(nil), // 1: separator.SeparateResponse
	(*AudioFileData)(nil),    // 2: separator.AudioFileData
	nil,                      // 3: separator.SeparateResponse.StemsEntry
}
var file_separator_proto_depIdxs = []int32{
	2, // 0: separator.SeparateRequest.audio_data:type_name -> separator.AudioFileData
	3, // 1: separator.SeparateResponse.stems:type_name -> separator.SeparateResponse.StemsEntry
	2, // 2: separator.SeparateResponse.StemsEntry.value:type_name -> separator.AudioFileData
	0, // 3: separator.AudioSeparator.SeparateAudio:input_type -> separator.SeparateRequest
	1, // 4: separator.AudioSeparator.SeparateAudio:output_type -> separator.SeparateResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_separator_proto_init() }
func file_separator_proto_init() {
	if File_separator_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_separator_proto_rawDesc), len(file_separator_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_separator_proto_goTypes,
		DependencyIndexes: file_separator_proto_depIdxs,
		MessageInfos:      file_separator_proto_msgTypes,
	}.Build()
	File_separator_proto = out.File
	file_separator_proto_goTypes = nil
	file_separator_proto_depIdxs = nil
}
