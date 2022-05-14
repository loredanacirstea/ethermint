package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	"github.com/tharsis/ethermint/x/inter-tx/types"
)

// Hooks wrapper struct for fees keeper
type Hooks struct {
	k Keeper
}

var _ evmtypes.EvmHooks = Hooks{}

// Hooks return the wrapper hooks struct for the Keeper
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// PostTxProcessing implements EvmHooks.PostTxProcessing. After each successful
// interaction with a registered contract, the contract deployer receives
// a share from the transaction fees paid by the user.
func (h Hooks) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	// check if the fees are globally enabled
	// params := h.k.GetParams(ctx)
	// if !params.EnableInterTx {
	// 	return nil
	// }

	// fmt.Println("-PostTxProcessing-receipt-", receipt)
	// fmt.Println("-PostTxProcessing-logs-", receipt.Logs)

	sendTxTopic := common.HexToHash("0xc87c5b688cdecc4e40f06a394de49a3215bd4a1a1256113bbe0426f9b80d563e")
	// fmt.Println("-PostTxProcessing-sendTxTopic-", sendTxTopic)

	var to common.Address
	var data []byte

	for _, log := range receipt.Logs {
		// fmt.Println("-PostTxProcessing-Topics-", len(log.Topics), log.Topics)
		if len(log.Topics) > 1 && log.Topics[0] == sendTxTopic {
			fmt.Println("---SendTransaction event topic detected: ", log.Topics[1])
			to = common.BytesToAddress(log.Topics[1].Bytes())
			data = log.Data

			fmt.Println("---to", log.Topics[1].Bytes(), to)
			// fmt.Println("---data", common.Bytes2Hex(data))

			value := new(big.Int).SetBytes(data[0:32])
			gasLimit := new(big.Int).SetBytes(data[32:64]).Uint64()
			dataLen := new(big.Int).SetBytes(data[96:128]).Uint64()
			data = data[128 : dataLen+128]
			fmt.Println("---data", common.Bytes2Hex(data))

			chainid := h.k.EvmKeeper.ChainID()
			nonce := uint64(0)
			// value := big.NewInt(0)
			// gasLimit := uint64(200000)
			gasPrice := big.NewInt(20)
			gasFeeCap := big.NewInt(20)
			gasTipCap := big.NewInt(20)
			accesses := &ethtypes.AccessList{}
			ethtx := evmtypes.NewTx(chainid, nonce, &to, value, gasLimit, gasPrice, gasFeeCap, gasTipCap, data, accesses)

			sender := msg.From().Hex()

			fmsg, err := types.NewMsgForwardEthereumTx(ethtx, sender)
			if err != nil {
				return err
			}

			res, err := h.k.ForwardEthereumTx(sdk.WrapSDKContext(ctx), fmsg)
			fmt.Println("-PostTxProcessing-res-", res, err)
		}
	}

	return nil
}
