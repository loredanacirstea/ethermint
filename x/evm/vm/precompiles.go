package extendedVM

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	keeper "github.com/tharsis/ethermint/x/controibc/keeper"
	types "github.com/tharsis/ethermint/x/controibc/types"
)

// Ethermint additions
// var PrecompiledContracts = map[common.Address]vm.PrecompiledContract{
// 	common.BytesToAddress([]byte{25}): &ibcPrecompile{},
// }

func GetPrecompiles(ctx sdk.Context, vmIbcKeeper keeper.Keeper) map[common.Address]vm.PrecompiledContract {
	return map[common.Address]vm.PrecompiledContract{
		common.BytesToAddress([]byte{25}): &ibcPrecompile{ctx, vmIbcKeeper},
	}
}

// ibcPrecompile implemented as a native contract.
type ibcPrecompile struct {
	ctx         sdk.Context
	vmIbcKeeper keeper.Keeper
}

func (c *ibcPrecompile) RequiredGas(input []byte) uint64 {
	return 3000
}

func (c *ibcPrecompile) Run(evm *vm.EVM, caller vm.ContractRef, input []byte) ([]byte, error) {
	signature := common.Bytes2Hex(input[0:4])
	callInput := input[4:]
	var result []byte
	var err error

	fmt.Println("---ibcPrecompile-Run--", signature, callInput)

	switch signature {
	case "1fe143a5": // sendMessage(string,string,address,bytes)
		result, err = sendMessage(c, caller, callInput)
	case "cf68263a": // getMessage(string,string,address,uint256)
		result, err = getMessage(c, caller, callInput)
	case "e161b98f": // countMessages(string,string,address)
		result, err = countMessages(c, caller, callInput)
	default:
		return nil, errors.New("invalid ibcPrecompile function")
	}
	if err != nil {
		fmt.Println("--ibcPrecompile err--", err)
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

	fmt.Println("--ibcPrecompile result--", encodedResult)
	return encodedResult, err
}

// function sendMessage(string memory portId, string memory channelId, address target, bytes memory data) external returns(bool success);

func sendMessage(c *ibcPrecompile, caller vm.ContractRef, input []byte) ([]byte, error) {
	portOffset := new(big.Int).SetBytes(input[0:32]).Uint64()
	channelOffset := new(big.Int).SetBytes(input[32:64]).Uint64()
	targetAddress := common.BytesToAddress(input[64:96])
	dataOffset := new(big.Int).SetBytes(input[96:128]).Uint64()

	portEnd := new(big.Int).SetBytes(input[portOffset : portOffset+32]).Uint64()
	portId := string(input[portOffset+32 : portOffset+32+portEnd])

	channelEnd := new(big.Int).SetBytes(input[channelOffset : channelOffset+32]).Uint64()
	channelId := string(input[channelOffset+32 : channelOffset+32+channelEnd])

	dataEnd := new(big.Int).SetBytes(input[dataOffset : dataOffset+32]).Uint64()
	data := input[dataOffset+32 : dataOffset+32+dataEnd]

	// 60sec
	timeoutTimestamp := uint64(time.Now().UnixNano() + 60*1000000000)
	// currentHeight := clienttypes.GetSelfHeight(c.ctx)
	// timeoutHeight := clienttypes.NewHeight(currentHeight.RevisionNumber, currentHeight.RevisionHeight+50)

	// TODO fix timeoutHeight
	timeoutHeight := clienttypes.NewHeight(2, 100000)

	packetData := types.VmibcMessagePacketData{
		// SourcePortId:    types.PortID,
		// SourceChannelId: channelId,
		// SourceAddress:   caller.Address().String(),
		// TargetPortId: portId,
		// TargetChannelId: channelId,
		// TargetAddress:   targetAddress.String(),
		// Timestamp:       uint64(time.Now().Unix()),
		Body: targetAddress.String() + string(data),
		// Body: "hello",
	}
	fmt.Println("---ibcPrecompile-packetData--", packetData)

	// channelCap := endpoint.Chain.GetChannelCapability(packet.GetSourcePort(), packet.GetSourceChannel())

	err := c.vmIbcKeeper.TransmitVmibcMessagePacket(c.ctx, packetData, portId, channelId, timeoutHeight, timeoutTimestamp)
	if err != nil {
		return make([]byte, 32), err
	}
	return new(big.Int).SetUint64(uint64(1)).FillBytes(make([]byte, 32)), nil
}

func getMessage(c *ibcPrecompile, caller vm.ContractRef, input []byte) ([]byte, error) {
	msg, success := c.vmIbcKeeper.GetVmIbcMessage(c.ctx, caller.Address().Bytes())
	if !success {
		return nil, errors.New("message not found")
	}
	return []byte(msg.Body), nil
}

func countMessages(c *ibcPrecompile, caller vm.ContractRef, input []byte) ([]byte, error) {
	return nil, nil
}
