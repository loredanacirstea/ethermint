package types

import (
	context "context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	cronjobstypes "github.com/tharsis/ethermint/x/cronjobs/types"
	feemarkettypes "github.com/tharsis/ethermint/x/feemarket/types"
	intertxtypes "github.com/tharsis/ethermint/x/inter-tx/types"
)

// AccountKeeper defines the expected account keeper interface
type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetAllAccounts(ctx sdk.Context) (accounts []authtypes.AccountI)
	IterateAccounts(ctx sdk.Context, cb func(account authtypes.AccountI) bool)
	GetSequence(sdk.Context, sdk.AccAddress) (uint64, error)
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, account authtypes.AccountI)
	RemoveAccount(ctx sdk.Context, account authtypes.AccountI)
	GetParams(ctx sdk.Context) (params authtypes.Params)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	// SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
}

// StakingKeeper returns the historical headers kept in store.
type StakingKeeper interface {
	GetHistoricalInfo(ctx sdk.Context, height int64) (stakingtypes.HistoricalInfo, bool)
	GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (validator stakingtypes.Validator, found bool)
}

// FeeMarketKeeper
type FeeMarketKeeper interface {
	GetBaseFee(ctx sdk.Context) *big.Int
	GetParams(ctx sdk.Context) feemarkettypes.Params
}

type InterTxKeeper interface {
	GetResponse(ctx sdk.Context, txKey []byte) []byte
	GetError(ctx sdk.Context, txKey []byte) []byte
	InterchainAccountFromAddress(goCtx context.Context, req *intertxtypes.QueryInterchainAccountFromAddressRequest) (*intertxtypes.QueryInterchainAccountFromAddressResponse, error)
	SubmitTx(goCtx context.Context, msg *intertxtypes.MsgSubmitTx) (*intertxtypes.MsgSubmitTxResponse, error)
	SubmitEthereumTx(goCtx context.Context, msg *intertxtypes.MsgSubmitEthereumTx) (*intertxtypes.MsgSubmitTxResponse, error)
	ForwardEthereumTx(goCtx context.Context, msg *intertxtypes.MsgForwardEthereumTx) (*intertxtypes.MsgSubmitTxResponse, error)
	RegisterAccount(goCtx context.Context, msg *intertxtypes.MsgRegisterAccount) (*intertxtypes.MsgRegisterAccountResponse, error)
	RegisterAbstractAccount(goCtx context.Context, msg *intertxtypes.MsgRegisterAbstractAccount) (*intertxtypes.MsgRegisterAccountResponse, error)
}

// Event Hooks
// These can be utilized to customize evm transaction processing.

// EvmHooks event hooks for evm tx processing
type EvmHooks interface {
	// Must be called after tx is processed successfully, if return an error, the whole transaction is reverted.
	PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error
}

type CronjobsKeeper interface {
	RegisterCronjob(goCtx context.Context, msg *cronjobstypes.MsgRegisterCronjob) (*cronjobstypes.MsgRegisterCronjobResponse, error)
	CancelCronjob(goCtx context.Context, msg *cronjobstypes.MsgCancelCronjob) (*cronjobstypes.MsgCancelCronjobResponse, error)
}
