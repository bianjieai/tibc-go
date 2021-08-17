package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

var _ exported.ClientState = (*ClientState)(nil)

func (m *ClientState) ClientType() string {
	return exported.BSC
}

func (m *ClientState) GetLatestHeight() exported.Height {
	return m.Header.Height
}

func (m *ClientState) Validate() error {
	panic("implement me")
}

func (m *ClientState) GetDelayTime() uint64 {
	return uint64((2*len(m.Validators)/3 + 1)) * m.Period
}

func (m *ClientState) GetDelayBlock() uint64 {
	return uint64(2*len(m.Validators)/3 + 1)
}

func (m *ClientState) GetPrefix() exported.Prefix {
	panic("implement me")
}

func (m *ClientState) Initialize(
	ctx sdk.Context,
	cdc codec.BinaryMarshaler,
	store sdk.KVStore,
	state exported.ConsensusState,
) error {
	if m.Header.Height.RevisionHeight%m.Epoch != 0 {
		return sdkerrors.Wrapf(ErrInvalidGenesisBlock, "block: %d is not epoch block", m.Header.Height.RevisionHeight)
	}
	validatorBytes := m.Header.Extra[extraVanity : len(m.Header.Extra)-extraSeal]
	validators, err := ParseValidators(validatorBytes)
	if err != nil {
		return err
	}
	//TODO
	m.Validators = validators
	return nil
}

func (m *ClientState) Status(
	ctx sdk.Context,
	clientStore sdk.KVStore,
	cdc codec.BinaryMarshaler,
) exported.Status {
	onsState, err := GetConsensusState(clientStore, cdc, m.GetLatestHeight())
	if err != nil {
		return exported.Unknown
	}
	if onsState.Timestamp+m.TrustingPeriod > uint64(ctx.BlockTime().Nanosecond()) {
		return exported.Expired
	}
	return exported.Active
}

func (m *ClientState) ExportMetadata(store sdk.KVStore) []exported.GenesisMetadata {
	return nil
}

func (m *ClientState) CheckHeaderAndUpdateState(
	ctx sdk.Context,
	marshaler codec.BinaryMarshaler,
	store sdk.KVStore,
	header exported.Header,
) (exported.ClientState, exported.ConsensusState, error) {
	panic("implement me")
}

func (m *ClientState) VerifyPacketCommitment(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	proof []byte,
	sourceChain, destChain string,
	sequence uint64,
	commitmentBytes []byte,
) error {
	panic("implement me")
}

func (m *ClientState) VerifyPacketAcknowledgement(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	proof []byte,
	sourceChain, destChain string,
	sequence uint64,
	acknowledgement []byte,
) error {
	panic("implement me")
}

func (m *ClientState) VerifyPacketCleanCommitment(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	proof []byte,
	sourceChain, destChain string,
	sequence uint64,
) error {
	panic("implement me")
}
