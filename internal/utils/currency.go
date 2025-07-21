package utils

type Currency string

const (
	CurrencyUSD  Currency = "USD"
	CurrencyEURO Currency = "EURO"
	CurrencyPLN  Currency = "PLN"
)

func IsCurrency(currency Currency) bool {
	switch currency {
	case CurrencyUSD, CurrencyEURO, CurrencyPLN:
		return true
	default:
		return false
	}
}
