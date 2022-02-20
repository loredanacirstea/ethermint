package extendedVM

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	keeper "github.com/tharsis/ethermint/x/controibc/keeper"
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

	var defaultTimeoutHeight = clienttypes.NewHeight(0, 100000)

	packet := channeltypes.NewPacket([]byte("hello precompile"), 1, portId, channelId, portId, channelId, defaultTimeoutHeight, 1645134730571546000)

	c.vmIbcKeeper.ChannelKeeper.SendPacket(c.ctx, channelCap, packet)

	return []byte("hello"), nil
}
