package testutils

import "strings"

const AlphaSize int = 26

var Alphabet []byte = make([]byte, AlphaSize)

func init() {
	for i := range Alphabet {
		Alphabet[i] = 'a' + byte(i)
	}
}

func CapitalizeWord(s string) string {
	if len(s) < 2 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
