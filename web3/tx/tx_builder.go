package tx

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type EthTxClient interface {
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	ChainID(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
}

type LegacyEthTxClient interface {
	EthTxClient
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
}

type DynamicFeeEthTxClient interface {
	LegacyEthTxClient
	HeaderByNumber(context.Context, *big.Int) (*types.Header, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
}

type TxBuilder struct {
	to      *common.Address
	data    []byte
	value   *big.Int
	nonce   *uint64
	chainID *big.Int

	from     *common.Address
	gasPrice *big.Int

	gasTipCap *big.Int
	gasFeeCap *big.Int

	gasMultiplier uint64
	gasDivisor    uint64
}

func NewBuilder(to *common.Address) *TxBuilder {
	return &TxBuilder{
		to:            to,
		value:         big.NewInt(0),
		gasMultiplier: 105,
		gasDivisor:    100,
	}
}

func (b *TxBuilder) Value(value *big.Int) *TxBuilder {
	b.value = value
	return b
}

func (b *TxBuilder) Data(data []byte) *TxBuilder {
	b.data = data
	return b
}

func (b *TxBuilder) Nonce(nonce uint64) *TxBuilder {
	b.nonce = &nonce
	return b
}

func (b *TxBuilder) ChainID(chainID *big.Int) *TxBuilder {
	b.chainID = chainID
	return b
}

func (b *TxBuilder) GasPrice(gasPrice *big.Int) *TxBuilder {
	b.gasPrice = gasPrice
	return b
}

func (b *TxBuilder) GasTipCap(gasTipCap *big.Int) *TxBuilder {
	b.gasTipCap = gasTipCap
	return b
}

func (b *TxBuilder) GasFeeCap(gasFeeCap *big.Int) *TxBuilder {
	b.gasFeeCap = gasFeeCap
	return b
}

func (b *TxBuilder) GasLimitMultiplier(multiplier, divisor uint64) *TxBuilder {
	b.gasMultiplier = multiplier
	b.gasDivisor = divisor
	return b
}

func adjustGasLimit(gas uint64, multiplier, divisor uint64) uint64 {
	return gas * multiplier / divisor
}

func (b *TxBuilder) EstimateGasFrom(ctx context.Context, client EthTxClient, from common.Address) (uint64, error) {
	return client.EstimateGas(ctx, ethereum.CallMsg{
		From:  from,
		To:    b.to,
		Data:  b.data,
		Value: b.value,
	})
}

func (b *TxBuilder) SignLegacy(ctx context.Context, client LegacyEthTxClient, prv *ecdsa.PrivateKey) (*types.Transaction, error) {

	fromAddress := crypto.PubkeyToAddress(prv.PublicKey)
	nonce, err := b.getNonce(ctx, client, fromAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, err := b.getGasPrice(ctx, client)
	if err != nil {
		return nil, err
	}

	gas, err := b.EstimateGasFrom(ctx, client, fromAddress)
	if err != nil {
		return nil, err
	}

	// Increase gas by multiplier/divisor
	gas = adjustGasLimit(gas, b.gasMultiplier, b.gasDivisor)

	return types.SignNewTx(prv, types.HomesteadSigner{}, &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gas,
		To:       b.to,
		Value:    b.value,
		Data:     b.data,
	})
}

func (b *TxBuilder) SignDynamicFee(ctx context.Context, client DynamicFeeEthTxClient, prv *ecdsa.PrivateKey) (*types.Transaction, error) {
	chainID, err := b.getChainID(ctx, client)
	if err != nil {
		return nil, err
	}

	fromAddress := crypto.PubkeyToAddress(prv.PublicKey)
	nonce, err := b.getNonce(ctx, client, fromAddress)
	if err != nil {
		return nil, err
	}

	gas, err := b.EstimateGasFrom(ctx, client, fromAddress)
	if err != nil {
		return nil, err
	}
	// Increase gas by multiplier/divisor
	gas = adjustGasLimit(gas, b.gasMultiplier, b.gasDivisor)

	tipCap, err := b.getGasTipCap(ctx, client)
	if err != nil {
		return nil, err
	}

	feeCap, err := b.getGasFeeCap(ctx, client)
	if err != nil {
		return nil, err
	}

	if feeCap.Cmp(tipCap) < 0 {
		return nil, fmt.Errorf("gasFeeCap must be greater than or equal to gasTipCap")
	}

	return types.SignNewTx(prv, types.NewCancunSigner(chainID), &types.DynamicFeeTx{
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		Gas:       gas,
		To:        b.to,
		Value:     b.value,
		Data:      b.data,
	})
}

// ----

func (b *TxBuilder) getGasTipCap(ctx context.Context, client DynamicFeeEthTxClient) (*big.Int, error) {
	if b.gasTipCap != nil {
		return b.gasTipCap, nil
	}

	if tipCap, err := client.SuggestGasTipCap(ctx); err != nil {
		return nil, err
	} else {
		b.gasTipCap = tipCap
	}

	return b.gasTipCap, nil
}

func (b *TxBuilder) getGasFeeCap(ctx context.Context, client DynamicFeeEthTxClient) (*big.Int, error) {
	if b.gasFeeCap != nil {
		return b.gasFeeCap, nil
	}

	header, err := client.HeaderByNumber(ctx, nil) // If `number` is nil, the latest known block header is returned.
	if err != nil {
		return nil, err
	}

	baseFee := header.BaseFee

	gasTipCap, err := b.getGasTipCap(ctx, client)
	if err != nil {
		return nil, err
	}

	// GasFeeCap = 2 * BaseFee + GasTipCap
	estimatedGasFeeCap := new(big.Int).Mul(baseFee, big.NewInt(2))
	b.gasFeeCap = estimatedGasFeeCap.Add(estimatedGasFeeCap, gasTipCap)

	return b.gasFeeCap, nil
}

func (b *TxBuilder) getGasPrice(ctx context.Context, client LegacyEthTxClient) (*big.Int, error) {
	if b.gasPrice != nil {
		return b.gasPrice, nil
	}

	if price, err := client.SuggestGasPrice(ctx); err != nil {
		return nil, err
	} else {
		b.gasPrice = price
	}

	return b.gasPrice, nil
}

func (b *TxBuilder) getChainID(ctx context.Context, cli EthTxClient) (*big.Int, error) {
	if b.chainID != nil {
		return b.chainID, nil
	}

	var err error
	b.chainID, err = cli.ChainID(ctx)

	return b.chainID, err
}

func (b *TxBuilder) getNonce(ctx context.Context, cli EthTxClient, from common.Address) (uint64, error) {
	if b.nonce != nil {
		return *b.nonce, nil
	}

	// The nonce cannot be cached because the from address is not stored directly in the Tx Builder, but is passed in as a parameter each time.
	return cli.PendingNonceAt(ctx, from)
}
