package types

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

func (m ClientState) CheckHeaderAndUpdateState(
	ctx sdk.Context,
	cdc codec.BinaryCodec,
	store sdk.KVStore,
	header exported.Header,
) (
	exported.ClientState,
	exported.ConsensusState,
	error,
) {
	bscHeader, ok := header.(*Header)
	if !ok {
		return nil, nil, sdkerrors.Wrapf(
			clienttypes.ErrInvalidHeader, "expected type %T, got %T", &Header{}, header,
		)
	}

	// get consensus state from clientStore
	bscConsState, err := GetConsensusState(store, cdc, m.GetLatestHeight())
	if err != nil {
		return nil, nil, sdkerrors.Wrapf(
			err, "could not get consensus state from clientstore at TrustedHeight: %s", m.GetLatestHeight(),
		)
	}

	if err := checkValidity(cdc, store, &m, bscConsState, *bscHeader); err != nil {
		return nil, nil, err
	}
	newClientState, consensusState, err := update(cdc, store, &m, bscHeader)
	if err != nil {
		return nil, nil, err
	}
	return newClientState, consensusState, nil
}

// checkValidity checks if the bsc header is valid.
func checkValidity(
	cdc codec.BinaryCodec,
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
func update(
	cdc codec.BinaryCodec,
	store sdk.KVStore,
	clientState *ClientState,
	header *Header,
) (
	*ClientState,
	*ConsensusState,
	error,
) {
	// The validator set change occurs at `header.Number % cs.Epoch == 0`
	number := header.Height.RevisionHeight
	if number%clientState.Epoch == 0 {
		validators, err := ParseValidators(header.Extra)
		if err != nil {
			return nil, nil, err
		}
		SetPendingValidators(store, cdc, validators)
	}

	// change validator set
	if number%clientState.Epoch == uint64(len(clientState.Validators)/2) {
		validators := GetPendingValidators(cdc, store).Validators
		newVals := make(map[common.Address]struct{}, len(validators))
		for _, val := range validators {
			newVals[common.BytesToAddress(val)] = struct{}{}
		}

		oldLimit := len(clientState.Validators)/2 + 1
		newLimit := len(newVals)/2 + 1
		if newLimit < oldLimit {
			for i := 0; i < oldLimit-newLimit; i++ {
				pruneHeight := clienttypes.NewHeight(header.Height.RevisionNumber, number-uint64(newLimit)-uint64(i))
				DeleteSigner(store, pruneHeight)
			}
		}
		clientState.Validators = validators
	}

	// update the recentSingers
	// Delete the oldest validator from the recent list to allow it signing again
	if limit := uint64(len(clientState.Validators)/2 + 1); number >= limit {
		pruneHeight := clienttypes.NewHeight(header.Height.RevisionNumber, number-limit)
		DeleteSigner(store, pruneHeight)
	}

	cs := &ConsensusState{
		Timestamp: header.Time,
		Number:    header.Height,
		Root:      header.Root,
	}

	clientState.Header = *header
	return clientState, cs, nil
}
