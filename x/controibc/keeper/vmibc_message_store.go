package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tharsis/ethermint/x/controibc/types"
)

// SetVmIbcMessage stores a message
func (k Keeper) SetVmIbcMessage(ctx sdk.Context, data types.VmibcMessagePacketData) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixVmIbcMessage)
	key := data.GetID()
	bz := k.cdc.MustMarshal(&data)
	store.Set(key, bz)
	fmt.Println("--SetVmIbcMessage--", key, bz)
}

// getVmIbcMessage - get registered message from the identifier
func (k Keeper) GetVmIbcMessage(ctx sdk.Context, id []byte) (types.VmibcMessagePacketData, bool) {
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

	fmt.Println("--GetVmIbcMessage--", data)

	// A retrieval triggers deletion
	k.DeleteVmIbcMessage(ctx, data)
	return data, true
}

// DeleteVmIbcMessage removes the message
func (k Keeper) DeleteVmIbcMessage(ctx sdk.Context, data types.VmibcMessagePacketData) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixVmIbcMessage)
	key := data.GetID()
	store.Delete(key)
}
