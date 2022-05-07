package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// errors
var (
	ErrInternalCronjobs = sdkerrors.Register(ModuleName, 2, "internal cronjobs error")
)
