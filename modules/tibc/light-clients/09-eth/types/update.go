package types

import (
	"bytes"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

func (m ClientState) CheckHeaderAndUpdateState(
	ctx sdk.Context,
	cdc codec.BinaryMarshaler,
	store sdk.KVStore,
	header exported.Header,
) (exported.ClientState, exported.ConsensusState, error) {
	ethHeader, ok := header.(*Header)
	if !ok {
		return nil, nil, sdkerrors.Wrapf(
			clienttypes.ErrInvalidHeader, "expected type %T, got %T", &Header{}, header,
		)
	}
	height := m.GetLatestHeight()
	// get consensus state from clientStore
	ethConsState, err := GetConsensusState(store, cdc, height)
	if err != nil {
		return nil, nil, sdkerrors.Wrapf(
			err, "could not get consensus state from clientstore at TrustedHeight: %s,please upgrade", m.GetLatestHeight(),
		)
	}
	if err := checkValidity(ctx, cdc, store, &m, ethConsState, *ethHeader); err != nil {
		return nil, nil, err
	}
	// Check the earliest consensus state to see if it is expired, if so then set the prune height
	// so that we can delete consensus state and all associated metadata.
	var (
		pruneHeight exported.Height
		pruneError  error
	)
	pruneCb := func(height exported.Height) bool {
		consState, err := GetConsensusState(store, cdc, height)
		// this error should never occur
		if err != nil {
			pruneError = err
			// this return just for get out of the func
			return true
		}
		bloclTime := uint64(ctx.BlockTime().Unix())
		if consState.Timestamp+m.TrustingPeriod < bloclTime {
			pruneHeight = height
		}
		return true
	}
	IterateConsensusStateAscending(store, pruneCb)
	if pruneError != nil {
		return nil, nil, pruneError
	}
	// if pruneHeight is set, delete consensus state and metadata
	if pruneHeight != nil {
		err = deleteConsensusState(cdc, store, pruneHeight)
		if err != nil {
			return nil, nil, err
		}
		deleteConsensusMetadata(store, pruneHeight)
	}

	newClientState, consensusState, err := update(ctx, cdc, store, &m, ethHeader)
	if err != nil {
		return nil, nil, err
	}

	// If  verify succeeds, save consensusState first . this store is header_index
	consensusStatetmp, err := cdc.MarshalInterface(consensusState)
	if err != nil {
		return nil, nil, sdkerrors.Wrapf(ErrUnmarshalInterface, "can not unmarshal ConsensusState interface in CheckHeaderAndUpdateState ")

	}
	store.Set(ConsensusStateIndexKey(consensusState.Header.Hash()), consensusStatetmp)
	//Check the bifurcation
	if bytes.Equal(ethConsState.Header.Hash().Bytes(), ethHeader.ParentHash) {
		// set all consensusState by struct (prefix+hash , consensusState)
		store.Set(host.ConsensusStateKey(ethHeader.Height), consensusStatetmp)
	} else {
		err = m.RestructChain(cdc, store, *ethHeader)
		if err != nil {
			return nil, nil, err
		}
	}
	m.Header = *ethHeader
	newClientState.Header = *ethHeader

	return newClientState, consensusState, nil
}

// checkValidity checks if the eth header is valid.
func checkValidity(
	ctx sdk.Context,
	cdc codec.BinaryMarshaler,
	store sdk.KVStore,
	clientState *ClientState,
	consState *ConsensusState,
	header Header,

) error {
	if err := header.ValidateBasic(); err != nil {
		return err
	}

	return verifyHeader(ctx, cdc, store, clientState, header)
}

// update the RecentSingers and the ConsensusState.
func update(ctx sdk.Context,
	cdc codec.BinaryMarshaler,
	store sdk.KVStore,
	clientState *ClientState,
	header *Header,
) (*ClientState, *ConsensusState, error) {

	cs := &ConsensusState{
		Timestamp: header.Time,
		Number:    header.Height,
		Root:      header.Root,
		Header:    *header,
	}
	setConsensusMetadata(ctx, store, header.GetHeight())
	return clientState, cs, nil
}

func (m ClientState) RestructChain(cdc codec.BinaryMarshaler, store sdk.KVStore, new Header) error {
	si, ti := m.Header.Height, new.Height
	var err error
	current := m.Header
	//si > ti
	if si.RevisionHeight > ti.RevisionHeight {
		currentTmp := store.Get(host.ConsensusStateKey(ti))
		if currentTmp == nil {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find consensus state for height %s in RestructChain", ti)
		}
		var currenttmp exported.ConsensusState
		if err = cdc.UnmarshalInterface(currentTmp, &currenttmp); err != nil {
			return sdkerrors.Wrapf(ErrUnmarshalInterface, "can not unmarshal ConsensusState interface in RestructChain ")

		}
		tmpconsensus, ok := currenttmp.(*ConsensusState)
		if !ok {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find consensus state for height %s in RestructChain", ti)
		}
		current = tmpconsensus.Header
		si = ti
	}
	newHashs := make([]common.Hash, 0)

	for ti.RevisionHeight > si.RevisionHeight {
		newHashs = append(newHashs, new.Hash())
		newTmp := store.Get(ConsensusStateIndexKey(new.ToEthHeader().ParentHash))
		if newTmp == nil {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find consensus state for hash %s in RestructChain in RestructChain", new.ToEthHeader().ParentHash,
			)
		}
		var currenttmp exported.ConsensusState
		if err := cdc.UnmarshalInterface(newTmp, &currenttmp); err != nil {
			return sdkerrors.Wrapf(ErrUnmarshalInterface, "can not unmarshal ConsensusState interface in RestructChain ")
		}
		tmpconsensus, ok := currenttmp.(*ConsensusState)
		if !ok {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find consensus state for hash %s in RestructChain", new.ToEthHeader().ParentHash)
		}
		new = tmpconsensus.Header
		ti.RevisionHeight--
	}
	// si.parent != ti.parent
	for !bytes.Equal(current.ParentHash, new.ParentHash) {
		newHashs = append(newHashs, new.Hash())
		newTmp := store.Get(ConsensusStateIndexKey(new.ToEthHeader().ParentHash))
		if newTmp == nil {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find consensus state for hash %s in RestructChain", new.ToEthHeader().ParentHash)
		}
		var currenttmp exported.ConsensusState
		if err = cdc.UnmarshalInterface(newTmp, &currenttmp); err != nil {
			return sdkerrors.Wrapf(ErrUnmarshalInterface, "can not unmarshal ConsensusState interface in RestructChain ")
		}
		tmpconsensus, ok := currenttmp.(*ConsensusState)
		if !ok {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find consensus state for hash %s in RestructChain", new.ToEthHeader().ParentHash)
		}
		new = tmpconsensus.Header
		ti.RevisionHeight--
		si.RevisionHeight--
		currentTmp := store.Get(host.ConsensusStateKey(si))
		if currentTmp == nil {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find consensus state for height %s in RestructChain", si)
		}
		if err = cdc.UnmarshalInterface(currentTmp, &currenttmp); err != nil {
			return sdkerrors.Wrapf(ErrUnmarshalInterface, "can not unmarshal ConsensusState interface in RestructChain ")
		}
		tmpconsensus = currenttmp.(*ConsensusState)
		if !ok {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not  consensus state for height %s in RestructChain", si)
		}
		current = tmpconsensus.Header
	}
	for i := len(newHashs) - 1; i >= 0; i-- {
		newTmp := store.Get(ConsensusStateIndexKey(newHashs[i]))
		if newTmp == nil {
			return sdkerrors.Wrapf(
				clienttypes.ErrInvalidConsensus, "can not find consensus state for hash %s in RestructChain", newHashs[i])

		}
		// set main_chain
		store.Set(host.ConsensusStateKey(ti), newTmp)
		ti.RevisionHeight++
	}
	return err
}
