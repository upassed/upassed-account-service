// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: group.proto

package client

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type FindStudentsInGroupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupId string `protobuf:"bytes,1,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
}

func (x *FindStudentsInGroupRequest) Reset() {
	*x = FindStudentsInGroupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindStudentsInGroupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindStudentsInGroupRequest) ProtoMessage() {}

func (x *FindStudentsInGroupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_group_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindStudentsInGroupRequest.ProtoReflect.Descriptor instead.
func (*FindStudentsInGroupRequest) Descriptor() ([]byte, []int) {
	return file_group_proto_rawDescGZIP(), []int{0}
}

func (x *FindStudentsInGroupRequest) GetGroupId() string {
	if x != nil {
		return x.GroupId
	}
	return ""
}

type FindStudentsInGroupResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StudentsInGroup []*StudentDTO `protobuf:"bytes,1,rep,name=students_in_group,json=studentsInGroup,proto3" json:"students_in_group,omitempty"`
}

func (x *FindStudentsInGroupResponse) Reset() {
	*x = FindStudentsInGroupResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_group_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindStudentsInGroupResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindStudentsInGroupResponse) ProtoMessage() {}

func (x *FindStudentsInGroupResponse) ProtoReflect() protoreflect.Message {
	mi := &file_group_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindStudentsInGroupResponse.ProtoReflect.Descriptor instead.
func (*FindStudentsInGroupResponse) Descriptor() ([]byte, []int) {
	return file_group_proto_rawDescGZIP(), []int{1}
}

func (x *FindStudentsInGroupResponse) GetStudentsInGroup() []*StudentDTO {
	if x != nil {
		return x.StudentsInGroup
	}
	return nil
}

var File_group_proto protoreflect.FileDescriptor

var file_group_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61,
	0x70, 0x69, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0d, 0x73, 0x74, 0x75,
	0x64, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x41, 0x0a, 0x1a, 0x46, 0x69,
	0x6e, 0x64, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x73, 0x49, 0x6e, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x08, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72,
	0x03, 0xb0, 0x01, 0x01, 0x52, 0x07, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x22, 0x5a, 0x0a,
	0x1b, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x73, 0x49, 0x6e, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x11,
	0x73, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x69, 0x6e, 0x5f, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x53, 0x74,
	0x75, 0x64, 0x65, 0x6e, 0x74, 0x44, 0x54, 0x4f, 0x52, 0x0f, 0x73, 0x74, 0x75, 0x64, 0x65, 0x6e,
	0x74, 0x73, 0x49, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x32, 0x61, 0x0a, 0x05, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x12, 0x58, 0x0a, 0x13, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e,
	0x74, 0x73, 0x49, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x1f, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x46, 0x69, 0x6e, 0x64, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x73, 0x49, 0x6e, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x46, 0x69, 0x6e, 0x64, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x73, 0x49, 0x6e, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1b, 0x5a, 0x19,
	0x75, 0x70, 0x61, 0x73, 0x73, 0x65, 0x64, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e,
	0x76, 0x31, 0x3b, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_group_proto_rawDescOnce sync.Once
	file_group_proto_rawDescData = file_group_proto_rawDesc
)

func file_group_proto_rawDescGZIP() []byte {
	file_group_proto_rawDescOnce.Do(func() {
		file_group_proto_rawDescData = protoimpl.X.CompressGZIP(file_group_proto_rawDescData)
	})
	return file_group_proto_rawDescData
}

var file_group_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_group_proto_goTypes = []any{
	(*FindStudentsInGroupRequest)(nil),  // 0: api.FindStudentsInGroupRequest
	(*FindStudentsInGroupResponse)(nil), // 1: api.FindStudentsInGroupResponse
	(*StudentDTO)(nil),                  // 2: api.StudentDTO
}
var file_group_proto_depIdxs = []int32{
	2, // 0: api.FindStudentsInGroupResponse.students_in_group:type_name -> api.StudentDTO
	0, // 1: api.Group.FindStudentsInGroup:input_type -> api.FindStudentsInGroupRequest
	1, // 2: api.Group.FindStudentsInGroup:output_type -> api.FindStudentsInGroupResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_group_proto_init() }
func file_group_proto_init() {
	if File_group_proto != nil {
		return
	}
	file_student_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_group_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*FindStudentsInGroupRequest); i {
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
		file_group_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*FindStudentsInGroupResponse); i {
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
			RawDescriptor: file_group_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_group_proto_goTypes,
		DependencyIndexes: file_group_proto_depIdxs,
		MessageInfos:      file_group_proto_msgTypes,
	}.Build()
	File_group_proto = out.File
	file_group_proto_rawDesc = nil
	file_group_proto_goTypes = nil
	file_group_proto_depIdxs = nil
}
