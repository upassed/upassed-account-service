// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: group.proto

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
	Group_FindStudentsInGroup_FullMethodName = "/api.Group/FindStudentsInGroup"
	Group_FindByID_FullMethodName            = "/api.Group/FindByID"
	Group_SearchByFilter_FullMethodName      = "/api.Group/SearchByFilter"
)

// GroupClient is the client API for Group service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GroupClient interface {
	FindStudentsInGroup(ctx context.Context, in *FindStudentsInGroupRequest, opts ...grpc.CallOption) (*FindStudentsInGroupResponse, error)
	FindByID(ctx context.Context, in *GroupFindByIDRequest, opts ...grpc.CallOption) (*GroupFindByIDResponse, error)
	SearchByFilter(ctx context.Context, in *GroupSearchByFilterRequest, opts ...grpc.CallOption) (*GroupSearchByFilterResponse, error)
}

type groupClient struct {
	cc grpc.ClientConnInterface
}

func NewGroupClient(cc grpc.ClientConnInterface) GroupClient {
	return &groupClient{cc}
}

func (c *groupClient) FindStudentsInGroup(ctx context.Context, in *FindStudentsInGroupRequest, opts ...grpc.CallOption) (*FindStudentsInGroupResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(FindStudentsInGroupResponse)
	err := c.cc.Invoke(ctx, Group_FindStudentsInGroup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupClient) FindByID(ctx context.Context, in *GroupFindByIDRequest, opts ...grpc.CallOption) (*GroupFindByIDResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GroupFindByIDResponse)
	err := c.cc.Invoke(ctx, Group_FindByID_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupClient) SearchByFilter(ctx context.Context, in *GroupSearchByFilterRequest, opts ...grpc.CallOption) (*GroupSearchByFilterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GroupSearchByFilterResponse)
	err := c.cc.Invoke(ctx, Group_SearchByFilter_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GroupServer is the server API for Group service.
// All implementations must embed UnimplementedGroupServer
// for forward compatibility.
type GroupServer interface {
	FindStudentsInGroup(context.Context, *FindStudentsInGroupRequest) (*FindStudentsInGroupResponse, error)
	FindByID(context.Context, *GroupFindByIDRequest) (*GroupFindByIDResponse, error)
	SearchByFilter(context.Context, *GroupSearchByFilterRequest) (*GroupSearchByFilterResponse, error)
	mustEmbedUnimplementedGroupServer()
}

// UnimplementedGroupServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGroupServer struct{}

func (UnimplementedGroupServer) FindStudentsInGroup(context.Context, *FindStudentsInGroupRequest) (*FindStudentsInGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindStudentsInGroup not implemented")
}
func (UnimplementedGroupServer) FindByID(context.Context, *GroupFindByIDRequest) (*GroupFindByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindByID not implemented")
}
func (UnimplementedGroupServer) SearchByFilter(context.Context, *GroupSearchByFilterRequest) (*GroupSearchByFilterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchByFilter not implemented")
}
func (UnimplementedGroupServer) mustEmbedUnimplementedGroupServer() {}
func (UnimplementedGroupServer) testEmbeddedByValue()               {}

// UnsafeGroupServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GroupServer will
// result in compilation errors.
type UnsafeGroupServer interface {
	mustEmbedUnimplementedGroupServer()
}

func RegisterGroupServer(s grpc.ServiceRegistrar, srv GroupServer) {
	// If the following call pancis, it indicates UnimplementedGroupServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Group_ServiceDesc, srv)
}

func _Group_FindStudentsInGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindStudentsInGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServer).FindStudentsInGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Group_FindStudentsInGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServer).FindStudentsInGroup(ctx, req.(*FindStudentsInGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Group_FindByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GroupFindByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServer).FindByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Group_FindByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServer).FindByID(ctx, req.(*GroupFindByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Group_SearchByFilter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GroupSearchByFilterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServer).SearchByFilter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Group_SearchByFilter_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServer).SearchByFilter(ctx, req.(*GroupSearchByFilterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Group_ServiceDesc is the grpc.ServiceDesc for Group service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Group_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Group",
	HandlerType: (*GroupServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindStudentsInGroup",
			Handler:    _Group_FindStudentsInGroup_Handler,
		},
		{
			MethodName: "FindByID",
			Handler:    _Group_FindByID_Handler,
		},
		{
			MethodName: "SearchByFilter",
			Handler:    _Group_SearchByFilter_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "group.proto",
}
