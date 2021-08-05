package cli

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	"github.com/bianjieai/tibc-go/modules/tibc/light-clients/07-tendermint/types"
)

// const (
// 	flagTrustLevel                   = "trust-level"
// 	flagProofSpecs                   = "proof-specs"
// 	flagUpgradePath                  = "upgrade-path"
// 	flagAllowUpdateAfterExpiry       = "allow_update_after_expiry"
// 	flagAllowUpdateAfterMisbehaviour = "allow_update_after_misbehaviour"
// )

// NewUpdateClientCmd defines the command to update a client as defined in
// https://github.com/cosmos/ics/tree/master/spec/ics-002-client-semantics#update
func NewUpdateClientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [client-id] [path/to/header.json]",
		Short: "update existing client with a header",
		Long:  "update existing tendermint client with a tendermint header",
		Example: fmt.Sprintf(
			"$ %s tx ibc %s update [client-id] [path/to/header.json] --from node0 --home ../node0/<app>cli --chain-id $CID",
			version.AppName, types.SubModuleName,
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			clientID := args[0]

			cdc := codec.NewProtoCodec(clientCtx.InterfaceRegistry)

			var header *types.Header
			if err := cdc.UnmarshalJSON([]byte(args[1]), header); err != nil {
				// check for file path if JSON input is not provided
				contents, err := ioutil.ReadFile(args[1])
				if err != nil {
					return errors.New("neither JSON input nor path to .json file were provided")
				}
				if err := cdc.UnmarshalJSON(contents, header); err != nil {
					return errors.Wrap(err, "error unmarshalling header file")
				}
			}

			msg, err := clienttypes.NewMsgUpdateClient(clientID, header, clientCtx.GetFromAddress())
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// func parseFraction(fraction string) (types.Fraction, error) {
// 	fr := strings.Split(fraction, "/")
// 	if len(fr) != 2 || fr[0] == fraction {
// 		return types.Fraction{}, fmt.Errorf("fraction must have format 'numerator/denominator' got %s", fraction)
// 	}

// 	numerator, err := strconv.ParseUint(fr[0], 10, 64)
// 	if err != nil {
// 		return types.Fraction{}, fmt.Errorf("invalid trust-level numerator: %w", err)
// 	}

// 	denominator, err := strconv.ParseUint(fr[1], 10, 64)
// 	if err != nil {
// 		return types.Fraction{}, fmt.Errorf("invalid trust-level denominator: %w", err)
// 	}

// 	return types.Fraction{
// 		Numerator:   numerator,
// 		Denominator: denominator,
// 	}, nil

// }
