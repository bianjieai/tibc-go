package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/client/utils"
	"github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
)

const (
	flagSequences = "sequences"
)

// GetCmdQueryPacketCommitments defines the command to query all packet commitments associated with
// source chain name and destination chain name
func GetCmdQueryPacketCommitments() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "packet-commitments [source-chain] [dest-chain]",
		Short: "Query all packet commitments associated with source",
		Long:  "Query all packet commitments associated with source chain name and destination chain name",
		Example: fmt.Sprintf(
			"%s query %s %s packet-commitments [source-chain] [dest-chain]",
			version.AppName, host.ModuleName, types.SubModuleName,
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			req := &types.QueryPacketCommitmentsRequest{
				SourceChain: args[0],
				DestChain:   args[1],
				Pagination:  pageReq,
			}

			res, err := queryClient.PacketCommitments(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "packet commitments associated with source chain name and destination chain name")

	return cmd
}

// GetCmdQueryPacketCommitment defines the command to query a packet commitment
func GetCmdQueryPacketCommitment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "packet-commitment [source-chain] [dest-chain] [sequence]",
		Short: "Query a packet commitment",
		Long:  "Query a packet commitment",
		Example: fmt.Sprintf(
			"%s query %s %s packet-commitment [source-chain] [dest-chain] [sequence]",
			version.AppName, host.ModuleName, types.SubModuleName,
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			sourceChain := args[0]
			destChain := args[1]
			prove, _ := cmd.Flags().GetBool(flags.FlagProve)

			seq, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			res, err := utils.QueryPacketCommitment(clientCtx, sourceChain, destChain, seq, prove)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().Bool(flags.FlagProve, true, "show proofs for the query results")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryPacketReceipt defines the command to query a packet receipt
func GetCmdQueryPacketReceipt() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "packet-receipt [source-chain] [dest-chain] [sequence]",
		Short: "Query a packet receipt",
		Long:  "Query a packet receipt",
		Example: fmt.Sprintf(
			"%s query %s %s packet-receipt [source-chain] [dest-chain] [sequence]",
			version.AppName, host.ModuleName, types.SubModuleName,
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			sourceChain := args[0]
			destChain := args[1]
			prove, _ := cmd.Flags().GetBool(flags.FlagProve)

			seq, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			res, err := utils.QueryPacketReceipt(clientCtx, sourceChain, destChain, seq, prove)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().Bool(flags.FlagProve, true, "show proofs for the query results")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryPacketAcknowledgement defines the command to query a packet acknowledgement
func GetCmdQueryPacketAcknowledgement() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "packet-ack [source-chain] [dest-chain] [sequence]",
		Short: "Query a packet acknowledgement",
		Long:  "Query a packet acknowledgement",
		Example: fmt.Sprintf(
			"%s query %s %s packet-ack [source-chain] [dest-chain] [sequence]",
			version.AppName, host.ModuleName, types.SubModuleName,
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			sourceChain := args[0]
			destChain := args[1]
			prove, _ := cmd.Flags().GetBool(flags.FlagProve)

			seq, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			res, err := utils.QueryPacketAcknowledgement(clientCtx, sourceChain, destChain, seq, prove)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().Bool(flags.FlagProve, true, "show proofs for the query results")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryUnreceivedPackets defines the command to query all the unreceived
// packets on the receiving chain
func GetCmdQueryUnreceivedPackets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unreceived-packets [source-chain] [dest-chain]",
		Short: "Query all the unreceived packets associated with source chain name and destination chain name",
		Long: `Determine if a packet, given a list of packet commitment sequences, is unreceived.

The return value represents:
- Unreceived packet commitments: no acknowledgement exists on receiving chain for the given packet commitment sequence on sending chain.
`,
		Example: fmt.Sprintf(
			"%s query %s %s unreceived-packets [source-chain] [dest-chain] --sequences=1,2,3",
			version.AppName, host.ModuleName, types.SubModuleName,
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			seqSlice, err := cmd.Flags().GetInt64Slice(flagSequences)
			if err != nil {
				return err
			}

			seqs := make([]uint64, len(seqSlice))
			for i := range seqSlice {
				seqs[i] = uint64(seqSlice[i])
			}

			req := &types.QueryUnreceivedPacketsRequest{
				SourceChain:               args[0],
				DestChain:                 args[1],
				PacketCommitmentSequences: seqs,
			}

			res, err := queryClient.UnreceivedPackets(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().Int64Slice(flagSequences, []int64{}, "comma separated list of packet sequence numbers")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryUnreceivedAcks defines the command to query all the unreceived acks on the original sending chain
func GetCmdQueryUnreceivedAcks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unreceived-acks [source-chain] [dest-chain]]",
		Short: "Query all the unreceived acks associated with source chain name and destination chain name",
		Long: `Given a list of acknowledgement sequences from counterparty, determine if an ack on the counterparty chain has been received on the executing chain.

The return value represents:
- Unreceived packet acknowledgement: packet commitment exists on original sending (executing) chain and ack exists on receiving chain.
`,
		Example: fmt.Sprintf(
			"%s query %s %s unreceived-acks [source-chain] [dest-chain] --sequences=1,2,3",
			version.AppName, host.ModuleName, types.SubModuleName,
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			seqSlice, err := cmd.Flags().GetInt64Slice(flagSequences)
			if err != nil {
				return err
			}

			seqs := make([]uint64, len(seqSlice))
			for i := range seqSlice {
				seqs[i] = uint64(seqSlice[i])
			}

			req := &types.QueryUnreceivedAcksRequest{
				SourceChain:        args[0],
				DestChain:          args[1],
				PacketAckSequences: seqs,
			}

			res, err := queryClient.UnreceivedAcks(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().Int64Slice(flagSequences, []int64{}, "comma separated list of packet sequence numbers")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryCleanPacketCommitment defines the command to query a packet commitment
func GetCmdQueryCleanPacketCommitment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean-packet-commitment [source-chain] [dest-chain]",
		Short: "Query the clean packet commitment",
		Long:  "Query the clean packet commitment",
		Example: fmt.Sprintf(
			"%s query %s %s clean-packet-commitment [source-chain] [dest-chain]",
			version.AppName, host.ModuleName, types.SubModuleName,
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			sourceChain := args[0]
			destChain := args[1]
			prove, _ := cmd.Flags().GetBool(flags.FlagProve)

			res, err := utils.QueryCleanPacketCommitment(clientCtx, sourceChain, destChain, prove)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().Bool(flags.FlagProve, true, "show proofs for the query results")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
