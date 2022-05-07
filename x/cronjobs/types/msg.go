package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgRegisterCronjob{}
	_ sdk.Msg = &MsgCancelCronjob{}
)

const (
	TypeMsgRegisterCronjob = "register_cronjob"
	TypeMsgCancelCronjob   = "cancel_cronjob"
)

// NewMsgRegisterCronjob creates new instance of MsgRegisterCronjob
func NewMsgRegisterCronjob(
	cronjob Cronjob,
	sender sdk.AccAddress,
) *MsgRegisterCronjob {
	return &MsgRegisterCronjob{
		Cronjob: cronjob,
		Sender:  sender.String(),
	}
}

// Route returns the name of the module
func (msg MsgRegisterCronjob) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgRegisterCronjob) Type() string { return TypeMsgRegisterCronjob }

// ValidateBasic runs stateless checks on the message
func (msg MsgRegisterCronjob) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(err, "invalid deployer address %s", msg.Sender)
	}
	// TODO validation cronjob fields
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgRegisterCronjob) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgRegisterCronjob) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{from}
}

// NewMsgCancelCronjob creates new instance of MsgRegisterCronjob
func NewMsgCancelCronjob(
	id string,
	sender sdk.AccAddress,
) *MsgCancelCronjob {
	return &MsgCancelCronjob{
		Identifier: id,
		Sender:     sender.String(),
	}
}

// Route returns the name of the module
func (msg MsgCancelCronjob) Route() string { return RouterKey }

// Type returns the the action
func (msg MsgCancelCronjob) Type() string { return TypeMsgCancelCronjob }

// ValidateBasic runs stateless checks on the message
func (msg MsgCancelCronjob) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(err, "invalid deployer address %s", msg.Sender)
	}
	// TODO validation id
	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgCancelCronjob) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCancelCronjob) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{from}
}
