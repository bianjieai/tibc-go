// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: tibc/core/client/v1/query.proto

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
	Query_ClientState_FullMethodName     = "/tibc.core.client.v1.Query/ClientState"
	Query_ClientStates_FullMethodName    = "/tibc.core.client.v1.Query/ClientStates"
	Query_ConsensusState_FullMethodName  = "/tibc.core.client.v1.Query/ConsensusState"
	Query_ConsensusStates_FullMethodName = "/tibc.core.client.v1.Query/ConsensusStates"
	Query_Relayers_FullMethodName        = "/tibc.core.client.v1.Query/Relayers"
)

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Query provides defines the gRPC querier service
type QueryClient interface {
	// ClientState queries an TIBC light client.
	ClientState(ctx context.Context, in *QueryClientStateRequest, opts ...grpc.CallOption) (*QueryClientStateResponse, error)
	// ClientStates queries all the TIBC light clients of a chain.
	ClientStates(ctx context.Context, in *QueryClientStatesRequest, opts ...grpc.CallOption) (*QueryClientStatesResponse, error)
	// ConsensusState queries a consensus state associated with a client state at
	// a given height.
	ConsensusState(ctx context.Context, in *QueryConsensusStateRequest, opts ...grpc.CallOption) (*QueryConsensusStateResponse, error)
	// ConsensusStates queries all the consensus state associated with a given
	// client.
	ConsensusStates(ctx context.Context, in *QueryConsensusStatesRequest, opts ...grpc.CallOption) (*QueryConsensusStatesResponse, error)
	// Relayers queries all the relayers associated with a given
	// client.
	Relayers(ctx context.Context, in *QueryRelayersRequest, opts ...grpc.CallOption) (*QueryRelayersResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) ClientState(ctx context.Context, in *QueryClientStateRequest, opts ...grpc.CallOption) (*QueryClientStateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryClientStateResponse)
	err := c.cc.Invoke(ctx, Query_ClientState_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ClientStates(ctx context.Context, in *QueryClientStatesRequest, opts ...grpc.CallOption) (*QueryClientStatesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryClientStatesResponse)
	err := c.cc.Invoke(ctx, Query_ClientStates_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ConsensusState(ctx context.Context, in *QueryConsensusStateRequest, opts ...grpc.CallOption) (*QueryConsensusStateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryConsensusStateResponse)
	err := c.cc.Invoke(ctx, Query_ConsensusState_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ConsensusStates(ctx context.Context, in *QueryConsensusStatesRequest, opts ...grpc.CallOption) (*QueryConsensusStatesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryConsensusStatesResponse)
	err := c.cc.Invoke(ctx, Query_ConsensusStates_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Relayers(ctx context.Context, in *QueryRelayersRequest, opts ...grpc.CallOption) (*QueryRelayersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryRelayersResponse)
	err := c.cc.Invoke(ctx, Query_Relayers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility.
//
// Query provides defines the gRPC querier service
type QueryServer interface {
	// ClientState queries an TIBC light client.
	ClientState(context.Context, *QueryClientStateRequest) (*QueryClientStateResponse, error)
	// ClientStates queries all the TIBC light clients of a chain.
	ClientStates(context.Context, *QueryClientStatesRequest) (*QueryClientStatesResponse, error)
	// ConsensusState queries a consensus state associated with a client state at
	// a given height.
	ConsensusState(context.Context, *QueryConsensusStateRequest) (*QueryConsensusStateResponse, error)
	// ConsensusStates queries all the consensus state associated with a given
	// client.
	ConsensusStates(context.Context, *QueryConsensusStatesRequest) (*QueryConsensusStatesResponse, error)
	// Relayers queries all the relayers associated with a given
	// client.
	Relayers(context.Context, *QueryRelayersRequest) (*QueryRelayersResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedQueryServer struct{}

func (UnimplementedQueryServer) ClientState(context.Context, *QueryClientStateRequest) (*QueryClientStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClientState not implemented")
}
func (UnimplementedQueryServer) ClientStates(context.Context, *QueryClientStatesRequest) (*QueryClientStatesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClientStates not implemented")
}
func (UnimplementedQueryServer) ConsensusState(context.Context, *QueryConsensusStateRequest) (*QueryConsensusStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConsensusState not implemented")
}
func (UnimplementedQueryServer) ConsensusStates(context.Context, *QueryConsensusStatesRequest) (*QueryConsensusStatesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConsensusStates not implemented")
}
func (UnimplementedQueryServer) Relayers(context.Context, *QueryRelayersRequest) (*QueryRelayersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Relayers not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}
func (UnimplementedQueryServer) testEmbeddedByValue()               {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	// If the following call pancis, it indicates UnimplementedQueryServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_ClientState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryClientStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ClientState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ClientState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ClientState(ctx, req.(*QueryClientStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ClientStates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryClientStatesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ClientStates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ClientStates_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ClientStates(ctx, req.(*QueryClientStatesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ConsensusState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryConsensusStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ConsensusState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ConsensusState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ConsensusState(ctx, req.(*QueryConsensusStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ConsensusStates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryConsensusStatesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ConsensusStates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_ConsensusStates_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ConsensusStates(ctx, req.(*QueryConsensusStatesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Relayers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRelayersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Relayers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Query_Relayers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Relayers(ctx, req.(*QueryRelayersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tibc.core.client.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ClientState",
			Handler:    _Query_ClientState_Handler,
		},
		{
			MethodName: "ClientStates",
			Handler:    _Query_ClientStates_Handler,
		},
		{
			MethodName: "ConsensusState",
			Handler:    _Query_ConsensusState_Handler,
		},
		{
			MethodName: "ConsensusStates",
			Handler:    _Query_ConsensusStates_Handler,
		},
		{
			MethodName: "Relayers",
			Handler:    _Query_Relayers_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tibc/core/client/v1/query.proto",
}
