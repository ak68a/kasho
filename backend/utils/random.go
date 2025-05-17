package utils

import (
	"fmt"
	"math/rand"
)

var alphabet string = ("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomString(r int) string {
	bits := []rune{}
	k := len(alphabet)

	for i := 0; i < r; i++ {
		index := rand.Intn(k)
		bits = append(bits, rune(alphabet[index]))
	}

	return string(bits)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@example.com", RandomString(6))
}