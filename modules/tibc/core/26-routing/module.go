package routing

import (
	"github.com/gogo/protobuf/grpc"
	"github.com/spf13/cobra"

	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/client/cli"
	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

// Name returns the TIBC packets ICS name.
func Name() string {
	return types.SubModuleName
}

// GetTxCmd returns the root tx command for TIBC packets.
func GetTxCmd() *cobra.Command {
	return nil
}

// GetQueryCmd returns the root query command for TIBC packets.
func GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterQueryService registers the gRPC query service for TIBC packets.
func RegisterQueryService(server grpc.Server, queryServer types.QueryServer) {
	types.RegisterQueryServer(server, queryServer)
}
