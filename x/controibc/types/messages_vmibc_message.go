package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSendVmibcMessage = "send_vmibc_message"

var _ sdk.Msg = &MsgSendVmibcMessage{}

func NewMsgSendVmibcMessage(
	creator string,
	port string,
	channelID string,
	timeoutTimestamp uint64,
	sourcePortId string,
	sourceChannelId string,
	sourceAddress string,
	targetPortId string,
	targetChannelId string,
	targetAddress string,
	index uint64,
	timestamp uint64,
	body string,
) *MsgSendVmibcMessage {
	return &MsgSendVmibcMessage{
		Creator:          creator,
		Port:             port,
		ChannelID:        channelID,
		TimeoutTimestamp: timeoutTimestamp,
		SourcePortId:     sourcePortId,
		SourceChannelId:  sourceChannelId,
		SourceAddress:    sourceAddress,
		TargetPortId:     targetPortId,
		TargetChannelId:  targetChannelId,
		TargetAddress:    targetAddress,
		Index:            index,
		Timestamp:        timestamp,
		Body:             body,
	}
}

func (msg *MsgSendVmibcMessage) Route() string {
	return RouterKey
}

func (msg *MsgSendVmibcMessage) Type() string {
	return TypeMsgSendVmibcMessage
}

func (msg *MsgSendVmibcMessage) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSendVmibcMessage) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSendVmibcMessage) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.Port == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid packet port")
	}
	if msg.ChannelID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid packet channel")
	}
	if msg.TimeoutTimestamp == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid packet timeout")
	}
	return nil
}
