package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/tharsis/ethermint/x/cronjobs/types"
	evmkeeper "github.com/tharsis/ethermint/x/evm/keeper"
	intertxkeeper "github.com/tharsis/ethermint/x/inter-tx/keeper"
)

// Keeper of this module maintains collections of fee information for contracts
// registered to receive transaction fees.
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramstore paramtypes.Subspace

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	EvmKeeper     *evmkeeper.Keeper
	// abstractAccountKeeper types.AbstractAccountKeeper
	AbstractAccountKeeper intertxkeeper.Keeper
	feeCollectorName      string
}

// NewKeeper creates new instances of the fees Keeper
func NewKeeper(
	storeKey sdk.StoreKey,
	cdc codec.BinaryCodec,
	ps paramtypes.Subspace,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	evmKeeper *evmkeeper.Keeper,
	// aak types.AbstractAccountKeeper,
	aak intertxkeeper.Keeper,
	feeCollector string,
) Keeper {
	// set KeyTable if it has not already been set
	// if !ps.HasKeyTable() {
	// 	ps = ps.WithKeyTable(types.ParamKeyTable())
	// }

	return Keeper{
		storeKey:              storeKey,
		cdc:                   cdc,
		paramstore:            ps,
		accountKeeper:         ak,
		bankKeeper:            bk,
		EvmKeeper:             evmKeeper,
		AbstractAccountKeeper: aak,
		feeCollectorName:      feeCollector,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
