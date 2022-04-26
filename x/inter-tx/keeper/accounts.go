package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// "github.com/ethereum/go-ethereum/crypto"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	"github.com/tharsis/ethermint/x/inter-tx/types"
)

func (k Keeper) GenerateAbstractAccount(ctx sdk.Context, ica sdk.AccAddress) (types.AbstractAccount, error) {
	privKey, err := ethsecp256k1.GenerateKey()
	if err != nil {
		return types.AbstractAccount{}, err
	}
	account := types.AbstractAccount{
		// PrivKey: crypto.FromECDSA(privKey),
		PrivKey: privKey.Bytes(),
		Type:    "eth_secp256k1",
	}
	return account, nil
}

func (k Keeper) SetAbstractAccount(ctx sdk.Context, ica sdk.AccAddress, account types.AbstractAccount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAbstractAccount)
	bz := k.cdc.MustMarshal(&account)
	store.Set(ica.Bytes(), bz)
}

func (k Keeper) GetAbstractAccount(ctx sdk.Context, ica sdk.AccAddress) (types.AbstractAccount, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixAbstractAccount)
	bz := store.Get(ica.Bytes())
	if len(bz) == 0 {
		return types.AbstractAccount{}, false
	}
	var account types.AbstractAccount
	k.cdc.MustUnmarshal(bz, &account)
	return account, true
}
