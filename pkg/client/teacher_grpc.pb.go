// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: teacher.proto

package client

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Teacher_FindByID_FullMethodName = "/api.Teacher/FindByID"
)

// TeacherClient is the client API for Teacher service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TeacherClient interface {
	FindByID(ctx context.Context, in *TeacherFindByIDRequest, opts ...grpc.CallOption) (*TeacherFindByIDResponse, error)
}

type teacherClient struct {
	cc grpc.ClientConnInterface
}

func NewTeacherClient(cc grpc.ClientConnInterface) TeacherClient {
	return &teacherClient{cc}
}

func (c *teacherClient) FindByID(ctx context.Context, in *TeacherFindByIDRequest, opts ...grpc.CallOption) (*TeacherFindByIDResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TeacherFindByIDResponse)
	err := c.cc.Invoke(ctx, Teacher_FindByID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TeacherServer is the server API for Teacher service.
// All implementations must embed UnimplementedTeacherServer
// for forward compatibility.
type TeacherServer interface {
	FindByID(context.Context, *TeacherFindByIDRequest) (*TeacherFindByIDResponse, error)
	mustEmbedUnimplementedTeacherServer()
}

// UnimplementedTeacherServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTeacherServer struct{}

func (UnimplementedTeacherServer) FindByID(context.Context, *TeacherFindByIDRequest) (*TeacherFindByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByID not implemented")
}
func (UnimplementedTeacherServer) mustEmbedUnimplementedTeacherServer() {}
func (UnimplementedTeacherServer) testEmbeddedByValue()                 {}

// UnsafeTeacherServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TeacherServer will
// result in compilation errors.
type UnsafeTeacherServer interface {
	mustEmbedUnimplementedTeacherServer()
}

func RegisterTeacherServer(s grpc.ServiceRegistrar, srv TeacherServer) {
	// If the following call pancis, it indicates UnimplementedTeacherServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Teacher_ServiceDesc, srv)
}

func _Teacher_FindByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TeacherFindByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TeacherServer).FindByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Teacher_FindByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TeacherServer).FindByID(ctx, req.(*TeacherFindByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Teacher_ServiceDesc is the grpc.ServiceDesc for Teacher service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Teacher_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Teacher",
	HandlerType: (*TeacherServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindByID",
			Handler:    _Teacher_FindByID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "teacher.proto",
}
