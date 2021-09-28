package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bianjieai/tibc-go/modules/tibc/core/26-routing/types"
)

var _ types.QueryServer = (*Keeper)(nil)

//RoutingRules implements the Query/RoutingRules gRPC method
func (q Keeper) RoutingRules(c context.Context, req *types.QueryRoutingRulesRequest) (*types.QueryRoutingRulesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	routingRules, found := q.GetRoutingRules(ctx)
	if !found {
		return nil, status.Error(
			codes.NotFound,
			sdkerrors.Wrap(types.ErrRoutingRulesNotFound, "routing rules not found").Error(),
		)
	}

	return &types.QueryRoutingRulesResponse{Rules: routingRules}, nil
}
