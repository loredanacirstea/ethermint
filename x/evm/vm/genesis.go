package vmibc

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tharsis/ethermint/x/evm/vm/keeper"
	"github.com/tharsis/ethermint/x/evm/vm/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, genState types.GenesisState) {

	fmt.Println("-------InitGenesis---------", genState)
	fmt.Println("-------InitGenesis-PortId--------", genState.PortId)

	// k.SetPort(ctx, genState.PortId)

	// Only try to bind to port if it is not already bound, since we may already own
	// port capability from capability InitGenesis
	if !keeper.IsBound(ctx, genState.PortId) {
		// module binds to desired ports on InitChain
		// and claims returned capabilities
		cap1 := keeper.PortKeeper().BindPort(ctx, genState.PortId)

		// NOTE: The module's scoped capability keeper must be private
		keeper.ScopedKeeper().ClaimCapability(ctx, cap1, "vmibc")
	}
	// k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	// genesis.Params = k.GetParams(ctx)

	genesis.PortId = k.GetPort(ctx)
	// genesis.MessageList = k.GetAllMessage(ctx)
	// genesis.MessageCount = k.GetMessageCount(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
