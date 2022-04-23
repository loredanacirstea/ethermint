package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type MsgEthereumIcaTx struct {
	*MsgEthereumTx
}

// message type and route constants
const (
	// TypeMsgEthereumIcaTx defines the type string of an Ethereum transaction
	TypeMsgEthereumIcaTx = "ethereum_ica_tx"
)

// NewIcaTx returns a reference to a new Ethereum transaction message.
func NewIcaTx(
	chainID *big.Int, nonce uint64, to *common.Address, amount *big.Int,
	gasLimit uint64, gasPrice, gasFeeCap, gasTipCap *big.Int, input []byte, accesses *ethtypes.AccessList,
) *MsgEthereumIcaTx {
	tx := NewTx(chainID, nonce, to, amount, gasLimit, gasPrice, gasFeeCap, gasTipCap, input, accesses)
	return &MsgEthereumIcaTx{MsgEthereumTx: tx}
}

// NewIcaTxContract returns a reference to a new Ethereum transaction
// message designated for contract creation.
func NewIcaTxContract(
	chainID *big.Int,
	nonce uint64,
	amount *big.Int,
	gasLimit uint64,
	gasPrice, gasFeeCap, gasTipCap *big.Int,
	input []byte,
	accesses *ethtypes.AccessList,
) *MsgEthereumIcaTx {
	tx := NewTxContract(chainID, nonce, amount, gasLimit, gasPrice, gasFeeCap, gasTipCap, input, accesses)
	return &MsgEthereumIcaTx{MsgEthereumTx: tx}
}

func newMsgEthereumIcaTx(
	chainID *big.Int, nonce uint64, to *common.Address, amount *big.Int,
	gasLimit uint64, gasPrice, gasFeeCap, gasTipCap *big.Int, input []byte, accesses *ethtypes.AccessList,
) *MsgEthereumIcaTx {
	tx := newMsgEthereumTx(chainID, nonce, to, amount, gasLimit, gasPrice, gasFeeCap, gasTipCap, input, accesses)
	return &MsgEthereumIcaTx{MsgEthereumTx: tx}
}

// GetSigners returns the expected signers for an Ethereum transaction message.
// For such a message, there should exist only a single 'signer'.
//
// NOTE: This method panics if 'Sign' hasn't been called first.
func (msg *MsgEthereumIcaTx) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{common.Hex2Bytes(msg.From)}
}
