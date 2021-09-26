package keeper

import (
	"context"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	packettypes "github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
	routingtypes "github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

// ClientState implements the TIBC QueryServer interface
func (q Keeper) ClientState(c context.Context, req *clienttypes.QueryClientStateRequest) (*clienttypes.QueryClientStateResponse, error) {
	return q.ClientKeeper.ClientState(c, req)
}

// ClientStates implements the TIBC QueryServer interface
func (q Keeper) ClientStates(c context.Context, req *clienttypes.QueryClientStatesRequest) (*clienttypes.QueryClientStatesResponse, error) {
	return q.ClientKeeper.ClientStates(c, req)
}

// ConsensusState implements the TIBC QueryServer interface
func (q Keeper) ConsensusState(c context.Context, req *clienttypes.QueryConsensusStateRequest) (*clienttypes.QueryConsensusStateResponse, error) {
	return q.ClientKeeper.ConsensusState(c, req)
}

// ConsensusStates implements the TIBC QueryServer interface
func (q Keeper) ConsensusStates(c context.Context, req *clienttypes.QueryConsensusStatesRequest) (*clienttypes.QueryConsensusStatesResponse, error) {
	return q.ClientKeeper.ConsensusStates(c, req)
}

// PacketCommitment implements the TIBC QueryServer interface
func (q Keeper) PacketCommitment(c context.Context, req *packettypes.QueryPacketCommitmentRequest) (*packettypes.QueryPacketCommitmentResponse, error) {
	return q.PacketKeeper.PacketCommitment(c, req)
}

// PacketCommitments implements the TIBC QueryServer interface
func (q Keeper) PacketCommitments(c context.Context, req *packettypes.QueryPacketCommitmentsRequest) (*packettypes.QueryPacketCommitmentsResponse, error) {
	return q.PacketKeeper.PacketCommitments(c, req)
}

// PacketReceipt implements the TIBC QueryServer interface
func (q Keeper) PacketReceipt(c context.Context, req *packettypes.QueryPacketReceiptRequest) (*packettypes.QueryPacketReceiptResponse, error) {
	return q.PacketKeeper.PacketReceipt(c, req)
}

// PacketAcknowledgement implements the TIBC QueryServer interface
func (q Keeper) PacketAcknowledgement(c context.Context, req *packettypes.QueryPacketAcknowledgementRequest) (*packettypes.QueryPacketAcknowledgementResponse, error) {
	return q.PacketKeeper.PacketAcknowledgement(c, req)
}

// PacketAcknowledgements implements the TIBC QueryServer interface
func (q Keeper) PacketAcknowledgements(c context.Context, req *packettypes.QueryPacketAcknowledgementsRequest) (*packettypes.QueryPacketAcknowledgementsResponse, error) {
	return q.PacketKeeper.PacketAcknowledgements(c, req)
}

// UnreceivedPackets implements the TIBC QueryServer interface
func (q Keeper) UnreceivedPackets(c context.Context, req *packettypes.QueryUnreceivedPacketsRequest) (*packettypes.QueryUnreceivedPacketsResponse, error) {
	return q.PacketKeeper.UnreceivedPackets(c, req)
}

// UnreceivedAcks implements the TIBC QueryServer interface
func (q Keeper) UnreceivedAcks(c context.Context, req *packettypes.QueryUnreceivedAcksRequest) (*packettypes.QueryUnreceivedAcksResponse, error) {
	return q.PacketKeeper.UnreceivedAcks(c, req)
}

// CleanPacketCommitment implements the TIBC QueryServer interface
func (q Keeper) CleanPacketCommitment(c context.Context, req *packettypes.QueryCleanPacketCommitmentRequest) (*packettypes.QueryCleanPacketCommitmentResponse, error) {
	return q.PacketKeeper.CleanPacketCommitment(c, req)
}

// RoutingRules implements the TIBC QueryServer interface
func (q Keeper) RoutingRules(c context.Context, req *routingtypes.QueryRoutingRulesRequest) (*routingtypes.QueryRoutingRulesResponse, error) {
	return q.RoutingKeeper.RoutingRules(c, req)
}

func (q Keeper) Relayers(c context.Context, req *clienttypes.QueryRelayersRequest) (*clienttypes.QueryRelayersResponse, error) {
	return q.ClientKeeper.Relayers(c, req)
}
