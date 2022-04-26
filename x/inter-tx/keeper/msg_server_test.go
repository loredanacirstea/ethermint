package keeper_test

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"

	// channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	// ibctesting "github.com/cosmos/ibc-go/v3/testing"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	"github.com/tharsis/ethermint/tests"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

func (suite *KeeperTestSuite) TestSubmitEthereumTx() {
	_, owner := generateKey()
	ica, err := icatypes.NewControllerPortID(owner.String())
	suite.Require().NoError(err)
	// ica, err := sdk.AccAddressFromBech32(icaString)
	suite.Require().NoError(err)

	ctx := suite.chainA.GetContext()
	keeper := suite.GetApp(suite.chainA).InterTxKeeper
	account, err := keeper.GenerateAbstractAccount(ctx)
	suite.Require().NoError(err)
	keeper.SetAbstractAccount(ctx, ica, account)

	_account, found := keeper.GetAbstractAccount(ctx, ica)
	suite.Require().True(found)
	suite.Require().Equal(account.PrivKey, _account.PrivKey, "wrong PrivKey")

	_, _, found = keeper.GetAbstractAccountHydrated(ctx, ica)
	suite.Require().True(found)

	// set active channel
	suite.GetApp(suite.chainA).ICAControllerKeeper.SetActiveChannelID(ctx, "connection-0", ica, "channel-0")
	// suite.GetApp(suite.chainA).ScopedICAHostKeeper.ClaimCapability(ctx, )

	nonce := uint64(0)
	value := big.NewInt(0)
	to := common.HexToAddress("0x0000000000000000000000000000000000000000")
	gasLimit := uint64(300000)
	gasPrice := big.NewInt(20)
	gasFeeCap := big.NewInt(20)
	gasTipCap := big.NewInt(20)
	data := make([]byte, 0)
	accesses := &ethtypes.AccessList{}
	chainId := suite.GetApp(suite.chainA).EvmKeeper.ChainID()
	ethtx := evmtypes.NewTx(chainId, nonce, &to, value, gasLimit, gasPrice, gasFeeCap, gasTipCap, data, accesses)
	ethtx.From = ica
	res, err := keeper.SubmitEthereumTx(sdk.WrapSDKContext(ctx), ethtx, owner, "connection-0")
	suite.Require().NoError(err)
	fmt.Println("res", res)
}

func generateKey() (*ethsecp256k1.PrivKey, sdk.AccAddress) {
	address, priv := tests.NewAddrKey()
	return priv.(*ethsecp256k1.PrivKey), sdk.AccAddress(address.Bytes())
}

// func (suite *KeeperTestSuite) TestRegisterInterchainAccount() {
// 	var (
// 		owner string
// 		path  *ibctesting.Path
// 	)

// 	testCases := []struct {
// 		name     string
// 		malleate func()
// 		expPass  bool
// 	}{
// 		{
// 			"success", func() {}, true,
// 		},
// 		{
// 			"port is already bound",
// 			func() {
// 				suite.GetApp(suite.chainA).IBCKeeper.PortKeeper.BindPort(suite.chainA.GetContext(), TestPortID)
// 			},
// 			false,
// 		},
// 		{
// 			"fails to generate port-id",
// 			func() {
// 				owner = ""
// 			},
// 			false,
// 		},
// 		{
// 			"MsgChanOpenInit fails - channel is already active",
// 			func() {
// 				portID, err := icatypes.NewControllerPortID(owner)
// 				suite.Require().NoError(err)

// 				channel := channeltypes.NewChannel(
// 					channeltypes.OPEN,
// 					channeltypes.ORDERED,
// 					channeltypes.NewCounterparty(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID),
// 					[]string{path.EndpointA.ConnectionID},
// 					path.EndpointA.ChannelConfig.Version,
// 				)
// 				suite.GetApp(suite.chainA).IBCKeeper.ChannelKeeper.SetChannel(suite.chainA.GetContext(), portID, ibctesting.FirstChannelID, channel)

// 				suite.GetApp(suite.chainA).ICAControllerKeeper.SetActiveChannelID(suite.chainA.GetContext(), ibctesting.FirstConnectionID, portID, ibctesting.FirstChannelID)
// 			},
// 			false,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc

// 		suite.Run(tc.name, func() {
// 			suite.SetupTest()

// 			owner = TestOwnerAddress // must be explicitly changed

// 			path = NewICAPath(suite.chainA, suite.chainB)
// 			suite.coordinator.SetupConnections(path)

// 			tc.malleate() // malleate mutates test data

// 			msgSrv := suite.GetApp(suite.chainA).InterTxKeeper
// 			msg := types.NewMsgRegisterAccount(owner, path.EndpointA.ConnectionID)

// 			res, err := msgSrv.RegisterAccount(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

// 			if tc.expPass {
// 				suite.Require().NoError(err)
// 				suite.Require().NotNil(res)
// 			} else {
// 				suite.Require().Error(err)
// 				suite.Require().Nil(res)
// 			}

// 		})
// 	}
// }
