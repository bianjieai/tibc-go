package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

// GetCmdQueryRoutingRulesCommitment defines the command to query routing rules
func GetCmdQueryRoutingRulesCommitment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routing-rules",
		Short: "Query routing rules commitment",
		Long:  "Query routing rules commitment",
		Example: fmt.Sprintf(
			"%s query %s %s routing-rules", version.AppName, host.ModuleName, types.SubModuleName,
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			if err != nil {
				return err
			}
			req := &types.QueryRoutingRulesRequest{}

			res, err := queryClient.RoutingRules(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
