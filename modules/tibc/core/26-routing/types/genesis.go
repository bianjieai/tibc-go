package types

import (
	"regexp"

	errorsmod "cosmossdk.io/errors"
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
			return errorsmod.Wrap(ErrInvalidRule, "invalid rule")
		}
	}
	return nil
}
