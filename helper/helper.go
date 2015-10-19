package helper

import (
	"math"
)

func Truncate(f float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	return float64(int(f*multiplier)) / multiplier
}
