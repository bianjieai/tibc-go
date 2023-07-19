// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: tibc/core/packet/v1/tx.proto

package packetv1

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

const (
	Msg_RecvPacket_FullMethodName      = "/tibc.core.packet.v1.Msg/RecvPacket"
	Msg_Acknowledgement_FullMethodName = "/tibc.core.packet.v1.Msg/Acknowledgement"
	Msg_CleanPacket_FullMethodName     = "/tibc.core.packet.v1.Msg/CleanPacket"
	Msg_RecvCleanPacket_FullMethodName = "/tibc.core.packet.v1.Msg/RecvCleanPacket"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	// RecvPacket defines a rpc handler method for MsgRecvPacket.
	RecvPacket(ctx context.Context, in *MsgRecvPacket, opts ...grpc.CallOption) (*MsgRecvPacketResponse, error)
	// Acknowledgement defines a rpc handler method for MsgAcknowledgement.
	Acknowledgement(ctx context.Context, in *MsgAcknowledgement, opts ...grpc.CallOption) (*MsgAcknowledgementResponse, error)
	// CleanPacket defines a rpc handler method for MsgCleanPacket.
	CleanPacket(ctx context.Context, in *MsgCleanPacket, opts ...grpc.CallOption) (*MsgCleanPacketResponse, error)
	// RecvCleanPacket defines a rpc handler method for MsgRecvCleanPacket.
	RecvCleanPacket(ctx context.Context, in *MsgRecvCleanPacket, opts ...grpc.CallOption) (*MsgRecvCleanPacketResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) RecvPacket(ctx context.Context, in *MsgRecvPacket, opts ...grpc.CallOption) (*MsgRecvPacketResponse, error) {
	out := new(MsgRecvPacketResponse)
	err := c.cc.Invoke(ctx, Msg_RecvPacket_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Acknowledgement(ctx context.Context, in *MsgAcknowledgement, opts ...grpc.CallOption) (*MsgAcknowledgementResponse, error) {
	out := new(MsgAcknowledgementResponse)
	err := c.cc.Invoke(ctx, Msg_Acknowledgement_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) CleanPacket(ctx context.Context, in *MsgCleanPacket, opts ...grpc.CallOption) (*MsgCleanPacketResponse, error) {
	out := new(MsgCleanPacketResponse)
	err := c.cc.Invoke(ctx, Msg_CleanPacket_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) RecvCleanPacket(ctx context.Context, in *MsgRecvCleanPacket, opts ...grpc.CallOption) (*MsgRecvCleanPacketResponse, error) {
	out := new(MsgRecvCleanPacketResponse)
	err := c.cc.Invoke(ctx, Msg_RecvCleanPacket_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	// RecvPacket defines a rpc handler method for MsgRecvPacket.
	RecvPacket(context.Context, *MsgRecvPacket) (*MsgRecvPacketResponse, error)
	// Acknowledgement defines a rpc handler method for MsgAcknowledgement.
	Acknowledgement(context.Context, *MsgAcknowledgement) (*MsgAcknowledgementResponse, error)
	// CleanPacket defines a rpc handler method for MsgCleanPacket.
	CleanPacket(context.Context, *MsgCleanPacket) (*MsgCleanPacketResponse, error)
	// RecvCleanPacket defines a rpc handler method for MsgRecvCleanPacket.
	RecvCleanPacket(context.Context, *MsgRecvCleanPacket) (*MsgRecvCleanPacketResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) RecvPacket(context.Context, *MsgRecvPacket) (*MsgRecvPacketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecvPacket not implemented")
}
func (UnimplementedMsgServer) Acknowledgement(context.Context, *MsgAcknowledgement) (*MsgAcknowledgementResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Acknowledgement not implemented")
}
func (UnimplementedMsgServer) CleanPacket(context.Context, *MsgCleanPacket) (*MsgCleanPacketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CleanPacket not implemented")
}
func (UnimplementedMsgServer) RecvCleanPacket(context.Context, *MsgRecvCleanPacket) (*MsgRecvCleanPacketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecvCleanPacket not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_RecvPacket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRecvPacket)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RecvPacket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_RecvPacket_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RecvPacket(ctx, req.(*MsgRecvPacket))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Acknowledgement_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAcknowledgement)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Acknowledgement(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_Acknowledgement_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Acknowledgement(ctx, req.(*MsgAcknowledgement))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CleanPacket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCleanPacket)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CleanPacket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_CleanPacket_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CleanPacket(ctx, req.(*MsgCleanPacket))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RecvCleanPacket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRecvCleanPacket)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RecvCleanPacket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_RecvCleanPacket_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RecvCleanPacket(ctx, req.(*MsgRecvCleanPacket))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tibc.core.packet.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RecvPacket",
			Handler:    _Msg_RecvPacket_Handler,
		},
		{
			MethodName: "Acknowledgement",
			Handler:    _Msg_Acknowledgement_Handler,
		},
		{
			MethodName: "CleanPacket",
			Handler:    _Msg_CleanPacket_Handler,
		},
		{
			MethodName: "RecvCleanPacket",
			Handler:    _Msg_RecvCleanPacket_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tibc/core/packet/v1/tx.proto",
}
