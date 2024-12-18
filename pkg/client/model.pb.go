// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: model.proto

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

type GroupDTO struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                 string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	SpecializationCode string `protobuf:"bytes,2,opt,name=specialization_code,json=specializationCode,proto3" json:"specialization_code,omitempty"`
	GroupNumber        string `protobuf:"bytes,3,opt,name=group_number,json=groupNumber,proto3" json:"group_number,omitempty"`
}

func (x *GroupDTO) Reset() {
	*x = GroupDTO{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GroupDTO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GroupDTO) ProtoMessage() {}

func (x *GroupDTO) ProtoReflect() protoreflect.Message {
	mi := &file_model_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GroupDTO.ProtoReflect.Descriptor instead.
func (*GroupDTO) Descriptor() ([]byte, []int) {
	return file_model_proto_rawDescGZIP(), []int{0}
}

func (x *GroupDTO) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GroupDTO) GetSpecializationCode() string {
	if x != nil {
		return x.SpecializationCode
	}
	return ""
}

func (x *GroupDTO) GetGroupNumber() string {
	if x != nil {
		return x.GroupNumber
	}
	return ""
}

type StudentDTO struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id               string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	FirstName        string    `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName         string    `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	MiddleName       string    `protobuf:"bytes,4,opt,name=middle_name,json=middleName,proto3" json:"middle_name,omitempty"`
	EducationalEmail string    `protobuf:"bytes,5,opt,name=educational_email,json=educationalEmail,proto3" json:"educational_email,omitempty"`
	Username         string    `protobuf:"bytes,6,opt,name=username,proto3" json:"username,omitempty"`
	Group            *GroupDTO `protobuf:"bytes,7,opt,name=group,proto3" json:"group,omitempty"`
}

func (x *StudentDTO) Reset() {
	*x = StudentDTO{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StudentDTO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StudentDTO) ProtoMessage() {}

func (x *StudentDTO) ProtoReflect() protoreflect.Message {
	mi := &file_model_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StudentDTO.ProtoReflect.Descriptor instead.
func (*StudentDTO) Descriptor() ([]byte, []int) {
	return file_model_proto_rawDescGZIP(), []int{1}
}

func (x *StudentDTO) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *StudentDTO) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *StudentDTO) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *StudentDTO) GetMiddleName() string {
	if x != nil {
		return x.MiddleName
	}
	return ""
}

func (x *StudentDTO) GetEducationalEmail() string {
	if x != nil {
		return x.EducationalEmail
	}
	return ""
}

func (x *StudentDTO) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *StudentDTO) GetGroup() *GroupDTO {
	if x != nil {
		return x.Group
	}
	return nil
}

type TeacherDTO struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	FirstName   string `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName    string `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	MiddleName  string `protobuf:"bytes,4,opt,name=middle_name,json=middleName,proto3" json:"middle_name,omitempty"`
	ReportEmail string `protobuf:"bytes,5,opt,name=report_email,json=reportEmail,proto3" json:"report_email,omitempty"`
	Username    string `protobuf:"bytes,6,opt,name=username,proto3" json:"username,omitempty"`
}

func (x *TeacherDTO) Reset() {
	*x = TeacherDTO{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TeacherDTO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TeacherDTO) ProtoMessage() {}

func (x *TeacherDTO) ProtoReflect() protoreflect.Message {
	mi := &file_model_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TeacherDTO.ProtoReflect.Descriptor instead.
func (*TeacherDTO) Descriptor() ([]byte, []int) {
	return file_model_proto_rawDescGZIP(), []int{2}
}

func (x *TeacherDTO) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TeacherDTO) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *TeacherDTO) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *TeacherDTO) GetMiddleName() string {
	if x != nil {
		return x.MiddleName
	}
	return ""
}

func (x *TeacherDTO) GetReportEmail() string {
	if x != nil {
		return x.ReportEmail
	}
	return ""
}

func (x *TeacherDTO) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

var File_model_proto protoreflect.FileDescriptor

var file_model_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61,
	0x70, 0x69, 0x22, 0x6e, 0x0a, 0x08, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x44, 0x54, 0x4f, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2f,
	0x0a, 0x13, 0x73, 0x70, 0x65, 0x63, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x73, 0x70, 0x65,
	0x63, 0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x64, 0x65, 0x12,
	0x21, 0x0a, 0x0c, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x4e, 0x75, 0x6d, 0x62,
	0x65, 0x72, 0x22, 0xe7, 0x01, 0x0a, 0x0a, 0x53, 0x74, 0x75, 0x64, 0x65, 0x6e, 0x74, 0x44, 0x54,
	0x4f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x69, 0x72, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a,
	0x0b, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2b,
	0x0a, 0x11, 0x65, 0x64, 0x75, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x5f, 0x65, 0x6d,
	0x61, 0x69, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x65, 0x64, 0x75, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x75,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x44, 0x54, 0x4f, 0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x22, 0xb8, 0x01, 0x0a,
	0x0a, 0x54, 0x65, 0x61, 0x63, 0x68, 0x65, 0x72, 0x44, 0x54, 0x4f, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x66,
	0x69, 0x72, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x61,
	0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c,
	0x61, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x6d, 0x69, 0x64, 0x64, 0x6c,
	0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6d, 0x69,
	0x64, 0x64, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x75,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x1b, 0x5a, 0x19, 0x75, 0x70, 0x61, 0x73, 0x73,
	0x65, 0x64, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x3b, 0x63, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_model_proto_rawDescOnce sync.Once
	file_model_proto_rawDescData = file_model_proto_rawDesc
)

func file_model_proto_rawDescGZIP() []byte {
	file_model_proto_rawDescOnce.Do(func() {
		file_model_proto_rawDescData = protoimpl.X.CompressGZIP(file_model_proto_rawDescData)
	})
	return file_model_proto_rawDescData
}

var file_model_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_model_proto_goTypes = []any{
	(*GroupDTO)(nil),   // 0: api.GroupDTO
	(*StudentDTO)(nil), // 1: api.StudentDTO
	(*TeacherDTO)(nil), // 2: api.TeacherDTO
}
var file_model_proto_depIdxs = []int32{
	0, // 0: api.StudentDTO.group:type_name -> api.GroupDTO
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_model_proto_init() }
func file_model_proto_init() {
	if File_model_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_model_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*GroupDTO); i {
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
		file_model_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*StudentDTO); i {
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
		file_model_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*TeacherDTO); i {
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
			RawDescriptor: file_model_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_model_proto_goTypes,
		DependencyIndexes: file_model_proto_depIdxs,
		MessageInfos:      file_model_proto_msgTypes,
	}.Build()
	File_model_proto = out.File
	file_model_proto_rawDesc = nil
	file_model_proto_goTypes = nil
	file_model_proto_depIdxs = nil
}
