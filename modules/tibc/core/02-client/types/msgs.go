package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

// message types for the TIBC client
const (
	TypeMsgUpdateClient string = "update_client"
)

var (
	_ sdk.Msg = &MsgCreateClient{}
	_ sdk.Msg = &MsgUpdateClient{}
	_ sdk.Msg = &MsgUpgradeClient{}
	_ sdk.Msg = &MsgRegisterRelayer{}

	_ codectypes.UnpackInterfacesMessage = MsgUpdateClient{}
	_ codectypes.UnpackInterfacesMessage = MsgCreateClient{}
	_ codectypes.UnpackInterfacesMessage = MsgUpgradeClient{}
)

// NewMsgUpdateClient creates a new MsgUpdateClient instance
//
//nolint:interfacer
func NewMsgUpdateClient(
	chainName string,
	header exported.Header,
	signer sdk.AccAddress,
) (*MsgUpdateClient, error) {
	anyHeader, err := PackHeader(header)
	if err != nil {
		return nil, err
	}

	return &MsgUpdateClient{
		ChainName: chainName,
		Header:    anyHeader,
		Signer:    signer.String(),
	}, nil
}

// ValidateBasic implements sdk.Msg
func (msg MsgUpdateClient) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"string could not be parsed as address: %v",
			err,
		)
	}
	header, err := UnpackHeader(msg.Header)
	if err != nil {
		return err
	}
	if err := header.ValidateBasic(); err != nil {
		return err
	}
	return host.ClientIdentifierValidator(msg.ChainName)
}

// GetSigners implements sdk.Msg
func (msg MsgUpdateClient) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{accAddr}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgUpdateClient) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var header exported.Header
	return unpacker.UnpackAny(msg.Header, &header)
}

// ValidateBasic implements sdk.Msg
func (msg MsgCreateClient) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"string could not be parsed as address: %v",
			err,
		)
	}
	context := CreateClientProposal{
		Title:          msg.Title,
		Description:    msg.Description,
		ChainName:      msg.ChainName,
		ClientState:    msg.ClientState,
		ConsensusState: msg.ConsensusState,
	}
	return context.ValidateBasic()
}

// GetSigners implements sdk.Msg
func (msg MsgCreateClient) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{accAddr}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgCreateClient) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if err := unpacker.UnpackAny(msg.ClientState, new(exported.ClientState)); err != nil {
		return err
	}

	if err := unpacker.UnpackAny(msg.ConsensusState, new(exported.ConsensusState)); err != nil {
		return err
	}
	return nil
}

// ValidateBasic implements sdk.Msg
func (msg MsgUpgradeClient) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"string could not be parsed as address: %v",
			err,
		)
	}
	context := UpgradeClientProposal{
		Title:          msg.Title,
		Description:    msg.Description,
		ChainName:      msg.ChainName,
		ClientState:    msg.ClientState,
		ConsensusState: msg.ConsensusState,
	}
	return context.ValidateBasic()
}

// GetSigners implements sdk.Msg
func (msg MsgUpgradeClient) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{accAddr}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgUpgradeClient) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if err := unpacker.UnpackAny(msg.ClientState, new(exported.ClientState)); err != nil {
		return err
	}

	if err := unpacker.UnpackAny(msg.ConsensusState, new(exported.ConsensusState)); err != nil {
		return err
	}
	return nil
}

// ValidateBasic implements sdk.Msg
func (msg MsgRegisterRelayer) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidAddress,
			"string could not be parsed as address: %v",
			err,
		)
	}
	context := RegisterRelayerProposal{
		Title:       msg.Title,
		Description: msg.Description,
		ChainName:   msg.ChainName,
		Relayers:    msg.Relayers,
	}
	return context.ValidateBasic()
}

// GetSigners implements sdk.Msg
func (msg MsgRegisterRelayer) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{accAddr}
}
