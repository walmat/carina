package util

import "math/rand"

func RangeIn(low, hi int) int {
	return low + rand.Intn(hi-low)
}
