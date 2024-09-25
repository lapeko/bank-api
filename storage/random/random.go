package random

import (
	"math/rand"
	"strings"
)

var (
	alphabet   = "abcdefghijklmnopqrstuvwxyz"
	currencies = [...]string{"USD", "EUR", "PLN"}
)

func Int64(min int, max int) int64 {
	return int64(min) + rand.Int63n(int64(max-min+1))
}

func Int(min int, max int) int {
	return min + rand.Intn(max-min+1)
}

func String(size int) string {
	buffer := strings.Builder{}
	alpSize := len(alphabet)

	for i := 0; i < size; i++ {
		buffer.WriteByte(alphabet[rand.Intn(alpSize)])
	}

	return buffer.String()
}

func Currency() string {
	return currencies[rand.Intn(len(currencies))]
}
