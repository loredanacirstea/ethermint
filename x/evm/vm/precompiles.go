package extendedVM

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	cronjobstypes "github.com/tharsis/ethermint/x/cronjobs/types"
	"github.com/tharsis/ethermint/x/evm/types"
	intertxtypes "github.com/tharsis/ethermint/x/inter-tx/types"
)

// Ethermint additions
func GetPrecompiles(ctx sdk.Context, intertxKeeper types.InterTxKeeper, cronjobsKeeper types.CronjobsKeeper) map[common.Address]vm.PrecompiledContract {
	return map[common.Address]vm.PrecompiledContract{
		common.BytesToAddress([]byte{25}): &ICAPrecompile{ctx, intertxKeeper},
		common.BytesToAddress([]byte{26}): &AbstractAccountPrecompile{ctx, intertxKeeper},
		common.BytesToAddress([]byte{27}): &CronjobsPrecompile{ctx, cronjobsKeeper},
	}
}

// ICAPrecompile implemented as a native contract.
type ICAPrecompile struct {
	ctx           sdk.Context
	intertxKeeper types.InterTxKeeper
}

func (c *ICAPrecompile) RequiredGas(input []byte) uint64 {
	return 3000
}

func (c *ICAPrecompile) Run(evm *vm.EVM, caller vm.ContractRef, input []byte) ([]byte, error) {
	signature := common.Bytes2Hex(input[0:4])
	callInput := input[4:]
	var result []byte
	var err error

	switch signature {
	case "8fca7148": // emitTx(address,uint256,uint256,bytes,address,uint256)
		result, err = emitTx(evm, c, caller, callInput)
	case "415a2bc1": // getResponse(bytes32)
		result, err = getTxResponse(c, caller, callInput)
	default:
		return nil, errors.New("invalid ICAPrecompile function")
	}
	fmt.Println("--ICAPrecompile result--", result, err)
	if err != nil {
		return nil, err
	}
	encodedResult := append(
		new(big.Int).SetUint64(32).FillBytes(make([]byte, 32)),
		new(big.Int).SetInt64(int64(len(result))).FillBytes(make([]byte, 32))...,
	)
	encodedResult = append(encodedResult, result...)

	padding := len(encodedResult) % 32
	if padding > 0 {
		encodedResult = append(encodedResult, make([]byte, 32-padding)...)
	}

	fmt.Println("--ICAPrecompile result--", encodedResult)
	return encodedResult, nil
}

func emitTx(evm *vm.EVM, c *ICAPrecompile, caller vm.ContractRef, input []byte) ([]byte, error) {
	connectionId := "connection-0"
	to := common.BytesToAddress(input[0:32])
	value := new(big.Int).SetBytes(input[32:64])
	gasLimit := new(big.Int).SetBytes(input[64:96]).Uint64()

	offsetInput := new(big.Int).SetBytes(input[96:128]).Uint64()
	calldataLength := new(big.Int).SetBytes(input[offsetInput : offsetInput+32]).Uint64()
	data := input[offsetInput+32 : offsetInput+32+calldataLength]
	owner := common.BytesToAddress(input[128:160])
	nonce := new(big.Int).SetBytes(input[160:192]).Uint64()
	gasPrice := big.NewInt(20)
	gasFeeCap := big.NewInt(20)
	gasTipCap := big.NewInt(20)

	accesses := &ethtypes.AccessList{}
	ethtx := types.NewTx(evm.ChainConfig().ChainID, nonce, &to, value, gasLimit, gasPrice, gasFeeCap, gasTipCap, data, accesses)
	msg, err := intertxtypes.NewMsgSubmitEthereumTx(ethtx, connectionId, owner.Hex())
	if err != nil {
		return nil, err
	}

	res, err := c.intertxKeeper.SubmitEthereumTx(sdk.WrapSDKContext(c.ctx), msg)
	fmt.Println("----SubmitTx res----", res, err)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func getTxResponse(c *ICAPrecompile, caller vm.ContractRef, key []byte) ([]byte, error) {
	err := c.intertxKeeper.GetError(c.ctx, key)
	if err != nil {
		return err, nil
	}
	res := c.intertxKeeper.GetResponse(c.ctx, key)
	return res, nil
}

// AbstractAccountPrecompile implemented as a native contract.
type AbstractAccountPrecompile struct {
	ctx           sdk.Context
	intertxKeeper types.InterTxKeeper
}

func (c *AbstractAccountPrecompile) RequiredGas(input []byte) uint64 {
	return 3000
}

func (c *AbstractAccountPrecompile) Run(evm *vm.EVM, caller vm.ContractRef, input []byte) ([]byte, error) {
	signature := common.Bytes2Hex(input[0:4])
	callInput := input[4:]
	var result []byte
	var err error

	switch signature {
	case "d9f226e9": // registerAccount()
		result, err = registerAccount(evm, c, caller)
	case "e6630840": // sendTx(address,uint256,uint256,bytes)
		result, err = forwardTx(evm, c, caller, callInput)
	default:
		return nil, errors.New("invalid AbstractAccountPrecompile function")
	}
	fmt.Println("--AbstractAccountPrecompile result--", result, err)
	if err != nil {
		return nil, err
	}
	encodedResult := append(
		new(big.Int).SetUint64(32).FillBytes(make([]byte, 32)),
		new(big.Int).SetInt64(int64(len(result))).FillBytes(make([]byte, 32))...,
	)
	encodedResult = append(encodedResult, result...)

	padding := len(encodedResult) % 32
	if padding > 0 {
		encodedResult = append(encodedResult, make([]byte, 32-padding)...)
	}

	fmt.Println("--AbstractAccountPrecompile result--", encodedResult)
	return encodedResult, err
}

func registerAccount(evm *vm.EVM, c *AbstractAccountPrecompile, caller vm.ContractRef) ([]byte, error) {
	_, err := c.intertxKeeper.RegisterAccount(sdk.WrapSDKContext(c.ctx), &intertxtypes.MsgRegisterAccount{
		Owner:        caller.Address().Hex(),
		ConnectionId: "connection-0", // TODO from input
	})
	if err != nil {
		return nil, err
	}
	return make([]byte, 0), nil
}

func sendTx(evm *vm.EVM, c *AbstractAccountPrecompile, caller vm.ContractRef, input []byte) ([]byte, error) {
	query := &intertxtypes.QueryInterchainAccountFromAddressRequest{}
	_, err := c.intertxKeeper.InterchainAccountFromAddress(sdk.WrapSDKContext(c.ctx), query)
	if err != nil {
		_, err := registerAccount(evm, c, caller)
		if err != nil {
			return nil, err
		}
		_, err = c.intertxKeeper.InterchainAccountFromAddress(sdk.WrapSDKContext(c.ctx), query)
		if err != nil {
			return nil, err
		}
	}

	fmt.Println("--AbstractAccountPrecompile submitTx-input-", input)
	to := common.BytesToAddress(input[0:32])
	gasLimit := new(big.Int).SetBytes(input[32:64]).Uint64()

	offsetInput := new(big.Int).SetBytes(input[64:96]).Uint64()
	calldataLength := new(big.Int).SetBytes(input[offsetInput : offsetInput+32]).Uint64()
	data := input[offsetInput+32 : offsetInput+32+calldataLength]
	fmt.Println("--AbstractAccountPrecompile data--", data)

	nonce := uint64(0)
	value := big.NewInt(0)
	gasPrice := big.NewInt(20)
	gasFeeCap := big.NewInt(20)
	gasTipCap := big.NewInt(20)
	accesses := &ethtypes.AccessList{}
	ethtx := types.NewTx(evm.ChainConfig().ChainID, nonce, &to, value, gasLimit, gasPrice, gasFeeCap, gasTipCap, data, accesses)
	fmt.Println("ethtx", ethtx)
	msg, err := intertxtypes.NewMsgForwardEthereumTx(ethtx, caller.Address().Hex())
	if err != nil {
		return nil, err
	}

	c.intertxKeeper.ForwardEthereumTx(sdk.WrapSDKContext(c.ctx), msg)
	return make([]byte, 0), nil
}

func forwardTx(evm *vm.EVM, c *AbstractAccountPrecompile, caller vm.ContractRef, input []byte) ([]byte, error) {
	to := common.BytesToAddress(input[0:32])
	value := new(big.Int).SetBytes(input[32:64])
	gasLimit := new(big.Int).SetBytes(input[64:96]).Uint64()

	offsetInput := new(big.Int).SetBytes(input[96:128]).Uint64()
	calldataLength := new(big.Int).SetBytes(input[offsetInput : offsetInput+32]).Uint64()
	data := input[offsetInput+32 : offsetInput+32+calldataLength]

	// TODO send nonce back to owner contract so it can keep track of it
	nonce := uint64(0)
	gasPrice := big.NewInt(20)
	gasFeeCap := big.NewInt(20)
	gasTipCap := big.NewInt(20)
	accesses := &ethtypes.AccessList{}
	ethtx := types.NewTx(evm.ChainConfig().ChainID, nonce, &to, value, gasLimit, gasPrice, gasFeeCap, gasTipCap, data, accesses)
	msg, err := intertxtypes.NewMsgForwardEthereumTx(ethtx, caller.Address().Hex())
	if err != nil {
		return nil, err
	}

	_, err = c.intertxKeeper.ForwardEthereumTx(sdk.WrapSDKContext(c.ctx), msg)
	if err != nil {
		return nil, err
	}
	return make([]byte, 0), nil
}

// CronjobsPrecompile implemented as a native contract.
type CronjobsPrecompile struct {
	ctx            sdk.Context
	cronjobsKeeper types.CronjobsKeeper
}

func (c *CronjobsPrecompile) RequiredGas(input []byte) uint64 {
	return 3000
}

var CronAbiJSON = `[{"inputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"string","name":"identifier","type":"string"}],"name":"cancelCron","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"identifier","type":"string"},{"internalType":"string","name":"epochIdentifier","type":"string"},{"internalType":"address","name":"contractAddress","type":"address"},{"internalType":"bytes","name":"input","type":"bytes"},{"internalType":"uint256","name":"value","type":"uint256"},{"internalType":"uint256","name":"gasLimit","type":"uint256"},{"internalType":"uint256","name":"gasPrice","type":"uint256"},{"internalType":"string","name":"sender","type":"string"}],"name":"registerCron","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"}]`

func (c *CronjobsPrecompile) Run(evm *vm.EVM, caller vm.ContractRef, input []byte) ([]byte, error) {
	signature := common.Bytes2Hex(input[0:4])
	callInput := input[4:]
	var result []byte
	var err error

	abi, err := abi.JSON(strings.NewReader(CronAbiJSON))
	if err != nil {
		return nil, err
	}
	fabi, err := abi.MethodById(input[0:4])
	if err != nil {
		return nil, err
	}

	switch signature {
	case "140fff69": // registerCron(string,string,address,bytes,uint256,uint256,uint256,string)
		result, err = registerCron(evm, c, caller, callInput, fabi)
	case "a323c1fe": // cancelCron(address,string)
		result, err = cancelCron(evm, c, caller, callInput, fabi)
	default:
		return nil, errors.New("invalid CronjobsPrecompile function")
	}
	// fmt.Println("--CronjobsPrecompile result--", result, err)
	if err != nil {
		return nil, err
	}
	encodedResult := append(
		new(big.Int).SetUint64(32).FillBytes(make([]byte, 32)),
		new(big.Int).SetInt64(int64(len(result))).FillBytes(make([]byte, 32))...,
	)
	encodedResult = append(encodedResult, result...)

	padding := len(encodedResult) % 32
	if padding > 0 {
		encodedResult = append(encodedResult, make([]byte, 32-padding)...)
	}

	// fmt.Println("--CronjobsPrecompile result--", encodedResult)
	return encodedResult, nil
}

func registerCron(evm *vm.EVM, c *CronjobsPrecompile, caller vm.ContractRef, input []byte, fAbi *abi.Method) ([]byte, error) {
	// fmt.Println("--CronjobsPrecompile registerCron--")

	// var values = []interface{}
	values, err := fAbi.Inputs.UnpackValues(input)
	if err != nil {
		return nil, err
	}
	// fmt.Println("--CronjobsPrecompile registerCron-values-", values)
	identifier, _ := (values[0]).(string)
	epochIdentifier, _ := (values[1]).(string)
	contract, _ := (values[2]).(common.Address)
	data, _ := (values[3]).([]byte)
	value, _ := (values[4]).(uint64)
	gasLimit, _ := (values[5]).(uint64)
	gasPrice, _ := (values[6]).(uint64)
	sender, _ := (values[7]).(string)

	// fmt.Sprintf("%T", values[2])
	// fmt.Sprintf("%T", values[5])
	// fmt.Sprintf("%T", values[6])

	// TODO sender must be caller

	// fmt.Println("--CronjobsPrecompile registerCron-values[0]-", values[0], identifier, sender)

	_sender, err := sdk.AccAddressFromHex(sender[2:])
	// fmt.Println("--CronjobsPrecompile registerCron-_sender-", _sender, err)
	if err != nil {
		return nil, err
	}
	// fmt.Println("--CronjobsPrecompile registerCron-_sender-", _sender)
	// fmt.Println("--CronjobsPrecompile registerCron-data-", data)
	// fmt.Println("--CronjobsPrecompile registerCron-contract-", contract, contract.Hex())
	// fmt.Println("--CronjobsPrecompile registerCron-gasPrice-", gasPrice)
	// fmt.Println("--CronjobsPrecompile registerCron-gasLimit-", gasLimit)

	cronjob := cronjobstypes.Cronjob{
		Identifier:      identifier,
		EpochIdentifier: epochIdentifier,
		ContractAddress: contract.Hex(),
		Input:           common.Bytes2Hex(data),
		Value:           value,
		GasLimit:        gasLimit,
		GasPrice:        gasPrice,
		Sender:          common.HexToAddress(sender).Hex(),
	}
	// fmt.Println("--CronjobsPrecompile registerCron-cronjob-", cronjob)
	msg := cronjobstypes.NewMsgRegisterCronjob(cronjob, _sender)
	c.cronjobsKeeper.RegisterCronjob(sdk.WrapSDKContext(c.ctx), msg)

	// fmt.Println("--CronjobsPrecompile-", c.cronjobsKeeper)

	return nil, nil
}

func cancelCron(evm *vm.EVM, c *CronjobsPrecompile, caller vm.ContractRef, input []byte, cAbi *abi.Method) ([]byte, error) {
	fmt.Println("--CronjobsPrecompile cancelCron--")

	return nil, nil
}
