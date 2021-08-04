package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgNftTransfer    = "nft_transfer"
)


// NewMsgNftTransfer creates a new NewMsgNftTransfer instance
func NewMsgNftTransfer(
	class, id, uri, sender, receiver string,
	awayFromOrigin bool,
	destChain, realayChain string,
) *MsgNftTransfer {
	return &MsgNftTransfer{
		Class:       class,
		Id:    id,
		Uri:            uri,
		Sender:           sender,
		Receiver:         receiver,
		AwayFromOrigin:    awayFromOrigin,
		DestChain: destChain,
		RealayChain: realayChain,
	}
}


// Route Implements Msg
func (msg MsgNftTransfer) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgNftTransfer) Type() string { return  TypeMsgNftTransfer}

// GetSigners Implements Msg.
func (msg MsgNftTransfer) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic Implements Msg.
func (msg MsgNftTransfer) ValidateBasic() error {
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgNftTransfer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}