package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrIBCAccountAlreadyExist      = sdkerrors.Register(ModuleName, 2, "interchain account already registered")
	ErrIBCAccountNotExist          = sdkerrors.Register(ModuleName, 3, "interchain account not exist")
	ErrAbstractAccountAlreadyExist = sdkerrors.Register(ModuleName, 4, "abstract account already registered")
	ErrAbstractAccountNotExist     = sdkerrors.Register(ModuleName, 5, "abstract account not exist")
	ErrInternalInterTx             = sdkerrors.Register(ModuleName, 6, "internal intertx error")
)
