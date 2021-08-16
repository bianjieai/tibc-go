package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
)

// GetQueryCmd returns the query commands for TIBC packets
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "tibc packet query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	queryCmd.AddCommand(
		GetCmdQueryPacketCommitment(),
		GetCmdQueryPacketCommitments(),
		GetCmdQueryPacketReceipt(),
		GetCmdQueryPacketAcknowledgement(),
		GetCmdQueryUnreceivedPackets(),
		GetCmdQueryUnreceivedAcks(),
	)

	return queryCmd
}
