package cli

import (
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

// NewCmdSubmitUpgradeProposal implements a command handler for submitting a client create proposal transaction.
func NewCmdCreateClientProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client-create [chain-name] [path/to/client_state.json] [path/to/consensus_state.json] [flags]",
		Args:  cobra.ExactArgs(3),
		Short: "Submit a client create proposal",
		Long:  "create a new IBC client with the specified client state and consensus state",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			chainName := args[0]

			cdc := codec.NewProtoCodec(clientCtx.InterfaceRegistry)
			// attempt to unmarshal client state argument
			var clientState exported.ClientState
			clientStateBz, err := ioutil.ReadFile(args[1])
			if err != nil {
				return errors.Wrap(err, "neither JSON input nor path to .json file for client state were provided")
			}

			if err := cdc.UnmarshalInterfaceJSON(clientStateBz, &clientState); err != nil {
				return errors.Wrap(err, "error unmarshalling client state file")
			}

			var consensusState exported.ConsensusState
			consensusStateBz, err := ioutil.ReadFile(args[2])
			if err != nil {
				return errors.Wrap(err, "neither JSON input nor path to .json file for consensus state were provided")
			}

			if err := cdc.UnmarshalInterfaceJSON(consensusStateBz, &consensusState); err != nil {
				return errors.Wrap(err, "error unmarshalling consensus state file")
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			content, err := types.NewCreateClientProposal(title, description, chainName, clientState, consensusState)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
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

// NewCmdUpgradeClientProposal implements a command handler for submitting a client upgrade proposal transaction.
func NewCmdUpgradeClientProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client-upgrade [chain-name] [path/to/client_state.json] [path/to/consensus_state.json] [flags]",
		Args:  cobra.ExactArgs(3),
		Short: "Submit a client upgrade proposal",
		Long:  "upgrade a IBC client with the specified client state and consensus state",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			chainName := args[0]

			cdc := codec.NewProtoCodec(clientCtx.InterfaceRegistry)
			// attempt to unmarshal client state argument
			var clientState exported.ClientState
			clientStateBz, err := ioutil.ReadFile(args[1])
			if err != nil {
				return errors.Wrap(err, "neither JSON input nor path to .json file for client state were provided")
			}

			if err := cdc.UnmarshalInterfaceJSON(clientStateBz, &clientState); err != nil {
				return errors.Wrap(err, "error unmarshalling client state file")
			}

			var consensusState exported.ConsensusState
			consensusStateBz, err := ioutil.ReadFile(args[2])
			if err != nil {
				return errors.Wrap(err, "neither JSON input nor path to .json file for consensus state were provided")
			}

			if err := cdc.UnmarshalInterfaceJSON(consensusStateBz, &consensusState); err != nil {
				return errors.Wrap(err, "error unmarshalling consensus state file")
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			content, err := types.NewUpgradeClientProposal(title, description, chainName, clientState, consensusState)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
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

// NewCmdRegisterRelayerProposal implements a command handler for submitting a relayer register proposal transaction.
func NewCmdRegisterRelayerProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relayer-register [chain-name] [relayers-address] [flags]",
		Args:  cobra.ExactArgs(2),
		Short: "Submit a relayer register proposal",
		Long:  "Submit a relayer register proposal for the specified client",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(cli.FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(cli.FlagDescription)
			if err != nil {
				return err
			}

			relayers := strings.Split(args[1], ",")
			content := types.NewRegisterRelayerProposal(title, description, args[0], relayers)

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
