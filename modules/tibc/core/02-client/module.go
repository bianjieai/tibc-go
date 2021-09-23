package client

import (
	"github.com/gogo/protobuf/grpc"
	"github.com/spf13/cobra"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/client/cli"
	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
)

// Name returns the TIBC client name.
func Name() string {
	return types.SubModuleName
}

// GetQueryCmd returns no root query command for the TIBC client.
func GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// GetTxCmd returns the root tx command for 02-client.
func GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// RegisterQueryService registers the gRPC query service for TIBC client.
func RegisterQueryService(server grpc.Server, queryServer types.QueryServer) {
	types.RegisterQueryServer(server, queryServer)
}
