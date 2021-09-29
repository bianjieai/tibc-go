package keeper

import (
	"encoding/hex"

	"github.com/armon/go-metrics"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

// CreateClient creates a new client state and populates it with a given
// client state and consensus state
func (k Keeper) CreateClient(
	ctx sdk.Context, chainName string, clientState exported.ClientState, consensusState exported.ConsensusState,
) error {
	k.SetClientState(ctx, chainName, clientState)
	// verifies initial consensus state against client state and initializes client store with any client-specific metadata
	// e.g. set ProcessedTime in Tendermint clients
	if err := clientState.Initialize(ctx, k.cdc, k.ClientStore(ctx, chainName), consensusState); err != nil {
		return err
	}

	// check if consensus state is nil in case the created client is Localhost
	k.SetClientConsensusState(ctx, chainName, clientState.GetLatestHeight(), consensusState)
	k.Logger(ctx).Info("client created at height", "chain-name", chainName, "height", clientState.GetLatestHeight().String())

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tibc", "client", "create"},
			1,
			[]metrics.Label{telemetry.NewLabel(types.LabelClientType, clientState.ClientType())},
		)
	}()

	return nil
}

// UpdateClient updates the consensus state and the state root from a provided header.
func (k Keeper) UpdateClient(ctx sdk.Context, chainName string, header exported.Header) error {
	clientState, found := k.GetClientState(ctx, chainName)
	if !found {
		return sdkerrors.Wrapf(types.ErrClientNotFound, "cannot update client %s", chainName)
	}

	clientStore := k.ClientStore(ctx, chainName)
	if status := clientState.Status(ctx, clientStore, k.cdc); status != exported.Active {
		return sdkerrors.Wrapf(types.ErrClientNotActive, "cannot update client (%s) with status %s", chainName, status)
	}

	// Any writes made in CheckHeaderAndUpdateState are persisted on both valid updates
	// Light client implementations are responsible for writing the correct metadata (if any) in either case.
	newClientState, newConsensusState, err := clientState.CheckHeaderAndUpdateState(ctx, k.cdc, k.ClientStore(ctx, chainName), header)
	if err != nil {
		return sdkerrors.Wrapf(err, "cannot update client %s", chainName)
	}

	// set new client state regardless of if update is valid update
	k.SetClientState(ctx, chainName, newClientState)

	// set new consensus state regardless of if update is valid update
	var consensusHeight = header.GetHeight()
	k.SetClientConsensusState(ctx, chainName, header.GetHeight(), newConsensusState)
	k.Logger(ctx).Info("client state updated",
		"chain-name", chainName,
		"height", consensusHeight.String(),
	)

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tibc", "client", "update"},
			1,
			[]metrics.Label{
				telemetry.NewLabel(types.LabelClientType, clientState.ClientType()),
				telemetry.NewLabel(types.LabelChainName, chainName),
				telemetry.NewLabel(types.LabelUpdateType, "msg"),
			},
		)
	}()

	// Marshal the Header as an Any and encode the resulting bytes to hex.
	// This prevents the event value from containing invalid UTF-8 characters
	// which may cause data to be lost when JSON encoding/decoding.
	var headerStr = hex.EncodeToString(types.MustMarshalHeader(k.cdc, header))

	// emitting events in the keeper emits for both begin block and handler client updates
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateClient,
			sdk.NewAttribute(types.AttributeKeyChainName, chainName),
			sdk.NewAttribute(types.AttributeKeyClientType, clientState.ClientType()),
			sdk.NewAttribute(types.AttributeKeyConsensusHeight, consensusHeight.String()),
			sdk.NewAttribute(types.AttributeKeyHeader, headerStr),
		),
	)

	return nil
}

// UpgradeClient upgrades the client to a new client state.
func (k Keeper) UpgradeClient(ctx sdk.Context, chainName string, upgradedClientState exported.ClientState, upgradedConsState exported.ConsensusState) error {
	clientState, found := k.GetClientState(ctx, chainName)
	if !found {
		return sdkerrors.Wrapf(types.ErrClientNotFound, "cannot update client %s", chainName)
	}

	if clientState.ClientType() != upgradedClientState.ClientType() {
		return sdkerrors.Wrapf(types.ErrInvalidClientType, "cannot update client %s, client-type not match", chainName)
	}

	k.SetClientState(ctx, chainName, upgradedClientState)
	k.SetClientConsensusState(ctx, chainName, upgradedClientState.GetLatestHeight(), upgradedConsState)

	k.Logger(ctx).Info("client state upgraded", "chain-name", chainName, "height", upgradedClientState.GetLatestHeight().String())

	defer func() {
		telemetry.IncrCounterWithLabels(
			[]string{"tibc", "client", "upgrade"},
			1,
			[]metrics.Label{
				telemetry.NewLabel(types.LabelClientType, upgradedClientState.ClientType()),
				telemetry.NewLabel(types.LabelChainName, chainName),
			},
		)
	}()

	// emitting events in the keeper emits for client upgrades
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpgradeClientProposal,
			sdk.NewAttribute(types.AttributeKeyChainName, chainName),
			sdk.NewAttribute(types.AttributeKeyClientType, upgradedClientState.ClientType()),
			sdk.NewAttribute(types.AttributeKeyConsensusHeight, upgradedClientState.GetLatestHeight().String()),
		),
	)
	return nil
}
