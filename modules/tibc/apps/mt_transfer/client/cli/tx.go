package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/types"
)

// NewTransferTxCmd returns the command to create a NewMsgTransfer transaction
func NewTransferTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [dest-chain] [receiver] [class] [id] [amount]",
		Short: "Transfer a mt through TIBC",
		Example: fmt.Sprintf(
			"%s tx tibc-mt-transfer transfer <dest-chain-name> <receiver> <denom-id> <mt-id> <amount> "+
				"--relay-chain=<relay-chain-name> "+
				"--dest-contract=<receive-the-contract-address-of-mt>",
			version.AppName,
		),
		Args: cobra.ExactArgs(5),
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
			amount := args[4]

			relayChain, err := cmd.Flags().GetString(FlagRelayChain)
			if err != nil {
				return err
			}

			destContract, err := cmd.Flags().GetString(FlagDestContract)
			if err != nil {
				return err
			}

			amt, err := strconv.ParseUint(amount, 10, 64)
			if err != nil {
				return err
			}

			msg := types.NewMsgMtTransfer(
				class, id, sender, receiver,
				destChain, relayChain, destContract, amt,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(FsMtTransfer)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
