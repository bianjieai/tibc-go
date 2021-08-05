package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

const (
	// ProposalTypeClientUpdate defines the type for a ClientUpdateProposal
	ProposalTypeClientUpdate = "CreateClient"
)

var (
	_ govtypes.Content = &CreateClientProposal{}
	_ govtypes.Content = &UpgradeClientProposal{}
	_ govtypes.Content = &RegisterRelayerProposal{}
)

// NewCreateClientProposal creates a new client proposal.
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
func (cup *CreateClientProposal) ProposalType() string { return ProposalTypeClientUpdate }

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

// GetTitle returns the title of a client update proposal.
func (cup *UpgradeClientProposal) GetTitle() string { return cup.Title }

// GetDescription returns the description of a client update proposal.
func (cup *UpgradeClientProposal) GetDescription() string { return cup.Description }

// ProposalRoute returns the routing key of a client update proposal.
func (cup *UpgradeClientProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a client update proposal.
func (cup *UpgradeClientProposal) ProposalType() string { return ProposalTypeClientUpdate }

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

// GetTitle returns the title of a client update proposal.
func (rrp *RegisterRelayerProposal) GetTitle() string { return rrp.Title }

// GetDescription returns the description of a client update proposal.
func (rrp *RegisterRelayerProposal) GetDescription() string { return rrp.Description }

// ProposalRoute returns the routing key of a client update proposal.
func (rrp *RegisterRelayerProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a client update proposal.
func (rrp *RegisterRelayerProposal) ProposalType() string { return ProposalTypeClientUpdate }

// ValidateBasic runs basic stateless validity checks
func (rrp *RegisterRelayerProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(rrp)
	if err != nil {
		return err
	}

	if err := host.ClientIdentifierValidator(rrp.ChainName); err != nil {
		return err
	}

	if len(rrp.Relayers) == 0 {
		return govtypes.ErrInvalidLengthGov
	}

	for _, relayer := range rrp.Relayers {
		_, err := sdk.AccAddressFromBech32(relayer)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
		}
	}
	return nil
}
