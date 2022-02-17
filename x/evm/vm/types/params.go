package types

import (
	"fmt"
	"time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	DefaultClaimsDenom        = "aevmos"
	DefaultDurationUntilDecay = 2629800 * time.Second         // 1 month = 30.4375 days
	DefaultDurationOfDecay    = 2 * DefaultDurationUntilDecay // 2 months
)

// Parameter store key
var (
	ParamStoreKeyMessageContent = []byte("MessageContent")
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyMessageContent, &p.MessageContent, validateBytes),
	}
}

// NewParams creates a new Params object
func NewParams(
	messageContent []byte,
) Params {
	return Params{
		MessageContent: messageContent,
	}
}

// DefaultParams creates a parameter instance with default values
// for the claims module.
func DefaultParams() Params {
	return Params{
		MessageContent: make([]byte, 32),
	}
}

func (p Params) Validate() error {
	return nil
}

func validateBytes(i interface{}) error {
	_, ok := i.([]byte)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}
