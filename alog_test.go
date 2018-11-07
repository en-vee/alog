package alog

import (
	"testing"
)

func TestLogConfigInLocal(t *testing.T) {

}

func TestLogConfigInEnvtVar(t *testing.T) {

}

func TestLogLevels(t *testing.T) {

	Trace("This is a TRACE message.")
	Debug("This is a DEBUG message.")
	Info("This is an INFO message.")
	Warn("This is a WARN message.")
	Error("This is an ERROR message.")
	Critical("This is a CRITICAL message.")
}
