package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
	"strconv"

	"github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
)

// NewSendCleanPacketCmd defines the command to send clean packet.
func NewSendCleanPacketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "send-clean-packet [dest-chain-name] [sequence] [flags]",
		Short:   "send a clean packet",
		Long:    "send a clean packet",
		Example: fmt.Sprintf("%s tx ibc %s send-clean-packet [dest-chain-name] [sequence] --source-chain-name test1 --relay-chain-name test2 --from node0", version.AppName, types.SubModuleName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sourceChain, err := cmd.Flags().GetString(FlagSourceChain)
			if err != nil {
				return err
			}

			relayChain, err := cmd.Flags().GetString(FlagRelayChain)
			if err != nil {
				return err
			}
			destChain := args[0];
			sequence, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			cleanPacket := types.CleanPacket{
				Sequence:         sequence,
				SourceChain:      sourceChain,
				DestinationChain: destChain,
				RelayChain:       relayChain,
			}
			
			msg:= types.NewMsgCleanPacket(cleanPacket, clientCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(FsSendCleanPacket)
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
