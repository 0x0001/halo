// Package wei provides functions for wei unit conversions
package wei

import (
	"math/big"

	"github.com/0x0001/halo/web3/convert"
)

// ToEther converts wei to ether
func ToEther(wei *big.Int) *big.Float {
	weiFloat := new(big.Float).SetInt(wei)
	return new(big.Float).Quo(weiFloat, big.NewFloat(convert.WeiPerEther))
}

// ToGwei converts wei to gwei
func ToGwei(wei *big.Int) *big.Float {
	weiFloat := new(big.Float).SetInt(wei)
	return new(big.Float).Quo(weiFloat, big.NewFloat(convert.WeiPerGwei))
}

// Convenience functions for basic data types using generics

// ToEtherFrom converts numeric wei values to ether
func ToEtherFrom[T convert.NumericType](wei T) *big.Float {
	weiBig := big.NewInt(int64(wei))
	return ToEther(weiBig)
}

// ToGweiFrom converts numeric wei values to gwei
func ToGweiFrom[T convert.NumericType](wei T) *big.Float {
	weiBig := big.NewInt(int64(wei))
	return ToGwei(weiBig)
}
