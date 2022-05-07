package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/tharsis/ethermint/x/cronjobs/types"
)

var _ types.MsgServer = &Keeper{}

// RegisterDevFeeInfo registers a contract to receive transaction fees
func (k Keeper) RegisterCronjob(
	goCtx context.Context,
	msg *types.MsgRegisterCronjob,
) (*types.MsgRegisterCronjobResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// if !k.isEnabled(ctx) {
	// 	return nil, sdkerrors.Wrapf(types.ErrInternalCronjobs, "cronjobs module is not enabled")
	// }

	id := k.GetId(msg.Cronjob.Sender, msg.Cronjob.Identifier)

	if k.IsCronjobRegistered(ctx, id) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "cronjob is already registered %s", id)
	}

	k.SetCronjob(ctx, id, msg.Cronjob)
	k.Logger(ctx).Debug(
		"registering cronjob for transaction fees",
		"id", id, "sender", msg.Cronjob.Sender,
	)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeRegisterCronjob,
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Cronjob.Sender),
				sdk.NewAttribute(types.AttributeKeyId, id),
				sdk.NewAttribute(types.AttributeKeyContract, msg.Cronjob.ContractAddress),
			),
		},
	)

	return &types.MsgRegisterCronjobResponse{}, nil
}

// DeleteCronjob deletes the fee for a contract
func (k Keeper) CancelCronjob(
	goCtx context.Context,
	msg *types.MsgCancelCronjob,
) (*types.MsgCancelCronjobResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.isEnabled(ctx) {
		return nil, sdkerrors.Wrapf(types.ErrInternalCronjobs, "cronjobs module is not enabled")
	}

	id := k.GetId(msg.Sender, msg.Identifier)

	if !k.IsCronjobRegistered(ctx, id) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "cronjob is not registered %s", id)
	}

	k.DeleteCronjob(ctx, id)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeDeleteCronjob,
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
				sdk.NewAttribute(types.AttributeKeyId, id),
			),
		},
	)

	return &types.MsgCancelCronjobResponse{}, nil
}

func (k Keeper) GetId(sender string, identifier string) string {
	return sender + "_" + identifier
}
