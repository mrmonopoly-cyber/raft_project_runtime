// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.2
// source: copy_state.proto

package protobuf

import (
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

type CopyState struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      *string  `protobuf:"bytes,1,req,name=Id" json:"Id,omitempty"`
	Term    *uint64  `protobuf:"varint,2,req,name=Term" json:"Term,omitempty"`
	Index   *uint64  `protobuf:"varint,3,req,name=Index" json:"Index,omitempty"`
	Voting  *bool    `protobuf:"varint,4,req,name=Voting" json:"Voting,omitempty"`
	Entries []*Entry `protobuf:"bytes,5,rep,name=Entries" json:"Entries,omitempty"`
}

func (x *CopyState) Reset() {
	*x = CopyState{}
	if protoimpl.UnsafeEnabled {
		mi := &file_copy_state_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CopyState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CopyState) ProtoMessage() {}

func (x *CopyState) ProtoReflect() protoreflect.Message {
	mi := &file_copy_state_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CopyState.ProtoReflect.Descriptor instead.
func (*CopyState) Descriptor() ([]byte, []int) {
	return file_copy_state_proto_rawDescGZIP(), []int{0}
}

func (x *CopyState) GetId() string {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return ""
}

func (x *CopyState) GetTerm() uint64 {
	if x != nil && x.Term != nil {
		return *x.Term
	}
	return 0
}

func (x *CopyState) GetIndex() uint64 {
	if x != nil && x.Index != nil {
		return *x.Index
	}
	return 0
}

func (x *CopyState) GetVoting() bool {
	if x != nil && x.Voting != nil {
		return *x.Voting
	}
	return false
}

func (x *CopyState) GetEntries() []*Entry {
	if x != nil {
		return x.Entries
	}
	return nil
}

var File_copy_state_proto protoreflect.FileDescriptor

var file_copy_state_proto_rawDesc = []byte{
	0x0a, 0x10, 0x63, 0x6f, 0x70, 0x79, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x1a, 0x0b, 0x65, 0x6e,
	0x74, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x88, 0x01, 0x0a, 0x09, 0x43, 0x6f,
	0x70, 0x79, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x02, 0x28, 0x09, 0x52, 0x02, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x18,
	0x02, 0x20, 0x02, 0x28, 0x04, 0x52, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x02, 0x28, 0x04, 0x52, 0x05, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x12, 0x16, 0x0a, 0x06, 0x56, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x04, 0x20, 0x02, 0x28,
	0x08, 0x52, 0x06, 0x56, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x29, 0x0a, 0x07, 0x45, 0x6e, 0x74,
	0x72, 0x69, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x45, 0x6e, 0x74,
	0x72, 0x69, 0x65, 0x73, 0x42, 0x0b, 0x5a, 0x09, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f,
}

var (
	file_copy_state_proto_rawDescOnce sync.Once
	file_copy_state_proto_rawDescData = file_copy_state_proto_rawDesc
)

func file_copy_state_proto_rawDescGZIP() []byte {
	file_copy_state_proto_rawDescOnce.Do(func() {
		file_copy_state_proto_rawDescData = protoimpl.X.CompressGZIP(file_copy_state_proto_rawDescData)
	})
	return file_copy_state_proto_rawDescData
}

var file_copy_state_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_copy_state_proto_goTypes = []interface{}{
	(*CopyState)(nil), // 0: protobuf.CopyState
	(*Entry)(nil),     // 1: protobuf.Entry
}
var file_copy_state_proto_depIdxs = []int32{
	1, // 0: protobuf.CopyState.Entries:type_name -> protobuf.Entry
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_copy_state_proto_init() }
func file_copy_state_proto_init() {
	if File_copy_state_proto != nil {
		return
	}
	file_entry_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_copy_state_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CopyState); i {
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
			RawDescriptor: file_copy_state_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_copy_state_proto_goTypes,
		DependencyIndexes: file_copy_state_proto_depIdxs,
		MessageInfos:      file_copy_state_proto_msgTypes,
	}.Build()
	File_copy_state_proto = out.File
	file_copy_state_proto_rawDesc = nil
	file_copy_state_proto_goTypes = nil
	file_copy_state_proto_depIdxs = nil
}
