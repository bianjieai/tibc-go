package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgNftTransfer = "tibc_nft_transfer"
)

var _ sdk.Msg = &MsgNftTransfer{}

// NewMsgNftTransfer creates a new NewMsgNftTransfer instance
func NewMsgNftTransfer(class, id, sender, receiver, destChain, realayChain, destContract string) *MsgNftTransfer {
	return &MsgNftTransfer{
		Class:        class,
		Id:           id,
		Sender:       sender,
		Receiver:     receiver,
		DestChain:    destChain,
		RealayChain:  realayChain,
		DestContract: destContract,
	}
}

// GetSigners implements sdk.Msg
func (msg MsgNftTransfer) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// ValidateBasic Implements Msg.
func (msg MsgNftTransfer) ValidateBasic() error {
	return nil
}
