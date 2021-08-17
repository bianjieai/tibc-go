package keeper

import (
	"context"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

// ClientState implements the IBC QueryServer interface
func (q Keeper) ClientState(c context.Context, req *clienttypes.QueryClientStateRequest) (*clienttypes.QueryClientStateResponse, error) {
	return q.ClientKeeper.ClientState(c, req)
}

// ClientStates implements the IBC QueryServer interface
func (q Keeper) ClientStates(c context.Context, req *clienttypes.QueryClientStatesRequest) (*clienttypes.QueryClientStatesResponse, error) {
	return q.ClientKeeper.ClientStates(c, req)
}

// ConsensusState implements the IBC QueryServer interface
func (q Keeper) ConsensusState(c context.Context, req *clienttypes.QueryConsensusStateRequest) (*clienttypes.QueryConsensusStateResponse, error) {
	return q.ClientKeeper.ConsensusState(c, req)
}

// ConsensusStates implements the IBC QueryServer interface
func (q Keeper) ConsensusStates(c context.Context, req *clienttypes.QueryConsensusStatesRequest) (*clienttypes.QueryConsensusStatesResponse, error) {
	return q.ClientKeeper.ConsensusStates(c, req)
}

// PacketCommitment implements the IBC QueryServer interface
func (q Keeper) PacketCommitment(c context.Context, req *packettypes.QueryPacketCommitmentRequest) (*packettypes.QueryPacketCommitmentResponse, error) {
	return q.Packetkeeper.PacketCommitment(c, req)
}

// PacketCommitments implements the IBC QueryServer interface
func (q Keeper) PacketCommitments(c context.Context, req *packettypes.QueryPacketCommitmentsRequest) (*packettypes.QueryPacketCommitmentsResponse, error) {
	return q.Packetkeeper.PacketCommitments(c, req)
}

// PacketReceipt implements the IBC QueryServer interface
func (q Keeper) PacketReceipt(c context.Context, req *packettypes.QueryPacketReceiptRequest) (*packettypes.QueryPacketReceiptResponse, error) {
	return q.Packetkeeper.PacketReceipt(c, req)
}

// PacketAcknowledgement implements the IBC QueryServer interface
func (q Keeper) PacketAcknowledgement(c context.Context, req *packettypes.QueryPacketAcknowledgementRequest) (*packettypes.QueryPacketAcknowledgementResponse, error) {
	return q.Packetkeeper.PacketAcknowledgement(c, req)
}

// PacketAcknowledgements implements the IBC QueryServer interface
func (q Keeper) PacketAcknowledgements(c context.Context, req *packettypes.QueryPacketAcknowledgementsRequest) (*packettypes.QueryPacketAcknowledgementsResponse, error) {
	return q.Packetkeeper.PacketAcknowledgements(c, req)
}

// UnreceivedPackets implements the IBC QueryServer interface
func (q Keeper) UnreceivedPackets(c context.Context, req *packettypes.QueryUnreceivedPacketsRequest) (*packettypes.QueryUnreceivedPacketsResponse, error) {
	return q.Packetkeeper.UnreceivedPackets(c, req)
}

// UnreceivedAcks implements the IBC QueryServer interface
func (q Keeper) UnreceivedAcks(c context.Context, req *packettypes.QueryUnreceivedAcksRequest) (*packettypes.QueryUnreceivedAcksResponse, error) {
	return q.Packetkeeper.UnreceivedAcks(c, req)
}

// CleanPacketCommitment implements the IBC QueryServer interface
func (q Keeper) CleanPacketCommitment(c context.Context, req *packettypes.QueryCleanPacketCommitmentRequest) (*packettypes.QueryCleanPacketCommitmentResponse, error) {
	return q.Packetkeeper.CleanPacketCommitment(c, req)
}

// RoutingRules implements the IBC QueryServer interface
func (q Keeper) RoutingRules(c context.Context, req *routingtypes.QueryRoutingRulesRequest) (*routingtypes.QueryRoutingRulesResponse, error) {
	return q.RoutingKeeper.RoutingRules(c, req)
}
