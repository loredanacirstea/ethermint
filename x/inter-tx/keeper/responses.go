package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tharsis/ethermint/x/inter-tx/types"
)

// SetResponse stores the response for a submitted EVM transaction
func (k Keeper) SetResponse(ctx sdk.Context, txKey []byte, response []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixResponse)
	store.Set(txKey, response)
}

// SetError stores the response for a submitted EVM transaction
func (k Keeper) SetError(ctx sdk.Context, txKey []byte, err string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixError)
	store.Set(txKey, []byte(err))
}

// GetResponse retrieves the response for a submitted EVM transaction
func (k Keeper) GetResponse(ctx sdk.Context, txKey []byte) []byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixResponse)
	return store.Get(txKey)
}

// GetError retrieves the response for a submitted EVM transaction
func (k Keeper) GetError(ctx sdk.Context, txKey []byte) []byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixError)
	return store.Get(txKey)
}
