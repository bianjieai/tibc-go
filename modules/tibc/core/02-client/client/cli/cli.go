package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
)

// GetQueryCmd returns a root CLI command handler for all tibc/client query commands.
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "TIBC client query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdQueryClientStates(),
		GetCmdQueryClientState(),
		GetCmdQueryConsensusStates(),
		GetCmdQueryConsensusState(),
		GetCmdQueryHeader(),
		GetCmdNodeConsensusState(),
		GetCmdQueryRelayers(),
	)

	return queryCmd
}

// NewTxCmd returns a root CLI command handler for all tibc/client transaction commands.
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "TIBC client transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewUpdateClientCmd(),
	)

	return txCmd
}
