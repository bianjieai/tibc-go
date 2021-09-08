package types

func DefaultGenesisState() GenesisState{
	return GenesisState{
		Rules: []string{},
	}
}