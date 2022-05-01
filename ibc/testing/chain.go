package ibctesting

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	ibcgotesting "github.com/cosmos/ibc-go/v3/testing"
	"github.com/cosmos/ibc-go/v3/testing/mock"

	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	ethermint "github.com/tharsis/ethermint/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"
)

// ChainIDPrefix defines the default chain ID prefix for Evmos test chains
var ChainIDPrefix = "evmos_9000-"

func init() {
	ibcgotesting.ChainIDPrefix = ChainIDPrefix
}

// NewTestChain initializes a new TestChain instance with a single validator set using a
// generated private key. It also creates a sender account to be used for delivering transactions.
//
// The first block height is committed to state in order to allow for client creations on
// counterparty chains. The TestChain will return with a block height starting at 2.
//
// Time management is handled by the Coordinator in order to ensure synchrony between chains.
// Each update of any chain increments the block header time for all chains by 5 seconds.
func NewTestChain(t *testing.T, coord *ibcgotesting.Coordinator, chainID string) *ibcgotesting.TestChain {
	// generate validator private/public key
	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})
	signersByAddress := make(map[string]tmtypes.PrivValidator, 1)
	signersByAddress[pubKey.Address().String()] = privVal

	// generate genesis account
	senderPrivKey, err := ethsecp256k1.GenerateKey()
	if err != nil {
		panic(err)
	}

	baseAcc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)

	acc := &ethermint.EthAccount{
		BaseAccount: baseAcc,
		CodeHash:    common.BytesToHash(evmtypes.EmptyCodeHash).Hex(),
	}

	amount := sdk.TokensFromConsensusPower(1, ethermint.PowerReduction)

	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, amount)),
	}

	app := SetupWithGenesisValSet(t, valSet, []authtypes.GenesisAccount{acc}, chainID, balance)

	// create current header and call begin block
	header := tmproto.Header{
		ChainID: chainID,
		Height:  1,
		Time:    coord.CurrentTime.UTC(),
	}

	txConfig := app.GetTxConfig()

	// create an account to send transactions from
	chain := &ibcgotesting.TestChain{
		T:             t,
		Coordinator:   coord,
		ChainID:       chainID,
		App:           app,
		CurrentHeader: header,
		QueryServer:   app.GetIBCKeeper(),
		TxConfig:      txConfig,
		Codec:         app.AppCodec(),
		Vals:          valSet,
		NextVals:      valSet,
		Signers:       signersByAddress,
		SenderPrivKey: senderPrivKey,
		SenderAccount: acc,
	}

	coord.CommitBlock(chain)

	return chain
}

func NewTransferPath(chainA, chainB *ibcgotesting.TestChain) *ibcgotesting.Path {
	path := ibcgotesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = ibcgotesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = ibcgotesting.TransferPort

	path.EndpointA.ChannelConfig.Order = channeltypes.UNORDERED
	path.EndpointB.ChannelConfig.Order = channeltypes.UNORDERED
	path.EndpointA.ChannelConfig.Version = "ics20-1"
	path.EndpointB.ChannelConfig.Version = "ics20-1"

	return path
}

// package ibctesting

// import (
// 	"testing"

// 	"github.com/stretchr/testify/require"

// 	"github.com/ethereum/go-ethereum/common"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
// 	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

// 	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
// 	tmtypes "github.com/tendermint/tendermint/types"

// 	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
// 	ibcgotesting "github.com/cosmos/ibc-go/v3/testing"
// 	"github.com/cosmos/ibc-go/v3/testing/mock"

// 	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
// 	ethermint "github.com/tharsis/ethermint/types"
// 	evmtypes "github.com/tharsis/ethermint/x/evm/types"
// )

// // ChainIDPrefix defines the default chain ID prefix for Evmos test chains
// var ChainIDPrefix = "evmos_9000-"
// var MaxAccounts = 1

// func init() {
// 	ibcgotesting.ChainIDPrefix = ChainIDPrefix
// }

// // NewTestChain initializes a new TestChain instance with a single validator set using a
// // generated private key. It also creates a sender account to be used for delivering transactions.
// //
// // The first block height is committed to state in order to allow for client creations on
// // counterparty chains. The TestChain will return with a block height starting at 2.
// //
// // Time management is handled by the Coordinator in order to ensure synchrony between chains.
// // Each update of any chain increments the block header time for all chains by 5 seconds.
// func NewTestChain(t *testing.T, coord *ibcgotesting.Coordinator, chainID string) *ibcgotesting.TestChain {
// 	// generate validators private/public key
// 	var (
// 		validatorsPerChain = 4
// 		validators         []*tmtypes.Validator
// 		signersByAddress   = make(map[string]tmtypes.PrivValidator, validatorsPerChain)
// 	)

// 	for i := 0; i < validatorsPerChain; i++ {
// 		privVal := mock.NewPV()
// 		pubKey, err := privVal.GetPubKey()
// 		require.NoError(t, err)
// 		validators = append(validators, tmtypes.NewValidator(pubKey, 1))
// 		signersByAddress[pubKey.Address().String()] = privVal
// 	}

// 	valSet := tmtypes.NewValidatorSet(validators)

// 	return NewTestChainWithValSet(t, coord, chainID, valSet, signersByAddress)
// }

// func NewTestChainWithValSet(t *testing.T, coord *ibcgotesting.Coordinator, chainID string, valSet *tmtypes.ValidatorSet, signers map[string]tmtypes.PrivValidator) *ibcgotesting.TestChain {
// 	genAccs := []authtypes.GenesisAccount{}
// 	genBals := []banktypes.Balance{}
// 	senderAccs := []ibcgotesting.SenderAccount{}

// 	// generate genesis accounts
// 	for i := 0; i < MaxAccounts; i++ {
// 		// generate genesis account
// 		senderPrivKey, err := ethsecp256k1.GenerateKey()
// 		if err != nil {
// 			panic(err)
// 		}

// 		baseAcc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), uint64(i), 0)

// 		acc := &ethermint.EthAccount{
// 			BaseAccount: baseAcc,
// 			CodeHash:    common.BytesToHash(evmtypes.EmptyCodeHash).Hex(),
// 		}

// 		amount, ok := sdk.NewIntFromString("10000000000000000000")
// 		require.True(t, ok)

// 		balance := banktypes.Balance{
// 			Address: acc.GetAddress().String(),
// 			Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, amount)),
// 		}

// 		genAccs = append(genAccs, acc)
// 		genBals = append(genBals, balance)

// 		senderAcc := ibcgotesting.SenderAccount{
// 			SenderAccount: acc,
// 			SenderPrivKey: senderPrivKey,
// 		}

// 		senderAccs = append(senderAccs, senderAcc)
// 	}

// 	app := SetupWithGenesisValSet(t, valSet, genAccs, chainID, genBals...)

// 	// create current header and call begin block
// 	header := tmproto.Header{
// 		ChainID: chainID,
// 		Height:  1,
// 		Time:    coord.CurrentTime.UTC(),
// 	}

// 	txConfig := app.GetTxConfig()

// 	// create an account to send transactions from
// 	chain := &ibcgotesting.TestChain{
// 		T:              t,
// 		Coordinator:    coord,
// 		ChainID:        chainID,
// 		App:            app,
// 		CurrentHeader:  header,
// 		QueryServer:    app.GetIBCKeeper(),
// 		TxConfig:       txConfig,
// 		Codec:          app.AppCodec(),
// 		Vals:           valSet,
// 		NextVals:       valSet,
// 		Signers:        signers,
// 		SenderPrivKey:  senderAccs[0].SenderPrivKey,
// 		SenderAccount:  senderAccs[0].SenderAccount,
// 		SenderAccounts: senderAccs,
// 	}

// 	coord.CommitBlock(chain)

// 	return chain
// }

// func NewTransferPath(chainA, chainB *ibcgotesting.TestChain) *ibcgotesting.Path {
// 	path := ibcgotesting.NewPath(chainA, chainB)
// 	path.EndpointA.ChannelConfig.PortID = ibcgotesting.TransferPort
// 	path.EndpointB.ChannelConfig.PortID = ibcgotesting.TransferPort

// 	path.EndpointA.ChannelConfig.Order = channeltypes.UNORDERED
// 	path.EndpointB.ChannelConfig.Order = channeltypes.UNORDERED
// 	path.EndpointA.ChannelConfig.Version = "ics20-1"
// 	path.EndpointB.ChannelConfig.Version = "ics20-1"

// 	return path
// }
