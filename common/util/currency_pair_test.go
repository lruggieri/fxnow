package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPairFromCurrencies(t *testing.T) {
	assert.Equal(t, "USD_JPY", PairFromCurrencies("USD", "JPY"))
}

func TestCurrenciesFromPair(t *testing.T) {
	t.Run("happy-path", func(t *testing.T) {
		from, to := CurrenciesFromPair("USD_JPY")
		assert.Equal(t, "USD", from)
		assert.Equal(t, "JPY", to)
	})

	t.Run("invalid-pair", func(t *testing.T) {
		from, to := CurrenciesFromPair("USDJPY")
		assert.Equal(t, "", from)
		assert.Equal(t, "", to)
	})
}
