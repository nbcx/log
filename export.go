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

// WithPath 设置日志输出路径
func (l *Logger) WithLevel(level string) {
	l.SetLevel(getLevel(level))
}

// 设置日志格式器
func (l *Logger) SetFormatter(level string) {
	l.SetLevel(getLevel(level))
}

// 设置日志输出
func (l *Logger) SetWriter(w io.Writer) {
	l.writer = w
}

// 为指定等级的日志设置额外的输出
// 通常用于需要特别关注的紧急日志
func SetLevelWriter(level string, w ...io.Writer) {
	std.SetLevelWriter(level, w...)
}

func (l *Logger) SetLevelWriter(level string, w ...io.Writer) {
	lv := getLevel(level)
	l.extra[lv] = append(l.extra[lv], w...)
}

// 设置日志显示格式
// 需要注意，对于一些自定义的formatter，它并不是绝对生效的
func (l *Logger) SetFlag(flag int) {
	l.opt.SetFlags(flag)
}

// 通常我们需要把日志传递给第三方模块使用，并想要标记是第三方模块
// 那么可以使用此函数获取一个带有指定标记的日志实例
func (l *Logger) SetPrefix(flag int) *Logger {
	l.opt.SetFlags(flag)

	// todo: wait do
	return l
}
