// Package ether provides functions for ether unit conversions
package ether

import (
	"math/big"

	"github.com/0x0001/halo/web3/convert"
)

// ToWei converts ether to wei
func ToWei(ether *big.Float) *big.Int {
	wei := new(big.Float).Mul(ether, big.NewFloat(convert.WeiPerEther))
	weiInt := new(big.Int)
	wei.Int(weiInt)
	return weiInt
}

// FromWei converts wei to ether
func FromWei(wei *big.Int) *big.Float {
	weiFloat := new(big.Float).SetInt(wei)
	return new(big.Float).Quo(weiFloat, big.NewFloat(convert.WeiPerEther))
}

// ToGwei converts ether to gwei
func ToGwei(ether *big.Float) *big.Float {
	return new(big.Float).Mul(ether, big.NewFloat(convert.GweiPerEther))
}

// FromGwei converts gwei to ether
func FromGwei(gwei *big.Float) *big.Float {
	return new(big.Float).Quo(gwei, big.NewFloat(convert.GweiPerEther))
}

// Convenience functions for basic data types using generics

// ToWeiFrom converts numeric ether values to wei
func ToWeiFrom[T convert.NumericType](ether T) *big.Int {
	etherFloat := big.NewFloat(float64(ether))
	return ToWei(etherFloat)
}

// FromWeiFloat converts wei to float64 ether value
func FromWeiFloat(wei *big.Int) float64 {
	ether := FromWei(wei)
	etherFloat, _ := ether.Float64()
	return etherFloat
}
