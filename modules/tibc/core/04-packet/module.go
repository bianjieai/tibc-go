package packet

import (
	"github.com/gogo/protobuf/grpc"
	"github.com/spf13/cobra"

	"github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/client/cli"
	"github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
)

// Name returns the TIBC packet TICS name.
func Name() string {
	return types.SubModuleName
}

// GetTxCmd returns the root tx command for TIBC packet.
func GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// GetQueryCmd returns the root query command for TIBC packet.
func GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterQueryService registers the gRPC query service for TIBC packet.
func RegisterQueryService(server grpc.Server, queryServer types.QueryServer) {
	types.RegisterQueryServer(server, queryServer)
}
