package extendedVM

import (
	"errors"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/tharsis/ethermint/x/evm/types"
	intertxtypes "github.com/tharsis/ethermint/x/inter-tx/types"
)

// Ethermint additions
func GetPrecompiles(ctx sdk.Context, intertxKeeper types.InterTxKeeper) map[common.Address]vm.PrecompiledContract {
	return map[common.Address]vm.PrecompiledContract{
		common.BytesToAddress([]byte{25}): &ICAPrecompile{ctx, intertxKeeper},
		common.BytesToAddress([]byte{26}): &AbstractAccountPrecompile{ctx, intertxKeeper},
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
	case "8937a3a7": // emitTx(bytes32,bool,bytes)
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
	return encodedResult, err
}

func emitTx(evm *vm.EVM, c *ICAPrecompile, caller vm.ContractRef, input []byte) ([]byte, error) {
	fmt.Println("----emitTx---caller-", caller.Address().Hash())
	fmt.Println("----emitTx----", input)
	owner := sdk.AccAddress(caller.Address().Bytes())
	fmt.Println("----emitTx-owner---", owner, owner.String())
	// req := &intertxtypes.QueryInterchainAccountFromAddressRequest{
	// 	Owner:        owner.String(),
	// 	ConnectionId: "connection-0",
	// }
	// fmt.Println("----emitTx-req---", req)
	// icaRes, err := c.intertxKeeper.InterchainAccountFromAddressInner(c.ctx, req)
	// fmt.Println("----InterchainAccountFromAddress----", icaRes, err)
	// if err != nil {
	// 	msgRegister := &intertxtypes.MsgRegisterAccount{
	// 		Owner:        req.Owner,
	// 		ConnectionId: req.ConnectionId,
	// 	}
	// 	_, err := c.intertxKeeper.RegisterAccount(c.ctx, msgRegister)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	fmt.Println("----RegisterAccount----", err)
	// 	icaRes, err = c.intertxKeeper.InterchainAccountFromAddressInner(c.ctx, req)
	// 	fmt.Println("----InterchainAccountFromAddress2----", icaRes, err)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	// ica := icaRes.InterchainAccountAddress
	// fmt.Println("----ica----", ica)
	nonce := uint64(0)
	value := big.NewInt(0)
	to := common.BytesToAddress(caller.Address().Bytes())
	gasLimit := uint64(300000)
	gasPrice := big.NewInt(20)
	gasFeeCap := big.NewInt(20)
	gasTipCap := big.NewInt(20)
	data := make([]byte, 0)
	accesses := &ethtypes.AccessList{}
	ethtx := types.NewTx(evm.ChainConfig().ChainID, nonce, &to, value, gasLimit, gasPrice, gasFeeCap, gasTipCap, data, accesses)

	res, err := c.intertxKeeper.SubmitEthereumTx(sdk.WrapSDKContext(c.ctx), ethtx, sdk.AccAddress(caller.Address().Bytes()), "connection-0")
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
	case "8937a3a7": // registerAccount()
		result, err = registerAccount(evm, c, caller)
	case "415a2bc1": // submitTx(bytes)
		result, err = submitTx(evm, c, caller, callInput)
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

func submitTx(evm *vm.EVM, c *AbstractAccountPrecompile, caller vm.ContractRef, input []byte) ([]byte, error) {
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

	nonce := uint64(0)
	value := big.NewInt(0)
	to := common.BytesToAddress(caller.Address().Bytes())
	gasLimit := uint64(300000)
	gasPrice := big.NewInt(20)
	gasFeeCap := big.NewInt(20)
	gasTipCap := big.NewInt(20)
	data := make([]byte, 0)
	accesses := &ethtypes.AccessList{}
	ethtx := types.NewTx(evm.ChainConfig().ChainID, nonce, &to, value, gasLimit, gasPrice, gasFeeCap, gasTipCap, data, accesses)

	c.intertxKeeper.SubmitEthereumTx(sdk.WrapSDKContext(c.ctx), ethtx, sdk.AccAddress(caller.Address().Bytes()), "connection-0")
	return make([]byte, 0), nil
}
