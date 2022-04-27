package keeper

import (
	"context"
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

// SubmitTx implements the Msg/SubmitTx interface
func (k Keeper) SubmitEthereumTx(goCtx context.Context, msgEthereumTx *evmtypes.MsgEthereumTx, owner sdk.AccAddress, connectionID string) (*types.MsgSubmitTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	ica := msgEthereumTx.From
	account, found := k.GetAbstractAccount(ctx, ica)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrAbstractAccountNotExist, "failed to retrieve abstract account for interchain account %s", ica)
	}
	priv := &ethsecp256k1.PrivKey{Key: account.PrivKey}
	address := common.BytesToAddress(priv.PubKey().Address().Bytes())
	msgEthereumTx.From = address.Hex()
	msgEthereumTx, err := k.signEthereumTx(priv, msgEthereumTx)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed to sign ethereum transaction with abstract account %s", address.Hex())
	}

	msg, err := types.NewMsgSubmitTx(msgEthereumTx, connectionID, owner.String())
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "failed build MsgSubmitTx")
	}
	return k.SubmitTx(goCtx, msg)
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
