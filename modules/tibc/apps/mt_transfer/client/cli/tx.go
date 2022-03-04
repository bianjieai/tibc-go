package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
)

// NewTransferTxCmd returns the command to create a NewMsgTransfer transaction
func NewTransferTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [dest-chain] [receiver] [class] [id]",
		Short: "Transfer a non fungible token through TIBC",
		Example: fmt.Sprintf(
			"%s tx tibc-nft-transfer transfer <dest-chain-name> <receiver> <denom-id> <nft-id> "+
				"--relay-chain=<relay-chain-name> "+
				"--dest-contract=<receive-the-contract-address-of-nft>",
			version.AppName,
		),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			sender := clientCtx.GetFromAddress().String()
			destChain := args[0]
			receiver := args[1]
			class := args[2]
			id := args[3]

			relayChain, err := cmd.Flags().GetString(FlagRelayChain)
			if err != nil {
				return err
			}

			destContract, err := cmd.Flags().GetString(FlagDestContract)
			if err != nil {
				return err
			}

			msg := types.NewMsgNftTransfer(
				class, id, sender, receiver,
				destChain, relayChain, destContract,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsNftTransfer)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
