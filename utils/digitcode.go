package utils

import (
	"math/rand"
	"time"
)

const digits = "0123456789"

// DigitCode .
func DigitCode(n int) (code string) {
	rand.Seed(time.Now().UnixNano())
	for len(code) < n {
		num := rand.Intn(len(digits))
		code += string(digits[num])
	}
	return
}
