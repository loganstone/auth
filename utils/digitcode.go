package utils

import (
	"math/rand"
	"time"
)

const digits = "0123456789"

// DigitCode .
func DigitCode(n int) string {
	rand.Seed(time.Now().UnixNano())
	code := make([]rune, n)
	for i := 0; i < n; i++ {
		idx := rand.Intn(len(digits))
		code[i] = rune(digits[idx])
	}
	return string(code)
}

// DigitCodes .
func DigitCodes(c, n int) []string {
	codes := make([]string, c)
	for i := 0; i < c; i++ {
		codes[i] = DigitCode(n)
	}
	return codes
}
