// Package convert provides utility functions for cryptocurrency unit conversions
package convert

// NumericType represents all supported numeric types for unit conversions
type NumericType interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

// Common cryptocurrency units
const (
	WeiPerGwei   = 1000000000
	WeiPerEther  = 1000000000000000000
	GweiPerEther = 1000000000
)
