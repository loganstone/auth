package configs

import "os"

// EnvMode indicates environment name for mode
const EnvMode = EnvPrefix + "MODE"

const (
	// DebugMode indicates mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates mode is release.
	ReleaseMode = "release"
	// TestMode indicates mode is test.
	TestMode = "test"
)

const (
	debugCode = iota
	releaseCode
	testCode
)

var modeCode = debugCode
var modeName = DebugMode

func init() {
	mode := os.Getenv(EnvMode)
	SetMode(mode)
}

// SetMode sets mode according to input string.
func SetMode(value string) {
	switch value {
	case DebugMode, "":
		modeCode = debugCode
	case ReleaseMode:
		modeCode = releaseCode
	case TestMode:
		modeCode = testCode
	default:
		panic("mode unknown: " + value)
	}
	if value == "" {
		value = DebugMode
	}
	modeName = value
}

// Mode returns currently mode.
func Mode() string {
	return modeName
}
