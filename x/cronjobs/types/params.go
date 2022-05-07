package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store key
var (
	DefaultEnableCronjobs       = true
	ParamStoreKeyEnableCronjobs = []byte("EnableCronjobs")
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(
	enableCronjobs bool,
) Params {
	return Params{
		EnableCronjobs: enableCronjobs,
	}
}

func DefaultParams() Params {
	return Params{
		EnableCronjobs: DefaultEnableCronjobs,
	}
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyEnableCronjobs, &p.EnableCronjobs, validateBool),
	}
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func (p Params) Validate() error {
	if err := validateBool(p.EnableCronjobs); err != nil {
		return err
	}
	return nil
}
