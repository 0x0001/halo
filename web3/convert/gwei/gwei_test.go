// Package gwei provides functions for gwei unit conversions
package gwei

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToWei(t *testing.T) {
	// Test integer value
	gwei := big.NewFloat(1) // 1 Gwei
	wei := ToWei(gwei)
	expected := big.NewInt(1000000000) // 1 Gwei = 10^9 wei
	assert.Equal(t, expected, wei)

	// Test decimal value
	gwei = big.NewFloat(0.5) // 0.5 Gwei
	wei = ToWei(gwei)
	expected = big.NewInt(500000000) // 0.5 Gwei = 5 * 10^8 wei
	assert.Equal(t, expected, wei)

	// Test large value
	gwei = big.NewFloat(1000000000) // 1 billion Gwei
	wei = ToWei(gwei)
	expected = new(big.Int).Mul(big.NewInt(1000000000), big.NewInt(1000000000))
	assert.Equal(t, expected, wei)
}

func TestFromWei(t *testing.T) {
	// Test exact conversion
	wei := big.NewInt(1000000000) // 10^9 wei
	gwei := FromWei(wei)
	expected := big.NewFloat(1)
	assert.Equal(t, 0, expected.Cmp(gwei))

	// Test decimal conversion
	wei = big.NewInt(500000000) // 5 * 10^8 wei
	gwei = FromWei(wei)
	expected = big.NewFloat(0.5)
	assert.Equal(t, 0, expected.Cmp(gwei))

	// Test small value
	wei = big.NewInt(1) // 1 wei
	gwei = FromWei(wei)
	expected1, _ := new(big.Float).SetString("0.000000001")
	assert.Equal(t, 0, expected1.Cmp(gwei))
}

func TestToEther(t *testing.T) {
	// Test integer value
	gwei := big.NewFloat(1000000000) // 1 billion Gwei
	ether := ToEther(gwei)
	expected := big.NewFloat(1) // 10^9 Gwei = 1 ETH
	assert.Equal(t, 0, expected.Cmp(ether))

	// Test decimal value
	gwei = big.NewFloat(500000000) // 500 million Gwei
	ether = ToEther(gwei)
	expected = big.NewFloat(0.5) // 5 * 10^8 Gwei = 0.5 ETH
	assert.Equal(t, 0, expected.Cmp(ether))
}

func TestFromEther(t *testing.T) {
	// Test integer value
	ether := big.NewFloat(1) // 1 ETH
	gwei := FromEther(ether)
	expected := big.NewFloat(1000000000) // 1 ETH = 10^9 Gwei
	assert.Equal(t, expected, gwei)

	// Test decimal value
	ether = big.NewFloat(0.5) // 0.5 ETH
	gwei = FromEther(ether)
	expected = big.NewFloat(500000000) // 0.5 ETH = 5 * 10^8 Gwei
	assert.Equal(t, expected, gwei)
	assert.Equal(t, 0, expected.Cmp(gwei))
}

func TestToWeiFrom(t *testing.T) {
	// Test integer value
	wei := ToWeiFrom(uint64(1))        // 1 Gwei
	expected := big.NewInt(1000000000) // 1 Gwei = 10^9 wei
	assert.Equal(t, expected, wei)

	// Test larger value
	wei = ToWeiFrom(uint64(5))        // 5 Gwei
	expected = big.NewInt(5000000000) // 5 Gwei = 5 * 10^9 wei
	assert.Equal(t, expected, wei)

	// Test with different types
	wei = ToWeiFrom(int64(1))         // 1 Gwei
	expected = big.NewInt(1000000000) // 1 Gwei = 10^9 wei
	assert.Equal(t, expected, wei)

	wei = ToWeiFrom(float64(1.0)) // 1 Gwei
	assert.Equal(t, expected, wei)

	wei = ToWeiFrom(float64(1.0)) // 1 Gwei
	assert.Equal(t, expected, wei)
}

func TestFromWeiFloat(t *testing.T) {
	// Test exact conversion
	gwei := FromWeiFloat(big.NewInt(1000000000)) // 1 Gwei in wei
	assert.Equal(t, 1.0, gwei)

	// Test decimal conversion
	gwei = FromWeiFloat(big.NewInt(500000000)) // 0.5 Gwei in wei
	assert.Equal(t, 0.5, gwei)
}
