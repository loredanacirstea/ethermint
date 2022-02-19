package keeper

import (
	"github.com/tharsis/ethermint/x/controibc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.PortId(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}


// PortId returns the PortId param
func (k Keeper) PortId(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyPortId, &res)
	return
}
