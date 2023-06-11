//
//Copyright IBM Corp. All Rights Reserved.
//
//SPDX-License-Identifier: Apache-2.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.4
// source: simplewal/simplewal.proto

package simplewal

import (
	eventpb "github.com/filecoin-project/mir/pkg/pb/eventpb"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type WALEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RetentionIndex uint64         `protobuf:"varint,1,opt,name=retention_index,json=retentionIndex,proto3" json:"retention_index,omitempty"`
	Event          *eventpb.Event `protobuf:"bytes,2,opt,name=event,proto3" json:"event,omitempty"`
}

func (x *WALEntry) Reset() {
	*x = WALEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_simplewal_simplewal_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WALEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WALEntry) ProtoMessage() {}

func (x *WALEntry) ProtoReflect() protoreflect.Message {
	mi := &file_simplewal_simplewal_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WALEntry.ProtoReflect.Descriptor instead.
func (*WALEntry) Descriptor() ([]byte, []int) {
	return file_simplewal_simplewal_proto_rawDescGZIP(), []int{0}
}

func (x *WALEntry) GetRetentionIndex() uint64 {
	if x != nil {
		return x.RetentionIndex
	}
	return 0
}

func (x *WALEntry) GetEvent() *eventpb.Event {
	if x != nil {
		return x.Event
	}
	return nil
}

var File_simplewal_simplewal_proto protoreflect.FileDescriptor

var file_simplewal_simplewal_proto_rawDesc = []byte{
	0x0a, 0x19, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x77, 0x61, 0x6c, 0x2f, 0x73, 0x69, 0x6d, 0x70,
	0x6c, 0x65, 0x77, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x69, 0x6d,
	0x70, 0x6c, 0x65, 0x77, 0x61, 0x6c, 0x1a, 0x15, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x2f,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x59, 0x0a,
	0x08, 0x57, 0x41, 0x4c, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x27, 0x0a, 0x0f, 0x72, 0x65, 0x74,
	0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0e, 0x72, 0x65, 0x74, 0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x12, 0x24, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0e, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x70, 0x62, 0x2e, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x69, 0x6c, 0x65, 0x63, 0x6f, 0x69, 0x6e, 0x2d,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x2f, 0x6d, 0x69, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x77, 0x61, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_simplewal_simplewal_proto_rawDescOnce sync.Once
	file_simplewal_simplewal_proto_rawDescData = file_simplewal_simplewal_proto_rawDesc
)

func file_simplewal_simplewal_proto_rawDescGZIP() []byte {
	file_simplewal_simplewal_proto_rawDescOnce.Do(func() {
		file_simplewal_simplewal_proto_rawDescData = protoimpl.X.CompressGZIP(file_simplewal_simplewal_proto_rawDescData)
	})
	return file_simplewal_simplewal_proto_rawDescData
}

var file_simplewal_simplewal_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_simplewal_simplewal_proto_goTypes = []interface{}{
	(*WALEntry)(nil),      // 0: simplewal.WALEntry
	(*eventpb.Event)(nil), // 1: eventpb.Event
}
var file_simplewal_simplewal_proto_depIdxs = []int32{
	1, // 0: simplewal.WALEntry.event:type_name -> eventpb.Event
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_simplewal_simplewal_proto_init() }
func file_simplewal_simplewal_proto_init() {
	if File_simplewal_simplewal_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_simplewal_simplewal_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WALEntry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_simplewal_simplewal_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_simplewal_simplewal_proto_goTypes,
		DependencyIndexes: file_simplewal_simplewal_proto_depIdxs,
		MessageInfos:      file_simplewal_simplewal_proto_msgTypes,
	}.Build()
	File_simplewal_simplewal_proto = out.File
	file_simplewal_simplewal_proto_rawDesc = nil
	file_simplewal_simplewal_proto_goTypes = nil
	file_simplewal_simplewal_proto_depIdxs = nil
}
