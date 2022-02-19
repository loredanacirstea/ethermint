package keeper

import (
	"github.com/tharsis/ethermint/x/controibc/types"
)

var _ types.QueryServer = Keeper{}
