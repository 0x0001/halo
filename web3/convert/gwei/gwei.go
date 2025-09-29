// Package gwei provides functions for gwei unit conversions
package gwei

import (
	"math/big"

	"github.com/0x0001/halo/web3/convert"
)

// ToWei converts gwei to wei
func ToWei(gwei *big.Float) *big.Int {
	wei := new(big.Float).Mul(gwei, big.NewFloat(convert.WeiPerGwei))
	weiInt := new(big.Int)
	wei.Int(weiInt)
	return weiInt
}

// FromWei converts wei to gwei
func FromWei(wei *big.Int) *big.Float {
	weiFloat := new(big.Float).SetInt(wei)
	return new(big.Float).Quo(weiFloat, big.NewFloat(convert.WeiPerGwei))
}

// ToEther converts gwei to ether
func ToEther(gwei *big.Float) *big.Float {
	return new(big.Float).Quo(gwei, big.NewFloat(convert.GweiPerEther))
}

// FromEther converts ether to gwei
func FromEther(ether *big.Float) *big.Float {
	return new(big.Float).Mul(ether, big.NewFloat(convert.GweiPerEther))
}

// Convenience functions for basic data types using generics

// ToWeiFrom converts numeric gwei values to wei
func ToWeiFrom[T convert.NumericType](gwei T) *big.Int {
	gweiFloat := big.NewFloat(float64(gwei))
	return ToWei(gweiFloat)
}

// FromWeiFloat converts wei to float64 gwei value
func FromWeiFloat(wei *big.Int) float64 {
	gwei := FromWei(wei)
	gweiFloat, _ := gwei.Float64()
	return gweiFloat
}
