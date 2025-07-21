package utils

import "math/rand/v2"

func RandIntInRange(min int, max int) int {
	return min + rand.IntN(max-min+1)
}

func RandInt64InRange(min int, max int) int64 {
	return int64(RandIntInRange(min, max))
}
