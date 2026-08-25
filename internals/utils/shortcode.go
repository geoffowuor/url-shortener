package utils

import (
	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortCode(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, length)

	for i := range code {
		code[i] = letters[r.Intn(len(letters))]
	}

	return string(code)
}
