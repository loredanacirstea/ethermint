package keeper_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibcgotesting "github.com/cosmos/ibc-go/v3/testing"
	ibctesting "github.com/tharsis/ethermint/ibc/testing"

	app "github.com/tharsis/ethermint/app"
)

var (
	// TestAccAddress defines a resuable bech32 address for testing purposes
	// TODO: update crypto.AddressHash() when sdk uses address.Module()
	TestAccAddress = icatypes.GenerateAddress(sdk.AccAddress(crypto.AddressHash([]byte(icatypes.ModuleName))), ibcgotesting.FirstConnectionID, TestPortID)
	// TestOwnerAddress defines a reusable bech32 address for testing purposes
	TestOwnerAddress = "cosmos17dtl0mjt3t77kpuhg2edqzjpszulwhgzuj9ljs"
	// TestPortID defines a resuable port identifier for testing purposes
	TestPortID, _ = icatypes.NewControllerPortID(TestOwnerAddress)
	// TestVersion defines a resuable interchainaccounts version string for testing purposes
	TestVersion = string(icatypes.ModuleCdc.MustMarshalJSON(&icatypes.Metadata{
		Version:                icatypes.Version,
		ControllerConnectionId: ibcgotesting.FirstConnectionID,
		HostConnectionId:       ibcgotesting.FirstConnectionID,
	}))
)

func init() {
	ibcgotesting.DefaultTestingAppInit = SetupICATestingApp
}

func SetupICATestingApp() (ibcgotesting.TestingApp, map[string]json.RawMessage) {
	return app.Setup(false, nil), app.NewDefaultGenesisState()
}

// KeeperTestSuite is a testing suite to test keeper functions
type KeeperTestSuite struct {
	suite.Suite

	coordinator *ibcgotesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibcgotesting.TestChain
	chainB *ibcgotesting.TestChain
}

func (suite *KeeperTestSuite) GetApp(chain *ibcgotesting.TestChain) *app.EthermintApp {
	app, ok := chain.App.(*app.EthermintApp)
	if !ok {
		panic("not ica app")
	}

	return app
}

// TestKeeperTestSuite runs all the tests within this package.
func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// SetupTest creates a coordinator with 2 test chains.
func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2, 0)
	suite.chainA = suite.coordinator.GetChain(ibcgotesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibcgotesting.GetChainID(2))
}

func NewICAPath(chainA, chainB *ibcgotesting.TestChain) *ibcgotesting.Path {
	path := ibcgotesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = icatypes.PortID
	path.EndpointB.ChannelConfig.PortID = icatypes.PortID
	path.EndpointA.ChannelConfig.Order = channeltypes.ORDERED
	path.EndpointB.ChannelConfig.Order = channeltypes.ORDERED
	path.EndpointA.ChannelConfig.Version = TestVersion
	path.EndpointB.ChannelConfig.Version = TestVersion

	return path
}
