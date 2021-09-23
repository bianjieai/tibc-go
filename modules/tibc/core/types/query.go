package types

import (
	"github.com/gogo/protobuf/grpc"

	client "github.com/bianjieai/tibc-go/modules/tibc/core/02-client"
	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	packet "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet"
	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	routing "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

// QueryServer defines the TIBC interfaces that the gRPC query server must implement
type QueryServer interface {
	clienttypes.QueryServer
	packettypes.QueryServer
	routingtypes.QueryServer
}

// RegisterQueryService registers each individual TIBC submodule query service
func RegisterQueryService(server grpc.Server, queryService QueryServer) {
	client.RegisterQueryService(server, queryService)
	packet.RegisterQueryService(server, queryService)
	routing.RegisterQueryService(server, queryService)
}
