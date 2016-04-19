package helper

import (
	"math"
	"time"
)

func Truncate(f float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	return float64(int(f*multiplier)) / multiplier
}

func ParseTime(timestamp string) (submittime time.Time, err error) {
	submittime, err = time.Parse("2006-01-02T15:04:05.999999999Z07:00", timestamp)
	if err != nil {
		submittime, err = time.Parse("2006-01-02 15:04:05", timestamp)
		if err != nil {
			submittime, err = time.Parse("2006-01-02 15:04:05.999999999Z07:00", timestamp)
			return
		}
		submittime = submittime.Add(8 * time.Hour)
	}
	return
}
