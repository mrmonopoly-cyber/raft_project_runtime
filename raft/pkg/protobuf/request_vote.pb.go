// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.2
// source: request_vote.proto

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

type RequestVote struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Term         *uint64 `protobuf:"varint,1,req,name=Term" json:"Term,omitempty"`
	CandidateId  *string `protobuf:"bytes,2,req,name=CandidateId" json:"CandidateId,omitempty"`
	LastLogIndex *uint64 `protobuf:"varint,3,req,name=LastLogIndex" json:"LastLogIndex,omitempty"`
	LastLogTerm  *uint64 `protobuf:"varint,4,req,name=LastLogTerm" json:"LastLogTerm,omitempty"`
}

func (x *RequestVote) Reset() {
	*x = RequestVote{}
	if protoimpl.UnsafeEnabled {
		mi := &file_request_vote_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestVote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestVote) ProtoMessage() {}

func (x *RequestVote) ProtoReflect() protoreflect.Message {
	mi := &file_request_vote_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestVote.ProtoReflect.Descriptor instead.
func (*RequestVote) Descriptor() ([]byte, []int) {
	return file_request_vote_proto_rawDescGZIP(), []int{0}
}

func (x *RequestVote) GetTerm() uint64 {
	if x != nil && x.Term != nil {
		return *x.Term
	}
	return 0
}

func (x *RequestVote) GetCandidateId() string {
	if x != nil && x.CandidateId != nil {
		return *x.CandidateId
	}
	return ""
}

func (x *RequestVote) GetLastLogIndex() uint64 {
	if x != nil && x.LastLogIndex != nil {
		return *x.LastLogIndex
	}
	return 0
}

func (x *RequestVote) GetLastLogTerm() uint64 {
	if x != nil && x.LastLogTerm != nil {
		return *x.LastLogTerm
	}
	return 0
}

var File_request_vote_proto protoreflect.FileDescriptor

var file_request_vote_proto_rawDesc = []byte{
	0x0a, 0x12, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x76, 0x6f, 0x74, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0x89,
	0x01, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x02, 0x28, 0x04, 0x52, 0x04, 0x54, 0x65,
	0x72, 0x6d, 0x12, 0x20, 0x0a, 0x0b, 0x43, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x02, 0x28, 0x09, 0x52, 0x0b, 0x43, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c, 0x4c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x02, 0x28, 0x04, 0x52, 0x0c, 0x4c, 0x61, 0x73, 0x74,
	0x4c, 0x6f, 0x67, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x20, 0x0a, 0x0b, 0x4c, 0x61, 0x73, 0x74,
	0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x18, 0x04, 0x20, 0x02, 0x28, 0x04, 0x52, 0x0b, 0x4c,
	0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x54, 0x65, 0x72, 0x6d, 0x42, 0x0b, 0x5a, 0x09, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
}

var (
	file_request_vote_proto_rawDescOnce sync.Once
	file_request_vote_proto_rawDescData = file_request_vote_proto_rawDesc
)

func file_request_vote_proto_rawDescGZIP() []byte {
	file_request_vote_proto_rawDescOnce.Do(func() {
		file_request_vote_proto_rawDescData = protoimpl.X.CompressGZIP(file_request_vote_proto_rawDescData)
	})
	return file_request_vote_proto_rawDescData
}

var file_request_vote_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_request_vote_proto_goTypes = []interface{}{
	(*RequestVote)(nil), // 0: protobuf.RequestVote
}
var file_request_vote_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_request_vote_proto_init() }
func file_request_vote_proto_init() {
	if File_request_vote_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_request_vote_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestVote); i {
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
			RawDescriptor: file_request_vote_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_request_vote_proto_goTypes,
		DependencyIndexes: file_request_vote_proto_depIdxs,
		MessageInfos:      file_request_vote_proto_msgTypes,
	}.Build()
	File_request_vote_proto = out.File
	file_request_vote_proto_rawDesc = nil
	file_request_vote_proto_goTypes = nil
	file_request_vote_proto_depIdxs = nil
}
