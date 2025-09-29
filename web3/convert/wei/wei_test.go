// Package wei provides functions for wei unit conversions
package wei

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToEther(t *testing.T) {
	// Test exact conversion
	wei := big.NewInt(1000000000000000000) // 10^18 wei
	ether := ToEther(wei)
	expected := big.NewFloat(1)
	assert.Equal(t, 0, expected.Cmp(ether))

	// Test decimal conversion
	wei = big.NewInt(500000000000000000) // 5 * 10^17 wei
	ether = ToEther(wei)
	expected = big.NewFloat(0.5)
	assert.Equal(t, 0, expected.Cmp(ether))

	// Test small value
	wei = big.NewInt(1) // 1 wei
	ether = ToEther(wei)
	expected1, _ := new(big.Float).SetString("0.000000000000000001")
	assert.Equal(t, 0, expected1.Cmp(ether))
}

func TestToGwei(t *testing.T) {
	// Test exact conversion
	wei := big.NewInt(1000000000) // 10^9 wei
	gwei := ToGwei(wei)
	expected := big.NewFloat(1)
	assert.Equal(t, 0, expected.Cmp(gwei))

	// Test decimal conversion
	wei = big.NewInt(500000000) // 5 * 10^8 wei
	gwei = ToGwei(wei)
	expected = big.NewFloat(0.5)
	assert.Equal(t, 0, expected.Cmp(gwei))

	// Test small value
	wei = big.NewInt(1) // 1 wei
	gwei = ToGwei(wei)
	expected1, _ := new(big.Float).SetString("0.000000001")
	assert.Equal(t, 0, expected1.Cmp(gwei))
}

func TestToEtherFrom(t *testing.T) {
	// Test exact conversion
	ether := ToEtherFrom(uint64(1000000000000000000)) // 1 ETH in wei
	expected := big.NewFloat(1.0)
	assert.Equal(t, 0, expected.Cmp(ether))

	// Test decimal conversion
	ether = ToEtherFrom(uint64(500000000000000000)) // 0.5 ETH in wei
	expected = big.NewFloat(0.5)
	assert.Equal(t, 0, expected.Cmp(ether))

	// Test with different types
	ether = ToEtherFrom(int64(1000000000000000000)) // 1 ETH in wei
	expected = big.NewFloat(1.0)
	assert.Equal(t, 0, expected.Cmp(ether))

	ether = ToEtherFrom(float64(1000000000000000000.0)) // 1 ETH in wei
	assert.Equal(t, 0, expected.Cmp(ether))
}

func TestToGweiFrom(t *testing.T) {
	// Test exact conversion
	gwei := ToGweiFrom(uint64(1000000000)) // 1 Gwei in wei
	expected := big.NewFloat(1.0)
	assert.Equal(t, 0, expected.Cmp(gwei))

	// Test decimal conversion
	gwei = ToGweiFrom(uint64(500000000)) // 0.5 Gwei in wei
	expected = big.NewFloat(0.5)
	assert.Equal(t, 0, expected.Cmp(gwei))

	// Test with different types
	gwei = ToGweiFrom(int64(1000000000)) // 1 Gwei in wei
	expected = big.NewFloat(1.0)
	assert.Equal(t, 0, expected.Cmp(gwei))

	gwei = ToGweiFrom(float64(1000000000.0)) // 1 Gwei in wei
	assert.Equal(t, 0, expected.Cmp(gwei))
}
