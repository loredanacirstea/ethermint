package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v3/modules/core/exported"
	"github.com/tharsis/ethermint/x/evm/vm/types"
)

// Keeper of this module maintains collections of erc20.
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   sdk.StoreKey
	paramstore paramtypes.Subspace

	// evmKeeper   *evmkeeper.Keeper // TODO: use interface
	ics4Wrapper     transfertypes.ICS4Wrapper
	scopedIBCKeeper capabilitykeeper.ScopedKeeper
}

// NewKeeper creates new instances of the erc20 Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	ps paramtypes.Subspace,
	// evmKeeper *evmkeeper.Keeper,
	ics4Wrapper transfertypes.ICS4Wrapper,
	scopedIBCKeeper capabilitykeeper.ScopedKeeper,
) Keeper {
	// set KeyTable if it has not already been set
	// if !ps.HasKeyTable() {
	// 	ps = ps.WithKeyTable(types.ParamKeyTable())
	// }

	return Keeper{
		storeKey:   storeKey,
		cdc:        cdc,
		paramstore: ps,
		// evmKeeper:   evmKeeper,
		ics4Wrapper:     ics4Wrapper,
		scopedIBCKeeper: scopedIBCKeeper,
	}
}

func (k Keeper) Ics4Wrapper() transfertypes.ICS4Wrapper {
	return k.ics4Wrapper
}

func (k Keeper) ScopedKeeper() capabilitykeeper.ScopedKeeper {
	return k.scopedIBCKeeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// OnRecvPacket performs an IBC receive callback. It performs a no-op if
// claims are inactive
func (k Keeper) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	ack exported.Acknowledgement,
) exported.Acknowledgement {
	params := k.GetParams(ctx)

	fmt.Println("----ibcprecompile OnRecvPacket--", params)

	// return the original success acknowledgement
	return ack
}

// OnAcknowledgementPacket claims the amount from the `ActionIBCTransfer` for
// the sender of the IBC transfer.
// The function performs a no-op if claims are disabled globally,
// acknowledgment failed, or if sender the sender has no claims record.
func (k Keeper) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
) error {
	params := k.GetParams(ctx)

	fmt.Println("----ibcprecompile OnAcknowledgementPacket--", params)
	return nil
}
