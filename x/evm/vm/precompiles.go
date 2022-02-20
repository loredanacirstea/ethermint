package extendedVM

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
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
	// channelCap := endpoint.Chain.GetChannelCapability(packet.GetSourcePort(), packet.GetSourceChannel())

	portId := "controibc"
	channelId := "channel-0"
	_, ok := c.vmIbcKeeper.ScopedKeeper.GetCapability(c.ctx, host.PortPath("controibc"))
	fmt.Println("----ibcPrecompile--GetCapability-----", portId, ok)

	channelCap, ok := c.vmIbcKeeper.ScopedKeeper.GetCapability(c.ctx, host.ChannelCapabilityPath(portId, channelId))
	fmt.Println("---ibcPrecompile-channelCap---", channelCap, ok)
	if !ok {
		return nil, sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	// 60sec
	timeoutTimestamp := uint64(time.Now().UnixNano() + 60*1000000000)
	// currentHeight := clienttypes.GetSelfHeight(c.ctx)

	// timeoutHeight := clienttypes.NewHeight(currentHeight.RevisionNumber, currentHeight.RevisionHeight+50)
	timeoutHeight := clienttypes.NewHeight(2, 100000)
	sequence, found := c.vmIbcKeeper.ChannelKeeper.GetNextSequenceSend(c.ctx, portId, channelId)
	if !found {
		sequence = 1
	}


	// packet := channeltypes.NewPacket([]byte("hello precompile"), sequence, portId, channelId, portId, channelId, timeoutHeight, timeoutTimestamp)

	// err := c.vmIbcKeeper.ChannelKeeper.SendPacket(c.ctx, channelCap, packet)
	// fmt.Println("----ibcPrecompile--err-----", err)

	// msgServer.SendVmibcMessage(c.ctx, msg)
	// ControibcPacketData.GetVmibcMessagePacket()

	packetData := types.VmibcMessagePacketData{
		Body: "hello",
	}
	err :=
		c.vmIbcKeeper.TransmitVmibcMessagePacket(c.ctx, packetData, portId, channelId, timeoutHeight, timeoutTimestamp)


	return []byte("hello"), nil
}
