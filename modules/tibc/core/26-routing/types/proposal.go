package types

import (
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

const (
	ProposalTypeSetRoutingRules = "SetRoutingRules"
)

var (
	_ govv1beta1.Content = &SetRoutingRulesProposal{}
)

func init() {
	govv1beta1.RegisterProposalType(ProposalTypeSetRoutingRules)
}

// NewSetRoutingRulesProposal creates a new setting rules proposal.
func NewSetRoutingRulesProposal(title, description string, rules []string) (*SetRoutingRulesProposal, error) {
	return &SetRoutingRulesProposal{
		Title:       title,
		Description: description,
		Rules:       rules,
	}, nil
}

// GetTitle returns the title of a setting rules proposal.
func (cup *SetRoutingRulesProposal) GetTitle() string { return cup.Title }

// GetDescription returns the description of a setting rules proposal.
func (cup *SetRoutingRulesProposal) GetDescription() string { return cup.Description }

// ProposalRoute returns the routing key of a setting rules proposal.
func (cup *SetRoutingRulesProposal) ProposalRoute() string { return host.RouterKey }

// ProposalType returns the type of a setting rules proposal.
func (cup *SetRoutingRulesProposal) ProposalType() string { return ProposalTypeSetRoutingRules }

// ValidateBasic runs basic stateless validity checks
func (cup *SetRoutingRulesProposal) ValidateBasic() error {
	err := govv1beta1.ValidateAbstract(cup)
	if err != nil {
		return err
	}

	if err := host.RoutingRulesValidator(cup.Rules); err != nil {
		return err
	}
	return nil
}
