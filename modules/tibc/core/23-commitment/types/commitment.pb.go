// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: tibc/core/commitment/v1/commitment.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	_go "github.com/cosmos/ics23/go"
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

// MerkleRoot defines a merkle root hash.
// In the Cosmos SDK, the AppHash of a block header becomes the root.
type MerkleRoot struct {
	Hash []byte `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (m *MerkleRoot) Reset()         { *m = MerkleRoot{} }
func (m *MerkleRoot) String() string { return proto.CompactTextString(m) }
func (*MerkleRoot) ProtoMessage()    {}
func (*MerkleRoot) Descriptor() ([]byte, []int) {
	return fileDescriptor_e76044a26d8b9539, []int{0}
}
func (m *MerkleRoot) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MerkleRoot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MerkleRoot.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MerkleRoot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MerkleRoot.Merge(m, src)
}
func (m *MerkleRoot) XXX_Size() int {
	return m.Size()
}
func (m *MerkleRoot) XXX_DiscardUnknown() {
	xxx_messageInfo_MerkleRoot.DiscardUnknown(m)
}

var xxx_messageInfo_MerkleRoot proto.InternalMessageInfo

// MerklePrefix is merkle path prefixed to the key.
// The constructed key from the Path and the key will be append(Path.KeyPath,
// append(Path.KeyPrefix, key...))
type MerklePrefix struct {
	KeyPrefix []byte `protobuf:"bytes,1,opt,name=key_prefix,json=keyPrefix,proto3" json:"key_prefix,omitempty" yaml:"key_prefix"`
}

func (m *MerklePrefix) Reset()         { *m = MerklePrefix{} }
func (m *MerklePrefix) String() string { return proto.CompactTextString(m) }
func (*MerklePrefix) ProtoMessage()    {}
func (*MerklePrefix) Descriptor() ([]byte, []int) {
	return fileDescriptor_e76044a26d8b9539, []int{1}
}
func (m *MerklePrefix) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MerklePrefix) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MerklePrefix.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MerklePrefix) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MerklePrefix.Merge(m, src)
}
func (m *MerklePrefix) XXX_Size() int {
	return m.Size()
}
func (m *MerklePrefix) XXX_DiscardUnknown() {
	xxx_messageInfo_MerklePrefix.DiscardUnknown(m)
}

var xxx_messageInfo_MerklePrefix proto.InternalMessageInfo

func (m *MerklePrefix) GetKeyPrefix() []byte {
	if m != nil {
		return m.KeyPrefix
	}
	return nil
}

// MerklePath is the path used to verify commitment proofs, which can be an
// arbitrary structured object (defined by a commitment type).
// MerklePath is represented from root-to-leaf
type MerklePath struct {
	KeyPath []string `protobuf:"bytes,1,rep,name=key_path,json=keyPath,proto3" json:"key_path,omitempty" yaml:"key_path"`
}

func (m *MerklePath) Reset()      { *m = MerklePath{} }
func (*MerklePath) ProtoMessage() {}
func (*MerklePath) Descriptor() ([]byte, []int) {
	return fileDescriptor_e76044a26d8b9539, []int{2}
}
func (m *MerklePath) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MerklePath) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MerklePath.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MerklePath) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MerklePath.Merge(m, src)
}
func (m *MerklePath) XXX_Size() int {
	return m.Size()
}
func (m *MerklePath) XXX_DiscardUnknown() {
	xxx_messageInfo_MerklePath.DiscardUnknown(m)
}

var xxx_messageInfo_MerklePath proto.InternalMessageInfo

func (m *MerklePath) GetKeyPath() []string {
	if m != nil {
		return m.KeyPath
	}
	return nil
}

// MerkleProof is a wrapper type over a chain of CommitmentProofs.
// It demonstrates membership or non-membership for an element or set of
// elements, verifiable in conjunction with a known commitment root. Proofs
// should be succinct.
// MerkleProofs are ordered from leaf-to-root
type MerkleProof struct {
	Proofs []*_go.CommitmentProof `protobuf:"bytes,1,rep,name=proofs,proto3" json:"proofs,omitempty"`
}

func (m *MerkleProof) Reset()         { *m = MerkleProof{} }
func (m *MerkleProof) String() string { return proto.CompactTextString(m) }
func (*MerkleProof) ProtoMessage()    {}
func (*MerkleProof) Descriptor() ([]byte, []int) {
	return fileDescriptor_e76044a26d8b9539, []int{3}
}
func (m *MerkleProof) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MerkleProof) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MerkleProof.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MerkleProof) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MerkleProof.Merge(m, src)
}
func (m *MerkleProof) XXX_Size() int {
	return m.Size()
}
func (m *MerkleProof) XXX_DiscardUnknown() {
	xxx_messageInfo_MerkleProof.DiscardUnknown(m)
}

var xxx_messageInfo_MerkleProof proto.InternalMessageInfo

func (m *MerkleProof) GetProofs() []*_go.CommitmentProof {
	if m != nil {
		return m.Proofs
	}
	return nil
}

func init() {
	proto.RegisterType((*MerkleRoot)(nil), "tibc.core.commitment.v1.MerkleRoot")
	proto.RegisterType((*MerklePrefix)(nil), "tibc.core.commitment.v1.MerklePrefix")
	proto.RegisterType((*MerklePath)(nil), "tibc.core.commitment.v1.MerklePath")
	proto.RegisterType((*MerkleProof)(nil), "tibc.core.commitment.v1.MerkleProof")
}

func init() {
	proto.RegisterFile("tibc/core/commitment/v1/commitment.proto", fileDescriptor_e76044a26d8b9539)
}

var fileDescriptor_e76044a26d8b9539 = []byte{
	// 350 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x91, 0xcf, 0x4e, 0xea, 0x40,
	0x14, 0xc6, 0xdb, 0x5c, 0xc2, 0xbd, 0x0c, 0x24, 0x37, 0x56, 0x8d, 0x86, 0x98, 0x42, 0xba, 0x30,
	0x6c, 0x98, 0x09, 0xe0, 0xc2, 0xb0, 0xac, 0x26, 0xae, 0x4c, 0x48, 0x97, 0xc6, 0xc4, 0x4c, 0xeb,
	0xd0, 0x8e, 0x50, 0x4e, 0xd3, 0x19, 0x88, 0x7d, 0x03, 0x97, 0x2e, 0x5d, 0xfa, 0x38, 0x2e, 0x59,
	0xba, 0x22, 0x06, 0xde, 0x80, 0x27, 0x30, 0x33, 0x23, 0xd2, 0xdd, 0x39, 0x39, 0xbf, 0xf3, 0xe7,
	0xfb, 0x0e, 0xea, 0x48, 0x1e, 0x46, 0x24, 0x82, 0x9c, 0x91, 0x08, 0xd2, 0x94, 0xcb, 0x94, 0xcd,
	0x24, 0x59, 0xf4, 0x4a, 0x19, 0xce, 0x72, 0x90, 0xe0, 0x9c, 0x28, 0x12, 0x2b, 0x12, 0x97, 0x6a,
	0x8b, 0x5e, 0xf3, 0x28, 0x86, 0x18, 0x34, 0x43, 0x54, 0x64, 0xf0, 0xe6, 0x59, 0x04, 0x22, 0x05,
	0x41, 0x78, 0x24, 0xfa, 0x03, 0x35, 0x30, 0xcb, 0x01, 0xc6, 0xc2, 0x54, 0xbd, 0x73, 0x84, 0x6e,
	0x59, 0x3e, 0x99, 0xb2, 0x00, 0x40, 0x3a, 0x0e, 0xaa, 0x24, 0x54, 0x24, 0xa7, 0x76, 0xdb, 0xee,
	0x34, 0x02, 0x1d, 0x0f, 0x2b, 0x2f, 0xef, 0x2d, 0xcb, 0xbb, 0x46, 0x0d, 0xc3, 0x8d, 0x72, 0x36,
	0xe6, 0xcf, 0xce, 0x05, 0x42, 0x13, 0x56, 0x3c, 0x64, 0x3a, 0x33, 0xbc, 0x7f, 0xbc, 0x5d, 0xb5,
	0x0e, 0x0a, 0x9a, 0x4e, 0x87, 0xde, 0xbe, 0xe6, 0x05, 0xb5, 0x09, 0x2b, 0x4c, 0x97, 0xe7, 0xef,
	0xb6, 0x8d, 0xa8, 0x4c, 0x1c, 0x8c, 0xfe, 0x69, 0x8e, 0x4a, 0xb5, 0xf1, 0x4f, 0xa7, 0xe6, 0x1f,
	0x6e, 0x57, 0xad, 0xff, 0xa5, 0x09, 0x54, 0x26, 0x5e, 0xf0, 0x57, 0xf5, 0x53, 0x99, 0x0c, 0x2b,
	0x6f, 0xea, 0x92, 0x1b, 0x54, 0xdf, 0x5d, 0x02, 0x30, 0x76, 0x2e, 0x51, 0xd5, 0x08, 0xd2, 0x23,
	0xea, 0xfd, 0x36, 0x36, 0x7a, 0xb1, 0xd6, 0x8b, 0x17, 0x3d, 0x7c, 0xf5, 0x6b, 0x92, 0xee, 0x08,
	0x7e, 0x78, 0xff, 0xfe, 0x63, 0xed, 0xda, 0xcb, 0xb5, 0x6b, 0x7f, 0xad, 0x5d, 0xfb, 0x75, 0xe3,
	0x5a, 0xcb, 0x8d, 0x6b, 0x7d, 0x6e, 0x5c, 0xeb, 0xce, 0x8f, 0xb9, 0x4c, 0xe6, 0xa1, 0xb2, 0x97,
	0x84, 0x9c, 0xce, 0x9e, 0x38, 0xa3, 0x9c, 0x28, 0xdb, 0xbb, 0x31, 0x90, 0x14, 0x1e, 0xe7, 0x53,
	0x26, 0xc8, 0xfe, 0x61, 0xfd, 0x41, 0xb7, 0xf4, 0x33, 0x59, 0x64, 0x4c, 0x84, 0x55, 0xed, 0xef,
	0xe0, 0x3b, 0x00, 0x00, 0xff, 0xff, 0xd3, 0xa7, 0x04, 0x13, 0xd8, 0x01, 0x00, 0x00,
}

func (m *MerkleRoot) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MerkleRoot) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MerkleRoot) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintCommitment(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MerklePrefix) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MerklePrefix) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MerklePrefix) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.KeyPrefix) > 0 {
		i -= len(m.KeyPrefix)
		copy(dAtA[i:], m.KeyPrefix)
		i = encodeVarintCommitment(dAtA, i, uint64(len(m.KeyPrefix)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MerklePath) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MerklePath) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MerklePath) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.KeyPath) > 0 {
		for iNdEx := len(m.KeyPath) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.KeyPath[iNdEx])
			copy(dAtA[i:], m.KeyPath[iNdEx])
			i = encodeVarintCommitment(dAtA, i, uint64(len(m.KeyPath[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *MerkleProof) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MerkleProof) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MerkleProof) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Proofs) > 0 {
		for iNdEx := len(m.Proofs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Proofs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintCommitment(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintCommitment(dAtA []byte, offset int, v uint64) int {
	offset -= sovCommitment(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MerkleRoot) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovCommitment(uint64(l))
	}
	return n
}

func (m *MerklePrefix) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.KeyPrefix)
	if l > 0 {
		n += 1 + l + sovCommitment(uint64(l))
	}
	return n
}

func (m *MerklePath) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.KeyPath) > 0 {
		for _, s := range m.KeyPath {
			l = len(s)
			n += 1 + l + sovCommitment(uint64(l))
		}
	}
	return n
}

func (m *MerkleProof) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Proofs) > 0 {
		for _, e := range m.Proofs {
			l = e.Size()
			n += 1 + l + sovCommitment(uint64(l))
		}
	}
	return n
}

func sovCommitment(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCommitment(x uint64) (n int) {
	return sovCommitment(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MerkleRoot) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommitment
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
			return fmt.Errorf("proto: MerkleRoot: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MerkleRoot: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommitment
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthCommitment
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthCommitment
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = append(m.Hash[:0], dAtA[iNdEx:postIndex]...)
			if m.Hash == nil {
				m.Hash = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommitment(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommitment
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
func (m *MerklePrefix) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommitment
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
			return fmt.Errorf("proto: MerklePrefix: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MerklePrefix: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field KeyPrefix", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommitment
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthCommitment
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthCommitment
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.KeyPrefix = append(m.KeyPrefix[:0], dAtA[iNdEx:postIndex]...)
			if m.KeyPrefix == nil {
				m.KeyPrefix = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommitment(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommitment
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
func (m *MerklePath) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommitment
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
			return fmt.Errorf("proto: MerklePath: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MerklePath: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field KeyPath", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommitment
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
				return ErrInvalidLengthCommitment
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCommitment
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.KeyPath = append(m.KeyPath, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommitment(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommitment
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
func (m *MerkleProof) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCommitment
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
			return fmt.Errorf("proto: MerkleProof: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MerkleProof: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Proofs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCommitment
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
				return ErrInvalidLengthCommitment
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCommitment
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Proofs = append(m.Proofs, &_go.CommitmentProof{})
			if err := m.Proofs[len(m.Proofs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCommitment(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCommitment
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
func skipCommitment(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCommitment
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
					return 0, ErrIntOverflowCommitment
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
					return 0, ErrIntOverflowCommitment
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
				return 0, ErrInvalidLengthCommitment
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupCommitment
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthCommitment
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthCommitment        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCommitment          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupCommitment = fmt.Errorf("proto: unexpected end of group")
)
