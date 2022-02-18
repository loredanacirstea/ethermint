package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	coretypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
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
	portKeeper      transfertypes.PortKeeper
}

// NewKeeper creates new instances of the erc20 Keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	ps paramtypes.Subspace,
	// evmKeeper *evmkeeper.Keeper,
	ics4Wrapper transfertypes.ICS4Wrapper,
	scopedIBCKeeper capabilitykeeper.ScopedKeeper,
	portKeeper transfertypes.PortKeeper,
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
		portKeeper:      portKeeper,
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

// IsBound checks a given port ID is already bounded.
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedIBCKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// // BindPort binds to a port and returns the associated capability.
// // Ports must be bound statically when the chain starts in `app.go`.
// // The capability must then be passed to a module which will need to pass
// // it as an extra parameter when calling functions on the IBC module.
// func (k *Keeper) BindPort(ctx sdk.Context, portID string) *capabilitytypes.Capability {
// 	if err := host.PortIdentifierValidator(portID); err != nil {
// 		panic(err.Error())
// 	}

// 	if k.IsBound(ctx, portID) {
// 		panic(fmt.Sprintf("port %s is already bound", portID))
// 	}

// 	key, err := k.scopedIBCKeeper.NewCapability(ctx, host.PortPath(portID))
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	k.Logger(ctx).Info("port binded", "port", portID)
// 	return key
// }

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	cap := k.portKeeper.BindPort(ctx, portID)
	return k.scopedIBCKeeper.ClaimCapability(ctx, cap, host.PortPath(portID))
}

// GetPort returns the portID for the transfer module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.PortKey))
}

// SetPort sets the portID for the transfer module. Used in InitGenesis
func (k Keeper) SetPort(ctx sdk.Context, portID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PortKey, []byte(portID))
}

// Authenticate authenticates a capability key against a port ID
// by checking if the memory address of the capability was previously
// generated and bound to the port (provided as a parameter) which the capability
// is being authenticated against.
func (k Keeper) Authenticate(ctx sdk.Context, key *capabilitytypes.Capability, portID string) bool {
	if err := host.PortIdentifierValidator(portID); err != nil {
		panic(err.Error())
	}

	return k.scopedIBCKeeper.AuthenticateCapability(ctx, key, host.PortPath(portID))
}

// LookupModuleByPort will return the IBCModule along with the capability associated with a given portID
func (k Keeper) LookupModuleByPort(ctx sdk.Context, portID string) (string, *capabilitytypes.Capability, error) {
	modules, cap, err := k.scopedIBCKeeper.LookupModules(ctx, host.PortPath(portID))
	if err != nil {
		return "", nil, err
	}

	return coretypes.GetModuleOwner(modules), cap, nil
}
