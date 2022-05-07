package keeper

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/tharsis/ethermint/x/cronjobs/types"
)

var _ types.QueryServer = Keeper{}

// DevFeeInfos returns all registered contracts for fee distribution
func (k Keeper) Cronjobs(
	c context.Context,
	req *types.QueryCronjobsRequest,
) (*types.QueryCronjobsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var cronjobs []types.Cronjob
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCronjob)

	pageRes, err := query.Paginate(
		store,
		req.Pagination,
		func(key, value []byte) error {
			var cronjob types.Cronjob
			k.cdc.MustUnmarshal(value, &cronjob)
			cronjobs = append(cronjobs, cronjob)
			return nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryCronjobsResponse{
		Cronjobs:   cronjobs,
		Pagination: pageRes,
	}, nil
}

// DevFeeInfo returns a given registered contract
func (k Keeper) Cronjob(
	c context.Context,
	req *types.QueryCronjobRequest,
) (*types.QueryCronjobResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	if strings.TrimSpace(req.Identifier) == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"identifier is empty",
		)
	}

	cronjob, found := k.GetCronjob(ctx, req.Identifier)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"cronjob id '%s'",
			req.Identifier,
		)
	}

	return &types.QueryCronjobResponse{Cronjob: cronjob}, nil
}

// Params returns the fees module params
func (k Keeper) Params(
	c context.Context,
	_ *types.QueryParamsRequest,
) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}
