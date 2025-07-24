package utils

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/lapeko/udemy__backend-master-class-golang-postgresql-kubernetes/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func GenRandIntInRange(min int, max int) int {
	return min + rand.IntN(max-min+1)
}

func GenRandInt64InRange(min int, max int) int64 {
	return int64(GenRandIntInRange(min, max))
}

func GenRandString(size int) string {
	var sb strings.Builder
	for i := 0; i < size; i++ {
		sb.WriteByte(Alphabet[rand.IntN(AlphaSize)])
	}
	return sb.String()
}

func GenRandFullName() string {
	firstName := CapitalizeWord(GenRandString(GenRandIntInRange(2, 10)))
	lastName := CapitalizeWord(GenRandString(GenRandIntInRange(2, 10)))
	return fmt.Sprintf("%s %s", firstName, lastName)
}

func GenRandCurrency() utils.Currency {
	switch rand.IntN(3) {
	case 0:
		return utils.CurrencyEURO
	case 1:
		return utils.CurrencyPLN
	default:
		return utils.CurrencyUSD
	}
}

func GetRandEmail() string {
	domain3 := GenRandString(GenRandIntInRange(2, 10))
	domain2 := GenRandString(GenRandIntInRange(2, 10))
	domain1 := GenRandString(GenRandIntInRange(2, 3))
	return fmt.Sprintf("%s@%s.%s", domain3, domain2, domain1)
}

func GenRandHashedPassword() string {
	pass, _ := bcrypt.GenerateFromPassword([]byte(GenRandString(10)), bcrypt.MinCost)
	return string(pass)
}
