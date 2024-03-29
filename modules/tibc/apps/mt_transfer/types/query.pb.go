// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: tibc/apps/mt_transfer/v1/query.proto

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
	return fileDescriptor_00405040fdfa1b2c, []int{0}
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
	return fileDescriptor_00405040fdfa1b2c, []int{1}
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
	return fileDescriptor_00405040fdfa1b2c, []int{2}
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
	return fileDescriptor_00405040fdfa1b2c, []int{3}
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
	proto.RegisterType((*QueryClassTraceRequest)(nil), "tibc.apps.mt_transfer.v1.QueryClassTraceRequest")
	proto.RegisterType((*QueryClassTraceResponse)(nil), "tibc.apps.mt_transfer.v1.QueryClassTraceResponse")
	proto.RegisterType((*QueryClassTracesRequest)(nil), "tibc.apps.mt_transfer.v1.QueryClassTracesRequest")
	proto.RegisterType((*QueryClassTracesResponse)(nil), "tibc.apps.mt_transfer.v1.QueryClassTracesResponse")
}

func init() {
	proto.RegisterFile("tibc/apps/mt_transfer/v1/query.proto", fileDescriptor_00405040fdfa1b2c)
}

var fileDescriptor_00405040fdfa1b2c = []byte{
	// 466 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0xc1, 0x8a, 0x13, 0x31,
	0x18, 0xc7, 0x9b, 0xaa, 0x0b, 0xa6, 0xe2, 0x21, 0x88, 0x96, 0x22, 0xb3, 0xcb, 0xb0, 0xac, 0x8b,
	0xb8, 0x89, 0x53, 0xf1, 0x2e, 0x2b, 0xea, 0x55, 0x8b, 0x20, 0xec, 0x45, 0x33, 0x63, 0x4c, 0x23,
	0xed, 0x64, 0x76, 0xbe, 0x4c, 0x61, 0x11, 0x2f, 0x3e, 0x81, 0xe0, 0x2b, 0x78, 0x10, 0x9f, 0xc1,
	0x07, 0xd8, 0xe3, 0x82, 0x17, 0xbd, 0xa8, 0xb4, 0x3e, 0x88, 0x24, 0x99, 0xb1, 0x91, 0x3a, 0x6c,
	0xf7, 0xf6, 0x31, 0xf3, 0xff, 0xbe, 0xff, 0xef, 0xfb, 0x27, 0xc1, 0xdb, 0x46, 0xa5, 0x19, 0xe3,
	0x45, 0x01, 0x6c, 0x6a, 0x9e, 0x9b, 0x92, 0xe7, 0xf0, 0x4a, 0x94, 0x6c, 0x96, 0xb0, 0xc3, 0x4a,
	0x94, 0x47, 0xb4, 0x28, 0xb5, 0xd1, 0xa4, 0x6f, 0x55, 0xd4, 0xaa, 0x68, 0xa0, 0xa2, 0xb3, 0x64,
	0x70, 0x45, 0x6a, 0xa9, 0x9d, 0x88, 0xd9, 0xca, 0xeb, 0x07, 0x37, 0x33, 0x0d, 0x53, 0x0d, 0x2c,
	0xe5, 0x20, 0xfc, 0x20, 0x36, 0x4b, 0x52, 0x61, 0x78, 0xc2, 0x0a, 0x2e, 0x55, 0xce, 0x8d, 0xd2,
	0x79, 0xa3, 0x6d, 0x25, 0x08, 0xad, 0xbc, 0xf6, 0xba, 0xd4, 0x5a, 0x4e, 0x04, 0xe3, 0x85, 0x62,
	0x3c, 0xcf, 0xb5, 0x71, 0x83, 0xc0, 0xff, 0x8d, 0x6f, 0xe1, 0xab, 0x4f, 0xac, 0xd7, 0xfd, 0x09,
	0x07, 0x78, 0x5a, 0xf2, 0x4c, 0x8c, 0xc4, 0x61, 0x25, 0xc0, 0x10, 0x82, 0xcf, 0x8f, 0x39, 0x8c,
	0xfb, 0x68, 0x0b, 0xed, 0x5e, 0x1c, 0xb9, 0x3a, 0x7e, 0x81, 0xaf, 0xad, 0xa8, 0xa1, 0xd0, 0x39,
	0x08, 0xf2, 0x00, 0xf7, 0x32, 0xfb, 0xd5, 0xda, 0x67, 0xc2, 0x75, 0xf5, 0x86, 0xdb, 0xb4, 0x2d,
	0x04, 0x1a, 0x8c, 0xc0, 0xd9, 0xdf, 0x3a, 0xe6, 0x2b, 0x0e, 0xd0, 0x00, 0x3d, 0xc4, 0x78, 0x19,
	0x44, 0x6d, 0xb0, 0x43, 0x7d, 0x6a, 0xd4, 0xa6, 0x46, 0x7d, 0xfc, 0x75, 0x6a, 0xf4, 0x31, 0x97,
	0xcd, 0x32, 0xa3, 0xa0, 0x33, 0xfe, 0x82, 0x70, 0x7f, 0xd5, 0xa3, 0x5e, 0xe3, 0x19, 0xbe, 0x14,
	0xac, 0x01, 0x7d, 0xb4, 0x75, 0x6e, 0xdd, 0x3d, 0xf6, 0x2f, 0x1f, 0xff, 0xd8, 0xec, 0x7c, 0xfe,
	0xb9, 0xb9, 0x51, 0xcf, 0xec, 0x2d, 0xf7, 0x02, 0xf2, 0xe8, 0x1f, 0xfa, 0xae, 0xa3, 0xbf, 0x71,
	0x2a, 0xbd, 0xa7, 0x0a, 0xf1, 0x87, 0xdf, 0xbb, 0xf8, 0x82, 0xc3, 0x27, 0x9f, 0x10, 0xc6, 0x4b,
	0x7b, 0x72, 0xbb, 0x1d, 0xf2, 0xff, 0x47, 0x3c, 0x48, 0xce, 0xd0, 0xe1, 0x49, 0xe2, 0xbb, 0xef,
	0xbe, 0xfe, 0xfe, 0xd0, 0x65, 0x64, 0x8f, 0xb5, 0x5e, 0xc1, 0x30, 0x3f, 0xf6, 0xc6, 0xde, 0x9b,
	0xb7, 0xe4, 0x23, 0xc2, 0xbd, 0x20, 0x6e, 0xb2, 0xbe, 0x73, 0x73, 0xfc, 0x83, 0xe1, 0x59, 0x5a,
	0x6a, 0x5a, 0xea, 0x68, 0x77, 0xc9, 0xce, 0x7a, 0xb4, 0xfb, 0x07, 0xc7, 0xf3, 0x08, 0x9d, 0xcc,
	0x23, 0xf4, 0x6b, 0x1e, 0xa1, 0xf7, 0x8b, 0xa8, 0x73, 0xb2, 0x88, 0x3a, 0xdf, 0x16, 0x51, 0xe7,
	0xe0, 0x9e, 0x54, 0x66, 0x5c, 0xa5, 0x34, 0xd3, 0x53, 0x96, 0x2a, 0x9e, 0xbf, 0x56, 0x82, 0x2b,
	0x37, 0x75, 0x4f, 0x6a, 0x36, 0xd5, 0x2f, 0xab, 0x89, 0x80, 0x16, 0x17, 0x73, 0x54, 0x08, 0x48,
	0x37, 0xdc, 0x83, 0xbb, 0xf3, 0x27, 0x00, 0x00, 0xff, 0xff, 0xe3, 0x4e, 0xcb, 0xfd, 0x3e, 0x04,
	0x00, 0x00,
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
	err := c.cc.Invoke(ctx, "/tibc.apps.mt_transfer.v1.Query/ClassTrace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) ClassTraces(ctx context.Context, in *QueryClassTracesRequest, opts ...grpc.CallOption) (*QueryClassTracesResponse, error) {
	out := new(QueryClassTracesResponse)
	err := c.cc.Invoke(ctx, "/tibc.apps.mt_transfer.v1.Query/ClassTraces", in, out, opts...)
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
		FullMethod: "/tibc.apps.mt_transfer.v1.Query/ClassTrace",
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
		FullMethod: "/tibc.apps.mt_transfer.v1.Query/ClassTraces",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ClassTraces(ctx, req.(*QueryClassTracesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "tibc.apps.mt_transfer.v1.Query",
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
	Metadata: "tibc/apps/mt_transfer/v1/query.proto",
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
