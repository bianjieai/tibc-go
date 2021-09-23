package utils

import (
	"context"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	commitmenttypes "github.com/bianjieai/tibc-go/modules/tibc/core/23-commitment/types"
	host "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"
	tibcclient "github.com/bianjieai/tibc-go/modules/tibc/core/client"
	"github.com/bianjieai/tibc-go/modules/tibc/core/exported"
	tibctmtypes "github.com/bianjieai/tibc-go/modules/tibc/light-clients/07-tendermint/types"
)

// QueryClientState returns a client state. If prove is true, it performs an ABCI store query
// in order to retrieve the merkle proof. Otherwise, it uses the gRPC query client.
func QueryClientState(
	clientCtx client.Context, chainName string, prove bool,
) (*types.QueryClientStateResponse, error) {
	if prove {
		return QueryClientStateABCI(clientCtx, chainName)
	}

	queryClient := types.NewQueryClient(clientCtx)
	req := &types.QueryClientStateRequest{
		ChainName: chainName,
	}

	return queryClient.ClientState(context.Background(), req)
}

// QueryClientStateABCI queries the store to get the light client state and a merkle proof.
func QueryClientStateABCI(
	clientCtx client.Context, chainName string,
) (*types.QueryClientStateResponse, error) {
	key := host.FullClientStateKey(chainName)

	value, proofBz, proofHeight, err := tibcclient.QueryTendermintProof(clientCtx, key)
	if err != nil {
		return nil, err
	}

	// check if client exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrap(types.ErrClientNotFound, chainName)
	}

	cdc := codec.NewProtoCodec(clientCtx.InterfaceRegistry)

	clientState, err := types.UnmarshalClientState(cdc, value)
	if err != nil {
		return nil, err
	}

	anyClientState, err := types.PackClientState(clientState)
	if err != nil {
		return nil, err
	}

	clientStateRes := types.NewQueryClientStateResponse(anyClientState, proofBz, proofHeight)
	return clientStateRes, nil
}

// QueryConsensusState returns a consensus state. If prove is true, it performs an ABCI store
// query in order to retrieve the merkle proof. Otherwise, it uses the gRPC query client.
func QueryConsensusState(
	clientCtx client.Context, chainName string, height exported.Height, prove, latestHeight bool,
) (*types.QueryConsensusStateResponse, error) {
	cdc := codec.NewProtoCodec(clientCtx.InterfaceRegistry)
	if height == nil || height.IsZero() {
		res, err := QueryClientState(clientCtx, chainName, false)
		if err != nil {
			return nil, err
		}
		var clientState exported.ClientState
		if err := cdc.UnpackAny(res.ClientState, &clientState); err != nil {
			return nil, err
		}
		height = clientState.GetLatestHeight()
	}
	if prove {
		return QueryConsensusStateABCI(clientCtx, chainName, height)
	}

	queryClient := types.NewQueryClient(clientCtx)
	req := &types.QueryConsensusStateRequest{
		ChainName:      chainName,
		RevisionNumber: height.GetRevisionNumber(),
		RevisionHeight: height.GetRevisionHeight(),
		LatestHeight:   latestHeight,
	}

	return queryClient.ConsensusState(context.Background(), req)
}

// QueryConsensusStateABCI queries the store to get the consensus state of a light client and a
// merkle proof of its existence or non-existence.
func QueryConsensusStateABCI(
	clientCtx client.Context, chainName string, height exported.Height,
) (*types.QueryConsensusStateResponse, error) {
	cdc := codec.NewProtoCodec(clientCtx.InterfaceRegistry)
	key := host.FullConsensusStateKey(chainName, height)

	value, proofBz, proofHeight, err := tibcclient.QueryTendermintProof(clientCtx, key)
	if err != nil {
		return nil, err
	}

	// check if consensus state exists
	if len(value) == 0 {
		return nil, sdkerrors.Wrap(types.ErrConsensusStateNotFound, chainName)
	}

	cs, err := types.UnmarshalConsensusState(cdc, value)
	if err != nil {
		return nil, err
	}

	anyConsensusState, err := types.PackConsensusState(cs)
	if err != nil {
		return nil, err
	}

	return types.NewQueryConsensusStateResponse(anyConsensusState, proofBz, proofHeight), nil
}

// QueryTendermintHeader takes a client context and returns the appropriate
// tendermint header
func QueryTendermintHeader(clientCtx client.Context) (tibctmtypes.Header, int64, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return tibctmtypes.Header{}, 0, err
	}

	info, err := node.ABCIInfo(context.Background())
	if err != nil {
		return tibctmtypes.Header{}, 0, err
	}

	var height int64
	if clientCtx.Height != 0 {
		height = clientCtx.Height
	} else {
		height = info.Response.LastBlockHeight
	}

	commit, err := node.Commit(context.Background(), &height)
	if err != nil {
		return tibctmtypes.Header{}, 0, err
	}

	page := 1
	count := 10_000

	validators, err := node.Validators(context.Background(), &height, &page, &count)
	if err != nil {
		return tibctmtypes.Header{}, 0, err
	}

	protoCommit := commit.SignedHeader.ToProto()
	protoValset, err := tmtypes.NewValidatorSet(validators.Validators).ToProto()
	if err != nil {
		return tibctmtypes.Header{}, 0, err
	}

	header := tibctmtypes.Header{
		SignedHeader: protoCommit,
		ValidatorSet: protoValset,
	}

	return header, height, nil
}

// QueryNodeConsensusState takes a client context and returns the appropriate
// tendermint consensus state
func QueryNodeConsensusState(clientCtx client.Context) (*tibctmtypes.ConsensusState, int64, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return &tibctmtypes.ConsensusState{}, 0, err
	}

	info, err := node.ABCIInfo(context.Background())
	if err != nil {
		return &tibctmtypes.ConsensusState{}, 0, err
	}

	var height int64
	if clientCtx.Height != 0 {
		height = clientCtx.Height
	} else {
		height = info.Response.LastBlockHeight
	}

	commit, err := node.Commit(context.Background(), &height)
	if err != nil {
		return &tibctmtypes.ConsensusState{}, 0, err
	}

	page := 1
	count := 10_000

	nextHeight := height + 1
	nextVals, err := node.Validators(context.Background(), &nextHeight, &page, &count)
	if err != nil {
		return &tibctmtypes.ConsensusState{}, 0, err
	}

	state := &tibctmtypes.ConsensusState{
		Timestamp:          commit.Time,
		Root:               commitmenttypes.NewMerkleRoot(commit.AppHash),
		NextValidatorsHash: tmtypes.NewValidatorSet(nextVals.Validators).Hash(),
	}

	return state, height, nil
}
