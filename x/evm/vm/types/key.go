package types

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
)

// constants
const (
	// module name
	ModuleName = "vmibc"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_vmibc"

	// Version defines the current version the IBC module supports
	Version = "vmibc-1"

	// PortID is the default port id that module binds to
	PortID = ModuleName
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("vmibc-port-")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// ModuleAddress is the native module address for EVM
var ModuleAddress common.Address

func init() {
	ModuleAddress = common.BytesToAddress(authtypes.NewModuleAddress(ModuleName).Bytes())
}

// prefix bytes for the EVM persistent store
const (
	prefixMessageRecord = iota + 1
)

// KVStore key prefixes
var (
	KeyPrefixMessageRecord = []byte{prefixMessageRecord}
)
