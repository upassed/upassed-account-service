// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: teacher.proto

package client

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

type TeacherCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FirstName   string `protobuf:"bytes,1,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName    string `protobuf:"bytes,2,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	MiddleName  string `protobuf:"bytes,3,opt,name=middle_name,json=middleName,proto3" json:"middle_name,omitempty"`
	ReportEmail string `protobuf:"bytes,4,opt,name=report_email,json=reportEmail,proto3" json:"report_email,omitempty"`
	Username    string `protobuf:"bytes,5,opt,name=username,proto3" json:"username,omitempty"`
}

func (x *TeacherCreateRequest) Reset() {
	*x = TeacherCreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_teacher_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TeacherCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TeacherCreateRequest) ProtoMessage() {}

func (x *TeacherCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_teacher_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TeacherCreateRequest.ProtoReflect.Descriptor instead.
func (*TeacherCreateRequest) Descriptor() ([]byte, []int) {
	return file_teacher_proto_rawDescGZIP(), []int{0}
}

func (x *TeacherCreateRequest) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *TeacherCreateRequest) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *TeacherCreateRequest) GetMiddleName() string {
	if x != nil {
		return x.MiddleName
	}
	return ""
}

func (x *TeacherCreateRequest) GetReportEmail() string {
	if x != nil {
		return x.ReportEmail
	}
	return ""
}

func (x *TeacherCreateRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

type TeacherCreateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TeacherId string `protobuf:"bytes,4,opt,name=teacher_id,json=teacherId,proto3" json:"teacher_id,omitempty"`
}

func (x *TeacherCreateResponse) Reset() {
	*x = TeacherCreateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_teacher_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TeacherCreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TeacherCreateResponse) ProtoMessage() {}

func (x *TeacherCreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_teacher_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TeacherCreateResponse.ProtoReflect.Descriptor instead.
func (*TeacherCreateResponse) Descriptor() ([]byte, []int) {
	return file_teacher_proto_rawDescGZIP(), []int{1}
}

func (x *TeacherCreateResponse) GetTeacherId() string {
	if x != nil {
		return x.TeacherId
	}
	return ""
}

var File_teacher_proto protoreflect.FileDescriptor

var file_teacher_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x74, 0x65, 0x61, 0x63, 0x68, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x03, 0x61, 0x70, 0x69, 0x22, 0xb2, 0x01, 0x0a, 0x14, 0x54, 0x65, 0x61, 0x63, 0x68, 0x65, 0x72,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a,
	0x0a, 0x66, 0x69, 0x72, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09,
	0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x6d, 0x69, 0x64,
	0x64, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a,
	0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x36, 0x0a, 0x15, 0x54, 0x65, 0x61,
	0x63, 0x68, 0x65, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x65, 0x61, 0x63, 0x68, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x65, 0x61, 0x63, 0x68, 0x65, 0x72, 0x49,
	0x64, 0x32, 0x4a, 0x0a, 0x07, 0x54, 0x65, 0x61, 0x63, 0x68, 0x65, 0x72, 0x12, 0x3f, 0x0a, 0x06,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x54, 0x65, 0x61,
	0x63, 0x68, 0x65, 0x72, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x54, 0x65, 0x61, 0x63, 0x68, 0x65, 0x72, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1b, 0x5a,
	0x19, 0x75, 0x70, 0x61, 0x73, 0x73, 0x65, 0x64, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x2e, 0x76, 0x31, 0x3b, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_teacher_proto_rawDescOnce sync.Once
	file_teacher_proto_rawDescData = file_teacher_proto_rawDesc
)

func file_teacher_proto_rawDescGZIP() []byte {
	file_teacher_proto_rawDescOnce.Do(func() {
		file_teacher_proto_rawDescData = protoimpl.X.CompressGZIP(file_teacher_proto_rawDescData)
	})
	return file_teacher_proto_rawDescData
}

var file_teacher_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_teacher_proto_goTypes = []any{
	(*TeacherCreateRequest)(nil),  // 0: api.TeacherCreateRequest
	(*TeacherCreateResponse)(nil), // 1: api.TeacherCreateResponse
}
var file_teacher_proto_depIdxs = []int32{
	0, // 0: api.Teacher.Create:input_type -> api.TeacherCreateRequest
	1, // 1: api.Teacher.Create:output_type -> api.TeacherCreateResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_teacher_proto_init() }
func file_teacher_proto_init() {
	if File_teacher_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_teacher_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*TeacherCreateRequest); i {
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
		file_teacher_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*TeacherCreateResponse); i {
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
			RawDescriptor: file_teacher_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_teacher_proto_goTypes,
		DependencyIndexes: file_teacher_proto_depIdxs,
		MessageInfos:      file_teacher_proto_msgTypes,
	}.Build()
	File_teacher_proto = out.File
	file_teacher_proto_rawDesc = nil
	file_teacher_proto_goTypes = nil
	file_teacher_proto_depIdxs = nil
}
