package extendedVM

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/holiman/uint256"
)

// Ethermint additions
var PrecompiledContracts = map[common.Address]vm.PrecompiledContract{
	common.BytesToAddress([]byte{25}): &evmPrecompile{},
}

// evmPrecompile implemented as a native contract.
type evmPrecompile struct{}

func (c *evmPrecompile) RequiredGas(input []byte) uint64 {
	return 3000
}

func (c *evmPrecompile) Run(evm *vm.EVM, caller vm.ContractRef, input []byte) ([]byte, error) {
	signature := common.Bytes2Hex(input[0:4])
	callInput := input[4:]
	var result []byte
	var err error

	fmt.Println("---evmPrecompile--", signature, callInput)

	switch signature {
	case "2b9416a8": // interpret(bytes,bytes,uint256,uint256)
		result, err = evmPrecompileInterpret(c, evm, caller, callInput)
	case "f750e68a": // analyze(bytes,bytes,uint256,uint256)
		result, err = evmPrecompileAnalyze(c, evm, caller, callInput)
	case "96bb50b3": // analyzeFrag(bytes,bytes32[],uint256,uint256)
		result, err = evmPrecompileAnalyzeFrag(c, evm, caller, callInput)
	case "349d0211": // part(bytes32,bytes32,bytes,bytes,uint256,uint256)
		result, err = nil, nil
	case "59e09d70": // partFrag(bytes32,bytes32,bytes,bytes32[],uint256,uint256)
		result, err = nil, nil
	default:
		return nil, errors.New("invalid evmPrecompile function")
	}
	if err != nil {
		return nil, err
	}

	fmt.Println("--evmPrecompile result--", common.Bytes2Hex(result))
	return result, err
}

// function interpret(bytes memory bytecode, bytes memory input, uint256 gas, uint256 value) view external returns(bytes memory result);
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

	ret, leftOverGas, err := innerEvm.CallWithBytecode(caller, contractAddress, bytecode, calldata, gas, value, nil)

	// TODO leftOverGas
	fmt.Println("------evmPrecompile--2--", ret, leftOverGas, nil)

	encodedResult := append(
		new(big.Int).SetUint64(32).FillBytes(make([]byte, 32)),
		new(big.Int).SetInt64(int64(len(ret))).FillBytes(make([]byte, 32))...,
	)
	encodedResult = append(encodedResult, ret...)

	padding := 32 - len(encodedResult)%32
	encodedResult = append(encodedResult, make([]byte, padding)...)

	return encodedResult, err
}

// function analyze(bytes memory bytecode, bytes memory input, uint256 gas, uint256 value) view external returns(uint256 pc, uint256 reads, uint256 writes, uint256 calls, uint256 gasUsed, bytes memory result);
func evmPrecompileAnalyze(c *evmPrecompile, evm *vm.EVM, caller vm.ContractRef, input []byte) ([]byte, error) {
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

	ret, leftOverGas, err := innerEvm.CallWithBytecode(caller, contractAddress, bytecode, calldata, gas, value, nil)

	interpreterContext := innerEvm.Interpreter().InterpretContext

	// uint256 pc, uint256 reads, uint256 writes, uint256 calls, uint256 msize, uint256 gasUsed, bytes memory result

	encodedResult := append(
		new(big.Int).SetUint64(interpreterContext.ReturnPc).FillBytes(make([]byte, 32)),
		new(big.Int).SetUint64(interpreterContext.Reads).FillBytes(make([]byte, 32))...,
	)
	encodedResult = append(encodedResult,
		new(big.Int).SetUint64(interpreterContext.Writes).FillBytes(make([]byte, 32))...,
	)
	encodedResult = append(encodedResult,
		new(big.Int).SetUint64(interpreterContext.Calls).FillBytes(make([]byte, 32))...,
	)
	encodedResult = append(encodedResult,
		new(big.Int).SetUint64(interpreterContext.Memsize).FillBytes(make([]byte, 32))...,
	)
	encodedResult = append(encodedResult,
		new(big.Int).SetUint64(leftOverGas).FillBytes(make([]byte, 32))...,
	)
	encodedResult = append(encodedResult,
		new(big.Int).SetUint64(uint64(224)).FillBytes(make([]byte, 32))...,
	)
	encodedResult = append(encodedResult,
		new(big.Int).SetInt64(int64(len(ret))).FillBytes(make([]byte, 32))...,
	)
	encodedResult = append(encodedResult, ret...)

	padding := len(encodedResult) % 32
	if padding > 0 {
		encodedResult = append(encodedResult, make([]byte, 32-padding)...)
	}

	return encodedResult, err
}

// function analyzeFrag(bytes memory bytecodeFragment, bytes32[] memory stack, uint256 gas, uint256 value) view external returns(uint256 pc, uint256 reads, uint256 writes, uint256 calls, uint256 memsize, uint256 gasUsed, bytes32[] memory stackOutput);
func evmPrecompileAnalyzeFrag(c *evmPrecompile, evm *vm.EVM, caller vm.ContractRef, input []byte) ([]byte, error) {
	offsetBytecode := new(big.Int).SetBytes(input[0:32]).Uint64()
	stackOffset := new(big.Int).SetBytes(input[32:64]).Uint64()
	gas := new(big.Int).SetBytes(input[64:96]).Uint64()
	value := new(big.Int).SetBytes(input[96:128])

	bytecodeLength := new(big.Int).SetBytes(input[offsetBytecode : offsetBytecode+32]).Uint64()
	bytecode := input[offsetBytecode+32 : offsetBytecode+32+bytecodeLength]

	stackCount := new(big.Int).SetBytes(input[stackOffset : stackOffset+32]).Uint64()
	var stackData []uint256.Int
	for i := uint64(0); i < stackCount; i++ {
		offset := stackOffset + 32 + i*32
		elem, overflow := uint256.FromBig(new(big.Int).SetBytes(input[offset : offset+32]))
		if overflow {
			return nil, errors.New("evmPrecompile stack overflow")
		}
		stackData = append(stackData, *elem)
	}

	innerEvm := vm.NewEVM(evm.Context, evm.TxContext, evm.StateDB, evm.ChainConfig(), evm.Config, evm.Precompiles)

	contractAddress := common.HexToAddress("0x0000000000000000000000000000000000000019")

	_, leftOverGas, err := innerEvm.CallWithBytecode(caller, contractAddress, bytecode, make([]byte, 0), gas, value, &stackData)

	interpreterContext := innerEvm.Interpreter().InterpretContext

	// uint256 pc, uint256 reads, uint256 writes, uint256 calls, uint256 memsize, uint256 gasUsed, bytes32[] memory stackOutput

	encodedResult := append(
		new(big.Int).SetUint64(interpreterContext.ReturnPc).FillBytes(make([]byte, 32)),
		new(big.Int).SetUint64(interpreterContext.Reads).FillBytes(make([]byte, 32))...,
	)
	encodedResult = append(encodedResult,
		new(big.Int).SetUint64(interpreterContext.Writes).FillBytes(make([]byte, 32))...,
	)
	encodedResult = append(encodedResult,
		new(big.Int).SetUint64(interpreterContext.Calls).FillBytes(make([]byte, 32))...,
	)
	encodedResult = append(encodedResult,
		new(big.Int).SetUint64(interpreterContext.Memsize).FillBytes(make([]byte, 32))...,
	)
	encodedResult = append(encodedResult,
		new(big.Int).SetUint64(leftOverGas).FillBytes(make([]byte, 32))...,
	)
	encodedResult = append(encodedResult,
		new(big.Int).SetUint64(uint64(224)).FillBytes(make([]byte, 32))...,
	)

	encodedResult = append(encodedResult,
		new(big.Int).SetUint64(uint64(len(interpreterContext.Stack))).FillBytes(make([]byte, 32))...,
	)
	for i := 0; i < len(interpreterContext.Stack); i++ {
		value := interpreterContext.Stack[i].Bytes32()
		encodedResult = append(encodedResult,
			value[:]...,
		)
	}

	return encodedResult, err
}
