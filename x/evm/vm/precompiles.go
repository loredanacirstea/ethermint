package vmibc

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	keeper "github.com/tharsis/ethermint/x/evm/vm/keeper"
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

	fmt.Println("---ibcPrecompile-0---", c.vmIbcKeeper)
	fmt.Println("---ibcPrecompile-1---", c.vmIbcKeeper.Ics4Wrapper())
	fmt.Println("---ibcPrecompile--2--", c.vmIbcKeeper.Ics4Wrapper().SendPacket)

	// channelCap := endpoint.Chain.GetChannelCapability(packet.GetSourcePort(), packet.GetSourceChannel())

	// c.vmIbcKeeper.Ics4Wrapper().SendPacket(c.ctx, channelCap, packet)

	// channelCap := vmIbcKeeper.GetCapability(ctx, channelCapName)

	// packet := channeltypes.NewPacket(ibctesting.MockPacketData, 1, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, defaultTimeoutHeight, 0)

	// scopedKeeper cosmosibckeeper.ScopedKeeper,
	// channelCap := c.vmIbcKeeper.GetCapability(c.ctx, channelCapName)

	// c.vmIbcKeeper.ScopedKeeper().NewCapability(c.ctx, "")

	channelCap, ok := c.vmIbcKeeper.ScopedKeeper().GetCapability(c.ctx, host.ChannelCapabilityPath("chat", "channel-0"))
	if !ok {
		return nil, sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	var defaultTimeoutHeight = clienttypes.NewHeight(0, 100000)

	packet := channeltypes.NewPacket([]byte("hello precompile"), 1, "chat", "channel-0", "chat", "channel-0", defaultTimeoutHeight, 1645134730571546000)

	c.vmIbcKeeper.Ics4Wrapper().SendPacket(c.ctx, channelCap, packet)

	return []byte("hello"), nil
}
