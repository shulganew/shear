// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.12.4
// source: internal/grpcs/proto/short.proto

package proto

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

// GetURL messages.
type GetURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Brief string `protobuf:"bytes,1,opt,name=brief,proto3" json:"brief,omitempty"`
}

func (x *GetURLRequest) Reset() {
	*x = GetURLRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_grpcs_proto_short_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetURLRequest) ProtoMessage() {}

func (x *GetURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_grpcs_proto_short_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetURLRequest.ProtoReflect.Descriptor instead.
func (*GetURLRequest) Descriptor() ([]byte, []int) {
	return file_internal_grpcs_proto_short_proto_rawDescGZIP(), []int{0}
}

func (x *GetURLRequest) GetBrief() string {
	if x != nil {
		return x.Brief
	}
	return ""
}

type GetURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Origin string `protobuf:"bytes,1,opt,name=origin,proto3" json:"origin,omitempty"`
}

func (x *GetURLResponse) Reset() {
	*x = GetURLResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_grpcs_proto_short_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetURLResponse) ProtoMessage() {}

func (x *GetURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_grpcs_proto_short_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetURLResponse.ProtoReflect.Descriptor instead.
func (*GetURLResponse) Descriptor() ([]byte, []int) {
	return file_internal_grpcs_proto_short_proto_rawDescGZIP(), []int{1}
}

func (x *GetURLResponse) GetOrigin() string {
	if x != nil {
		return x.Origin
	}
	return ""
}

// SetURL messages.
type SetURLRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Origin string `protobuf:"bytes,1,opt,name=origin,proto3" json:"origin,omitempty"`
}

func (x *SetURLRequest) Reset() {
	*x = SetURLRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_grpcs_proto_short_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetURLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetURLRequest) ProtoMessage() {}

func (x *SetURLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_grpcs_proto_short_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetURLRequest.ProtoReflect.Descriptor instead.
func (*SetURLRequest) Descriptor() ([]byte, []int) {
	return file_internal_grpcs_proto_short_proto_rawDescGZIP(), []int{2}
}

func (x *SetURLRequest) GetOrigin() string {
	if x != nil {
		return x.Origin
	}
	return ""
}

type SetURLResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Brief string `protobuf:"bytes,1,opt,name=brief,proto3" json:"brief,omitempty"`
}

func (x *SetURLResponse) Reset() {
	*x = SetURLResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_grpcs_proto_short_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetURLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetURLResponse) ProtoMessage() {}

func (x *SetURLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_grpcs_proto_short_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetURLResponse.ProtoReflect.Descriptor instead.
func (*SetURLResponse) Descriptor() ([]byte, []int) {
	return file_internal_grpcs_proto_short_proto_rawDescGZIP(), []int{3}
}

func (x *SetURLResponse) GetBrief() string {
	if x != nil {
		return x.Brief
	}
	return ""
}

// GetUserURLs messages.
type GetURLs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetURLs) Reset() {
	*x = GetURLs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_grpcs_proto_short_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetURLs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetURLs) ProtoMessage() {}

func (x *GetURLs) ProtoReflect() protoreflect.Message {
	mi := &file_internal_grpcs_proto_short_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetURLs.ProtoReflect.Descriptor instead.
func (*GetURLs) Descriptor() ([]byte, []int) {
	return file_internal_grpcs_proto_short_proto_rawDescGZIP(), []int{4}
}

type Short struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Brief  string `protobuf:"bytes,1,opt,name=brief,proto3" json:"brief,omitempty"`
	Origin string `protobuf:"bytes,2,opt,name=origin,proto3" json:"origin,omitempty"`
}

func (x *Short) Reset() {
	*x = Short{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_grpcs_proto_short_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Short) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Short) ProtoMessage() {}

func (x *Short) ProtoReflect() protoreflect.Message {
	mi := &file_internal_grpcs_proto_short_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Short.ProtoReflect.Descriptor instead.
func (*Short) Descriptor() ([]byte, []int) {
	return file_internal_grpcs_proto_short_proto_rawDescGZIP(), []int{5}
}

func (x *Short) GetBrief() string {
	if x != nil {
		return x.Brief
	}
	return ""
}

func (x *Short) GetOrigin() string {
	if x != nil {
		return x.Origin
	}
	return ""
}

type GetUserURLsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Short []*Short `protobuf:"bytes,1,rep,name=short,proto3" json:"short,omitempty"`
}

func (x *GetUserURLsResponse) Reset() {
	*x = GetUserURLsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_grpcs_proto_short_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserURLsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserURLsResponse) ProtoMessage() {}

func (x *GetUserURLsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_grpcs_proto_short_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserURLsResponse.ProtoReflect.Descriptor instead.
func (*GetUserURLsResponse) Descriptor() ([]byte, []int) {
	return file_internal_grpcs_proto_short_proto_rawDescGZIP(), []int{6}
}

func (x *GetUserURLsResponse) GetShort() []*Short {
	if x != nil {
		return x.Short
	}
	return nil
}

var File_internal_grpcs_proto_short_proto protoreflect.FileDescriptor

var file_internal_grpcs_proto_short_proto_rawDesc = []byte{
	0x0a, 0x20, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x73,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x09, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x67, 0x72, 0x70, 0x63, 0x22, 0x25, 0x0a,
	0x0d, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14,
	0x0a, 0x05, 0x62, 0x72, 0x69, 0x65, 0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x62,
	0x72, 0x69, 0x65, 0x66, 0x22, 0x28, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x22, 0x27,
	0x0a, 0x0d, 0x53, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x22, 0x26, 0x0a, 0x0e, 0x53, 0x65, 0x74, 0x55, 0x52,
	0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x72, 0x69,
	0x65, 0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x62, 0x72, 0x69, 0x65, 0x66, 0x22,
	0x09, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x73, 0x22, 0x35, 0x0a, 0x05, 0x53, 0x68,
	0x6f, 0x72, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x72, 0x69, 0x65, 0x66, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x62, 0x72, 0x69, 0x65, 0x66, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x72, 0x69,
	0x67, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x72, 0x69, 0x67, 0x69,
	0x6e, 0x22, 0x3d, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x55, 0x52, 0x4c, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x05, 0x73, 0x68, 0x6f, 0x72,
	0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x67,
	0x72, 0x70, 0x63, 0x2e, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x52, 0x05, 0x73, 0x68, 0x6f, 0x72, 0x74,
	0x32, 0xc8, 0x01, 0x0a, 0x05, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12, 0x3d, 0x0a, 0x06, 0x47, 0x65,
	0x74, 0x55, 0x52, 0x4c, 0x12, 0x18, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19,
	0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x52,
	0x4c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a, 0x06, 0x53, 0x65, 0x74,
	0x55, 0x52, 0x4c, 0x12, 0x18, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x67, 0x72, 0x70, 0x63, 0x2e,
	0x53, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e,
	0x73, 0x68, 0x6f, 0x72, 0x74, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x53, 0x65, 0x74, 0x55, 0x52, 0x4c,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x55, 0x52, 0x4c, 0x73, 0x12, 0x12, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x67,
	0x72, 0x70, 0x63, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x52, 0x4c, 0x73, 0x1a, 0x1e, 0x2e, 0x73, 0x68,
	0x6f, 0x72, 0x74, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x55,
	0x52, 0x4c, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x11, 0x5a, 0x0f, 0x73,
	0x68, 0x6f, 0x72, 0x74, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_grpcs_proto_short_proto_rawDescOnce sync.Once
	file_internal_grpcs_proto_short_proto_rawDescData = file_internal_grpcs_proto_short_proto_rawDesc
)

func file_internal_grpcs_proto_short_proto_rawDescGZIP() []byte {
	file_internal_grpcs_proto_short_proto_rawDescOnce.Do(func() {
		file_internal_grpcs_proto_short_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_grpcs_proto_short_proto_rawDescData)
	})
	return file_internal_grpcs_proto_short_proto_rawDescData
}

var file_internal_grpcs_proto_short_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_internal_grpcs_proto_short_proto_goTypes = []interface{}{
	(*GetURLRequest)(nil),       // 0: shortgrpc.GetURLRequest
	(*GetURLResponse)(nil),      // 1: shortgrpc.GetURLResponse
	(*SetURLRequest)(nil),       // 2: shortgrpc.SetURLRequest
	(*SetURLResponse)(nil),      // 3: shortgrpc.SetURLResponse
	(*GetURLs)(nil),             // 4: shortgrpc.GetURLs
	(*Short)(nil),               // 5: shortgrpc.Short
	(*GetUserURLsResponse)(nil), // 6: shortgrpc.GetUserURLsResponse
}
var file_internal_grpcs_proto_short_proto_depIdxs = []int32{
	5, // 0: shortgrpc.GetUserURLsResponse.short:type_name -> shortgrpc.Short
	0, // 1: shortgrpc.Users.GetURL:input_type -> shortgrpc.GetURLRequest
	2, // 2: shortgrpc.Users.SetURL:input_type -> shortgrpc.SetURLRequest
	4, // 3: shortgrpc.Users.GetUserURLs:input_type -> shortgrpc.GetURLs
	1, // 4: shortgrpc.Users.GetURL:output_type -> shortgrpc.GetURLResponse
	3, // 5: shortgrpc.Users.SetURL:output_type -> shortgrpc.SetURLResponse
	6, // 6: shortgrpc.Users.GetUserURLs:output_type -> shortgrpc.GetUserURLsResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_internal_grpcs_proto_short_proto_init() }
func file_internal_grpcs_proto_short_proto_init() {
	if File_internal_grpcs_proto_short_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_grpcs_proto_short_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetURLRequest); i {
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
		file_internal_grpcs_proto_short_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetURLResponse); i {
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
		file_internal_grpcs_proto_short_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetURLRequest); i {
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
		file_internal_grpcs_proto_short_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetURLResponse); i {
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
		file_internal_grpcs_proto_short_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetURLs); i {
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
		file_internal_grpcs_proto_short_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Short); i {
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
		file_internal_grpcs_proto_short_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserURLsResponse); i {
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
			RawDescriptor: file_internal_grpcs_proto_short_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_grpcs_proto_short_proto_goTypes,
		DependencyIndexes: file_internal_grpcs_proto_short_proto_depIdxs,
		MessageInfos:      file_internal_grpcs_proto_short_proto_msgTypes,
	}.Build()
	File_internal_grpcs_proto_short_proto = out.File
	file_internal_grpcs_proto_short_proto_rawDesc = nil
	file_internal_grpcs_proto_short_proto_goTypes = nil
	file_internal_grpcs_proto_short_proto_depIdxs = nil
}
