package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)


var (
	KeyPortId = []byte("PortId")
	// TODO: Determine the default value
	DefaultPortId string = "port_id"
)


// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	portId string,
) Params {
	return Params{
        PortId: portId,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
        DefaultPortId,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyPortId, &p.PortId, validatePortId),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
   	if err := validatePortId(p.PortId); err != nil {
   		return err
   	}
   	
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// validatePortId validates the PortId param
func validatePortId(v interface{}) error {
	portId, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// TODO implement validation
	_ = portId

	return nil
}
