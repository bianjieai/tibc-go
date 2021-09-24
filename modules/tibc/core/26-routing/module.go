package routing

import (
	"github.com/gogo/protobuf/grpc"
	"github.com/spf13/cobra"

	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/client/cli"
	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

// Name returns the TIBC routing TICS name.
func Name() string {
	return types.SubModuleName
}

// GetTxCmd returns the root tx command for TIBC routing.
func GetTxCmd() *cobra.Command {
	return nil
}

// GetQueryCmd returns the root query command for TIBC routing.
func GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterQueryService registers the gRPC query service for TIBC routing.
func RegisterQueryService(server grpc.Server, queryServer types.QueryServer) {
	types.RegisterQueryServer(server, queryServer)
}
