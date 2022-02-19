package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	"github.com/tharsis/ethermint/x/controibc/types"
)

func (k msgServer) SendVmibcMessage(goCtx context.Context, msg *types.MsgSendVmibcMessage) (*types.MsgSendVmibcMessageResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: logic before transmitting the packet

	// Construct the packet
	var packet types.VmibcMessagePacketData

	packet.Body = msg.Body

	// Transmit the packet
	err := k.TransmitVmibcMessagePacket(
		ctx,
		packet,
		msg.Port,
		msg.ChannelID,
		clienttypes.ZeroHeight(),
		msg.TimeoutTimestamp,
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendVmibcMessageResponse{}, nil
}
