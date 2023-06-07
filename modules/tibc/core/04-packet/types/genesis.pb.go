// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: tibc/core/packet/v1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

// GenesisState defines the tibc channel submodule's genesis state.
type GenesisState struct {
	Acknowledgements []PacketState    `protobuf:"bytes,2,rep,name=acknowledgements,proto3" json:"acknowledgements"`
	Commitments      []PacketState    `protobuf:"bytes,3,rep,name=commitments,proto3" json:"commitments"`
	Receipts         []PacketState    `protobuf:"bytes,4,rep,name=receipts,proto3" json:"receipts"`
	SendSequences    []PacketSequence `protobuf:"bytes,5,rep,name=send_sequences,json=sendSequences,proto3" json:"send_sequences" yaml:"send_sequences"`
	RecvSequences    []PacketSequence `protobuf:"bytes,6,rep,name=recv_sequences,json=recvSequences,proto3" json:"recv_sequences" yaml:"recv_sequences"`
	AckSequences     []PacketSequence `protobuf:"bytes,7,rep,name=ack_sequences,json=ackSequences,proto3" json:"ack_sequences" yaml:"ack_sequences"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_0ac7d6194d8bf289, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetAcknowledgements() []PacketState {
	if m != nil {
		return m.Acknowledgements
	}
	return nil
}

func (m *GenesisState) GetCommitments() []PacketState {
	if m != nil {
		return m.Commitments
	}
	return nil
}

func (m *GenesisState) GetReceipts() []PacketState {
	if m != nil {
		return m.Receipts
	}
	return nil
}

func (m *GenesisState) GetSendSequences() []PacketSequence {
	if m != nil {
		return m.SendSequences
	}
	return nil
}

func (m *GenesisState) GetRecvSequences() []PacketSequence {
	if m != nil {
		return m.RecvSequences
	}
	return nil
}

func (m *GenesisState) GetAckSequences() []PacketSequence {
	if m != nil {
		return m.AckSequences
	}
	return nil
}

// PacketSequence defines the genesis type necessary to retrieve and store
// next send and receive sequences.
type PacketSequence struct {
	SourceChain      string `protobuf:"bytes,1,opt,name=source_chain,json=sourceChain,proto3" json:"source_chain,omitempty" yaml:"source_chain"`
	DestinationChain string `protobuf:"bytes,2,opt,name=destination_chain,json=destinationChain,proto3" json:"destination_chain,omitempty" yaml:"destination_chain"`
	Sequence         uint64 `protobuf:"varint,3,opt,name=sequence,proto3" json:"sequence,omitempty"`
}

func (m *PacketSequence) Reset()         { *m = PacketSequence{} }
func (m *PacketSequence) String() string { return proto.CompactTextString(m) }
func (*PacketSequence) ProtoMessage()    {}
func (*PacketSequence) Descriptor() ([]byte, []int) {
	return fileDescriptor_0ac7d6194d8bf289, []int{1}
}
func (m *PacketSequence) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PacketSequence) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PacketSequence.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PacketSequence) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PacketSequence.Merge(m, src)
}
func (m *PacketSequence) XXX_Size() int {
	return m.Size()
}
func (m *PacketSequence) XXX_DiscardUnknown() {
	xxx_messageInfo_PacketSequence.DiscardUnknown(m)
}

var xxx_messageInfo_PacketSequence proto.InternalMessageInfo

func (m *PacketSequence) GetSourceChain() string {
	if m != nil {
		return m.SourceChain
	}
	return ""
}

func (m *PacketSequence) GetDestinationChain() string {
	if m != nil {
		return m.DestinationChain
	}
	return ""
}

func (m *PacketSequence) GetSequence() uint64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "tibc.core.packet.v1.GenesisState")
	proto.RegisterType((*PacketSequence)(nil), "tibc.core.packet.v1.PacketSequence")
}

func init() { proto.RegisterFile("tibc/core/packet/v1/genesis.proto", fileDescriptor_0ac7d6194d8bf289) }

var fileDescriptor_0ac7d6194d8bf289 = []byte{
	// 454 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0xc1, 0x6e, 0xd3, 0x30,
	0x18, 0xc7, 0x9b, 0xb5, 0x8c, 0xe1, 0xb6, 0xd3, 0xc8, 0x86, 0x88, 0xaa, 0x91, 0x86, 0x70, 0xe9,
	0x65, 0x31, 0x03, 0x4e, 0x3b, 0x70, 0x08, 0x07, 0xe0, 0x86, 0xcc, 0x05, 0x71, 0x99, 0x5c, 0xf7,
	0x23, 0x33, 0x6d, 0xec, 0x12, 0xbb, 0x45, 0x7b, 0x0b, 0x9e, 0x84, 0x0b, 0x2f, 0xb1, 0xe3, 0x8e,
	0x9c, 0x2a, 0xd4, 0xbe, 0x41, 0x9f, 0x00, 0xd9, 0xce, 0xba, 0x54, 0x9b, 0x10, 0xdd, 0x2d, 0xfe,
	0xfc, 0xff, 0xff, 0x7e, 0x91, 0xa5, 0x0f, 0x3d, 0xd5, 0xbc, 0xcf, 0x30, 0x93, 0x05, 0xe0, 0x31,
	0x65, 0x43, 0xd0, 0x78, 0x7a, 0x8c, 0x33, 0x10, 0xa0, 0xb8, 0x4a, 0xc6, 0x85, 0xd4, 0xd2, 0xdf,
	0x37, 0x91, 0xc4, 0x44, 0x12, 0x17, 0x49, 0xa6, 0xc7, 0x9d, 0x83, 0x4c, 0x66, 0xd2, 0xde, 0x63,
	0xf3, 0xe5, 0xa2, 0x9d, 0xe8, 0x36, 0x5a, 0x59, 0xb2, 0x89, 0xf8, 0x67, 0x03, 0xb5, 0xde, 0x3a,
	0xfc, 0x47, 0x4d, 0x35, 0xf8, 0x04, 0xed, 0x51, 0x36, 0x14, 0xf2, 0xfb, 0x08, 0x06, 0x19, 0xe4,
	0x20, 0xb4, 0x0a, 0xb6, 0xa2, 0x7a, 0xaf, 0xf9, 0x22, 0x4a, 0x6e, 0x11, 0x27, 0x1f, 0xec, 0x97,
	0xed, 0xa6, 0x8d, 0x8b, 0x59, 0xb7, 0x46, 0x6e, 0xf4, 0xfd, 0x77, 0xa8, 0xc9, 0x64, 0x9e, 0x73,
	0xed, 0x70, 0xf5, 0x8d, 0x70, 0xd5, 0xaa, 0x9f, 0xa2, 0x9d, 0x02, 0x18, 0xf0, 0xb1, 0x56, 0x41,
	0x63, 0x23, 0xcc, 0xaa, 0xe7, 0x73, 0xb4, 0xab, 0x40, 0x0c, 0x4e, 0x15, 0x7c, 0x9b, 0x80, 0x60,
	0xa0, 0x82, 0x7b, 0x96, 0xf4, 0xec, 0x5f, 0xa4, 0x32, 0x9b, 0x3e, 0x31, 0xb0, 0xe5, 0xac, 0xfb,
	0xe8, 0x9c, 0xe6, 0xa3, 0x93, 0x78, 0x1d, 0x14, 0x93, 0xb6, 0x19, 0x5c, 0x85, 0xad, 0xaa, 0x00,
	0x36, 0xad, 0xa8, 0xb6, 0xef, 0xac, 0x5a, 0x07, 0xc5, 0xa4, 0x6d, 0x06, 0xd7, 0xaa, 0x2f, 0xa8,
	0x4d, 0xd9, 0xb0, 0x62, 0xba, 0xff, 0xff, 0xa6, 0xc3, 0xd2, 0x74, 0xe0, 0x4c, 0x6b, 0x9c, 0x98,
	0xb4, 0x28, 0x1b, 0xae, 0x3c, 0xf1, 0x2f, 0x0f, 0xed, 0xae, 0xd7, 0xfd, 0x13, 0xd4, 0x52, 0x72,
	0x52, 0x30, 0x38, 0x65, 0x67, 0x94, 0x8b, 0xc0, 0x8b, 0xbc, 0xde, 0x83, 0xf4, 0xf1, 0x72, 0xd6,
	0xdd, 0x2f, 0x5f, 0xa9, 0x72, 0x1b, 0x93, 0xa6, 0x3b, 0xbe, 0x31, 0x27, 0xff, 0x3d, 0x7a, 0x38,
	0x00, 0xa5, 0xb9, 0xa0, 0x9a, 0x4b, 0x51, 0x02, 0xb6, 0x2c, 0xe0, 0x70, 0x39, 0xeb, 0x06, 0x0e,
	0x70, 0x23, 0x12, 0x93, 0xbd, 0xca, 0xcc, 0xa1, 0x3a, 0x68, 0xe7, 0xea, 0xaf, 0x83, 0x7a, 0xe4,
	0xf5, 0x1a, 0x64, 0x75, 0x4e, 0x3f, 0x5d, 0xcc, 0x43, 0xef, 0x72, 0x1e, 0x7a, 0x7f, 0xe6, 0xa1,
	0xf7, 0x63, 0x11, 0xd6, 0x2e, 0x17, 0x61, 0xed, 0xf7, 0x22, 0xac, 0x7d, 0x7e, 0x9d, 0x71, 0x7d,
	0x36, 0xe9, 0x27, 0x4c, 0xe6, 0xb8, 0xcf, 0xa9, 0xf8, 0xca, 0x81, 0x72, 0x6c, 0x1e, 0xed, 0x28,
	0x93, 0x38, 0x97, 0x83, 0xc9, 0x08, 0x14, 0xbe, 0xde, 0xa3, 0xe7, 0xaf, 0x8e, 0xca, 0x55, 0xd2,
	0xe7, 0x63, 0x50, 0xfd, 0x6d, 0xbb, 0x47, 0x2f, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0x23, 0x33,
	0x44, 0x2a, 0xb9, 0x03, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.AckSequences) > 0 {
		for iNdEx := len(m.AckSequences) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.AckSequences[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if len(m.RecvSequences) > 0 {
		for iNdEx := len(m.RecvSequences) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.RecvSequences[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if len(m.SendSequences) > 0 {
		for iNdEx := len(m.SendSequences) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SendSequences[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.Receipts) > 0 {
		for iNdEx := len(m.Receipts) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Receipts[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Commitments) > 0 {
		for iNdEx := len(m.Commitments) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Commitments[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Acknowledgements) > 0 {
		for iNdEx := len(m.Acknowledgements) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Acknowledgements[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	return len(dAtA) - i, nil
}

func (m *PacketSequence) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PacketSequence) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PacketSequence) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Sequence != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.Sequence))
		i--
		dAtA[i] = 0x18
	}
	if len(m.DestinationChain) > 0 {
		i -= len(m.DestinationChain)
		copy(dAtA[i:], m.DestinationChain)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.DestinationChain)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.SourceChain) > 0 {
		i -= len(m.SourceChain)
		copy(dAtA[i:], m.SourceChain)
		i = encodeVarintGenesis(dAtA, i, uint64(len(m.SourceChain)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Acknowledgements) > 0 {
		for _, e := range m.Acknowledgements {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Commitments) > 0 {
		for _, e := range m.Commitments {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Receipts) > 0 {
		for _, e := range m.Receipts {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.SendSequences) > 0 {
		for _, e := range m.SendSequences {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.RecvSequences) > 0 {
		for _, e := range m.RecvSequences {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.AckSequences) > 0 {
		for _, e := range m.AckSequences {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func (m *PacketSequence) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.SourceChain)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	l = len(m.DestinationChain)
	if l > 0 {
		n += 1 + l + sovGenesis(uint64(l))
	}
	if m.Sequence != 0 {
		n += 1 + sovGenesis(uint64(m.Sequence))
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Acknowledgements", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Acknowledgements = append(m.Acknowledgements, PacketState{})
			if err := m.Acknowledgements[len(m.Acknowledgements)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Commitments", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Commitments = append(m.Commitments, PacketState{})
			if err := m.Commitments[len(m.Commitments)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Receipts", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Receipts = append(m.Receipts, PacketState{})
			if err := m.Receipts[len(m.Receipts)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SendSequences", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SendSequences = append(m.SendSequences, PacketSequence{})
			if err := m.SendSequences[len(m.SendSequences)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RecvSequences", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RecvSequences = append(m.RecvSequences, PacketSequence{})
			if err := m.RecvSequences[len(m.RecvSequences)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AckSequences", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AckSequences = append(m.AckSequences, PacketSequence{})
			if err := m.AckSequences[len(m.AckSequences)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func (m *PacketSequence) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: PacketSequence: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PacketSequence: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SourceChain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SourceChain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DestinationChain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DestinationChain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sequence", wireType)
			}
			m.Sequence = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Sequence |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
