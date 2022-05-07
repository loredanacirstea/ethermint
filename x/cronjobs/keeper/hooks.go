package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tharsis/ethermint/x/cronjobs/types"
	epochstypes "github.com/tharsis/ethermint/x/epochs/types"
)

// BeforeEpochStart: noop, We don't need to do anything here
func (k Keeper) BeforeEpochStart(_ sdk.Context, _ string, _ int64) {
}

// AfterEpochEnd mints and distributes coins at the end of each epoch end
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	fmt.Println("---Epoch---", epochIdentifier, epochNumber)

	k.IterateCronjobsByEpoch(ctx, epochIdentifier, func(cronjob types.Cronjob) bool {
		fmt.Println("---Running cron---", cronjob.Identifier, cronjob.EpochIdentifier)
		k.ExecuteCron(ctx, cronjob)
		return false
	})
}

// ___________________________________________________________________________________________________

// Hooks wrapper struct for incentives keeper
type Hooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// epochs hooks
func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}
