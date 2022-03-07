package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
)

// NewTxCmd returns the transaction commands for TIBC multi token transfer
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        "tibc-mt-transfer",
		Short:                      "TIBC multi token transfer transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewTransferTxCmd(),
	)

	return txCmd
}

// GetQueryCmd returns the query commands for TIBC multi token transfer
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        "tibc-mt-transfer",
		Short:                      "TIBC multi token transfer query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}

	queryCmd.AddCommand(
		GetCmdQueryClassTrace(),
		GetCmdQueryClassTraces(),
	)

	return queryCmd
}
