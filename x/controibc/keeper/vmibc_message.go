package keeper

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	"github.com/tharsis/ethermint/x/controibc/types"
)

// TransmitVmibcMessagePacket transmits the packet over IBC with the specified source port and source channel
func (k Keeper) TransmitVmibcMessagePacket(
	ctx sdk.Context,
	packetData types.VmibcMessagePacketData,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
) error {
	fmt.Println("---TransmitVmibcMessagePacket---", packetData, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp)

	sourceChannelEnd, found := k.ChannelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", sourcePort, sourceChannel)
	}

	fmt.Println("---TransmitVmibcMessagePacket-sourceChannelEnd-found---", sourceChannelEnd, found)

	destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	fmt.Println("---TransmitVmibcMessagePacket-destinationPort-destinationChannel---", destinationPort, destinationChannel)

	// get the next sequence
	sequence, found := k.ChannelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrapf(
			channeltypes.ErrSequenceSendNotFound,
			"source port: %s, source channel: %s", sourcePort, sourceChannel,
		)
	}

	fmt.Println("---TransmitVmibcMessagePacket-sequence-found---", sequence, found)

	channelCap, ok := k.ScopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	fmt.Println("---TransmitVmibcMessagePacket-channelCap-ok---", channelCap, ok)

	packetBytes, err := packetData.GetBytes()
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, "cannot marshal the packet: "+err.Error())
	}

	fmt.Println("---TransmitVmibcMessagePacket-packetBytes---", packetBytes)

	packet := channeltypes.NewPacket(
		packetBytes,
		sequence,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
		timeoutHeight,
		timeoutTimestamp,
	)

	fmt.Println("---TransmitVmibcMessagePacket-packet---", packet)

	fmt.Println("---TransmitVmibcMessagePacket-post event---", types.EventTypeVmibcMessagePacket, sdk.EventTypeMessage)

	err = k.ics4Wrapper.SendPacket(ctx, channelCap, packet)

	// fmt.Println("----", channeltypes.AttributeValueCategory)

	// ctx.EventManager().EmitEvents(sdk.Events{
	// 	sdk.NewEvent(
	// 		types.EventTypeVmibcMessagePacket,
	// 		sdk.NewAttribute(sdk.AttributeKeySender, "sender1"),
	// 		// sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver),
	// 		sdk.NewAttribute("receiver", "receiver1"),
	// 	),
	// 	sdk.NewEvent(
	// 		sdk.EventTypeMessage,
	// 		sdk.NewAttribute(sdk.AttributeKeyModule, channeltypes.AttributeValueCategory),
	// 	),
	// })

	fmt.Println("---TransmitVmibcMessagePacket-post SendPacket---", err)

	if err != nil {
		return err
	}

	return nil
}

// OnRecvVmibcMessagePacket processes packet reception
func (k Keeper) OnRecvVmibcMessagePacket(ctx sdk.Context, packet channeltypes.Packet, data types.VmibcMessagePacketData) (packetAck types.VmibcMessagePacketAck, err error) {
	fmt.Println("----ibcPrecompile--OnRecvVmibcMessagePacket-----", packet, data)
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	// TODO: packet reception logic
	k.SetVmIbcMessage(ctx, data)

	return packetAck, nil
}

// OnAcknowledgementVmibcMessagePacket responds to the the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementVmibcMessagePacket(ctx sdk.Context, packet channeltypes.Packet, data types.VmibcMessagePacketData, ack channeltypes.Acknowledgement) error {
	fmt.Println("----ibcPrecompile--OnAcknowledgementVmibcMessagePacket-----", packet, data)
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:

		// TODO: failed acknowledgement logic
		_ = dispatchedAck.Error

		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.VmibcMessagePacketAck

		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		// TODO: successful acknowledgement logic

		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return errors.New("invalid acknowledgment format")
	}
}

// OnTimeoutVmibcMessagePacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutVmibcMessagePacket(ctx sdk.Context, packet channeltypes.Packet, data types.VmibcMessagePacketData) error {
	fmt.Println("----ibcPrecompile--OnTimeoutVmibcMessagePacket-----", packet, data)
	// TODO: packet timeout logic

	return nil
}
