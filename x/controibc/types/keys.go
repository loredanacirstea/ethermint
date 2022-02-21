package types

const (
	// ModuleName defines the module name
	ModuleName = "controibc"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_controibc"

	// Version defines the current version the IBC module supports
	Version = "controibc-1"

	// PortID is the default port id that module binds to
	PortID = "controibc"
)

// prefix bytes for the EVM persistent store
const (
	prefixVmIbcMessage = iota + 1
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey               = KeyPrefix("controibc-port-")
	KeyPrefixVmIbcMessage = []byte{prefixVmIbcMessage}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
