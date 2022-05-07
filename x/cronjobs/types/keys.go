package types

// constants
const (
	// module name
	ModuleName = "cronjob"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for message routing
	RouterKey = ModuleName
)

// prefix bytes for the fees persistent store
const (
	prefixCronjob = iota + 1
	prefixCronjobIdentifier
)

// KVStore key prefixes
var (
	KeyPrefixCronjob           = []byte{prefixCronjob}
	KeyPrefixCronjobIdentifier = []byte{prefixCronjobIdentifier}
)
