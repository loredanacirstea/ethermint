package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tharsis/ethermint/x/controibc/types"
)

// SetTokenPair stores a message
func (k Keeper) SetVmIbcMessage(ctx sdk.Context, data types.VmibcMessagePacketData) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixVmIbcMessage)
	key := data.GetID()
	bz := k.cdc.MustMarshal(&data)
	store.Set(key, bz)
}

// getVmIbcMessage - get registered message from the identifier
func (k Keeper) getVmIbcMessage(ctx sdk.Context, id []byte) (types.VmibcMessagePacketData, bool) {
	if id == nil {
		return types.VmibcMessagePacketData{}, false
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixVmIbcMessage)
	var data types.VmibcMessagePacketData
	bz := store.Get(id)
	if len(bz) == 0 {
		return types.VmibcMessagePacketData{}, false
	}

	k.cdc.MustUnmarshal(bz, &data)
	return data, true
}

// DeleteVmIbcMessage removes the message
func (k Keeper) DeleteVmIbcMessage(ctx sdk.Context, data types.VmibcMessagePacketData) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixVmIbcMessage)
	key := data.GetID()
	store.Delete(key)
}
