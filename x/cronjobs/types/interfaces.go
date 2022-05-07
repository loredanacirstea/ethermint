package types

import (
	context "context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	// evmtypes "github.com/tharsis/ethermint/x/evm/types"
	abstractacctypes "github.com/tharsis/ethermint/x/inter-tx/types"
)

// AccountKeeper defines the expected interface needed to retrieve account info.
type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// // EVMKeeper defines the expected EVM keeper interface used on erc20
// type EVMKeeper interface {
// 	EVMConfig(ctx sdk.Context) (*evmtypes.EVMConfig, error)
// 	GetParams(ctx sdk.Context) evmtypes.Params
// 	GetAccountWithoutBalance(ctx sdk.Context, addr common.Address) *statedb.Account
// }

// AbstractAccountKeeper defines the expected keeper interface
type AbstractAccountKeeper interface {
	RegisterAbstractAccount(goCtx context.Context, msg *abstractacctypes.MsgRegisterAbstractAccount) (*abstractacctypes.MsgRegisterAccountResponse, error)
	ForwardEthereumTx(goCtx context.Context, msg *abstractacctypes.MsgForwardEthereumTx) (*abstractacctypes.MsgSubmitTxResponse, error)
}
