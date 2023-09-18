package util

import (
	"fmt"
	"strings"
)

func PairFromCurrencies(fromCurrency, toCurrency string) string {
	return fmt.Sprintf("%s_%s", fromCurrency, toCurrency)
}

func CurrenciesFromPair(currencyPair string) (fromCurrency, toCurrency string) {
	parts := strings.Split(currencyPair, "_")
	if len(parts) != 2 {
		return "", ""
	}

	return parts[0], parts[1]
}
