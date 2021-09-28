package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

const (
	ProposalTypeClientCreate    = "CreateClient"
	ProposalTypeClientUpgrade   = "UpgradeClient"
	ProposalTypeRelayerRegister = "RegisterRelayer"
)

var (
	_ govtypes.Content = &CreateClientProposal{}
	_ govtypes.Content = &UpgradeClientProposal{}
	_ govtypes.Content = &RegisterRelayerProposal{}

	_ codectypes.UnpackInterfacesMessage = &CreateClientProposal{}
	_ codectypes.UnpackInterfacesMessage = &UpgradeClientProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeClientCreate)
	govtypes.RegisterProposalType(ProposalTypeClientUpgrade)
	govtypes.RegisterProposalType(ProposalTypeRelayerRegister)
}

// NewCreateClientProposal creates a new creating client proposal.
func NewCreateClientProposal(title, description, chainName string, clientState exported.ClientState, consensusState exported.ConsensusState) (*CreateClientProposal, error) {
	clientStateAny, err := PackClientState(clientState)
	if err != nil {
		return nil, err
	}

	consensusStateAny, err := PackConsensusState(consensusState)
	if err != nil {
		return nil, err
	}

	return &CreateClientProposal{
		Title:          title,
		Description:    description,
		ChainName:      chainName,
		ClientState:    clientStateAny,
		ConsensusState: consensusStateAny,
	}, nil
}

// GetTitle returns the title of a client update proposal.
func (cup *CreateClientProposal) GetTitle() string { return cup.Title }

// GetDescription returns the description of a client update proposal.
func (cup *CreateClientProposal) GetDescription() string { return cup.Description }

// ProposalRoute returns the routing key of a client update proposal.
func (cup *CreateClientProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a client update proposal.
func (cup *CreateClientProposal) ProposalType() string { return ProposalTypeClientCreate }

// ValidateBasic runs basic stateless validity checks
func (cup *CreateClientProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(cup)
	if err != nil {
		return err
	}

	if err := host.ClientIdentifierValidator(cup.ChainName); err != nil {
		return err
	}

	clientState, err := UnpackClientState(cup.ClientState)
	if err != nil {
		return err
	}
	return clientState.Validate()
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (cup CreateClientProposal) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if err := unpacker.UnpackAny(cup.ClientState, new(exported.ClientState)); err != nil {
		return err
	}

	if err := unpacker.UnpackAny(cup.ConsensusState, new(exported.ConsensusState)); err != nil {
		return err
	}
	return nil
}

// NewUpgradeClientProposal create a upgrade client proposal.
func NewUpgradeClientProposal(
	title, description, chainName string,
	clientState exported.ClientState,
	consensusState exported.ConsensusState,
) (
	*UpgradeClientProposal, error,
) {
	clientStateAny, err := PackClientState(clientState)
	if err != nil {
		return nil, err
	}

	consensusStateAny, err := PackConsensusState(consensusState)
	if err != nil {
		return nil, err
	}

	return &UpgradeClientProposal{
		Title:          title,
		Description:    description,
		ChainName:      chainName,
		ClientState:    clientStateAny,
		ConsensusState: consensusStateAny,
	}, nil
}

// GetTitle returns the title of a client upgrade proposal.
func (cup *UpgradeClientProposal) GetTitle() string { return cup.Title }

// GetDescription returns the description of a client upgrade proposal.
func (cup *UpgradeClientProposal) GetDescription() string { return cup.Description }

// ProposalRoute returns the routing key of a client upgrade proposal.
func (cup *UpgradeClientProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a client upgrade proposal.
func (cup *UpgradeClientProposal) ProposalType() string { return ProposalTypeClientUpgrade }

// ValidateBasic runs basic stateless validity checks
func (cup *UpgradeClientProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(cup)
	if err != nil {
		return err
	}

	if err := host.ClientIdentifierValidator(cup.ChainName); err != nil {
		return err
	}

	clientState, err := UnpackClientState(cup.ClientState)
	if err != nil {
		return err
	}
	return clientState.Validate()
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (cup UpgradeClientProposal) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if err := unpacker.UnpackAny(cup.ClientState, new(exported.ClientState)); err != nil {
		return err
	}

	if err := unpacker.UnpackAny(cup.ConsensusState, new(exported.ConsensusState)); err != nil {
		return err
	}
	return nil
}

// NewRegisterRelayerProposal creates a new registering relayer proposal.
func NewRegisterRelayerProposal(title, description, chainName string, relayers []string) *RegisterRelayerProposal {
	return &RegisterRelayerProposal{
		Title:       title,
		Description: description,
		ChainName:   chainName,
		Relayers:    relayers,
	}
}

// GetTitle returns the title of a registering relayer proposal.
func (rrp *RegisterRelayerProposal) GetTitle() string { return rrp.Title }

// GetDescription returns the description of a registering relayer proposal.
func (rrp *RegisterRelayerProposal) GetDescription() string { return rrp.Description }

// ProposalRoute returns the routing key of a registering relayer proposal.
func (rrp *RegisterRelayerProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a client registering relayer proposal.
func (rrp *RegisterRelayerProposal) ProposalType() string { return ProposalTypeRelayerRegister }

// ValidateBasic runs basic stateless validity checks
func (rrp *RegisterRelayerProposal) ValidateBasic() error {
	if err := govtypes.ValidateAbstract(rrp); err != nil {
		return err
	}

	if err := host.ClientIdentifierValidator(rrp.ChainName); err != nil {
		return err
	}

	if len(rrp.Relayers) == 0 {
		return govtypes.ErrInvalidLengthGov
	}

	for _, relayer := range rrp.Relayers {
		if _, err := sdk.AccAddressFromBech32(relayer); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
		}
	}
	return nil
}
