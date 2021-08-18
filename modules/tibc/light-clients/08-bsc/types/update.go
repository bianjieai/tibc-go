package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

func (m ClientState) CheckHeaderAndUpdateState(
	ctx sdk.Context,
	cdc codec.BinaryMarshaler,
	store sdk.KVStore,
	header exported.Header,
) (exported.ClientState, exported.ConsensusState, error) {
	bscHeader, ok := header.(*Header)
	if !ok {
		return nil, nil, sdkerrors.Wrapf(
			clienttypes.ErrInvalidHeader, "expected type %T, got %T", &Header{}, header,
		)
	}

	// get consensus state from clientStore
	bscConsState, err := GetConsensusState(store, cdc, bscHeader.Height)
	if err != nil {
		return nil, nil, sdkerrors.Wrapf(
			err, "could not get consensus state from clientstore at TrustedHeight: %s", bscHeader.Height,
		)
	}

	if err := checkValidity(cdc, store, &m, bscConsState, *bscHeader); err != nil {
		return nil, nil, err
	}
	newClientState, consensusState := update(ctx, store, &m, bscHeader)
	return newClientState, consensusState, nil
}

// checkValidity checks if the bsc header is valid.
func checkValidity(
	cdc codec.BinaryMarshaler,
	store sdk.KVStore,
	clientState *ClientState,
	consState *ConsensusState,
	header Header,

) error {
	if err := header.ValidateBasic(); err != nil {
		return err
	}

	return verifyHeader(cdc, store, clientState, header)
}

// update the RecentSingers and the ConsensusState.
func update(ctx sdk.Context, clientStore sdk.KVStore, clientState *ClientState, header *Header) (*ClientState, *ConsensusState) {
	//height := header.GetHeight().(clienttypes.Height)
	// TODO
	return nil, nil
}
