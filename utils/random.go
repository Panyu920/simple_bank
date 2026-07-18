package utils

import (
	"math/rand"
	"strings"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyz"
)

var currency = [3]string{"USD", "EUR", "CAD"}

func RandomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func RandomOwner() string {
	nameLength := RandomInt(3, 6)
	var name strings.Builder
	name.Grow(nameLength)
	for range nameLength {
		name.WriteByte(letters[rand.Intn(26)])
	}

	return name.String()
}

func RandomBalance() int64 {
	return int64(RandomInt(100, 10000))
}

func RandomCurrency() string {
	return currency[rand.Intn(len(currency))]
}
