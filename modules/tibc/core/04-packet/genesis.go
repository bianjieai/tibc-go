package packet

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/keeper"
	"github.com/bianjieai/tibc-go/modules/tibc/core/04-packet/types"
)

// InitGenesis initializes the tibc packet submodule's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, gs types.GenesisState) {
	for _, ack := range gs.Acknowledgements {
		k.SetPacketAcknowledgement(ctx, ack.SourceChain, ack.DestinationChain, ack.Sequence, ack.Data)
	}
	for _, commitment := range gs.Commitments {
		k.SetPacketCommitment(ctx, commitment.SourceChain, commitment.DestinationChain, commitment.Sequence, commitment.Data)
	}
	for _, receipt := range gs.Receipts {
		k.SetPacketReceipt(ctx, receipt.SourceChain, receipt.DestinationChain, receipt.Sequence)
	}
	for _, ss := range gs.SendSequences {
		k.SetNextSequenceSend(ctx, ss.SourceChain, ss.DestinationChain, ss.Sequence)
	}
}

// ExportGenesis returns the tibc packet submodule's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	return types.GenesisState{
		Acknowledgements: k.GetAllPacketAcks(ctx),
		Commitments:      k.GetAllPacketCommitments(ctx),
		Receipts:         k.GetAllPacketReceipts(ctx),
		SendSequences:    k.GetAllPacketSendSeqs(ctx),
		RecvSequences:    k.GetAllPacketRecvSeqs(ctx),
		AckSequences:     k.GetAllPacketAckSeqs(ctx),
	}
}
