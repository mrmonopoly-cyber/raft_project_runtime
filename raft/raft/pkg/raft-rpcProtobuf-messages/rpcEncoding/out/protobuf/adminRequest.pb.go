// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v4.25.3
// source: adminCli/adminRequest.proto

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

type AdminOp int32

const (
	AdminOp_CHANGE_CONF_NEW AdminOp = 0
	AdminOp_CHANGE_CONF_ADD AdminOp = 1
	AdminOp_CHANGE_CONF_REM AdminOp = 2
)

// Enum value maps for AdminOp.
var (
	AdminOp_name = map[int32]string{
		0: "CHANGE_CONF_NEW",
		1: "CHANGE_CONF_ADD",
		2: "CHANGE_CONF_REM",
	}
	AdminOp_value = map[string]int32{
		"CHANGE_CONF_NEW": 0,
		"CHANGE_CONF_ADD": 1,
		"CHANGE_CONF_REM": 2,
	}
)

func (x AdminOp) Enum() *AdminOp {
	p := new(AdminOp)
	*p = x
	return p
}

func (x AdminOp) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AdminOp) Descriptor() protoreflect.EnumDescriptor {
	return file_adminCli_adminRequest_proto_enumTypes[0].Descriptor()
}

func (AdminOp) Type() protoreflect.EnumType {
	return &file_adminCli_adminRequest_proto_enumTypes[0]
}

func (x AdminOp) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AdminOp.Descriptor instead.
func (AdminOp) EnumDescriptor() ([]byte, []int) {
	return file_adminCli_adminRequest_proto_rawDescGZIP(), []int{0}
}

type ClusterConf struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Conf   []string `protobuf:"bytes,1,rep,name=conf,proto3" json:"conf,omitempty"`
	Leader *string  `protobuf:"bytes,2,opt,name=leader,proto3,oneof" json:"leader,omitempty"`
}

func (x *ClusterConf) Reset() {
	*x = ClusterConf{}
	if protoimpl.UnsafeEnabled {
		mi := &file_adminCli_adminRequest_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClusterConf) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClusterConf) ProtoMessage() {}

func (x *ClusterConf) ProtoReflect() protoreflect.Message {
	mi := &file_adminCli_adminRequest_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClusterConf.ProtoReflect.Descriptor instead.
func (*ClusterConf) Descriptor() ([]byte, []int) {
	return file_adminCli_adminRequest_proto_rawDescGZIP(), []int{0}
}

func (x *ClusterConf) GetConf() []string {
	if x != nil {
		return x.Conf
	}
	return nil
}

func (x *ClusterConf) GetLeader() string {
	if x != nil && x.Leader != nil {
		return *x.Leader
	}
	return ""
}

type ChangeConfReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Op   AdminOp      `protobuf:"varint,1,opt,name=op,proto3,enum=protobuf.AdminOp" json:"op,omitempty"`
	Conf *ClusterConf `protobuf:"bytes,2,opt,name=conf,proto3" json:"conf,omitempty"`
}

func (x *ChangeConfReq) Reset() {
	*x = ChangeConfReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_adminCli_adminRequest_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangeConfReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangeConfReq) ProtoMessage() {}

func (x *ChangeConfReq) ProtoReflect() protoreflect.Message {
	mi := &file_adminCli_adminRequest_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangeConfReq.ProtoReflect.Descriptor instead.
func (*ChangeConfReq) Descriptor() ([]byte, []int) {
	return file_adminCli_adminRequest_proto_rawDescGZIP(), []int{1}
}

func (x *ChangeConfReq) GetOp() AdminOp {
	if x != nil {
		return x.Op
	}
	return AdminOp_CHANGE_CONF_NEW
}

func (x *ChangeConfReq) GetConf() *ClusterConf {
	if x != nil {
		return x.Conf
	}
	return nil
}

type InfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReqType   AdminOp `protobuf:"varint,1,opt,name=reqType,proto3,enum=protobuf.AdminOp" json:"reqType,omitempty"`
	Timestamp string  `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *InfoRequest) Reset() {
	*x = InfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_adminCli_adminRequest_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InfoRequest) ProtoMessage() {}

func (x *InfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_adminCli_adminRequest_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InfoRequest.ProtoReflect.Descriptor instead.
func (*InfoRequest) Descriptor() ([]byte, []int) {
	return file_adminCli_adminRequest_proto_rawDescGZIP(), []int{2}
}

func (x *InfoRequest) GetReqType() AdminOp {
	if x != nil {
		return x.ReqType
	}
	return AdminOp_CHANGE_CONF_NEW
}

func (x *InfoRequest) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

type InfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata *InfoRequest `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Payload  []byte       `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *InfoResponse) Reset() {
	*x = InfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_adminCli_adminRequest_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InfoResponse) ProtoMessage() {}

func (x *InfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_adminCli_adminRequest_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InfoResponse.ProtoReflect.Descriptor instead.
func (*InfoResponse) Descriptor() ([]byte, []int) {
	return file_adminCli_adminRequest_proto_rawDescGZIP(), []int{3}
}

func (x *InfoResponse) GetMetadata() *InfoRequest {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *InfoResponse) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

type Redirect struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Leaderip string `protobuf:"bytes,1,opt,name=leaderip,proto3" json:"leaderip,omitempty"`
}

func (x *Redirect) Reset() {
	*x = Redirect{}
	if protoimpl.UnsafeEnabled {
		mi := &file_adminCli_adminRequest_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Redirect) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Redirect) ProtoMessage() {}

func (x *Redirect) ProtoReflect() protoreflect.Message {
	mi := &file_adminCli_adminRequest_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Redirect.ProtoReflect.Descriptor instead.
func (*Redirect) Descriptor() ([]byte, []int) {
	return file_adminCli_adminRequest_proto_rawDescGZIP(), []int{4}
}

func (x *Redirect) GetLeaderip() string {
	if x != nil {
		return x.Leaderip
	}
	return ""
}

var File_adminCli_adminRequest_proto protoreflect.FileDescriptor

var file_adminCli_adminRequest_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x43, 0x6c, 0x69, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0x49, 0x0a, 0x0b, 0x63, 0x6c, 0x75, 0x73, 0x74,
	0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x6e, 0x66, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x6e, 0x66, 0x12, 0x1b, 0x0a, 0x06, 0x6c, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x06, 0x6c, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x88, 0x01, 0x01, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x6c, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x22, 0x5d, 0x0a, 0x0d, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x43, 0x6f, 0x6e, 0x66,
	0x52, 0x65, 0x71, 0x12, 0x21, 0x0a, 0x02, 0x6f, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x64, 0x6d, 0x69, 0x6e,
	0x4f, 0x70, 0x52, 0x02, 0x6f, 0x70, 0x12, 0x29, 0x0a, 0x04, 0x63, 0x6f, 0x6e, 0x66, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x52, 0x04, 0x63, 0x6f, 0x6e,
	0x66, 0x22, 0x58, 0x0a, 0x0b, 0x69, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x2b, 0x0a, 0x07, 0x72, 0x65, 0x71, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x64, 0x6d,
	0x69, 0x6e, 0x4f, 0x70, 0x52, 0x07, 0x72, 0x65, 0x71, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x5b, 0x0a, 0x0c, 0x69,
	0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x08, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x69, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x18,
	0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x22, 0x26, 0x0a, 0x08, 0x72, 0x65, 0x64, 0x69,
	0x72, 0x65, 0x63, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x69, 0x70,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x69, 0x70,
	0x2a, 0x48, 0x0a, 0x07, 0x41, 0x64, 0x6d, 0x69, 0x6e, 0x4f, 0x70, 0x12, 0x13, 0x0a, 0x0f, 0x43,
	0x48, 0x41, 0x4e, 0x47, 0x45, 0x5f, 0x43, 0x4f, 0x4e, 0x46, 0x5f, 0x4e, 0x45, 0x57, 0x10, 0x00,
	0x12, 0x13, 0x0a, 0x0f, 0x43, 0x48, 0x41, 0x4e, 0x47, 0x45, 0x5f, 0x43, 0x4f, 0x4e, 0x46, 0x5f,
	0x41, 0x44, 0x44, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x43, 0x48, 0x41, 0x4e, 0x47, 0x45, 0x5f,
	0x43, 0x4f, 0x4e, 0x46, 0x5f, 0x52, 0x45, 0x4d, 0x10, 0x02, 0x42, 0x0b, 0x5a, 0x09, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_adminCli_adminRequest_proto_rawDescOnce sync.Once
	file_adminCli_adminRequest_proto_rawDescData = file_adminCli_adminRequest_proto_rawDesc
)

func file_adminCli_adminRequest_proto_rawDescGZIP() []byte {
	file_adminCli_adminRequest_proto_rawDescOnce.Do(func() {
		file_adminCli_adminRequest_proto_rawDescData = protoimpl.X.CompressGZIP(file_adminCli_adminRequest_proto_rawDescData)
	})
	return file_adminCli_adminRequest_proto_rawDescData
}

var file_adminCli_adminRequest_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_adminCli_adminRequest_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_adminCli_adminRequest_proto_goTypes = []interface{}{
	(AdminOp)(0),          // 0: protobuf.AdminOp
	(*ClusterConf)(nil),   // 1: protobuf.clusterConf
	(*ChangeConfReq)(nil), // 2: protobuf.changeConfReq
	(*InfoRequest)(nil),   // 3: protobuf.infoRequest
	(*InfoResponse)(nil),  // 4: protobuf.infoResponse
	(*Redirect)(nil),      // 5: protobuf.redirect
}
var file_adminCli_adminRequest_proto_depIdxs = []int32{
	0, // 0: protobuf.changeConfReq.op:type_name -> protobuf.AdminOp
	1, // 1: protobuf.changeConfReq.conf:type_name -> protobuf.clusterConf
	0, // 2: protobuf.infoRequest.reqType:type_name -> protobuf.AdminOp
	3, // 3: protobuf.infoResponse.metadata:type_name -> protobuf.infoRequest
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_adminCli_adminRequest_proto_init() }
func file_adminCli_adminRequest_proto_init() {
	if File_adminCli_adminRequest_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_adminCli_adminRequest_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClusterConf); i {
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
		file_adminCli_adminRequest_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangeConfReq); i {
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
		file_adminCli_adminRequest_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InfoRequest); i {
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
		file_adminCli_adminRequest_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InfoResponse); i {
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
		file_adminCli_adminRequest_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Redirect); i {
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
	file_adminCli_adminRequest_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_adminCli_adminRequest_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_adminCli_adminRequest_proto_goTypes,
		DependencyIndexes: file_adminCli_adminRequest_proto_depIdxs,
		EnumInfos:         file_adminCli_adminRequest_proto_enumTypes,
		MessageInfos:      file_adminCli_adminRequest_proto_msgTypes,
	}.Build()
	File_adminCli_adminRequest_proto = out.File
	file_adminCli_adminRequest_proto_rawDesc = nil
	file_adminCli_adminRequest_proto_goTypes = nil
	file_adminCli_adminRequest_proto_depIdxs = nil
}
