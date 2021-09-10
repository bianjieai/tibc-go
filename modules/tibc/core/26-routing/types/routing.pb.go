// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: tibc/core/routing/v1/routing.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

// SetRoutingRulesProposal defines a proposal to set routing rules
type SetRoutingRulesProposal struct {
	// the title of the update proposal
	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	// the description of the proposal
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// routing rules list
	Rules []string `protobuf:"bytes,3,rep,name=rules,proto3" json:"rules,omitempty"`
}

func (m *SetRoutingRulesProposal) Reset()         { *m = SetRoutingRulesProposal{} }
func (m *SetRoutingRulesProposal) String() string { return proto.CompactTextString(m) }
func (*SetRoutingRulesProposal) ProtoMessage()    {}
func (*SetRoutingRulesProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_e62ff9a6c7c9021e, []int{0}
}
func (m *SetRoutingRulesProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SetRoutingRulesProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SetRoutingRulesProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SetRoutingRulesProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetRoutingRulesProposal.Merge(m, src)
}
func (m *SetRoutingRulesProposal) XXX_Size() int {
	return m.Size()
}
func (m *SetRoutingRulesProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_SetRoutingRulesProposal.DiscardUnknown(m)
}

var xxx_messageInfo_SetRoutingRulesProposal proto.InternalMessageInfo

func init() {
	proto.RegisterType((*SetRoutingRulesProposal)(nil), "tibc.core.routing.v1.SetRoutingRulesProposal")
}

func init() {
	proto.RegisterFile("tibc/core/routing/v1/routing.proto", fileDescriptor_e62ff9a6c7c9021e)
}

var fileDescriptor_e62ff9a6c7c9021e = []byte{
	// 240 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2a, 0xc9, 0x4c, 0x4a,
	0xd6, 0x4f, 0xce, 0x2f, 0x4a, 0xd5, 0x2f, 0xca, 0x2f, 0x2d, 0xc9, 0xcc, 0x4b, 0xd7, 0x2f, 0x33,
	0x84, 0x31, 0xf5, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0x44, 0x40, 0x6a, 0xf4, 0x40, 0x6a, 0xf4,
	0x60, 0x12, 0x65, 0x86, 0x52, 0x22, 0xe9, 0xf9, 0xe9, 0xf9, 0x60, 0x05, 0xfa, 0x20, 0x16, 0x44,
	0xad, 0x52, 0x2e, 0x97, 0x78, 0x70, 0x6a, 0x49, 0x10, 0x44, 0x59, 0x50, 0x69, 0x4e, 0x6a, 0x71,
	0x40, 0x51, 0x7e, 0x41, 0x7e, 0x71, 0x62, 0x8e, 0x90, 0x08, 0x17, 0x6b, 0x49, 0x66, 0x49, 0x4e,
	0xaa, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x84, 0x23, 0xa4, 0xc0, 0xc5, 0x9d, 0x92, 0x5a,
	0x9c, 0x5c, 0x94, 0x59, 0x50, 0x92, 0x99, 0x9f, 0x27, 0xc1, 0x04, 0x96, 0x43, 0x16, 0x02, 0xe9,
	0x2b, 0x02, 0x19, 0x24, 0xc1, 0xac, 0xc0, 0x0c, 0xd2, 0x07, 0xe6, 0x58, 0xb1, 0x74, 0x2c, 0x90,
	0x67, 0x70, 0x8a, 0x3c, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18,
	0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1, 0xc6, 0x63, 0x39, 0x86, 0x28, 0xfb, 0xf4,
	0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0xfd, 0xa4, 0xcc, 0xc4, 0xbc, 0xac, 0xcc,
	0xd4, 0xc4, 0x4c, 0x7d, 0x90, 0x4f, 0x74, 0xd3, 0xf3, 0xf5, 0x73, 0xf3, 0x53, 0x40, 0xa6, 0xe8,
	0x23, 0x7c, 0x6f, 0x64, 0xa6, 0x0b, 0x0b, 0x80, 0x92, 0xca, 0x82, 0xd4, 0xe2, 0x24, 0x36, 0xb0,
	0x87, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xd5, 0xce, 0x09, 0xc0, 0x22, 0x01, 0x00, 0x00,
}

func (m *SetRoutingRulesProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SetRoutingRulesProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SetRoutingRulesProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Rules) > 0 {
		for iNdEx := len(m.Rules) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Rules[iNdEx])
			copy(dAtA[i:], m.Rules[iNdEx])
			i = encodeVarintRouting(dAtA, i, uint64(len(m.Rules[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintRouting(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintRouting(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintRouting(dAtA []byte, offset int, v uint64) int {
	offset -= sovRouting(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SetRoutingRulesProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovRouting(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovRouting(uint64(l))
	}
	if len(m.Rules) > 0 {
		for _, s := range m.Rules {
			l = len(s)
			n += 1 + l + sovRouting(uint64(l))
		}
	}
	return n
}

func sovRouting(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozRouting(x uint64) (n int) {
	return sovRouting(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SetRoutingRulesProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRouting
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
			return fmt.Errorf("proto: SetRoutingRulesProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SetRoutingRulesProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRouting
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
				return ErrInvalidLengthRouting
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRouting
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRouting
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
				return ErrInvalidLengthRouting
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRouting
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rules", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRouting
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
				return ErrInvalidLengthRouting
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRouting
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Rules = append(m.Rules, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRouting(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRouting
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
func skipRouting(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRouting
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
					return 0, ErrIntOverflowRouting
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
					return 0, ErrIntOverflowRouting
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
				return 0, ErrInvalidLengthRouting
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupRouting
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthRouting
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthRouting        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRouting          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupRouting = fmt.Errorf("proto: unexpected end of group")
)
