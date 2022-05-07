package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tharsis/ethermint/x/cronjobs/types"
)

// GetAllCronjobs - get all registered DevFeeInfo instances
func (k Keeper) GetAllCronjobs(ctx sdk.Context) []types.Cronjob {
	values := []types.Cronjob{}

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixCronjob)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var value types.Cronjob
		k.cdc.MustUnmarshal(iterator.Value(), &value)
		values = append(values, value)
	}

	return values
}

// GetAllCronjobs - get all registered DevFeeInfo instances
func (k Keeper) GetAllCronjobsByEpoch(ctx sdk.Context, epochIdentifier string) []types.Cronjob {
	values := []types.Cronjob{}
	ids := k.GetCronjobIdentifiersByEpoch(ctx, epochIdentifier)
	fmt.Println("--GetCronjobIdentifiersByEpoch", epochIdentifier, ids)
	for _, id := range ids.Identifiers {
		value, _ := k.GetCronjob(ctx, id)
		values = append(values, value)
	}
	return values
}

// IterateCronjobs iterates over all registered contracts and performs a
// callback with the corresponding DevFeeInfo.
func (k Keeper) IterateCronjobsByEpoch(
	ctx sdk.Context,
	epochIdentifier string,
	handlerFn func(cronjob types.Cronjob) (stop bool),
) {
	ids := k.GetCronjobIdentifiersByEpoch(ctx, epochIdentifier)
	for _, id := range ids.Identifiers {
		value, _ := k.GetCronjob(ctx, id)
		if handlerFn(value) {
			break
		}
	}
}

// GetCronjob
func (k Keeper) GetCronjob(ctx sdk.Context, id string) (types.Cronjob, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCronjob)
	bz := store.Get([]byte(id))
	if len(bz) == 0 {
		return types.Cronjob{}, false
	}

	var value types.Cronjob
	k.cdc.MustUnmarshal(bz, &value)
	return value, true
}

// GetCronjobIdentifiersByEpoch
func (k Keeper) GetCronjobIdentifiersByEpoch(ctx sdk.Context, epochIdentifier string) types.CronjobsByEpochIdentifier {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCronjobIdentifier)
	bz := store.Get([]byte(epochIdentifier))
	if len(bz) == 0 {
		return types.CronjobsByEpochIdentifier{Identifiers: make([]string, 0)}
	}

	var value types.CronjobsByEpochIdentifier
	k.cdc.MustUnmarshal(bz, &value)
	return value
}

// SetFee stores the developer fee information for a registered contract
func (k Keeper) SetCronjob(ctx sdk.Context, id string, cronjob types.Cronjob) {
	fmt.Println("---SetCronjob", id, cronjob)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCronjob)
	bz := k.cdc.MustMarshal(&cronjob)
	store.Set([]byte(id), bz)
	k.SetCronjobIdentifier(ctx, cronjob.EpochIdentifier, id)
}

// SetCronjobIdentifier
func (k Keeper) SetCronjobIdentifier(ctx sdk.Context, epochIdentifier string, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCronjobIdentifier)
	fmt.Println("---SetCronjobIdentifier id", id, epochIdentifier)
	ids := k.GetCronjobIdentifiersByEpoch(ctx, epochIdentifier)
	ids.Identifiers = append(ids.Identifiers, id)
	fmt.Println("---SetCronjobIdentifier ids.Identifiers", ids.Identifiers)
	bz := k.cdc.MustMarshal(&ids)
	store.Set([]byte(epochIdentifier), bz)
}

// DeleteFee removes a registered contract
func (k Keeper) DeleteCronjob(ctx sdk.Context, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCronjob)
	store.Delete([]byte(id))
}

// IsFeeRegistered checks if a contract was registered for receiving fees
func (k Keeper) IsCronjobRegistered(
	ctx sdk.Context,
	id string,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCronjob)
	return store.Has([]byte(id))
}
