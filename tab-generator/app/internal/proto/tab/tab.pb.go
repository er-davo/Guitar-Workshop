// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: tab.proto

package tab

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

type TabRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Audio         *AudioFileData         `protobuf:"bytes,1,opt,name=audio,proto3" json:"audio,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TabRequest) Reset() {
	*x = TabRequest{}
	mi := &file_tab_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TabRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TabRequest) ProtoMessage() {}

func (x *TabRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tab_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TabRequest.ProtoReflect.Descriptor instead.
func (*TabRequest) Descriptor() ([]byte, []int) {
	return file_tab_proto_rawDescGZIP(), []int{0}
}

func (x *TabRequest) GetAudio() *AudioFileData {
	if x != nil {
		return x.Audio
	}
	return nil
}

type TabResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Tab           string                 `protobuf:"bytes,1,opt,name=tab,proto3" json:"tab,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TabResponse) Reset() {
	*x = TabResponse{}
	mi := &file_tab_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TabResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TabResponse) ProtoMessage() {}

func (x *TabResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tab_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TabResponse.ProtoReflect.Descriptor instead.
func (*TabResponse) Descriptor() ([]byte, []int) {
	return file_tab_proto_rawDescGZIP(), []int{1}
}

func (x *TabResponse) GetTab() string {
	if x != nil {
		return x.Tab
	}
	return ""
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
	mi := &file_tab_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AudioFileData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AudioFileData) ProtoMessage() {}

func (x *AudioFileData) ProtoReflect() protoreflect.Message {
	mi := &file_tab_proto_msgTypes[2]
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
	return file_tab_proto_rawDescGZIP(), []int{2}
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

var File_tab_proto protoreflect.FileDescriptor

const file_tab_proto_rawDesc = "" +
	"\n" +
	"\ttab.proto\x12\x03tab\"6\n" +
	"\n" +
	"TabRequest\x12(\n" +
	"\x05audio\x18\x01 \x01(\v2\x12.tab.AudioFileDataR\x05audio\"\x1f\n" +
	"\vTabResponse\x12\x10\n" +
	"\x03tab\x18\x01 \x01(\tR\x03tab\"M\n" +
	"\rAudioFileData\x12\x1b\n" +
	"\tfile_name\x18\x01 \x01(\tR\bfileName\x12\x1f\n" +
	"\vaudio_bytes\x18\x02 \x01(\fR\n" +
	"audioBytes2?\n" +
	"\vTabGenerate\x120\n" +
	"\vGenerateTab\x12\x0f.tab.TabRequest\x1a\x10.tab.TabResponseB\x14Z\x12internal/proto/tabb\x06proto3"

var (
	file_tab_proto_rawDescOnce sync.Once
	file_tab_proto_rawDescData []byte
)

func file_tab_proto_rawDescGZIP() []byte {
	file_tab_proto_rawDescOnce.Do(func() {
		file_tab_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_tab_proto_rawDesc), len(file_tab_proto_rawDesc)))
	})
	return file_tab_proto_rawDescData
}

var file_tab_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_tab_proto_goTypes = []any{
	(*TabRequest)(nil),    // 0: tab.TabRequest
	(*TabResponse)(nil),   // 1: tab.TabResponse
	(*AudioFileData)(nil), // 2: tab.AudioFileData
}
var file_tab_proto_depIdxs = []int32{
	2, // 0: tab.TabRequest.audio:type_name -> tab.AudioFileData
	0, // 1: tab.TabGenerate.GenerateTab:input_type -> tab.TabRequest
	1, // 2: tab.TabGenerate.GenerateTab:output_type -> tab.TabResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_tab_proto_init() }
func file_tab_proto_init() {
	if File_tab_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_tab_proto_rawDesc), len(file_tab_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tab_proto_goTypes,
		DependencyIndexes: file_tab_proto_depIdxs,
		MessageInfos:      file_tab_proto_msgTypes,
	}.Build()
	File_tab_proto = out.File
	file_tab_proto_goTypes = nil
	file_tab_proto_depIdxs = nil
}
