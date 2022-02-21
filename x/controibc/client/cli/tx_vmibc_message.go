package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	channelutils "github.com/cosmos/ibc-go/v2/modules/core/04-channel/client/utils"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/tharsis/ethermint/x/controibc/types"
)

var _ = strconv.Itoa(0)

func CmdSendVmibcMessage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-vmibc-message [src-port] [src-channel] [source-port-id] [source-channel-id] [source-address] [target-port-id] [target-channel-id] [target-address] [index] [timestamp] [body]",
		Short: "Send a vmibcMessage over IBC",
		Args:  cobra.ExactArgs(11),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			creator := clientCtx.GetFromAddress().String()
			srcPort := args[0]
			srcChannel := args[1]

			argSourcePortId := args[2]
			argSourceChannelId := args[3]
			argSourceAddress := args[4]
			argTargetPortId := args[5]
			argTargetChannelId := args[6]
			argTargetAddress := args[7]
			argIndex, err := cast.ToUint64E(args[8])
			if err != nil {
				return err
			}
			argTimestamp, err := cast.ToUint64E(args[9])
			if err != nil {
				return err
			}
			argBody := args[10]

			// Get the relative timeout timestamp
			timeoutTimestamp, err := cmd.Flags().GetUint64(flagPacketTimeoutTimestamp)
			if err != nil {
				return err
			}
			consensusState, _, _, err := channelutils.QueryLatestConsensusState(clientCtx, srcPort, srcChannel)
			if err != nil {
				return err
			}
			if timeoutTimestamp != 0 {
				timeoutTimestamp = consensusState.GetTimestamp() + timeoutTimestamp
			}

			msg := types.NewMsgSendVmibcMessage(creator, srcPort, srcChannel, timeoutTimestamp, argSourcePortId, argSourceChannelId, argSourceAddress, argTargetPortId, argTargetChannelId, argTargetAddress, argIndex, argTimestamp, argBody)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().Uint64(flagPacketTimeoutTimestamp, DefaultRelativePacketTimeoutTimestamp, "Packet timeout timestamp in nanoseconds. Default is 10 minutes.")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
