package keeper

import (
	"context"
	"fmt"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	"github.com/tharsis/ethermint/tests"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	"github.com/tharsis/ethermint/x/inter-tx/types"
)

var _ types.MsgServer = &Keeper{}

// RegisterAccount implements the Msg/RegisterAccount interface
func (k Keeper) RegisterAccount(goCtx context.Context, msg *types.MsgRegisterAccount) (*types.MsgRegisterAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.icaControllerKeeper.RegisterInterchainAccount(ctx, msg.ConnectionId, msg.Owner); err != nil {
		return nil, err
	}

	portID, err := icatypes.NewControllerPortID(msg.Owner)
	if err != nil {
		return nil, err
	}

	account, err := k.GenerateAbstractAccount(ctx)
	if err != nil {
		return nil, err
	}
	k.SetAbstractAccount(ctx, portID, account)

	return &types.MsgRegisterAccountResponse{}, nil
}

// RegisterAccount implements the Msg/RegisterAbstractAccount interface
func (k Keeper) RegisterAbstractAccount(goCtx context.Context, msg *types.MsgRegisterAbstractAccount) (*types.MsgRegisterAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	account, err := k.GenerateAbstractAccount(ctx)
	if err != nil {
		return nil, err
	}
	k.SetAbstractAccount(ctx, msg.Owner, account)

	return &types.MsgRegisterAccountResponse{}, nil
}

// SubmitTx implements the Msg/SubmitTx interface
func (k Keeper) SubmitTx(goCtx context.Context, msg *types.MsgSubmitTx) (*types.MsgSubmitTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	portID, err := icatypes.NewControllerPortID(msg.Owner)
	if err != nil {
		return nil, err
	}

	channelID, found := k.icaControllerKeeper.GetActiveChannelID(ctx, msg.ConnectionId, portID)
	if !found {
		return nil, sdkerrors.Wrapf(icatypes.ErrActiveChannelNotFound, "failed to retrieve active channel for port %s", portID)
	}

	chanCap, found := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if !found {
		return nil, sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	data, err := icatypes.SerializeCosmosTx(k.cdc, []sdk.Msg{msg.GetTxMsg()})
	if err != nil {
		return nil, err
	}

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}

	// timeoutTimestamp set to max value with the unsigned bit shifted to sastisfy hermes timestamp conversion
	// it is the responsibility of the auth module developer to ensure an appropriate timeout timestamp
	timeoutTimestamp := time.Now().Add(time.Minute).UnixNano()
	_, err = k.icaControllerKeeper.SendTx(ctx, chanCap, msg.ConnectionId, portID, packetData, uint64(timeoutTimestamp))
	if err != nil {
		return nil, err
	}

	return &types.MsgSubmitTxResponse{}, nil
}

// SubmitEthereumTx implements the Msg/SubmitEthereumTx interface
func (k Keeper) SubmitEthereumTx(goCtx context.Context, msg *types.MsgSubmitEthereumTx) (*types.MsgSubmitTxResponse, error) {
	fmt.Println("---SubmitEthereumTx--")
	ctx := sdk.UnwrapSDKContext(goCtx)
	owner := sdk.AccAddress(common.HexToAddress(msg.Owner).Bytes())
	connectionID := msg.ConnectionId
	msgEthereumTx := msg.GetTxMsg().(*evmtypes.MsgEthereumTx)

	portID, err := icatypes.NewControllerPortID(owner.String())
	if err != nil {
		return nil, err
	}
	ica, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, msg.ConnectionId, portID)
	if !found {
		portID = icatypes.PortID
		ica, found = k.icaControllerKeeper.GetInterchainAccountAddress(ctx, msg.ConnectionId, portID)

		return nil, sdkerrors.Wrapf(icatypes.ErrInterchainAccountNotFound, "failed to retrieve interchain account for connection %s; portID %s", msg.ConnectionId, portID)
	}

	account, found := k.GetAbstractAccount(ctx, portID)
	if !found {
		registerMsg := &types.MsgRegisterAbstractAccount{
			Owner: msg.Owner,
		}
		_, err := k.RegisterAbstractAccount(goCtx, registerMsg)
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrAbstractAccountCouldNotBeCreated, "failed to create abstract account for %s", msg.Owner)
		}
		account, found = k.GetAbstractAccount(ctx, portID)
		if !found {
			return nil, sdkerrors.Wrapf(types.ErrAbstractAccountNotExist, "failed to retrieve abstract account for interchain account %s", ica)
		}
	}
	priv := &ethsecp256k1.PrivKey{Key: account.PrivKey}
	address := common.BytesToAddress(priv.PubKey().Address().Bytes())
	msgEthereumTx.From = address.Hex()
	msgEthereumTx, err = k.signEthereumTx(priv, msgEthereumTx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to sign ethereum transaction with abstract account %s", address.Hex())
	}

	// We wrap the EthereumTx message to pass through ICA module authorization
	// MsgWrappedEthereumTx returns the IcaAddress when `msg.GetSigners()` is called
	// We avoid checking the MsgEthereumTx signature in the ICA module, but this is checked
	// in the EVM module
	wrappedMsg, err := types.NewMsgWrappedEthereumTx(msgEthereumTx, ica)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed build MsgWrappedEthereumTx")
	}

	newmsg, err := types.NewMsgSubmitTx(wrappedMsg, connectionID, owner.String())
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed build MsgSubmitTx")
	}
	return k.SubmitTx(goCtx, newmsg)
}

func (k Keeper) UnwrapEthereumTx(goCtx context.Context, msg *types.MsgWrappedEthereumTx) (*types.MsgSubmitTxResponse, error) {
	fmt.Println("---UnwrapEthereumTx--")
	msgEthereumTx := msg.GetTxMsg().(*evmtypes.MsgEthereumTx)

	// Unwrap the EthereumTx and send it to the EVM module
	_, err := k.EvmKeeper.EthereumTx(goCtx, msgEthereumTx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to forward transaction")
	}
	return &types.MsgSubmitTxResponse{}, nil
}

// ForwardEthereumTx implements the Msg/ForwardEthereumTx interface
// It forwards a transaction from a contract account to be signed with
// the contract's abstract account and sent to the EVM module
func (k Keeper) ForwardEthereumTx(goCtx context.Context, msg *types.MsgForwardEthereumTx) (*types.MsgSubmitTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	msgEthereumTx := msg.GetTxMsg().(*evmtypes.MsgEthereumTx)
	owner := sdk.AccAddress(common.HexToAddress(msg.Owner).Bytes())
	account, found := k.GetAbstractAccount(ctx, msg.Owner)
	if !found {
		registerMsg := &types.MsgRegisterAbstractAccount{
			Owner: msg.Owner,
		}
		_, err := k.RegisterAbstractAccount(goCtx, registerMsg)
		if err != nil {
			return nil, sdkerrors.Wrapf(types.ErrAbstractAccountCouldNotBeCreated, "failed to create abstract account for %s", msg.Owner)
		}
		account, found = k.GetAbstractAccount(ctx, msg.Owner)
		if !found {
			return nil, sdkerrors.Wrapf(types.ErrAbstractAccountNotExist, "failed to retrieve abstract account for interchain account %s", msg.Owner)
		}
	}
	priv := &ethsecp256k1.PrivKey{Key: account.PrivKey}
	address := common.BytesToAddress(priv.PubKey().Address().Bytes())
	msgEthereumTx.From = address.Hex()
	msgEthereumTx, err := k.signEthereumTx(priv, msgEthereumTx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to sign ethereum transaction with abstract account %s", address.Hex())
	}

	var cost *big.Int
	txData, err := evmtypes.UnpackTxData(msgEthereumTx.Data)
	if err != nil {
		cost = txData.Cost()
	} else {
		cost = big.NewInt(int64(msgEthereumTx.GetGas() * msgEthereumTx.GetFee().Uint64()))
	}

	coins := sdk.Coins{{Denom: k.EvmKeeper.GetParams(ctx).EvmDenom, Amount: sdk.NewIntFromBigInt(cost)}}
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, types.ModuleName, coins)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to send transaction fees from owner %s to module %s", address.Hex(), types.ModuleName)
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(address.Bytes()), coins)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to send transaction fees from module %s to address %s", types.ModuleName, address.Hex())
	}

	_, err = k.EvmKeeper.EthereumTx(goCtx, msgEthereumTx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to forward transaction")
	}
	return &types.MsgSubmitTxResponse{}, nil
}

func (k Keeper) signEthereumTx(priv *ethsecp256k1.PrivKey, msgEthereumTx *evmtypes.MsgEthereumTx) (*evmtypes.MsgEthereumTx, error) {
	ethSigner := ethtypes.LatestSignerForChainID(k.EvmKeeper.ChainID())
	keyringSigner := tests.NewSigner(priv)
	err := msgEthereumTx.Sign(ethSigner, keyringSigner)
	if err != nil {
		return nil, err
	}
	return msgEthereumTx, nil
}
