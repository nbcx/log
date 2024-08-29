package log

import (
	"io"
	"os"
	"path/filepath"
)

type Options func(g *Logger)

// Set log
func Set(options ...Options) { // level string, flag int,
	// if isInit {
	// 	return
	// }
	// isInit = true
	// gLogger.level, eLogger.level = getLevel(level), getLevel(level)
	// gLogger.flag, eLogger.flag = flag, flag
	for _, op := range options {
		op(std)
	}
	// gLogger.baseLogger = log.New(gLogger.out, "", gLogger.flag)
	// eLogger.baseLogger = log.New(eLogger.out, "", eLogger.flag)
}

// WithShowFuncName 设置日志是否显示函数名
func WithShowFuncName() Options {
	return func(g *Logger) {
		g.showFuncName = true
	}
}

// WithPath 设置日志输出路径
func WithPath(logDir string) Options {
	return func(g *Logger) {
		std.baseLogger.SetOutput(g.combineFileWithStdWriter(filepath.Join(logDir, globalFileName), os.Stdout))
	}
}

// WithWriter 设置日志输出Writer
func WithWriter(w io.Writer) Options {
	return func(g *Logger) {
		g.baseLogger.SetOutput(w)
	}
}

// WithPath 设置日志输出路径
func WithLevel(level string) Options {
	return func(g *Logger) {
		g.SetLevel(getLevel(level))
	}
}

func WithFlag(flag int) Options {
	return func(g *Logger) {
		g.flag = flag
		g.baseLogger.SetFlags(flag)
	}
}

// WithShowFuncName 设置日志是否显示函数名
func (l *Logger) WithShowFuncName() {
	l.showFuncName = true
}

// WithPath 设置日志输出路径
func (l *Logger) WithPath(logDir string) {
	l.baseLogger.SetOutput(l.combineFileWithStdWriter(filepath.Join(logDir, globalFileName), os.Stdout))
}

// WithWriter 设置日志输出Writer
func (l *Logger) WithWriter(w io.Writer) {
	l.baseLogger.SetOutput(w)
}

// WithPath 设置日志输出路径
func (l *Logger) WithLevel(level string) {
	l.SetLevel(getLevel(level))
}

func (l *Logger) WithFlag(flag int) {
	l.flag = flag
	l.baseLogger.SetFlags(flag)
}
