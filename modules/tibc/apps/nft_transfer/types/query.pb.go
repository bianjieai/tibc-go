// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: tibc/apps/nft_transfer/v1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	query "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// QueryClassTraceRequest is the request type for the Query/ClassTrace RPC
// method
type QueryClassTraceRequest struct {
	// hash (in hex format) of the class trace information.
	Hash string `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (m *QueryClassTraceRequest) Reset()         { *m = QueryClassTraceRequest{} }
func (m *QueryClassTraceRequest) String() string { return proto.CompactTextString(m) }
func (*QueryClassTraceRequest) ProtoMessage()    {}
func (*QueryClassTraceRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1bc5b7eb47b9a64, []int{0}
}
func (m *QueryClassTraceRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryClassTraceRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryClassTraceRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryClassTraceRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryClassTraceRequest.Merge(m, src)
}
func (m *QueryClassTraceRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryClassTraceRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryClassTraceRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryClassTraceRequest proto.InternalMessageInfo

func (m *QueryClassTraceRequest) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

// QueryClassTraceResponse is the response type for the Query/ClassTrace RPC
// method.
type QueryClassTraceResponse struct {
	// class_trace returns the requested class trace information.
	ClassTrace *ClassTrace `protobuf:"bytes,1,opt,name=class_trace,json=classTrace,proto3" json:"class_trace,omitempty"`
}

func (m *QueryClassTraceResponse) Reset()         { *m = QueryClassTraceResponse{} }
func (m *QueryClassTraceResponse) String() string { return proto.CompactTextString(m) }
func (*QueryClassTraceResponse) ProtoMessage()    {}
func (*QueryClassTraceResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1bc5b7eb47b9a64, []int{1}
}
func (m *QueryClassTraceResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryClassTraceResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryClassTraceResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryClassTraceResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryClassTraceResponse.Merge(m, src)
}
func (m *QueryClassTraceResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryClassTraceResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryClassTraceResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryClassTraceResponse proto.InternalMessageInfo

func (m *QueryClassTraceResponse) GetClassTrace() *ClassTrace {
	if m != nil {
		return m.ClassTrace
	}
	return nil
}

// QueryConnectionsRequest is the request type for the Query/ClassTraces RPC
// method
type QueryClassTracesRequest struct {
	// pagination defines an optional pagination for the request.
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryClassTracesRequest) Reset()         { *m = QueryClassTracesRequest{} }
func (m *QueryClassTracesRequest) String() string { return proto.CompactTextString(m) }
func (*QueryClassTracesRequest) ProtoMessage()    {}
func (*QueryClassTracesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1bc5b7eb47b9a64, []int{2}
}
func (m *QueryClassTracesRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryClassTracesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryClassTracesRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryClassTracesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryClassTracesRequest.Merge(m, src)
}
func (m *QueryClassTracesRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryClassTracesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryClassTracesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryClassTracesRequest proto.InternalMessageInfo

func (m *QueryClassTracesRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

// QueryConnectionsResponse is the response type for the Query/ClassTraces RPC
// method.
type QueryClassTracesResponse struct {
	// class_traces returns all class trace information.
	ClassTraces Traces `protobuf:"bytes,1,rep,name=class_traces,json=classTraces,proto3,castrepeated=Traces" json:"class_traces"`
	// pagination defines the pagination in the response.
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryClassTracesResponse) Reset()         { *m = QueryClassTracesResponse{} }
func (m *QueryClassTracesResponse) String() string { return proto.CompactTextString(m) }
func (*QueryClassTracesResponse) ProtoMessage()    {}
func (*QueryClassTracesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f1bc5b7eb47b9a64, []int{3}
}
func (m *QueryClassTracesResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryClassTracesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryClassTracesResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryClassTracesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryClassTracesResponse.Merge(m, src)
}
func (m *QueryClassTracesResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryClassTracesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryClassTracesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryClassTracesResponse proto.InternalMessageInfo

func (m *QueryClassTracesResponse) GetClassTraces() Traces {
	if m != nil {
		return m.ClassTraces
	}
	return nil
}

func (m *QueryClassTracesResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryClassTraceRequest)(nil), "tibc.apps.nft_transfer.v1.QueryClassTraceRequest")
	proto.RegisterType((*QueryClassTraceResponse)(nil), "tibc.apps.nft_transfer.v1.QueryClassTraceResponse")
	proto.RegisterType((*QueryClassTracesRequest)(nil), "tibc.apps.nft_transfer.v1.QueryClassTracesRequest")
	proto.RegisterType((*QueryClassTracesResponse)(nil), "tibc.apps.nft_transfer.v1.QueryClassTracesResponse")
}

func init() {
	proto.RegisterFile("tibc/apps/nft_transfer/v1/query.proto", fileDescriptor_f1bc5b7eb47b9a64)
}

var fileDescriptor_f1bc5b7eb47b9a64 = []byte{
	// 467 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0x41, 0x6b, 0x13, 0x41,
	0x14, 0xc7, 0x33, 0x51, 0x0b, 0x4e, 0xc4, 0xc3, 0x20, 0x1a, 0x83, 0x6c, 0xcb, 0x42, 0xad, 0x4a,
	0x9d, 0x71, 0x53, 0xf0, 0x6e, 0x85, 0x7a, 0xd5, 0xe0, 0x41, 0xf4, 0x20, 0xb3, 0xeb, 0x74, 0x32,
	0x92, 0xcc, 0x6c, 0xf7, 0xcd, 0x06, 0x8a, 0x78, 0xf1, 0x13, 0x08, 0x7e, 0x08, 0x41, 0x3f, 0x83,
	0xf7, 0x1e, 0x0b, 0x5e, 0x3c, 0x59, 0x49, 0xfc, 0x20, 0x32, 0x33, 0x1b, 0xb3, 0xa5, 0x5d, 0x4d,
	0x6e, 0x8f, 0xe4, 0xff, 0xde, 0xff, 0xf7, 0xfe, 0xf3, 0x16, 0x6f, 0x5a, 0x95, 0x66, 0x8c, 0xe7,
	0x39, 0x30, 0xbd, 0x6f, 0x5f, 0xdb, 0x82, 0x6b, 0xd8, 0x17, 0x05, 0x9b, 0x24, 0xec, 0xa0, 0x14,
	0xc5, 0x21, 0xcd, 0x0b, 0x63, 0x0d, 0xb9, 0xe9, 0x64, 0xd4, 0xc9, 0x68, 0x5d, 0x46, 0x27, 0x49,
	0xef, 0x9a, 0x34, 0xd2, 0x78, 0x15, 0x73, 0x55, 0x68, 0xe8, 0xdd, 0xcb, 0x0c, 0x8c, 0x0d, 0xb0,
	0x94, 0x83, 0x08, 0x93, 0xd8, 0x24, 0x49, 0x85, 0xe5, 0x09, 0xcb, 0xb9, 0x54, 0x9a, 0x5b, 0x65,
	0x74, 0xa5, 0xdd, 0x6e, 0x66, 0x38, 0x65, 0x16, 0xd4, 0xb7, 0xa4, 0x31, 0x72, 0x24, 0x18, 0xcf,
	0x15, 0xe3, 0x5a, 0x1b, 0xeb, 0x47, 0x41, 0xf8, 0x37, 0xde, 0xc6, 0xd7, 0x9f, 0x39, 0xb7, 0xc7,
	0x23, 0x0e, 0xf0, 0xbc, 0xe0, 0x99, 0x18, 0x88, 0x83, 0x52, 0x80, 0x25, 0x04, 0x5f, 0x1c, 0x72,
	0x18, 0x76, 0xd1, 0x06, 0xba, 0x73, 0x79, 0xe0, 0xeb, 0x98, 0xe3, 0x1b, 0x67, 0xd4, 0x90, 0x1b,
	0x0d, 0x82, 0xec, 0xe1, 0x4e, 0xe6, 0x7e, 0x75, 0xf6, 0x99, 0xf0, 0x5d, 0x9d, 0xfe, 0x26, 0x6d,
	0xcc, 0x81, 0xd6, 0x66, 0xe0, 0xec, 0x6f, 0x7d, 0x8e, 0x05, 0xcc, 0x89, 0xf6, 0x30, 0x5e, 0x64,
	0x51, 0x39, 0xdc, 0xa6, 0x21, 0x38, 0xea, 0x82, 0xa3, 0xe1, 0x09, 0xaa, 0xe0, 0xe8, 0x53, 0x2e,
	0xe7, 0xdb, 0x0c, 0x6a, 0x9d, 0xf1, 0x37, 0x84, 0xbb, 0x67, 0x3d, 0xaa, 0x3d, 0x5e, 0xe0, 0x2b,
	0xb5, 0x3d, 0xa0, 0x8b, 0x36, 0x2e, 0x2c, 0xbd, 0xc8, 0xee, 0xd5, 0xa3, 0x9f, 0xeb, 0xad, 0x2f,
	0x27, 0xeb, 0x6b, 0xd5, 0xd0, 0xce, 0x62, 0x31, 0x20, 0x4f, 0x4e, 0xe1, 0xb7, 0x3d, 0xfe, 0xd6,
	0x7f, 0xf1, 0x03, 0x56, 0x9d, 0xbf, 0x7f, 0xd2, 0xc6, 0x97, 0x3c, 0x3f, 0xf9, 0x8a, 0x30, 0x5e,
	0xd8, 0x93, 0xe4, 0x1f, 0x94, 0xe7, 0xbf, 0x72, 0xaf, 0xbf, 0x4a, 0x4b, 0x60, 0x89, 0x1f, 0x7e,
	0xf8, 0xfe, 0xfb, 0x53, 0xfb, 0x01, 0xa1, 0xac, 0xf9, 0x10, 0xeb, 0x19, 0xb2, 0x77, 0xee, 0x78,
	0xde, 0x93, 0xcf, 0x08, 0x77, 0x6a, 0x91, 0x93, 0x15, 0xbc, 0xe7, 0x37, 0xd0, 0xdb, 0x59, 0xa9,
	0xa7, 0x02, 0x66, 0x1e, 0xf8, 0x2e, 0xd9, 0x5a, 0x12, 0x78, 0xf7, 0xd5, 0xd1, 0x34, 0x42, 0xc7,
	0xd3, 0x08, 0xfd, 0x9a, 0x46, 0xe8, 0xe3, 0x2c, 0x6a, 0x1d, 0xcf, 0xa2, 0xd6, 0x8f, 0x59, 0xd4,
	0x7a, 0xf9, 0x48, 0x2a, 0x3b, 0x2c, 0x53, 0x9a, 0x99, 0x31, 0x4b, 0x15, 0xd7, 0x6f, 0x95, 0xe0,
	0xca, 0x8f, 0xbd, 0x2f, 0x0d, 0x1b, 0x9b, 0x37, 0xe5, 0x48, 0x40, 0x93, 0x8d, 0x3d, 0xcc, 0x05,
	0xa4, 0x6b, 0xfe, 0xcb, 0xdb, 0xf9, 0x13, 0x00, 0x00, 0xff, 0xff, 0x12, 0x84, 0x97, 0x8a, 0x4b,
	0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// ClassTrace queries a class trace information.
	ClassTrace(ctx context.Context, in *QueryClassTraceRequest, opts ...grpc.CallOption) (*QueryClassTraceResponse, error)
	// ClassTraces queries all class traces.
	ClassTraces(ctx context.Context, in *QueryClassTracesRequest, opts ...grpc.CallOption) (*QueryClassTracesResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) ClassTrace(ctx context.Context, in *QueryClassTraceRequest, opts ...grpc.CallOption) (*QueryClassTraceResponse, error) {
	out := new(QueryClassTraceResponse)
	err := c.cc.Invoke(ctx, "/tibc.apps.nft_transfer.v1.Query/ClassTrace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ClassTraces(ctx context.Context, in *QueryClassTracesRequest, opts ...grpc.CallOption) (*QueryClassTracesResponse, error) {
	out := new(QueryClassTracesResponse)
	err := c.cc.Invoke(ctx, "/tibc.apps.nft_transfer.v1.Query/ClassTraces", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// ClassTrace queries a class trace information.
	ClassTrace(context.Context, *QueryClassTraceRequest) (*QueryClassTraceResponse, error)
	// ClassTraces queries all class traces.
	ClassTraces(context.Context, *QueryClassTracesRequest) (*QueryClassTracesResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) ClassTrace(ctx context.Context, req *QueryClassTraceRequest) (*QueryClassTraceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClassTrace not implemented")
}
func (*UnimplementedQueryServer) ClassTraces(ctx context.Context, req *QueryClassTracesRequest) (*QueryClassTracesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClassTraces not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_ClassTrace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryClassTraceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ClassTrace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tibc.apps.nft_transfer.v1.Query/ClassTrace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ClassTrace(ctx, req.(*QueryClassTraceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_ClassTraces_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryClassTracesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ClassTraces(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tibc.apps.nft_transfer.v1.Query/ClassTraces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ClassTraces(ctx, req.(*QueryClassTracesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "tibc.apps.nft_transfer.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ClassTrace",
			Handler:    _Query_ClassTrace_Handler,
		},
		{
			MethodName: "ClassTraces",
			Handler:    _Query_ClassTraces_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tibc/apps/nft_transfer/v1/query.proto",
}

func (m *QueryClassTraceRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryClassTraceRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryClassTraceRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryClassTraceResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryClassTraceResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryClassTraceResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ClassTrace != nil {
		{
			size, err := m.ClassTrace.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryClassTracesRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryClassTracesRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryClassTracesRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryClassTracesResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryClassTracesResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryClassTracesResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.ClassTraces) > 0 {
		for iNdEx := len(m.ClassTraces) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ClassTraces[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryClassTraceRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryClassTraceResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ClassTrace != nil {
		l = m.ClassTrace.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryClassTracesRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryClassTracesResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ClassTraces) > 0 {
		for _, e := range m.ClassTraces {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryClassTraceRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryClassTraceRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryClassTraceRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryClassTraceResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryClassTraceResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryClassTraceResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClassTrace", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ClassTrace == nil {
				m.ClassTrace = &ClassTrace{}
			}
			if err := m.ClassTrace.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryClassTracesRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryClassTracesRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryClassTracesRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageRequest{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryClassTracesResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryClassTracesResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryClassTracesResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClassTraces", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClassTraces = append(m.ClassTraces, ClassTrace{})
			if err := m.ClassTraces[len(m.ClassTraces)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageResponse{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
