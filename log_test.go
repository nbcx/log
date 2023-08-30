package log

import (
	"log"
	"testing"
)

func TestInitLogger(t *testing.T) {
	Init("debug", log.LstdFlags|log.Lmicroseconds, WithShowFuncName())
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
