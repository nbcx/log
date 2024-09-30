package log

import (
	"fmt"
	"os"
	"testing"
)

func TestInitLogger(t *testing.T) {
	Set(WithLevel("debug"))

	SetLevelWriter("error", os.Stderr)

	Debug("this is debug 1")
	Set(WithLevel("info"))
	Debug("this is debug 2")

	// Set(WithLevel("info"), WithWriter(os.Stdout)) // "debug", log.LstdFlags|log.Lmicroseconds, WithShowFuncName()
	Info("this is info %v", 111, 1113)
	Error("this is error")
	Warn(1213)
	// Panic("this is error")

	fmt.Println("hello")

}
