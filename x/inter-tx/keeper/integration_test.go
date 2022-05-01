package keeper_test

// "fmt"
// "math/big"

// sdk "github.com/cosmos/cosmos-sdk/types"
// icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"

// channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
// ibctesting "github.com/cosmos/ibc-go/v3/testing"
// "github.com/ethereum/go-ethereum/common"
// ethtypes "github.com/ethereum/go-ethereum/core/types"
// "github.com/tharsis/ethermint/crypto/ethsecp256k1"
// "github.com/tharsis/ethermint/tests"
// evmtypes "github.com/tharsis/ethermint/x/evm/types"
// "github.com/tharsis/ethermint/x/inter-tx/types"

// OnChanCloseInit on controller (chainA)
func (suite *KeeperTestSuite) TestOnChanCloseInit() {
	path := NewICAPath(suite.chainA, suite.chainB)
	suite.coordinator.SetupConnections(path)

	err := SetupICAPath(path, TestOwnerAddress)
	suite.Require().NoError(err)

	module, _, err := suite.chainA.App.GetIBCKeeper().PortKeeper.LookupModuleByPort(suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID)
	suite.Require().NoError(err)

	cbs, ok := suite.chainA.App.GetIBCKeeper().Router.GetRoute(module)
	suite.Require().True(ok)

	err = cbs.OnChanCloseInit(
		suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID,
	)

	suite.Require().Error(err)
}

// func (suite *KeeperTestSuite) TestSubmitEthereumTx() {
// 	_, owner := generateKey()
// 	ica, err := icatypes.NewControllerPortID(owner.String())
// 	suite.Require().NoError(err)
// 	// ica, err := sdk.AccAddressFromBech32(icaString)
// 	suite.Require().NoError(err)

// 	ctx := suite.chainA.GetContext()
// 	keeper := suite.GetApp(suite.chainA).InterTxKeeper
// 	account, err := keeper.GenerateAbstractAccount(ctx)
// 	suite.Require().NoError(err)
// 	keeper.SetAbstractAccount(ctx, ica, account)

// 	_account, found := keeper.GetAbstractAccount(ctx, ica)
// 	suite.Require().True(found)
// 	suite.Require().Equal(account.PrivKey, _account.PrivKey, "wrong PrivKey")

// 	_, _, found = keeper.GetAbstractAccountHydrated(ctx, ica)
// 	suite.Require().True(found)

// 	// set active channel
// 	suite.GetApp(suite.chainA).ICAControllerKeeper.SetActiveChannelID(ctx, "connection-0", ica, "channel-0")
// 	// suite.GetApp(suite.chainA).ScopedICAHostKeeper.ClaimCapability(ctx, )

// 	nonce := uint64(0)
// 	value := big.NewInt(0)
// 	to := common.HexToAddress("0x0000000000000000000000000000000000000000")
// 	gasLimit := uint64(300000)
// 	gasPrice := big.NewInt(20)
// 	gasFeeCap := big.NewInt(20)
// 	gasTipCap := big.NewInt(20)
// 	data := make([]byte, 0)
// 	accesses := &ethtypes.AccessList{}
// 	chainId := suite.GetApp(suite.chainA).EvmKeeper.ChainID()
// 	ethtx := evmtypes.NewTx(chainId, nonce, &to, value, gasLimit, gasPrice, gasFeeCap, gasTipCap, data, accesses)
// 	ethtx.From = ica
// 	msg, err := types.NewMsgSubmitEthereumTx(ethtx, "connection-0", common.BytesToAddress(owner.Bytes()).Hex())
// 	suite.Require().NoError(err)
// 	res, err := keeper.SubmitEthereumTx(sdk.WrapSDKContext(ctx), msg)
// 	suite.Require().NoError(err)
// 	fmt.Println("res", res)
// }

// func generateKey() (*ethsecp256k1.PrivKey, sdk.AccAddress) {
// 	address, priv := tests.NewAddrKey()
// 	return priv.(*ethsecp256k1.PrivKey), sdk.AccAddress(address.Bytes())
// }
