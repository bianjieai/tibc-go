package types

import (
	"encoding/base64"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	commitmenttypes "github.com/bianjieai/tibc-go/modules/tibc/core/23-commitment/types"
)

var _ sdk.Msg = &MsgRecvPacket{}
var _ sdk.Msg = &MsgAcknowledgement{}
var _ sdk.Msg = &MsgRecvCleanPacket{}

// NewMsgRecvPacket constructs new MsgRecvPacket
// nolint:interfacer
func NewMsgRecvPacket(
	packet Packet, proofCommitment []byte, proofHeight clienttypes.Height, signer sdk.AccAddress,
) *MsgRecvPacket {
	return &MsgRecvPacket{
		Packet:          packet,
		ProofCommitment: proofCommitment,
		ProofHeight:     proofHeight,
		Signer:          signer.String(),
	}
}

// ValidateBasic implements sdk.Msg
func (msg MsgRecvPacket) ValidateBasic() error {
	if len(msg.ProofCommitment) == 0 {
		return sdkerrors.Wrap(commitmenttypes.ErrInvalidProof, "cannot submit an empty proof")
	}
	if msg.ProofHeight.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, "proof height must be non-zero")
	}
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	}
	return msg.Packet.ValidateBasic()
}

// GetDataSignBytes returns the base64-encoded bytes used for the
// data field when signing the packet.
func (msg MsgRecvPacket) GetDataSignBytes() []byte {
	s := "\"" + base64.StdEncoding.EncodeToString(msg.Packet.Data) + "\""
	return []byte(s)
}

// GetSigners implements sdk.Msg
func (msg MsgRecvPacket) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// NewMsgAcknowledgement constructs a new MsgAcknowledgement
// nolint:interfacer
func NewMsgAcknowledgement(
	packet Packet,
	ack, proofAcked []byte,
	proofHeight clienttypes.Height,
	signer sdk.AccAddress,
) *MsgAcknowledgement {
	return &MsgAcknowledgement{
		Packet:          packet,
		Acknowledgement: ack,
		ProofAcked:      proofAcked,
		ProofHeight:     proofHeight,
		Signer:          signer.String(),
	}
}

// ValidateBasic implements sdk.Msg
func (msg MsgAcknowledgement) ValidateBasic() error {
	if len(msg.ProofAcked) == 0 {
		return sdkerrors.Wrap(commitmenttypes.ErrInvalidProof, "cannot submit an empty proof")
	}
	if msg.ProofHeight.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, "proof height must be non-zero")
	}
	if len(msg.Acknowledgement) == 0 {
		return sdkerrors.Wrap(ErrInvalidAcknowledgement, "ack bytes cannot be empty")
	}
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	}
	return msg.Packet.ValidateBasic()
}

// GetSigners implements sdk.Msg
func (msg MsgAcknowledgement) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

var _ sdk.Msg = &MsgCleanPacket{}

// NewMsgCleanPacket constructs new MsgCleanPacket
// nolint:interfacer
func NewMsgCleanPacket(packet CleanPacket, signer sdk.AccAddress) *MsgCleanPacket {
	return &MsgCleanPacket{
		CleanPacket: packet,
		Signer:      signer.String(),
	}
}

// ValidateBasic implements sdk.Msg
func (msg MsgCleanPacket) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	}
	return msg.CleanPacket.ValidateBasic()
}

// GetSigners implements sdk.Msg
func (msg MsgCleanPacket) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// NewMsgCleanPacket constructs new MsgCleanPacket
// nolint:interfacer
func NewMsgRecvCleanPacket(
	cleanPacket CleanPacket,
	proofCommitment []byte,
	proofHeight clienttypes.Height,
	signer sdk.AccAddress,
) *MsgRecvCleanPacket {
	return &MsgRecvCleanPacket{
		CleanPacket:     cleanPacket,
		ProofCommitment: proofCommitment,
		ProofHeight:     proofHeight,
		Signer:          signer.String(),
	}
}

// ValidateBasic implements sdk.Msg
func (msg MsgRecvCleanPacket) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	}
	return msg.CleanPacket.ValidateBasic()
}

// GetSigners implements sdk.Msg
func (msg MsgRecvCleanPacket) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}
