package utils

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestDigitCode(t *testing.T) {
	codeLen := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, n := range codeLen {
		code := DigitCode(n)
		assert.Equal(t, len(code), n)
	}

	var prev string
	for i := 0; i < 100; i++ {
		code := DigitCode(6)
		assert.NotEqual(t, prev, code)
		prev = code
	}
}
