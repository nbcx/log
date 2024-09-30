package log

import (
	"io"
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

type Options func(g *Log)

// Set log
func Set(options ...Options) { // level string, flag int,
	for _, op := range options {
		op(std)
	}
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

type Logger interface {
	Debug(msg any, a ...interface{})
	Info(msg any, a ...interface{})
	Warn(msg any, a ...interface{})
	Error(msg any, a ...interface{})
	Panic(msg any, a ...interface{})
	Fatal(msg any, a ...interface{})
}

// WithPath 设置日志输出路径
func WithLevel(level string) Options {
	return func(g *Log) {
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
