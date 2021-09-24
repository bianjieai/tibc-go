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
	_ sdk.Msg = &MsgUpdateClient{}

	_ codectypes.UnpackInterfacesMessage = MsgUpdateClient{}
)

// ValidateBasic implements sdk.Msg
// func (msg MsgCreateClient) ValidateBasic() error {
// 	_, err := sdk.AccAddressFromBech32(msg.Signer)
// 	if err != nil {
// 		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
// 	}
// 	clientState, err := UnpackClientState(msg.ClientState)
// 	if err != nil {
// 		return err
// 	}
// 	if err := clientState.Validate(); err != nil {
// 		return err
// 	}
// 	if clientState.ClientType() == exported.Localhost {
// 		return sdkerrors.Wrap(ErrInvalidClient, "localhost client can only be created on chain initialization")
// 	}
// 	consensusState, err := UnpackConsensusState(msg.ConsensusState)
// 	if err != nil {
// 		return err
// 	}
// 	if clientState.ClientType() != consensusState.ClientType() {
// 		return sdkerrors.Wrap(ErrInvalidClientType, "client type for client state and consensus state do not match")
// 	}
// 	if err := ValidateClientType(clientState.ClientType()); err != nil {
// 		return sdkerrors.Wrap(err, "client type does not meet naming constraints")
// 	}
// 	return consensusState.ValidateBasic()
// }

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
// func (msg MsgCreateClient) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
// 	var clientState exported.ClientState
// 	err := unpacker.UnpackAny(msg.ClientState, &clientState)
// 	if err != nil {
// 		return err
// 	}

// 	var consensusState exported.ConsensusState
// 	return unpacker.UnpackAny(msg.ConsensusState, &consensusState)
// }

// NewMsgUpdateClient creates a new MsgUpdateClient instance
//nolint:interfacer
func NewMsgUpdateClient(chainName string, header exported.Header, signer sdk.AccAddress) (*MsgUpdateClient, error) {
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
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
