package types

import (
	"bytes"
	"errors"

	"github.com/ethereum/go-ethereum/common"

	clienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	for {
		if err := checkValidity(cdc, store, &m, ethConsState, *ethHeader); err != nil {
			return nil, nil, err
		}
		break
	}
	newClientState, consensusState, err := update(cdc, store, &m, ethHeader)
	if err != nil {
		return nil, nil, err
	}

	// If  verify succeeds, save consensusState first . this store is header_index
	consensusStatetmp, err := cdc.MarshalInterface(consensusState)
	if err != nil {
		return nil, nil, err
	}
	store.Set(host.ConsensusStateIndexKey(string(consensusState.GetRoot().GetHash())), consensusStatetmp)

	//Check the bifurcation
	if bytes.Equal(ethConsState.Header.Hash().Bytes(), ethHeader.ParentHash) {
		// set all consensusState by struct (prefix+hash , consensusState)
		store.Set(host.ConsensusStateKey(ethHeader.Height), consensusStatetmp)
	} else {
		err := m.RestructChain(cdc, store, *ethHeader)
		if err != nil {
			return nil, nil, err
		}
	}
	return newClientState, consensusState, nil
}

// checkValidity checks if the eth header is valid.
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
	//if clientState.Header.Hash() != header.ToEthHeader().ParentHash {
	//	return errors.New("parent hash not same")
	//}
	return verifyHeader(cdc, store, clientState, header)
}

// update the RecentSingers and the ConsensusState.
func update(cdc codec.BinaryMarshaler,
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

	clientState.Header = *header
	return clientState, cs, nil
}
func (m ClientState) RestructChain(cdc codec.BinaryMarshaler, store sdk.KVStore, new Header) error {
	si, ti := m.Header.Height, new.Height
	var err error
	current := m.Header
	if si.RevisionHeight > ti.RevisionHeight {
		currentTmp := store.Get(host.ConsensusStateKey(ti))
		if currentTmp == nil {
			err = errors.New("no found ConsensusState")
			return err
		}
		err := cdc.UnmarshalInterface(currentTmp, current)
		if err != nil {
			return err
		}
		si = ti
	}
	newHashs := make([]common.Hash, 0)
	for ti.RevisionHeight > si.RevisionHeight {
		newHashs = append(newHashs, new.Hash())
		newTmp := store.Get(host.ConsensusStateIndexKey(string(new.ParentHash)))
		if newTmp == nil {
			err = errors.New("no found ConsensusState")
			return err
		}
		err := cdc.UnmarshalInterface(newTmp, new)
		if err != nil {
			return err
		}
		ti.RevisionHeight--
	}
	for bytes.Equal(current.ParentHash, new.ParentHash) {
		newHashs = append(newHashs, new.Hash())
		newTmp := store.Get(host.ConsensusStateIndexKey(string(new.ParentHash)))
		if newTmp == nil {
			err = errors.New("no found ConsensusState")
			return err
		}
		err := cdc.UnmarshalInterface(newTmp, new)
		if err != nil {
			return err
		}
		ti.RevisionHeight--
		si.RevisionHeight--
		currentTmp := store.Get(host.ConsensusStateKey(si))
		err = cdc.UnmarshalInterface(currentTmp, current)
		if err != nil {
			return err
		}
	}
	for i := len(newHashs) - 1; i >= 0; i-- {
		newTmp := store.Get(host.ConsensusStateIndexKey(newHashs[i].String()))
		if newTmp == nil {
			err = errors.New("no found ConsensusState")
			return err
		}
		store.Set(host.ConsensusStateKey(ti), newTmp)
		ti.RevisionHeight++
	}
	return err
}
