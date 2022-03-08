package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
)

// GetCmdQueryClassTrace defines the command to query a class trace from a given hash.
func GetCmdQueryClassTrace() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "class-trace [hash]",
		Short:   "Query the class trace info from a given trace hash",
		Long:    "Query the class trace info from a given trace hash",
		Example: fmt.Sprintf("%s query tibc-mt-transfer class-trace [hash]", version.AppName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryClassTraceRequest{
				Hash: args[0],
			}

			res, err := queryClient.ClassTrace(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryClassTraces defines the command to query all the class trace infos
// that this chain mantains.
func GetCmdQueryClassTraces() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "class-traces",
		Short:   "Query the trace info for all mt classes",
		Long:    "Query the trace info for all mt classes",
		Example: fmt.Sprintf("%s query tibc-mt-transfer class-traces", version.AppName),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryClassTracesRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.ClassTraces(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "class trace")

	return cmd
}
