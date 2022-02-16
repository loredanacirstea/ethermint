package extendedVM

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

// Ethermint additions
var PrecompiledContracts = map[common.Address]vm.PrecompiledContract{
	common.BytesToAddress([]byte{25}): &ibcPrecompile{},
}

// ibcPrecompile implemented as a native contract.
type ibcPrecompile struct{}

func (c *ibcPrecompile) RequiredGas(input []byte) uint64 {
	return 3000
}

func (c *ibcPrecompile) Run(evm *vm.EVM, caller vm.ContractRef, input []byte) ([]byte, error) {

	return []byte("hello"), nil
}
