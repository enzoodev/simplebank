package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	min := int64(0)
	max := int64(100)
	randomInt := RandomInt(min, max)
	require.True(t, randomInt >= min && randomInt <= max, "random int is out of bounds")
}

func TestRandomString(t *testing.T) {
	n := 6
	randomStr := RandomString(n)
	require.Len(t, randomStr, n, "random string length should be %d", n)
}

func TestRandomOwner(t *testing.T) {
	randomOwner := RandomOwner()
	require.Len(t, randomOwner, 6, "random owner length should be 6")
}

func TestRandomMoney(t *testing.T) {
	randomMoney := RandomMoney()
	require.True(t, randomMoney >= 0 && randomMoney <= 1000, "random money is out of bounds")
}

func TestRandomCurrency(t *testing.T) {
	currency := RandomCurrency()
	validCurrencies := map[string]bool{
		"USD": true,
		"EUR": true,
		"CAD": true,
	}
	require.True(t, validCurrencies[currency], "currency should be one of USD, EUR, or CAD")
}

func TestInt64ToNullInt64(t *testing.T) {
	value := int64(12345)
	nullInt64 := Int64ToNullInt64(value)

	require.Equal(t, value, nullInt64.Int64, "expected and actual values don't match")
	require.True(t, nullInt64.Valid, "sql.NullInt64 should be valid")
}

func TestInt64ToNullInt64Zero(t *testing.T) {
	value := int64(0)
	nullInt64 := Int64ToNullInt64(value)

	require.Equal(t, value, nullInt64.Int64, "expected 0 but got different value")
	require.True(t, nullInt64.Valid, "sql.NullInt64 should be valid when the value is 0")
}
