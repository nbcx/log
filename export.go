package log

import (
	"io"
	"os"

	"github.com/shiena/ansicolor"
)

var std *Logger

func init() {
	std = new(Logger)
	std.level = debugLevel
	std.opt = &Option{
		CallDepth:    4,
		ShowFuncName: true,
		Flag:         LstdFlags | Ltime | Lshortfile,
	}
	std.format = NewConsole()
	std.writer = ansicolor.NewAnsiColorWriter(os.Stdout)
}

type Options func(g *Logger)

// Set log
func Set(options ...Options) { // level string, flag int,
	for _, op := range options {
		op(std)
	}
}

// GetLogger default log
func GetLogger() *Logger {
	return std
}

// WithPath 设置日志输出路径
// func WithPath(logDir string) Options {
// 	return func(g *Logger) {
// 		std.handle.SetOutput(g.combineFileWithStdWriter(filepath.Join(logDir, globalFileName), os.Stdout))
// 	}
// }

// WithWriter 设置日志输出Writer
// func WithWriter(w io.Writer) Options {
// 	return func(g *Logger) {
// 		g.handle.SetOutput(w)
// 	}
// }

// WithPath 设置日志输出路径
func WithLevel(level string) Options {
	return func(g *Logger) {
		g.SetLevel(getLevel(level))
	}
}

// Debug Debug
func Debug(format any, a ...interface{}) {
	std.Debug(format, a...)
}

// Info Info
func Info(format any, a ...interface{}) {
	std.Info(format, a...)
}

// Warn Warn
func Warn(format any, a ...interface{}) {
	std.Warn(format, a...)
}

// Error Error
func Error(format any, a ...interface{}) {
	std.Error(format, a...)
}

// Error Error
func Panic(format any, a ...interface{}) {
	std.Panic(format, a...)
}

// Fatal Fatal
func Fatal(format any, a ...interface{}) {
	std.Fatal(format, a...)
}

// 为指定等级的日志设置额外的输出
// 通常用于需要特别关注的紧急日志
func SetLevelWriter(level string, w ...io.Writer) {
	std.SetLevelWriter(level, w...)
}
