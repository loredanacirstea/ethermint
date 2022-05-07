package cronjobs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tharsis/ethermint/x/cronjobs/keeper"
	"github.com/tharsis/ethermint/x/cronjobs/types"
)

// InitGenesis import module genesis
func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	data types.GenesisState,
) {
	// k.SetParams(ctx, data.Params)

	for _, cronjob := range data.Cronjobs {
		k.SetCronjob(ctx, cronjob.Identifier, cronjob)
	}
}

// ExportGenesis export module state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:   k.GetParams(ctx),
		Cronjobs: k.GetAllCronjobs(ctx),
	}
}
