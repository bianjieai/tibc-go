// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: tibc/core/client/v1/tx.proto

package clientv1

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
	Msg_CreateClient_FullMethodName    = "/tibc.core.client.v1.Msg/CreateClient"
	Msg_UpdateClient_FullMethodName    = "/tibc.core.client.v1.Msg/UpdateClient"
	Msg_UpgradeClient_FullMethodName   = "/tibc.core.client.v1.Msg/UpgradeClient"
	Msg_RegisterRelayer_FullMethodName = "/tibc.core.client.v1.Msg/RegisterRelayer"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Msg defines the tibc/client Msg service.
type MsgClient interface {
	// CreateClient defines a rpc handler method for MsgCreateClient.
	CreateClient(ctx context.Context, in *MsgCreateClient, opts ...grpc.CallOption) (*MsgCreateClientResponse, error)
	// UpdateClient defines a rpc handler method for MsgUpdateClient.
	UpdateClient(ctx context.Context, in *MsgUpdateClient, opts ...grpc.CallOption) (*MsgUpdateClientResponse, error)
	// UpgradeClient defines a rpc handler method for MsgUpgradeClient.
	UpgradeClient(ctx context.Context, in *MsgUpgradeClient, opts ...grpc.CallOption) (*MsgUpgradeClientResponse, error)
	// RegisterRelayer defines a rpc handler method for MsgRegisterRelayer.
	RegisterRelayer(ctx context.Context, in *MsgRegisterRelayer, opts ...grpc.CallOption) (*MsgRegisterRelayerResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) CreateClient(ctx context.Context, in *MsgCreateClient, opts ...grpc.CallOption) (*MsgCreateClientResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgCreateClientResponse)
	err := c.cc.Invoke(ctx, Msg_CreateClient_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateClient(ctx context.Context, in *MsgUpdateClient, opts ...grpc.CallOption) (*MsgUpdateClientResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgUpdateClientResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateClient_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpgradeClient(ctx context.Context, in *MsgUpgradeClient, opts ...grpc.CallOption) (*MsgUpgradeClientResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgUpgradeClientResponse)
	err := c.cc.Invoke(ctx, Msg_UpgradeClient_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) RegisterRelayer(ctx context.Context, in *MsgRegisterRelayer, opts ...grpc.CallOption) (*MsgRegisterRelayerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MsgRegisterRelayerResponse)
	err := c.cc.Invoke(ctx, Msg_RegisterRelayer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility.
//
// Msg defines the tibc/client Msg service.
type MsgServer interface {
	// CreateClient defines a rpc handler method for MsgCreateClient.
	CreateClient(context.Context, *MsgCreateClient) (*MsgCreateClientResponse, error)
	// UpdateClient defines a rpc handler method for MsgUpdateClient.
	UpdateClient(context.Context, *MsgUpdateClient) (*MsgUpdateClientResponse, error)
	// UpgradeClient defines a rpc handler method for MsgUpgradeClient.
	UpgradeClient(context.Context, *MsgUpgradeClient) (*MsgUpgradeClientResponse, error)
	// RegisterRelayer defines a rpc handler method for MsgRegisterRelayer.
	RegisterRelayer(context.Context, *MsgRegisterRelayer) (*MsgRegisterRelayerResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMsgServer struct{}

func (UnimplementedMsgServer) CreateClient(context.Context, *MsgCreateClient) (*MsgCreateClientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateClient not implemented")
}
func (UnimplementedMsgServer) UpdateClient(context.Context, *MsgUpdateClient) (*MsgUpdateClientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateClient not implemented")
}
func (UnimplementedMsgServer) UpgradeClient(context.Context, *MsgUpgradeClient) (*MsgUpgradeClientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpgradeClient not implemented")
}
func (UnimplementedMsgServer) RegisterRelayer(context.Context, *MsgRegisterRelayer) (*MsgRegisterRelayerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterRelayer not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}
func (UnimplementedMsgServer) testEmbeddedByValue()             {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	// If the following call pancis, it indicates UnimplementedMsgServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_CreateClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateClient)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_CreateClient_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateClient(ctx, req.(*MsgCreateClient))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateClient)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpdateClient_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateClient(ctx, req.(*MsgUpdateClient))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpgradeClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpgradeClient)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpgradeClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_UpgradeClient_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpgradeClient(ctx, req.(*MsgUpgradeClient))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RegisterRelayer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRegisterRelayer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RegisterRelayer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_RegisterRelayer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RegisterRelayer(ctx, req.(*MsgRegisterRelayer))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tibc.core.client.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateClient",
			Handler:    _Msg_CreateClient_Handler,
		},
		{
			MethodName: "UpdateClient",
			Handler:    _Msg_UpdateClient_Handler,
		},
		{
			MethodName: "UpgradeClient",
			Handler:    _Msg_UpgradeClient_Handler,
		},
		{
			MethodName: "RegisterRelayer",
			Handler:    _Msg_RegisterRelayer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tibc/core/client/v1/tx.proto",
}
