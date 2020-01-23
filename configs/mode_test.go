package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetMode(t *testing.T) {
	for mode, code := range map[string]int{
		DebugMode:   debugCode,
		ReleaseMode: releaseCode,
		TestMode:    testCode,
	} {
		SetMode(mode)
		assert.Equal(t, code, modeCode)
		assert.Equal(t, mode, Mode())
	}
	assert.Panics(t, func() { SetMode("panic") })
}
