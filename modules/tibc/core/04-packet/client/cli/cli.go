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
		GetCmdQueryCleanPacketCommitment(),
	)

	return queryCmd
}

// NewTxCmd returns the command to create and handle TIBC clients
func NewTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      "TIBC client transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewSendCleanPacketCmd(),
	)
	return txCmd
}
