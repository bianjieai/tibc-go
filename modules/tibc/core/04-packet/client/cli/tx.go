package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
)

// NewSendCleanPacketCmd defines the command to send clean packet.
func NewSendCleanPacketCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-clean-packet [dest-chain-name] [sequence] [flags]",
		Short: "send a clean packet",
		Long:  "send a clean packet",
		Example: fmt.Sprintf(
			"%s tx tibc %s send-clean-packet [dest-chain-name] [sequence] --relay-chain-name test2 --from node0",
			version.AppName, types.SubModuleName,
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			relayChain, err := cmd.Flags().GetString(FlagRelayChain)
			if err != nil {
				return err
			}
			destChain := args[0]
			sequence, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			cleanPacket := types.CleanPacket{
				Sequence:         sequence,
				DestinationChain: destChain,
				RelayChain:       relayChain,
			}

			msg := types.NewMsgCleanPacket(cleanPacket, clientCtx.GetFromAddress())
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
