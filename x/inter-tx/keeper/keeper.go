package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/keeper"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	abci "github.com/tendermint/tendermint/abci/types"
	evmkeeper "github.com/tharsis/ethermint/x/evm/keeper"
	"github.com/tharsis/ethermint/x/inter-tx/types"
)

type Keeper struct {
	cdc codec.Codec

	storeKey sdk.StoreKey

	scopedKeeper        capabilitykeeper.ScopedKeeper
	icaControllerKeeper icacontrollerkeeper.Keeper
	EvmKeeper           *evmkeeper.Keeper
	accountKeeper       authkeeper.AccountKeeper
	bankKeeper          types.BankKeeper
	feeCollector        string
	deliverTx           func(req abci.RequestDeliverTx) abci.ResponseDeliverTx
	commitFn            func() (res abci.ResponseCommit)
	moduleBasics        module.BasicManager // TODO use just use app.GetTxConfig()
	ClientCtx           *client.Context
}

func NewKeeper(cdc codec.Codec, storeKey sdk.StoreKey, iaKeeper icacontrollerkeeper.Keeper, scopedKeeper capabilitykeeper.ScopedKeeper, evmKeeper *evmkeeper.Keeper, bankKeeper types.BankKeeper, accountKeeper authkeeper.AccountKeeper, feeCollector string, deliverTx func(req abci.RequestDeliverTx) abci.ResponseDeliverTx, moduleBasics module.BasicManager, commitFn func() (res abci.ResponseCommit)) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,

		scopedKeeper:        scopedKeeper,
		icaControllerKeeper: iaKeeper,
		EvmKeeper:           evmKeeper,
		accountKeeper:       accountKeeper,
		bankKeeper:          bankKeeper,
		feeCollector:        feeCollector,
		deliverTx:           deliverTx,
		moduleBasics:        moduleBasics,
		commitFn:            commitFn,
	}
}

// Logger returns the application logger, scoped to the associated module
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s-%s", host.ModuleName, types.ModuleName))
}

// ClaimCapability claims the channel capability passed via the OnOpenChanInit callback
func (k *Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}
