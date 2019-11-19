package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestDigitCodes(t *testing.T) {
	codeLen := 6
	codesLen := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, c := range codesLen {
		codes := DigitCodes(c, codeLen)
		assert.Equal(t, len(codes), c)

		var prev string
		for _, code := range codes {
			assert.NotEqual(t, prev, code)
			prev = code
		}
	}
}
