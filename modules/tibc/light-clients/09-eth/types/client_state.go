package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/light"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"

	commitmenttypes "github.com/bianjieai/tibc-go/modules/tibc/core/23-commitment/types"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
)

var _ exported.ClientState = (*ClientState)(nil)

func (m ClientState) ClientType() string {
	return exported.ETH
}

func (m ClientState) GetLatestHeight() exported.Height {
	return m.Header.Height
}

func (m ClientState) Validate() error {
	return m.Header.ValidateBasic()
}

func (m ClientState) GetDelayTime() uint64 {
	return m.TimeDelay
}

func (m ClientState) GetDelayBlock() uint64 {
	return m.BlockDelay
}

func (m ClientState) GetPrefix() exported.Prefix {
	return commitmenttypes.MerklePrefix{}
}

func (m ClientState) Initialize(
	ctx sdk.Context,
	cdc codec.BinaryMarshaler,
	store sdk.KVStore,
	state exported.ConsensusState,
) error {
	marshalInterface, err := cdc.MarshalInterface(state)
	if err != nil {
		return sdkerrors.Wrap(ErrInvalidGenesisBlock, "marshal consensus to interface failed")
	}
	consensusState := state.(*ConsensusState)
	header := consensusState.Header.ToEthHeader()
	store.Set(ConsensusStateIndexKey(header.Hash()), marshalInterface)

	return nil
}

func (m ClientState) Status(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
) exported.Status {
	onsState, err := GetConsensusState(store, cdc, m.GetLatestHeight())
	if err != nil {
		return exported.Unknown
	}
	if onsState.Timestamp+m.GetDelayTime() < uint64(ctx.BlockTime().Nanosecond()) {
		return exported.Expired
	}
	return exported.Active
}

func (m ClientState) ExportMetadata(store sdk.KVStore) []exported.GenesisMetadata {
	//TODO
	return nil
}

func (m ClientState) VerifyPacketCommitment(
	ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	height exported.Height,
	proof []byte,
	sourceChain, destChain string,
	sequence uint64,
	commitment []byte,
) error {
	ethProof, consensusState, err := produceVerificationArgs(store, cdc, m, height, proof)
	if err != nil {
		return err
	}

	// check delay period has passed
	delayBlock := m.Header.Height.RevisionHeight - height.GetRevisionHeight()
	if delayBlock < m.GetDelayBlock() {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"delay block (%d) < client state delay block (%d)",
			delayBlock, m.GetDelayBlock(),
		)
	}
	// verify that the provided commitment has been stored
	return verifyMerkleProof(ethProof, consensusState, m.ContractAddress, commitment)
}

func (m ClientState) VerifyPacketAcknowledgement(
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

func (m ClientState) VerifyPacketCleanCommitment(
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

// produceVerificationArgs performs the basic checks on the arguments that are
// shared between the verification functions and returns the unmarshal
// merkle proof, the consensus state and an error if one occurred.
func produceVerificationArgs(
	store sdk.KVStore,
	cdc codec.BinaryMarshaler,
	cs ClientState,
	height exported.Height,
	proof []byte,
) (merkleProof Proof, consensusState *ConsensusState, err error) {
	if cs.GetLatestHeight().LT(height) {
		return Proof{}, nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"client state height < proof height (%d < %d)", cs.GetLatestHeight(), height,
		)
	}

	if proof == nil {
		return Proof{}, nil, sdkerrors.Wrap(ErrInvalidProof, "proof cannot be empty")
	}

	if err = json.Unmarshal(proof, &merkleProof); err != nil {
		return Proof{}, nil, sdkerrors.Wrap(ErrInvalidProof, "failed to unmarshal proof into proof")
	}

	consensusState, err = GetConsensusState(store, cdc, height)
	if err != nil {
		return Proof{}, nil, err
	}
	return merkleProof, consensusState, nil
}

func verifyMerkleProof(ethProof Proof,
	consensusState *ConsensusState,
	contractAddr []byte,
	commitment []byte,
) error {
	//1. prepare verify account
	nodeList := new(light.NodeList)

	for _, s := range ethProof.AccountProof {
		_ = nodeList.Put(nil, common.FromHex(s))
	}
	ns := nodeList.NodeSet()

	addr := common.FromHex(ethProof.Address)
	if !bytes.Equal(addr, contractAddr) {
		return fmt.Errorf("verifyMerkleProof, contract address is error, proof address: %s, side chain address: %s", ethProof.Address, hex.EncodeToString(contractAddr))
	}
	acctKey := crypto.Keccak256(addr)

	//2. verify account proof
	root := common.BytesToHash(consensusState.Root)
	acctVal, err := trie.VerifyProof(root, acctKey, ns)
	if err != nil {
		return fmt.Errorf("verifyMerkleProof, verify account proof error:%s", err)
	}

	storageHash := common.HexToHash(ethProof.StorageHash)
	codeHash := common.HexToHash(ethProof.CodeHash)
	nonce := common.HexToHash(ethProof.Nonce).Big()
	balance := common.HexToHash(ethProof.Balance).Big()

	acct := &ProofAccount{
		Nonce:    nonce,
		Balance:  balance,
		Storage:  storageHash,
		Codehash: codeHash,
	}

	accRlp, err := rlp.EncodeToBytes(acct)
	if err != nil {
		return err
	}

	if !bytes.Equal(accRlp, acctVal) {
		return fmt.Errorf("verifyMerkleProof, verify account proof failed, wanted:%v, get:%v", accRlp, acctVal)
	}

	//3.verify storage proof
	nodeList = new(light.NodeList)
	if len(ethProof.StorageProof) != 1 {
		return fmt.Errorf("verifyMerkleProof, invalid storage proof format")
	}

	sp := ethProof.StorageProof[0]
	storageKey := crypto.Keccak256(common.HexToHash(sp.Key).Bytes())

	for _, prf := range sp.Proof {
		_ = nodeList.Put(nil, common.FromHex(prf))
	}

	ns = nodeList.NodeSet()
	val, err := trie.VerifyProof(storageHash, storageKey, ns)
	if err != nil {
		return fmt.Errorf("verifyMerkleProof, verify storage proof error:%s", err)
	}

	if !checkProofResult(val, commitment) {
		return fmt.Errorf("verifyMerkleProof, verify storage result failed")
	}
	return nil
}

func checkProofResult(result, value []byte) bool {
	var tempBytes []byte
	err := rlp.DecodeBytes(result, &tempBytes)
	if err != nil {
		return false
	}
	//
	var s []byte
	for i := len(tempBytes); i < 32; i++ {
		s = append(s, 0)
	}
	s = append(s, tempBytes...)
	// TODO
	//hash := crypto.Keccak256(value)
	return bytes.Equal(s, value)
}
