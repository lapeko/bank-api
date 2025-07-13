package utils

import "math/rand/v2"

func RandIntInRange(min int, max int) int {
	return min + rand.IntN(max-min+1)
}
