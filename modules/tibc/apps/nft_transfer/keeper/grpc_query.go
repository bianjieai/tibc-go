package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
)

var _ types.QueryServer = Keeper{}

// ClassTrace implements the Query/ClassTrace gRPC method
func (q Keeper) ClassTrace(c context.Context, req *types.QueryClassTraceRequest) (*types.QueryClassTraceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	hash, err := types.ParseHexHash(req.Hash)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid class trace hash %s, %s", req.Hash, err))
	}

	ctx := sdk.UnwrapSDKContext(c)
	classTrace, found := q.GetClassTrace(ctx, hash)
	if !found {
		return nil, status.Error(
			codes.NotFound,
			sdkerrors.Wrap(types.ErrTraceNotFound, req.Hash).Error(),
		)
	}

	return &types.QueryClassTraceResponse{
		ClassTrace: &classTrace,
	}, nil
}

// ClassTraces implements the Query/ClassTraces gRPC method
func (q Keeper) ClassTraces(c context.Context, req *types.QueryClassTracesRequest) (*types.QueryClassTracesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	traces := types.Traces{}
	store := prefix.NewStore(ctx.KVStore(q.storeKey), types.ClassTraceKey)

	pageRes, err := query.Paginate(store, req.Pagination, func(_, value []byte) error {
		result, err := q.UnmarshalClassTrace(value)
		if err != nil {
			return err
		}

		traces = append(traces, result)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryClassTracesResponse{
		ClassTraces: traces.Sort(),
		Pagination:  pageRes,
	}, nil
}
