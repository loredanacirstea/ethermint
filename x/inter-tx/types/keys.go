package types

const (
	ModuleName = "intertx"

	StoreKey = ModuleName

	RouterKey = ModuleName

	QuerierRoute = ModuleName
)

// prefix bytes for the inter-tx persistent store
const (
	prefixResponse = iota + 1
	prefixError
	prefixAbstractAccount
)

// KVStore key prefixes
var (
	KeyPrefixResponse        = []byte{prefixResponse}
	KeyPrefixError           = []byte{prefixError}
	KeyPrefixAbstractAccount = []byte{prefixAbstractAccount}
)
