package random

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math/rand"
	"strings"
	"testing"
)

func TestInts(t *testing.T) {
	minNum, maxNum := rand.Int(), rand.Int()
	if minNum > maxNum {
		temp := minNum
		minNum = maxNum
		maxNum = temp
	}
	randInt64 := Int64(minNum, maxNum)
	fmt.Println("randInt64, maxNum", randInt64, maxNum)
	require.LessOrEqual(t, randInt64, int64(maxNum))
	require.GreaterOrEqual(t, randInt64, int64(minNum))

	randInt := Int(minNum, maxNum)
	require.LessOrEqual(t, randInt, maxNum)
	require.GreaterOrEqual(t, randInt, minNum)
}

func TestString(t *testing.T) {
	size := 10
	result := String(size)

	require.Equal(t, size, len(result))

	for _, char := range result {
		require.True(t, strings.Contains(alphabet, string(char)))
	}
}

func TestCurrency(t *testing.T) {
	result := Currency()

	validCurrencies := map[string]bool{
		"USD": true,
		"EUR": true,
		"PLN": true,
	}

	require.Contains(t, validCurrencies, result)
}
