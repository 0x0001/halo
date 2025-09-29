package tx

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEthTxClient is a mock implementation of the EthTxClient interface
type MockEthTxClient struct {
	mock.Mock
}

func (m *MockEthTxClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	args := m.Called(ctx, account)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockEthTxClient) ChainID(ctx context.Context) (*big.Int, error) {
	args := m.Called(ctx)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *MockEthTxClient) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	args := m.Called(ctx, msg)
	return args.Get(0).(uint64), args.Error(1)
}

// MockLegacyEthTxClient is a mock implementation of the LegacyEthTxClient interface
type MockLegacyEthTxClient struct {
	MockEthTxClient
}

func (m *MockLegacyEthTxClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	args := m.Called(ctx)
	return args.Get(0).(*big.Int), args.Error(1)
}

// MockDynamicFeeEthTxClient is a mock implementation of the DynamicFeeEthTxClient interface
type MockDynamicFeeEthTxClient struct {
	MockLegacyEthTxClient
}

func (m *MockDynamicFeeEthTxClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	args := m.Called(ctx, number)
	return args.Get(0).(*types.Header), args.Error(1)
}

func (m *MockDynamicFeeEthTxClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	args := m.Called(ctx)
	return args.Get(0).(*big.Int), args.Error(1)
}

// generateTestPrivateKey generates a private key for testing
func generateTestPrivateKey(t *testing.T) *ecdsa.PrivateKey {
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	return privateKey
}

func TestNewBuilder(t *testing.T) {
	toAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")
	builder := NewBuilder(&toAddress)

	assert.NotNil(t, builder)
	assert.Equal(t, &toAddress, builder.to)
	assert.Equal(t, big.NewInt(0), builder.value)
	assert.Equal(t, uint64(105), builder.gasMultiplier)
	assert.Equal(t, uint64(100), builder.gasDivisor)
}

func TestBuilderMethods(t *testing.T) {
	builder := NewBuilder(nil)
	testValue := big.NewInt(1000)
	testData := []byte("test data")
	testNonce := uint64(5)
	testChainID := big.NewInt(1)
	testGasPrice := big.NewInt(20000000000)
	testGasTipCap := big.NewInt(1000000000)
	testGasFeeCap := big.NewInt(30000000000)

	// Test method chaining
	result := builder.Value(testValue).
		Data(testData).
		Nonce(testNonce).
		ChainID(testChainID).
		GasPrice(testGasPrice).
		GasTipCap(testGasTipCap).
		GasFeeCap(testGasFeeCap).
		GasLimitMultiplier(110, 100)

	assert.Equal(t, testValue, builder.value)
	assert.Equal(t, testData, builder.data)
	assert.Equal(t, &testNonce, builder.nonce)
	assert.Equal(t, testChainID, builder.chainID)
	assert.Equal(t, testGasPrice, builder.gasPrice)
	assert.Equal(t, testGasTipCap, builder.gasTipCap)
	assert.Equal(t, testGasFeeCap, builder.gasFeeCap)
	assert.Equal(t, uint64(110), builder.gasMultiplier)
	assert.Equal(t, uint64(100), builder.gasDivisor)
	assert.Equal(t, builder, result) // Ensure it returns itself to support method chaining
}

func TestAdjustGasLimit(t *testing.T) {
	tests := []struct {
		gas        uint64
		multiplier uint64
		divisor    uint64
		expected   uint64
	}{
		{100000, 105, 100, 105000},
		{50000, 110, 100, 55000},
		{21000, 100, 100, 21000}, // 不增加
		{100000, 120, 100, 120000},
	}

	for _, tt := range tests {
		result := adjustGasLimit(tt.gas, tt.multiplier, tt.divisor)
		assert.Equal(t, tt.expected, result)
	}
}

func TestEstimateGasFrom(t *testing.T) {
	mockClient := new(MockEthTxClient)
	builder := NewBuilder(nil)
	fromAddress := common.HexToAddress("0x1234567890123456789012345678901234567890")

	expectedGas := uint64(21000)
	mockClient.On("EstimateGas", mock.Anything, mock.AnythingOfType("ethereum.CallMsg")).Return(expectedGas, nil)

	gas, err := builder.EstimateGasFrom(context.Background(), mockClient, fromAddress)

	assert.NoError(t, err)
	assert.Equal(t, expectedGas, gas)
	mockClient.AssertExpectations(t)
}

func TestSignLegacy(t *testing.T) {
	mockClient := new(MockLegacyEthTxClient)
	builder := NewBuilder(nil)
	privateKey := generateTestPrivateKey(t)

	// Set mock expectations
	mockClient.MockEthTxClient.On("PendingNonceAt", mock.Anything, mock.AnythingOfType("common.Address")).Return(uint64(1), nil)
	mockClient.On("SuggestGasPrice", mock.Anything).Return(big.NewInt(20000000000), nil)
	mockClient.MockEthTxClient.On("EstimateGas", mock.Anything, mock.AnythingOfType("ethereum.CallMsg")).Return(uint64(21000), nil)

	tx, err := builder.SignLegacy(context.Background(), mockClient, privateKey)

	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, uint64(1), tx.Nonce())
	assert.Equal(t, big.NewInt(20000000000), tx.GasPrice())
	// 21000 * 105 / 100 = 22050
	assert.Equal(t, uint64(22050), tx.Gas())
	mockClient.AssertExpectations(t)
}

func TestSignDynamicFee(t *testing.T) {
	mockClient := new(MockDynamicFeeEthTxClient)
	builder := NewBuilder(nil)
	privateKey := generateTestPrivateKey(t)
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Set mock expectations
	mockClient.MockEthTxClient.On("PendingNonceAt", mock.Anything, fromAddress).Return(uint64(1), nil)
	mockClient.MockEthTxClient.On("ChainID", mock.Anything).Return(big.NewInt(1), nil)
	mockClient.MockEthTxClient.On("EstimateGas", mock.Anything, mock.AnythingOfType("ethereum.CallMsg")).Return(uint64(21000), nil)
	mockClient.On("SuggestGasTipCap", mock.Anything).Return(big.NewInt(1000000000), nil)
	mockClient.On("HeaderByNumber", mock.Anything, mock.Anything).Return(&types.Header{
		BaseFee: big.NewInt(20000000000),
	}, nil)

	tx, err := builder.SignDynamicFee(context.Background(), mockClient, privateKey)

	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, uint64(1), tx.Nonce())
	assert.Equal(t, big.NewInt(1000000000), tx.GasTipCap())
	assert.Equal(t, big.NewInt(41000000000), tx.GasFeeCap()) // 2 * 20000000000 + 1000000000
	// 21000 * 105 / 100 = 22050
	assert.Equal(t, uint64(22050), tx.Gas())
	mockClient.AssertExpectations(t)
}

func TestSignDynamicFeeGasValidation(t *testing.T) {
	mockClient := new(MockDynamicFeeEthTxClient)
	builder := NewBuilder(nil)
	privateKey := generateTestPrivateKey(t)

	// Set mock expectations
	mockClient.MockEthTxClient.On("ChainID", mock.Anything).Return(big.NewInt(1), nil)
	mockClient.MockEthTxClient.On("PendingNonceAt", mock.Anything, mock.AnythingOfType("common.Address")).Return(uint64(1), nil)
	mockClient.MockEthTxClient.On("EstimateGas", mock.Anything, mock.AnythingOfType("ethereum.CallMsg")).Return(uint64(21000), nil)

	// Set gas fee cap less than tip cap to trigger validation error
	builder.GasFeeCap(big.NewInt(1000000000)) // 1 gwei
	builder.GasTipCap(big.NewInt(2000000000)) // 2 gwei (larger)

	_, err := builder.SignDynamicFee(context.Background(), mockClient, privateKey)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "gasFeeCap must be greater than or equal to gasTipCap")
}

func TestGettersWithCachedValues(t *testing.T) {
	builder := NewBuilder(nil)
	testNonce := uint64(5)
	testChainID := big.NewInt(1)
	testGasPrice := big.NewInt(20000000000)
	testGasTipCap := big.NewInt(1000000000)
	testGasFeeCap := big.NewInt(30000000000)

	// Pre-set values
	builder.nonce = &testNonce
	builder.chainID = testChainID
	builder.gasPrice = testGasPrice
	builder.gasTipCap = testGasTipCap
	builder.gasFeeCap = testGasFeeCap

	// Create mock client, but it won't be called because values are cached
	mockClient := new(MockDynamicFeeEthTxClient)

	nonce, err := builder.getNonce(context.Background(), mockClient, common.Address{})
	assert.NoError(t, err)
	assert.Equal(t, testNonce, nonce)

	chainID, err := builder.getChainID(context.Background(), mockClient)
	assert.NoError(t, err)
	assert.Equal(t, testChainID, chainID)

	gasPrice, err := builder.getGasPrice(context.Background(), mockClient)
	assert.NoError(t, err)
	assert.Equal(t, testGasPrice, gasPrice)

	gasTipCap, err := builder.getGasTipCap(context.Background(), mockClient)
	assert.NoError(t, err)
	assert.Equal(t, testGasTipCap, gasTipCap)

	gasFeeCap, err := builder.getGasFeeCap(context.Background(), mockClient)
	assert.NoError(t, err)
	assert.Equal(t, testGasFeeCap, gasFeeCap)
}
