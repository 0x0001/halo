// Package ether provides functions for ether unit conversions
package ether

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToWei(t *testing.T) {
	// Test integer value
	ether := big.NewFloat(1) // 1 ETH
	wei := ToWei(ether)
	expected := big.NewInt(1000000000000000000) // 1 ETH = 10^18 wei
	assert.Equal(t, expected, wei)

	// Test decimal value
	ether = big.NewFloat(0.5) // 0.5 ETH
	wei = ToWei(ether)
	expected = big.NewInt(500000000000000000) // 0.5 ETH = 5 * 10^17 wei
	assert.Equal(t, expected, wei)

	// Test large value
	ether = big.NewFloat(1000) // 1000 ETH
	wei = ToWei(ether)
	expected = new(big.Int).Mul(big.NewInt(1000), big.NewInt(1000000000000000000))
	assert.Equal(t, expected, wei)
}

func TestFromWei(t *testing.T) {
	// Test exact conversion
	wei := big.NewInt(1000000000000000000) // 10^18 wei
	ether := FromWei(wei)
	expected := big.NewFloat(1)
	assert.Equal(t, 0, expected.Cmp(ether))

	// Test decimal conversion
	wei = big.NewInt(500000000000000000) // 5 * 10^17 wei
	ether = FromWei(wei)
	expected = big.NewFloat(0.5)
	assert.Equal(t, 0, expected.Cmp(ether))

	// Test small value
	wei = big.NewInt(1) // 1 wei
	ether = FromWei(wei)
	expected1, _ := new(big.Float).SetString("0.000000000000000001")
	assert.Equal(t, 0, expected1.Cmp(ether))
}

func TestToGwei(t *testing.T) {
	// Test integer value
	ether := big.NewFloat(1) // 1 ETH
	gwei := ToGwei(ether)
	expected := big.NewFloat(1000000000) // 1 ETH = 10^9 Gwei
	assert.Equal(t, expected, gwei)

	// Test decimal value
	ether = big.NewFloat(0.5) // 0.5 ETH
	gwei = ToGwei(ether)
	expected = big.NewFloat(500000000) // 0.5 ETH = 5 * 10^8 Gwei
	assert.Equal(t, expected, gwei)
}

func TestFromGwei(t *testing.T) {
	// Test integer value
	gwei := big.NewFloat(1000000000) // 1 billion Gwei
	ether := FromGwei(gwei)
	expected := big.NewFloat(1) // 10^9 Gwei = 1 ETH
	assert.Equal(t, 0, expected.Cmp(ether))

	// Test decimal value
	gwei = big.NewFloat(500000000) // 500 million Gwei
	ether = FromGwei(gwei)
	expected = big.NewFloat(0.5) // 5 * 10^8 Gwei = 0.5 ETH
	assert.Equal(t, 0, expected.Cmp(ether))
}

func TestToWeiFrom(t *testing.T) {
	// Test integer value
	wei := ToWeiFrom(uint64(1))                 // 1 ETH
	expected := big.NewInt(1000000000000000000) // 1 ETH = 10^18 wei
	assert.Equal(t, expected, wei)

	// Test larger value
	wei = ToWeiFrom(uint64(5))                 // 5 ETH
	expected = big.NewInt(5000000000000000000) // 5 ETH = 5 * 10^18 wei
	assert.Equal(t, expected, wei)

	// Test with different types
	wei = ToWeiFrom(int64(1))                  // 1 ETH
	expected = big.NewInt(1000000000000000000) // 1 ETH = 10^18 wei
	assert.Equal(t, expected, wei)

	wei = ToWeiFrom(float64(1.0)) // 1 ETH
	assert.Equal(t, expected, wei)
}

func TestFromWeiFloat64(t *testing.T) {
	// Test exact conversion
	ether := FromWeiFloat(big.NewInt(1000000000000000000)) // 1 ETH in wei
	assert.Equal(t, 1.0, ether)

	// Test decimal conversion
	ether = FromWeiFloat(big.NewInt(500000000000000000)) // 0.5 ETH in wei
	assert.Equal(t, 0.5, ether)
}
