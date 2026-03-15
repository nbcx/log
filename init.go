package log

import (
	"os"

	"github.com/shiena/ansicolor"
)

var std *Log

func init() {
	std = new(Log)
	std.level = debugLevel
	std.opt = &Option{
		CallDepth:    4,
		ShowFuncName: true,
		Flag:         LstdFlags | Ltime | Lshortfile,
	}
	std.format = NewConsole()
	std.writer = ansicolor.NewAnsiColorWriter(os.Stdout)
}
