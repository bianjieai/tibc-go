package cli

import (
	"github.com/spf13/cobra"

	"github.com/bianjieai/tibc-go/modules/tibc/light-clients/07-tendermint/types"
)

// NewTxCmd returns a root CLI command handler for all x/ibc/light-clients/07-tendermint transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "Tendermint client transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	txCmd.AddCommand(
		NewUpdateClientCmd(),
	)
	return txCmd
}
