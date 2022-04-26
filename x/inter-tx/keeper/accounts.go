package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	"github.com/tharsis/ethermint/x/inter-tx/types"
)

func (k Keeper) GenerateAbstractAccount(ctx sdk.Context) (types.AbstractAccount, error) {
	privKey, err := ethsecp256k1.GenerateKey()
	if err != nil {
		return types.AbstractAccount{}, err
	}
	account := types.AbstractAccount{
		PrivKey: privKey.Bytes(),
		Type:    "eth_secp256k1",
	}
	return account, nil
}

func (k Keeper) SetAbstractAccount(ctx sdk.Context, ica string, account types.AbstractAccount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAbstractAccount)
	bz := k.cdc.MustMarshal(&account)
	store.Set([]byte(ica), bz)
}

func (k Keeper) GetAbstractAccount(ctx sdk.Context, ica string) (types.AbstractAccount, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAbstractAccount)
	bz := store.Get([]byte(ica))
	if len(bz) == 0 {
		return types.AbstractAccount{}, false
	}
	var account types.AbstractAccount
	k.cdc.MustUnmarshal(bz, &account)
	account.PubKey = account.PrivKey
	return account, true
}

func (k Keeper) GetAbstractAccountHydrated(ctx sdk.Context, ica string) (*ethsecp256k1.PrivKey, sdk.AccAddress, bool) {
	account, found := k.GetAbstractAccount(ctx, ica)
	if !found {
		return nil, nil, false
	}
	priv := &ethsecp256k1.PrivKey{
		Key: account.PrivKey,
	}
	address := sdk.AccAddress(priv.PubKey().Address().Bytes())

	return priv, address, true
}
