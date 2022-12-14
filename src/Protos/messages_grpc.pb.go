// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: Protos/messages.proto

package Client_Server_Token_Manager

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TokenMgmtClient is the client API for TokenMgmt service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TokenMgmtClient interface {
	CreateRequest(ctx context.Context, in *CreateToken, opts ...grpc.CallOption) (*CreateResponse, error)
	DropRequest(ctx context.Context, in *DropToken, opts ...grpc.CallOption) (*DropResponse, error)
	WriteRequest(ctx context.Context, in *WriteToken, opts ...grpc.CallOption) (*WriteResponse, error)
	ReadRequest(ctx context.Context, in *ReadToken, opts ...grpc.CallOption) (*ReadResponse, error)
	CreateReplicate(ctx context.Context, in *CreateToken, opts ...grpc.CallOption) (*CreateReplicateResponse, error)
	WriteReplicate(ctx context.Context, in *WriteToken, opts ...grpc.CallOption) (*WriteReplicateResponse, error)
	ReadWrite(ctx context.Context, in *ReadWriteToken, opts ...grpc.CallOption) (*ReadWriteResponse, error)
}

type tokenMgmtClient struct {
	cc grpc.ClientConnInterface
}

func NewTokenMgmtClient(cc grpc.ClientConnInterface) TokenMgmtClient {
	return &tokenMgmtClient{cc}
}

func (c *tokenMgmtClient) CreateRequest(ctx context.Context, in *CreateToken, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, "/Protos.TokenMgmt/CreateRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenMgmtClient) DropRequest(ctx context.Context, in *DropToken, opts ...grpc.CallOption) (*DropResponse, error) {
	out := new(DropResponse)
	err := c.cc.Invoke(ctx, "/Protos.TokenMgmt/DropRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenMgmtClient) WriteRequest(ctx context.Context, in *WriteToken, opts ...grpc.CallOption) (*WriteResponse, error) {
	out := new(WriteResponse)
	err := c.cc.Invoke(ctx, "/Protos.TokenMgmt/WriteRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenMgmtClient) ReadRequest(ctx context.Context, in *ReadToken, opts ...grpc.CallOption) (*ReadResponse, error) {
	out := new(ReadResponse)
	err := c.cc.Invoke(ctx, "/Protos.TokenMgmt/ReadRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenMgmtClient) CreateReplicate(ctx context.Context, in *CreateToken, opts ...grpc.CallOption) (*CreateReplicateResponse, error) {
	out := new(CreateReplicateResponse)
	err := c.cc.Invoke(ctx, "/Protos.TokenMgmt/CreateReplicate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenMgmtClient) WriteReplicate(ctx context.Context, in *WriteToken, opts ...grpc.CallOption) (*WriteReplicateResponse, error) {
	out := new(WriteReplicateResponse)
	err := c.cc.Invoke(ctx, "/Protos.TokenMgmt/WriteReplicate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tokenMgmtClient) ReadWrite(ctx context.Context, in *ReadWriteToken, opts ...grpc.CallOption) (*ReadWriteResponse, error) {
	out := new(ReadWriteResponse)
	err := c.cc.Invoke(ctx, "/Protos.TokenMgmt/ReadWrite", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TokenMgmtServer is the server API for TokenMgmt service.
// All implementations must embed UnimplementedTokenMgmtServer
// for forward compatibility
type TokenMgmtServer interface {
	CreateRequest(context.Context, *CreateToken) (*CreateResponse, error)
	DropRequest(context.Context, *DropToken) (*DropResponse, error)
	WriteRequest(context.Context, *WriteToken) (*WriteResponse, error)
	ReadRequest(context.Context, *ReadToken) (*ReadResponse, error)
	CreateReplicate(context.Context, *CreateToken) (*CreateReplicateResponse, error)
	WriteReplicate(context.Context, *WriteToken) (*WriteReplicateResponse, error)
	ReadWrite(context.Context, *ReadWriteToken) (*ReadWriteResponse, error)
	mustEmbedUnimplementedTokenMgmtServer()
}

// UnimplementedTokenMgmtServer must be embedded to have forward compatible implementations.
type UnimplementedTokenMgmtServer struct {
}

func (UnimplementedTokenMgmtServer) CreateRequest(context.Context, *CreateToken) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRequest not implemented")
}
func (UnimplementedTokenMgmtServer) DropRequest(context.Context, *DropToken) (*DropResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DropRequest not implemented")
}
func (UnimplementedTokenMgmtServer) WriteRequest(context.Context, *WriteToken) (*WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WriteRequest not implemented")
}
func (UnimplementedTokenMgmtServer) ReadRequest(context.Context, *ReadToken) (*ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadRequest not implemented")
}
func (UnimplementedTokenMgmtServer) CreateReplicate(context.Context, *CreateToken) (*CreateReplicateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateReplicate not implemented")
}
func (UnimplementedTokenMgmtServer) WriteReplicate(context.Context, *WriteToken) (*WriteReplicateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WriteReplicate not implemented")
}
func (UnimplementedTokenMgmtServer) ReadWrite(context.Context, *ReadWriteToken) (*ReadWriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadWrite not implemented")
}
func (UnimplementedTokenMgmtServer) mustEmbedUnimplementedTokenMgmtServer() {}

// UnsafeTokenMgmtServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TokenMgmtServer will
// result in compilation errors.
type UnsafeTokenMgmtServer interface {
	mustEmbedUnimplementedTokenMgmtServer()
}

func RegisterTokenMgmtServer(s grpc.ServiceRegistrar, srv TokenMgmtServer) {
	s.RegisterService(&TokenMgmt_ServiceDesc, srv)
}

func _TokenMgmt_CreateRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenMgmtServer).CreateRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Protos.TokenMgmt/CreateRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenMgmtServer).CreateRequest(ctx, req.(*CreateToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenMgmt_DropRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DropToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenMgmtServer).DropRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Protos.TokenMgmt/DropRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenMgmtServer).DropRequest(ctx, req.(*DropToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenMgmt_WriteRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenMgmtServer).WriteRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Protos.TokenMgmt/WriteRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenMgmtServer).WriteRequest(ctx, req.(*WriteToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenMgmt_ReadRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenMgmtServer).ReadRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Protos.TokenMgmt/ReadRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenMgmtServer).ReadRequest(ctx, req.(*ReadToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenMgmt_CreateReplicate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenMgmtServer).CreateReplicate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Protos.TokenMgmt/CreateReplicate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenMgmtServer).CreateReplicate(ctx, req.(*CreateToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenMgmt_WriteReplicate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenMgmtServer).WriteReplicate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Protos.TokenMgmt/WriteReplicate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenMgmtServer).WriteReplicate(ctx, req.(*WriteToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _TokenMgmt_ReadWrite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadWriteToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenMgmtServer).ReadWrite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Protos.TokenMgmt/ReadWrite",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenMgmtServer).ReadWrite(ctx, req.(*ReadWriteToken))
	}
	return interceptor(ctx, in, info, handler)
}

// TokenMgmt_ServiceDesc is the grpc.ServiceDesc for TokenMgmt service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TokenMgmt_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Protos.TokenMgmt",
	HandlerType: (*TokenMgmtServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRequest",
			Handler:    _TokenMgmt_CreateRequest_Handler,
		},
		{
			MethodName: "DropRequest",
			Handler:    _TokenMgmt_DropRequest_Handler,
		},
		{
			MethodName: "WriteRequest",
			Handler:    _TokenMgmt_WriteRequest_Handler,
		},
		{
			MethodName: "ReadRequest",
			Handler:    _TokenMgmt_ReadRequest_Handler,
		},
		{
			MethodName: "CreateReplicate",
			Handler:    _TokenMgmt_CreateReplicate_Handler,
		},
		{
			MethodName: "WriteReplicate",
			Handler:    _TokenMgmt_WriteReplicate_Handler,
		},
		{
			MethodName: "ReadWrite",
			Handler:    _TokenMgmt_ReadWrite_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Protos/messages.proto",
}
