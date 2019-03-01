package alog

import (
	"testing"
)

// TestConfigInLocal tests if logging works according to the configuration in alog.conf present in current directory
func TestLogConfigInLocal(t *testing.T) {
	// Check if
}

func TestLogConfigInEnvtVar(t *testing.T) {

}

func TestLogLevels(t *testing.T) {
	SetLogLevel(ERROR)
	Trace("This is a TRACE message.")
	Debug("This is a DEBUG message.")
	Info("This is an INFO message.")
	Warn("This is a WARN message.")
	Error("This is an ERROR message.")
	Critical("This is a CRITICAL message.")

}
