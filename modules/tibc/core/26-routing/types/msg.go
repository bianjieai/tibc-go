package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgSetRoutingRules{}
)

// ValidateBasic implements sdk.Msg
func (msg MsgSetRoutingRules) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"string could not be parsed as address: %v",
			err,
		)
	}
	context := SetRoutingRulesProposal{
		Title:       msg.Title,
		Description: msg.Description,
		Rules:       msg.Rules,
	}
	return context.ValidateBasic()
}

// GetSigners implements sdk.Msg
func (msg MsgSetRoutingRules) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{accAddr}
}
