package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/tharsis/ethermint/x/cronjobs/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
	intertxtypes "github.com/tharsis/ethermint/x/inter-tx/types"
)

// GetAllCronjobs - get all registered DevFeeInfo instances
func (k Keeper) ExecuteCron(ctx sdk.Context, cron types.Cronjob) (bool, error) {

	nonce := uint64(0)
	value := big.NewInt(int64(cron.Value))
	gasPrice := big.NewInt(20)
	gasFeeCap := big.NewInt(20)
	gasTipCap := big.NewInt(20)
	gasLimit := uint64(200000)
	data := common.Hex2Bytes(cron.Input)
	to := common.HexToAddress(cron.ContractAddress)
	chainid := k.EvmKeeper.ChainID()

	ethtx := evmtypes.NewTx(chainid, nonce, &to, value, gasLimit, gasPrice, gasFeeCap, gasTipCap, data, nil)

	msg, err := intertxtypes.NewMsgForwardEthereumTx(ethtx, cron.Sender)
	if err != nil {
		return false, err
	}

	_, err = k.AbstractAccountKeeper.RegisterAbstractAccount(sdk.WrapSDKContext(ctx), &intertxtypes.MsgRegisterAbstractAccount{
		Owner: cron.Sender,
	})

	_, err = k.AbstractAccountKeeper.ForwardEthereumTx(sdk.WrapSDKContext(ctx), msg)

	if err != nil {
		return false, err
	}
	return true, nil
}
