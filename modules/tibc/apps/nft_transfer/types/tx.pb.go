// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: tibc/apps/nft_transfer/v1/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
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

type MsgNftTransfer struct {
	// the class to which the NFT to be transferred belongs
	Class string `protobuf:"bytes,1,opt,name=class,proto3" json:"class,omitempty"`
	// the nft id
	Id string `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	// the nft sender
	Sender string `protobuf:"bytes,3,opt,name=sender,proto3" json:"sender,omitempty"`
	// the nft receiver
	Receiver string `protobuf:"bytes,4,opt,name=receiver,proto3" json:"receiver,omitempty"`
	// target chain of transmission
	DestChain string `protobuf:"bytes,5,opt,name=dest_chain,json=destChain,proto3" json:"dest_chain,omitempty"`
	// relay chain during transmission
	RealayChain string `protobuf:"bytes,6,opt,name=realay_chain,json=realayChain,proto3" json:"realay_chain,omitempty"`
	// the destination contract address to receive the nft
	DestContract string `protobuf:"bytes,7,opt,name=dest_contract,json=destContract,proto3" json:"dest_contract,omitempty"`
}

func (m *MsgNftTransfer) Reset()         { *m = MsgNftTransfer{} }
func (m *MsgNftTransfer) String() string { return proto.CompactTextString(m) }
func (*MsgNftTransfer) ProtoMessage()    {}
func (*MsgNftTransfer) Descriptor() ([]byte, []int) {
	return fileDescriptor_9963dad398b8a1b5, []int{0}
}
func (m *MsgNftTransfer) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgNftTransfer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgNftTransfer.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgNftTransfer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgNftTransfer.Merge(m, src)
}
func (m *MsgNftTransfer) XXX_Size() int {
	return m.Size()
}
func (m *MsgNftTransfer) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgNftTransfer.DiscardUnknown(m)
}

var xxx_messageInfo_MsgNftTransfer proto.InternalMessageInfo

// MsgTransferResponse defines the Msg/NftTransfer response type.
type MsgNftTransferResponse struct {
}

func (m *MsgNftTransferResponse) Reset()         { *m = MsgNftTransferResponse{} }
func (m *MsgNftTransferResponse) String() string { return proto.CompactTextString(m) }
func (*MsgNftTransferResponse) ProtoMessage()    {}
func (*MsgNftTransferResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9963dad398b8a1b5, []int{1}
}
func (m *MsgNftTransferResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgNftTransferResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgNftTransferResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgNftTransferResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgNftTransferResponse.Merge(m, src)
}
func (m *MsgNftTransferResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgNftTransferResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgNftTransferResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgNftTransferResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgNftTransfer)(nil), "tibc.apps.nft_transfer.v1.MsgNftTransfer")
	proto.RegisterType((*MsgNftTransferResponse)(nil), "tibc.apps.nft_transfer.v1.MsgNftTransferResponse")
}

func init() {
	proto.RegisterFile("tibc/apps/nft_transfer/v1/tx.proto", fileDescriptor_9963dad398b8a1b5)
}

var fileDescriptor_9963dad398b8a1b5 = []byte{
	// 349 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0x31, 0x4f, 0x32, 0x31,
	0x18, 0xc7, 0xef, 0xe0, 0x85, 0x17, 0x0a, 0x32, 0x34, 0x84, 0x54, 0x12, 0x0f, 0xc5, 0x45, 0x07,
	0xdb, 0xa0, 0x9b, 0x9b, 0x3a, 0xe3, 0x40, 0x9c, 0x74, 0x20, 0xbd, 0xbb, 0x72, 0x54, 0xa1, 0xbd,
	0xb4, 0x85, 0xc8, 0x37, 0x70, 0xf4, 0x23, 0xf0, 0x71, 0x1c, 0x19, 0x4d, 0x5c, 0x0c, 0xb7, 0xf8,
	0x31, 0xcc, 0xf5, 0x4e, 0x03, 0x03, 0x89, 0xdb, 0x3d, 0xbf, 0xe7, 0x77, 0xff, 0x27, 0x4f, 0x1f,
	0xd0, 0x35, 0xdc, 0x0f, 0x08, 0x8d, 0x63, 0x4d, 0xc4, 0xc8, 0x0c, 0x8d, 0xa2, 0x42, 0x8f, 0x98,
	0x22, 0xf3, 0x1e, 0x31, 0xcf, 0x38, 0x56, 0xd2, 0x48, 0xb8, 0x9f, 0x3a, 0x38, 0x75, 0xf0, 0xa6,
	0x83, 0xe7, 0xbd, 0x76, 0x33, 0x92, 0x91, 0xb4, 0x16, 0x49, 0xbf, 0xb2, 0x1f, 0xba, 0x1f, 0x2e,
	0x68, 0xf4, 0x75, 0x74, 0x3b, 0x32, 0x77, 0xb9, 0x0b, 0x9b, 0xa0, 0x14, 0x4c, 0xa8, 0xd6, 0xc8,
	0x3d, 0x74, 0x4f, 0xaa, 0x83, 0xac, 0x80, 0x0d, 0x50, 0xe0, 0x21, 0x2a, 0x58, 0x54, 0xe0, 0x21,
	0x6c, 0x81, 0xb2, 0x66, 0x22, 0x64, 0x0a, 0x15, 0x2d, 0xcb, 0x2b, 0xd8, 0x06, 0x15, 0xc5, 0x02,
	0xc6, 0xe7, 0x4c, 0xa1, 0x7f, 0xb6, 0xf3, 0x5b, 0xc3, 0x03, 0x00, 0x42, 0xa6, 0xcd, 0x30, 0x18,
	0x53, 0x2e, 0x50, 0xc9, 0x76, 0xab, 0x29, 0xb9, 0x49, 0x01, 0x3c, 0x02, 0x75, 0xc5, 0xe8, 0x84,
	0x2e, 0x72, 0xa1, 0x6c, 0x85, 0x5a, 0xc6, 0x32, 0xe5, 0x18, 0xec, 0x65, 0x09, 0x52, 0x18, 0x45,
	0x03, 0x83, 0xfe, 0x5b, 0xa7, 0x6e, 0x43, 0x72, 0x76, 0x59, 0x79, 0x59, 0x76, 0x9c, 0xaf, 0x65,
	0xc7, 0xe9, 0x22, 0xd0, 0xda, 0x5e, 0x6e, 0xc0, 0x74, 0x2c, 0x85, 0x66, 0xe7, 0x0a, 0x14, 0xfb,
	0x3a, 0x82, 0x4f, 0xa0, 0xb6, 0xb9, 0xfa, 0x29, 0xde, 0xf9, 0x7e, 0x78, 0x3b, 0xa8, 0xdd, 0xfb,
	0xb3, 0xfa, 0x33, 0xf3, 0xfa, 0xe1, 0x6d, 0xed, 0xb9, 0xab, 0xb5, 0xe7, 0x7e, 0xae, 0x3d, 0xf7,
	0x35, 0xf1, 0x9c, 0x55, 0xe2, 0x39, 0xef, 0x89, 0xe7, 0xdc, 0x5f, 0x45, 0xdc, 0x8c, 0x67, 0x3e,
	0x0e, 0xe4, 0x94, 0xf8, 0x9c, 0x8a, 0x47, 0xce, 0x28, 0x27, 0xe9, 0x80, 0xb3, 0x48, 0x92, 0xa9,
	0x0c, 0x67, 0x13, 0xa6, 0xc9, 0x8e, 0xfb, 0x9b, 0x45, 0xcc, 0xb4, 0x5f, 0xb6, 0xf7, 0xbc, 0xf8,
	0x0e, 0x00, 0x00, 0xff, 0xff, 0xb5, 0xf0, 0x13, 0x54, 0x26, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	// NftTransfer defines a rpc handler method for MsgNftTransfer.
	NftTransfer(ctx context.Context, in *MsgNftTransfer, opts ...grpc.CallOption) (*MsgNftTransferResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) NftTransfer(ctx context.Context, in *MsgNftTransfer, opts ...grpc.CallOption) (*MsgNftTransferResponse, error) {
	out := new(MsgNftTransferResponse)
	err := c.cc.Invoke(ctx, "/tibc.apps.nft_transfer.v1.Msg/NftTransfer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// NftTransfer defines a rpc handler method for MsgNftTransfer.
	NftTransfer(context.Context, *MsgNftTransfer) (*MsgNftTransferResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) NftTransfer(ctx context.Context, req *MsgNftTransfer) (*MsgNftTransferResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NftTransfer not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_NftTransfer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgNftTransfer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).NftTransfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tibc.apps.nft_transfer.v1.Msg/NftTransfer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).NftTransfer(ctx, req.(*MsgNftTransfer))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "tibc.apps.nft_transfer.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NftTransfer",
			Handler:    _Msg_NftTransfer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tibc/apps/nft_transfer/v1/tx.proto",
}

func (m *MsgNftTransfer) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgNftTransfer) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgNftTransfer) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.DestContract) > 0 {
		i -= len(m.DestContract)
		copy(dAtA[i:], m.DestContract)
		i = encodeVarintTx(dAtA, i, uint64(len(m.DestContract)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.RealayChain) > 0 {
		i -= len(m.RealayChain)
		copy(dAtA[i:], m.RealayChain)
		i = encodeVarintTx(dAtA, i, uint64(len(m.RealayChain)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.DestChain) > 0 {
		i -= len(m.DestChain)
		copy(dAtA[i:], m.DestChain)
		i = encodeVarintTx(dAtA, i, uint64(len(m.DestChain)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Receiver) > 0 {
		i -= len(m.Receiver)
		copy(dAtA[i:], m.Receiver)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Receiver)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Class) > 0 {
		i -= len(m.Class)
		copy(dAtA[i:], m.Class)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Class)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgNftTransferResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgNftTransferResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgNftTransferResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgNftTransfer) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Class)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Receiver)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.DestChain)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.RealayChain)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.DestContract)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgNftTransferResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgNftTransfer) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgNftTransfer: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgNftTransfer: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Class", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Class = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Receiver", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Receiver = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DestChain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DestChain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RealayChain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RealayChain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DestContract", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DestContract = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgNftTransferResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgNftTransferResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgNftTransferResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
