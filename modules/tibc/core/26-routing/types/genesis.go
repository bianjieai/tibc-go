package types

import (
	"regexp"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func DefaultGenesisState() GenesisState {
	return GenesisState{
		Rules: []string{},
	}
}
func NewGenesisState(rules []string) GenesisState {
	return GenesisState{
		Rules: rules,
	}
}

func (gs GenesisState) Validate() error {
	for _, rule := range gs.Rules {
		valid, _ := regexp.MatchString(RulePattern, rule)
		if !valid {
			return sdkerrors.Wrap(ErrInvalidRule, "invalid rule")
		}
	}
	return nil
}
