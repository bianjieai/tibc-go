package cli

import (
	"fmt"
	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
)

// NewTransferTxCmd returns the command to create a NewMsgTransfer transaction
func NewTransferTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nft transfer [dest-chain] [class] [id] [uri] [receiver]",
		Short: "Transfer a non fungible token through TIBC",
		Example: fmt.Sprintf("%s tx tibc-nft-transfer transfer <dest-chain> <class> <id> <uri> <receiver> "+
			"--away-from-origin=<away-from-origin>"+
			"--relay-chain=<dest-chain>",
			version.AppName),
		Args:    cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			sender := clientCtx.GetFromAddress().String()
			destChain := args[0]
			class := args[1]
			id := args[2]
			uri := args[3]
			receiver := args[4]
			relayChain, err := cmd.Flags().GetString(FlagRelayChain)
			if err != nil {
				return err
			}

			awayFromChain, err := cmd.Flags().GetBool(FlagAwayFromChain)
			if err != nil {
				return err
			}

			if err != nil {
				return err
			}

			msg := types.NewMsgNftTransfer(
				class, id, uri, sender, receiver, awayFromChain,
				destChain, relayChain,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsNftTransfer)
	_ = cmd.MarkFlagRequired(FlagAwayFromChain)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}