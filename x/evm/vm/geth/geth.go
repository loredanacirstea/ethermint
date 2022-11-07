package geth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"

	evm "github.com/evmos/ethermint/x/evm/vm"
)

var (
	_ evm.EVM         = (*EVM)(nil)
	_ evm.Constructor = NewEVM
)

// EVM is the wrapper for the go-ethereum EVM.
type EVM struct {
	*vm.EVM
	ctx                    sdk.Context
	getPrecompilesExtended func(ctx sdk.Context, evm *vm.EVM) evm.PrecompiledContracts
}

// NewEVM defines the constructor function for the go-ethereum (geth) EVM. It uses
// the default precompiled contracts and the EVM concrete implementation from
// geth.
func NewEVM(
	ctx sdk.Context,
	blockCtx vm.BlockContext,
	txCtx vm.TxContext,
	stateDB vm.StateDB,
	chainConfig *params.ChainConfig,
	config vm.Config,
	getPrecompilesExtended func(ctx sdk.Context, evm *vm.EVM) evm.PrecompiledContracts,
) evm.EVM {
	newEvm := &EVM{
		EVM:                    vm.NewEVM(blockCtx, txCtx, stateDB, chainConfig, config),
		ctx:                    ctx,
		getPrecompilesExtended: getPrecompilesExtended,
	}
	newEvm.EVM.GetActivePrecompiles = newEvm.GetActivePrecompiles
	return newEvm
}

// Context returns the EVM's Block Context
func (e EVM) Context() vm.BlockContext {
	return e.EVM.Context
}

// TxContext returns the EVM's Tx Context
func (e EVM) TxContext() vm.TxContext {
	return e.EVM.TxContext
}

// Config returns the configuration options for the EVM.
func (e EVM) Config() vm.Config {
	return e.EVM.Config
}

// // Precompile returns the precompiled contract associated with the given address
// // and the current chain configuration. If the contract cannot be found it returns
// // nil.
// func (e EVM) Precompile(addr common.Address) (p vm.PrecompiledContract, found bool) {
// 	p, found = e.precompiles[addr]
// 	return p, found
// }

func (e EVM) GetActivePrecompiles(rules params.Rules) map[common.Address]vm.PrecompiledContract {
	precompiles := vm.DefaultActivePrecompileMap(rules)
	customPrecompiles := e.getPrecompilesExtended(e.ctx, e.EVM)

	for k, v := range customPrecompiles {
		precompiles[k] = v
	}
	return precompiles
}

// ActivePrecompiles returns a list of all the active precompiled contract addresses
// for the current chain configuration.
func (e EVM) ActivePrecompiles(rules params.Rules) []common.Address {
	return e.EVM.ActivePrecompiles(rules)
}
