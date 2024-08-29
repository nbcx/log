package log

import (
	"testing"
)

func TestInitLogger(t *testing.T) {
	Set(WithLevel("debug"))
	Debug("this is debug 1")
	Set(WithLevel("info"))
	Debug("this is debug 2")

	Set(WithLevel("info")) // "debug", log.LstdFlags|log.Lmicroseconds, WithShowFuncName()
	Info("this is info")
	Error("this is error")
}

func Test_runFuncName(t *testing.T) {
	tests := []struct {
		name  string
		depth int
		want  string
	}{
		{"test", 1, "runFuncName"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runFuncName(tt.depth); got != tt.want {
				t.Errorf("runFuncName() = %v, want %v", got, tt.want)
			}
		})
	}
}
