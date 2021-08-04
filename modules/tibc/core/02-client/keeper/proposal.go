package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
)

// HandleCreateClientProposal will try to create the client with the new ClientState and ConsensusState if and only if the proposal passes.
func (k Keeper) HandleCreateClientProposal(ctx sdk.Context, p *types.CreateClientProposal) error {
	_, has := k.GetClientState(ctx, p.ChainName)
	if has {
		return sdkerrors.Wrapf(types.ErrClientExists, "chain-name: %s", p.ChainName)
	}

	clientState, err := types.UnpackClientState(p.ClientState)
	if err != nil {
		return err
	}

	consensusState, err := types.UnpackConsensusState(p.ConsensusState)
	if err != nil {
		return err
	}
	return k.CreateClient(ctx, p.ChainName, clientState, consensusState)
}

// CreateClientProposal will try to update the client with the new ClientState and ConsensusState if and only if the proposal passes.
func (k Keeper) HandleUpgradeClientProposal(ctx sdk.Context, p *types.UpgradeClientProposal) error {
	clientState, err := types.UnpackClientState(p.ClientState)
	if err != nil {
		return err
	}

	consensusState, err := types.UnpackConsensusState(p.ConsensusState)
	if err != nil {
		return err
	}
	return k.UpgradeClient(ctx, p.ChainName, clientState, consensusState)
}

// RegisterRelayerProposal will try to update the client with the new header if and only if
// the proposal passes. The localhost client is not allowed to be modified with a proposal.
func (k Keeper) HandleRegisterRelayerProposal(ctx sdk.Context, p *types.RegisterRelayerProposal) error {
	//TODO
	return nil
}
