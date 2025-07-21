package utils

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const AlphaSize int = 26

var Alphabet []byte = make([]byte, AlphaSize, AlphaSize)

func init() {
	for i := range Alphabet {
		Alphabet[i] = 'a' + byte(i)
	}
}

func GenRandString(size int) string {
	var sb strings.Builder
	for i := 0; i < size; i++ {
		sb.WriteByte(Alphabet[rand.IntN(AlphaSize)])
	}
	return sb.String()
}

func GenRandFullName() string {
	firstName := CapitalizeWord(GenRandString(RandIntInRange(2, 10)))
	lastName := CapitalizeWord(GenRandString(RandIntInRange(2, 10)))
	return fmt.Sprintf("%s %s", firstName, lastName)
}

func GenRandCurrency() Currency {
	switch rand.IntN(3) {
	case 0:
		return CurrencyEURO
	case 1:
		return CurrencyPLN
	default:
		return CurrencyUSD
	}
}

func GetRandEmail() string {
	domain3 := GenRandString(RandIntInRange(2, 10))
	domain2 := GenRandString(RandIntInRange(2, 10))
	domain1 := GenRandString(RandIntInRange(2, 3))
	return fmt.Sprintf("%s@%s.%s", domain3, domain2, domain1)
}

func GenRandHashedPassword() string {
	pass, _ := bcrypt.GenerateFromPassword([]byte(GenRandString(10)), bcrypt.MinCost)
	return string(pass)
}
