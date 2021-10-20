package keeper

import (
	"github.com/armon/go-metrics"

	"github.com/cosmos/cosmos-sdk/telemetry"
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

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tibc", "client", "create"},
			1,
			[]metrics.Label{
				telemetry.NewLabel(types.LabelClientType, clientState.ClientType()),
				telemetry.NewLabel(types.LabelChainName, p.ChainName),
			},
		)
	}()

	// emitting events in the keeper for proposal updates to clients
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCreateClientProposal,
			sdk.NewAttribute(types.AttributeKeyChainName, p.ChainName),
			sdk.NewAttribute(types.AttributeKeyClientType, clientState.ClientType()),
			sdk.NewAttribute(types.AttributeKeyConsensusHeight, clientState.GetLatestHeight().String()),
		),
	)
	return k.CreateClient(ctx, p.ChainName, clientState, consensusState)
}

// HandleUpgradeClientProposal will try to update the client with the new ClientState and ConsensusState if and only if the proposal passes.
func (k Keeper) HandleUpgradeClientProposal(ctx sdk.Context, p *types.UpgradeClientProposal) error {
	clientState, err := types.UnpackClientState(p.ClientState)
	if err != nil {
		return err
	}

	consensusState, err := types.UnpackConsensusState(p.ConsensusState)
	if err != nil {
		return err
	}

	k.Logger(ctx).Info("client updated after governance proposal passed", "client-name", p.ChainName, "height", clientState.GetLatestHeight().String())
	return k.UpgradeClient(ctx, p.ChainName, clientState, consensusState)
}

// HandleRegisterRelayerProposal will try to save the registered relayer address under the specified client
func (k Keeper) HandleRegisterRelayerProposal(ctx sdk.Context, p *types.RegisterRelayerProposal) error {
	// _, has := k.GetClientState(ctx, p.ChainName)
	// if !has {
	// 	return sdkerrors.Wrapf(types.ErrClientNotFound, "chain-name: %s", p.ChainName)
	// }
	k.RegisterRelayers(ctx, p.ChainName, p.Relayers)
	return nil
}
