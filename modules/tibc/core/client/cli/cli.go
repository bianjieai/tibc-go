package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	ibcclient "github.com/bianjieai/tibc-go/modules/tibc/core/02-client"
	packet "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	routing "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	ibcTxCmd := &cobra.Command{
		Use:                        host.ModuleName,
		Short:                      "IBC transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ibcTxCmd.AddCommand(
		ibcclient.GetTxCmd(),
		packet.GetTxCmd(),
	)

	return ibcTxCmd
}

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group tibc queries under a subcommand
	ibcQueryCmd := &cobra.Command{
		Use:                        host.ModuleName,
		Short:                      "Querying commands for the TIBC module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ibcQueryCmd.AddCommand(
		ibcclient.GetQueryCmd(),
		packet.GetQueryCmd(),
		routing.GetQueryCmd(),
	)

	return ibcQueryCmd
}
