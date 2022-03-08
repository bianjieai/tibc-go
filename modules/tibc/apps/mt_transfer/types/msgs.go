package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeMsgMtTransfer = "tibc_mt_transfer"
)

var _ sdk.Msg = &MsgMtTransfer{}

// NewMsgMtTransfer creates a new NewMsgMtTransfer instance
func NewMsgMtTransfer(class, id, sender, receiver, destChain, realayChain, destContract string, amount uint64) *MsgMtTransfer {
	return &MsgMtTransfer{
		Class:        class,
		Id:           id,
		Sender:       sender,
		Receiver:     receiver,
		DestChain:    destChain,
		RealayChain:  realayChain,
		DestContract: destContract,
		Amount:       amount,
	}
}

// GetSigners implements sdk.Msg
func (msg MsgMtTransfer) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

// ValidateBasic Implements Msg.
func (msg MsgMtTransfer) ValidateBasic() error {
	return nil
}
