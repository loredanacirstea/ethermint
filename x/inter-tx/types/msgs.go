package types

import (
	fmt "fmt"
	"strings"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	proto "github.com/gogo/protobuf/proto"
)

var (
	_ sdk.Msg = &MsgRegisterAccount{}
	_ sdk.Msg = &MsgSubmitTx{}
	_ sdk.Msg = &MsgSubmitEthereumTx{}
	_ sdk.Msg = &MsgForwardEthereumTx{}
	_ sdk.Msg = &MsgRegisterAbstractAccount{}

	_ codectypes.UnpackInterfacesMessage = MsgSubmitTx{}
)

// NewMsgRegisterAccount creates a new MsgRegisterAccount instance
func NewMsgRegisterAccount(owner, connectionID string) *MsgRegisterAccount {
	return &MsgRegisterAccount{
		Owner:        owner,
		ConnectionId: connectionID,
	}
}

// ValidateBasic implements sdk.Msg
func (msg MsgRegisterAccount) ValidateBasic() error {
	if strings.TrimSpace(msg.Owner) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "failed to parse address: %s", msg.Owner)
	}

	return nil
}

// GetSigners implements sdk.Msg
func (msg MsgRegisterAccount) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{accAddr}
}

// NewMsgRegisterAccount creates a new RegisterAbstractAccount instance
func NewRegisterAbstractAccount(owner string) *MsgRegisterAbstractAccount {
	return &MsgRegisterAbstractAccount{
		Owner: owner,
	}
}

// ValidateBasic implements sdk.Msg
func (msg MsgRegisterAbstractAccount) ValidateBasic() error {
	if strings.TrimSpace(msg.Owner) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "failed to parse address: %s", msg.Owner)
	}

	return nil
}

// GetSigners implements sdk.Msg
func (msg MsgRegisterAbstractAccount) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{accAddr}
}

// NewMsgSubmitTx creates and returns a new MsgSubmitTx instance
func NewMsgSubmitTx(sdkMsg sdk.Msg, connectionID, owner string) (*MsgSubmitTx, error) {
	any, err := PackTxMsgAny(sdkMsg)
	if err != nil {
		return nil, err
	}

	return &MsgSubmitTx{
		ConnectionId: connectionID,
		Owner:        owner,
		Msg:          any,
	}, nil
}

// PackTxMsgAny marshals the sdk.Msg payload to a protobuf Any type
func PackTxMsgAny(sdkMsg sdk.Msg) (*codectypes.Any, error) {
	msg, ok := sdkMsg.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("can't proto marshal %T", sdkMsg)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, err
	}

	return any, nil
}

// UnpackInterfaces implements codectypes.UnpackInterfacesMessage
func (msg MsgSubmitTx) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var (
		sdkMsg sdk.Msg
	)

	return unpacker.UnpackAny(msg.Msg, &sdkMsg)
}

// GetTxMsg fetches the cached any message
func (msg *MsgSubmitTx) GetTxMsg() sdk.Msg {
	sdkMsg, ok := msg.Msg.GetCachedValue().(sdk.Msg)
	if !ok {
		return nil
	}

	return sdkMsg
}

// GetSigners implements sdk.Msg
func (msg MsgSubmitTx) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{accAddr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgSubmitTx) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid owner address")
	}

	return nil
}

// NewMsgSubmitTx creates and returns a new MsgSubmitTx instance
func NewMsgSubmitEthereumTx(sdkMsg sdk.Msg, connectionID, owner string) (*MsgSubmitEthereumTx, error) {
	any, err := PackTxMsgAny(sdkMsg)
	if err != nil {
		return nil, err
	}

	return &MsgSubmitEthereumTx{
		ConnectionId: connectionID,
		Owner:        owner,
		Msg:          any,
	}, nil
}

// UnpackInterfaces implements codectypes.UnpackInterfacesMessage
func (msg MsgSubmitEthereumTx) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var (
		sdkMsg sdk.Msg
	)

	return unpacker.UnpackAny(msg.Msg, &sdkMsg)
}

// GetTxMsg fetches the cached any message
func (msg *MsgSubmitEthereumTx) GetTxMsg() sdk.Msg {
	sdkMsg, ok := msg.Msg.GetCachedValue().(sdk.Msg)
	if !ok {
		return nil
	}

	return sdkMsg
}

// GetSigners implements sdk.Msg
func (msg MsgSubmitEthereumTx) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{accAddr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgSubmitEthereumTx) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid owner address")
	}

	return nil
}

// NewMsgForwardEthereumTx creates and returns a new MsgSubmitTx instance
func NewMsgForwardEthereumTx(sdkMsg sdk.Msg, owner string) (*MsgForwardEthereumTx, error) {
	any, err := PackTxMsgAny(sdkMsg)
	if err != nil {
		return nil, err
	}

	return &MsgForwardEthereumTx{
		Owner: owner,
		Msg:   any,
	}, nil
}

// UnpackInterfaces implements codectypes.UnpackInterfacesMessage
func (msg MsgForwardEthereumTx) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var (
		sdkMsg sdk.Msg
	)

	return unpacker.UnpackAny(msg.Msg, &sdkMsg)
}

// GetTxMsg fetches the cached any message
func (msg *MsgForwardEthereumTx) GetTxMsg() sdk.Msg {
	sdkMsg, ok := msg.Msg.GetCachedValue().(sdk.Msg)
	if !ok {
		return nil
	}

	return sdkMsg
}

// GetSigners implements sdk.Msg
func (msg MsgForwardEthereumTx) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{accAddr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgForwardEthereumTx) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid owner address")
	}

	return nil
}

// NewMsgForwardEthereumTx creates and returns a new MsgSubmitTx instance
func NewMsgWrappedEthereumTx(sdkMsg sdk.Msg, icaAddress string) (*MsgWrappedEthereumTx, error) {
	any, err := PackTxMsgAny(sdkMsg)
	if err != nil {
		return nil, err
	}

	return &MsgWrappedEthereumTx{
		IcaAddress: icaAddress,
		Msg:        any,
	}, nil
}

// UnpackInterfaces implements codectypes.UnpackInterfacesMessage
func (msg MsgWrappedEthereumTx) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var (
		sdkMsg sdk.Msg
	)

	return unpacker.UnpackAny(msg.Msg, &sdkMsg)
}

// GetTxMsg fetches the cached any message
func (msg *MsgWrappedEthereumTx) GetTxMsg() sdk.Msg {
	sdkMsg, ok := msg.Msg.GetCachedValue().(sdk.Msg)
	if !ok {
		return nil
	}

	return sdkMsg
}

// GetSigners implements sdk.Msg
func (msg MsgWrappedEthereumTx) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.IcaAddress)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{accAddr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgWrappedEthereumTx) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.IcaAddress)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid owner address")
	}

	return nil
}
