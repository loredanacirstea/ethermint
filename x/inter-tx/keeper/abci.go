package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlock sets the sdk Context and EIP155 chain id to the Keeper.
func (k *Keeper) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {

}

// EndBlock also retrieves the bloom filter value from the transient store and commits it to the
// KVStore. The EVM end block logic doesn't update the validator set, thus it returns
// an empty slice.
func (k *Keeper) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	// abciev := ctx.EventManager().ABCIEvents()

	// fmt.Println("---intertx--EndBlock---ABCIEvents", abciev)

	// events, err := backend.AllTxLogsFromEvents(abciev)
	// fmt.Println("--intertx---EndBlock---events", events)
	// fmt.Println("---intertx--EndBlock---err", err)

	// attrs := abciev[0].GetAttributes()
	// fmt.Println("---intertx--EndBlock---GetAttributes", attrs)
	// // fmt.Println("---intertx--EndBlock---key", string(attrs[0].Key))
	// // fmt.Println("---intertx--EndBlock---GetAttributes", attrs[0])

	return []abci.ValidatorUpdate{}
}
