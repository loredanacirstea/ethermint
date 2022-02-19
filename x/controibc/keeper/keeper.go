package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/tharsis/ethermint/x/controibc/types"

	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	channelkeeper "github.com/cosmos/ibc-go/v3/modules/core/04-channel/keeper"
	portkeeper "github.com/cosmos/ibc-go/v3/modules/core/05-port/keeper"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
)

type (
	Keeper struct {
		portKey       []byte
		storeKey      sdk.StoreKey
		ChannelKeeper channelkeeper.Keeper
		PortKeeper    *portkeeper.Keeper
		ScopedKeeper  capabilitykeeper.ScopedKeeper
		cdc           codec.BinaryCodec
		// memKey        sdk.StoreKey
		paramstore paramtypes.Subspace
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	// memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	channelKeeper channelkeeper.Keeper,
	portKeeper *portkeeper.Keeper,
	scopedKeeper capabilitykeeper.ScopedKeeper,

) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		portKey:       types.PortKey,
		storeKey:      storeKey,
		ChannelKeeper: channelKeeper,
		PortKeeper:    portKeeper,
		ScopedKeeper:  scopedKeeper,
		cdc:           cdc,
		// memKey:        memKey,
		paramstore: ps,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// // ChanCloseInit defines a wrapper function for the channel Keeper's function
// func (k Keeper) ChanCloseInit(ctx sdk.Context, portID, channelID string) error {
// 	capName := host.ChannelCapabilityPath(portID, channelID)
// 	chanCap, ok := k.ScopedKeeper.GetCapability(ctx, capName)
// 	if !ok {
// 		return sdkerrors.Wrapf(channeltypes.ErrChannelCapabilityNotFound, "could not retrieve channel capability at: %s", capName)
// 	}
// 	return k.ChannelKeeper.ChanCloseInit(ctx, portID, channelID, chanCap)
// }

// IsBound checks if the module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.ScopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	cap := k.PortKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, cap, host.PortPath(portID))
}

// GetPort returns the portID for the module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(k.portKey))
}

// SetPort sets the portID for the module. Used in InitGenesis
func (k Keeper) SetPort(ctx sdk.Context, portID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(k.portKey, []byte(portID))
}

// AuthenticateCapability wraps the scopedKeeper's AuthenticateCapability function
func (k Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
	return k.ScopedKeeper.AuthenticateCapability(ctx, cap, name)
}

// ClaimCapability allows the module that can claim a capability that IBC module passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.ScopedKeeper.ClaimCapability(ctx, cap, name)
}
