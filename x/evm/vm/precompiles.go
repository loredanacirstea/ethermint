package extendedVM

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

// Ethermint additions
var PrecompiledContracts = map[common.Address]vm.PrecompiledContract{
	common.BytesToAddress([]byte{25}): &evmPrecompile{},
}

// ibcPrecompile implemented as a native contract.
type evmPrecompile struct{}

func (c *evmPrecompile) RequiredGas(input []byte) uint64 {
	return 3000
}

func (c *evmPrecompile) Run(evm *vm.EVM, caller vm.ContractRef, input []byte) ([]byte, error) {
	signature := common.Bytes2Hex(input[0:4])
	callInput := input[4:]
	var result []byte
	var err error

	fmt.Println("---ibcPrecompile--", signature, callInput)

	switch signature {
	case "2b9416a8": // interpret(bytes,bytes,uint256,uint256)
		result, err = evmPrecompileInterpret(c, evm, caller, callInput)
	default:
		return nil, errors.New("invalid ibcPrecompile function")
	}
	if err != nil {
		return nil, err
	}

	encodedResult := append(
		new(big.Int).SetUint64(32).FillBytes(make([]byte, 32)),
		new(big.Int).SetInt64(int64(len(result))).FillBytes(make([]byte, 32))...,
	)
	encodedResult = append(encodedResult, result...)

	padding := 32 - len(encodedResult)%32
	encodedResult = append(encodedResult, make([]byte, padding)...)

	fmt.Println("--ibcPrecompile result--", encodedResult)
	return encodedResult, err
}

func evmPrecompileInterpret(c *evmPrecompile, evm *vm.EVM, caller vm.ContractRef, input []byte) ([]byte, error) {
	offsetBytecode := new(big.Int).SetBytes(input[0:32]).Uint64()
	offsetInput := new(big.Int).SetBytes(input[32:64]).Uint64()
	gas := new(big.Int).SetBytes(input[64:96]).Uint64()
	value := new(big.Int).SetBytes(input[96:128])

	bytecodeLength := new(big.Int).SetBytes(input[offsetBytecode : offsetBytecode+32]).Uint64()
	bytecode := input[offsetBytecode+32 : offsetBytecode+32+bytecodeLength]

	calldataLength := new(big.Int).SetBytes(input[offsetInput : offsetInput+32]).Uint64()
	calldata := input[offsetInput+32 : offsetInput+32+calldataLength]

	innerEvm := vm.NewEVM(evm.Context, evm.TxContext, evm.StateDB, evm.ChainConfig(), evm.Config, evm.Precompiles)

	contractAddress := common.HexToAddress("0x0000000000000000000000000000000000000019")

	ret, leftOverGas, err := innerEvm.CallWithBytecode(caller, contractAddress, bytecode, calldata, gas, value)

	// TODO leftOverGas
	fmt.Println("------evmPrecompile--2--", ret, leftOverGas, nil)

	return ret, err
}
