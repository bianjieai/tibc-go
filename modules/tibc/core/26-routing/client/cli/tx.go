package cli

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

// NewSetRoutingRulesProposalCmd implements a command handler for submitting a setting rules proposal transaction.
func NewSetRoutingRulesProposalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-rules [path/to/routing_rules.json] [flags]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a rules set proposal",
		Long:  "set routing rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			routingRulesBz, err := ioutil.ReadFile(args[0])
			if err != nil {
				return errors.Wrap(err, "neither JSON input nor path to .json file for routing rules were provided")
			}

			var rules []string
			if err := json.Unmarshal(routingRulesBz, &rules); err != nil {
				return errors.Wrap(err, "error unmarshalling rules file")
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			content, err := types.NewSetRoutingRulesProposal(title, description, rules)
			if err != nil {
				return err
			}

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagTitle, "", "title of proposal")
	cmd.Flags().String(cli.FlagDescription, "", "description of proposal")
	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")
	return cmd
}
