package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

// GetQueryCmd returns the query commands for TIBC routing
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "tibc routing query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdQueryRoutingRulesCommitment(),
	)
	return queryCmd
}
