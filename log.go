package log

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"
)

// Logger Logger
type Logger struct {
	level  int32
	opt    *Option
	format Formatter
	writer io.Writer
	extra  [fatalLevel + 1][]io.Writer
}

// Close It's dangerous to call the method on logging
func (logger *Logger) Close() {
	// if logger.baseFile != nil {
	// 	_ = logger.baseFile.Close()
	// }

	// // logger.handle = nil
	// logger.baseFile = nil
}

func (l *Logger) doPrintf(level int32, format string, a ...interface{}) {
	if level < atomic.LoadInt32(&l.level) {
		return
	}
	if l.format == nil {
		panic("logger closed case format is nil")
	}

	s := l.format.Format(l.opt, level, format, a...)
	_, _ = l.writer.Write(s)
	// 额外的日志输出通道
	for _, v := range l.extra[level] {
		_, _ = v.Write(s)
	}

	if level == panicLevel {
		panic(fmt.Sprintf(format, a...))
	}

	if level == fatalLevel {
		os.Exit(1)
	}
}

// SetLevel SetLevel
func (l *Logger) SetLevel(level int32) {
	atomic.StoreInt32(&l.level, level)
}

// Debug Debug
func (l *Logger) Debug(format string, a ...interface{}) {
	l.doPrintf(debugLevel, format, a...)
}

// Info Info
func (l *Logger) Info(format string, a ...interface{}) {
	l.doPrintf(infoLevel, format, a...)
}

// Warn Warn
func (l *Logger) Warn(format string, a ...interface{}) {
	l.doPrintf(warnLevel, format, a...)
}

// Error Error
func (l *Logger) Error(format string, a ...interface{}) {
	l.doPrintf(errorLevel, format, a...)
}

// Panic Panic
func (l *Logger) Panic(format string, a ...interface{}) {
	l.doPrintf(panicLevel, format, a...)
}

// Fatal Fatal
func (l *Logger) Fatal(format string, a ...interface{}) {
	l.doPrintf(fatalLevel, format, a...)
}

// Debug Debug
func Debug(format string, a ...interface{}) {
	std.Debug(format, a...)
}

// Info Info
func Info(format string, a ...interface{}) {
	std.Info(format, a...)
}

// Warn Warn
func Warn(format string, a ...interface{}) {
	std.Warn(format, a...)
}

// Error Error
func Error(format string, a ...interface{}) {
	std.Error(format, a...)
}

// Error Error
func Panic(format string, a ...interface{}) {
	std.Panic(format, a...)
}

// Fatal Fatal
func Fatal(format string, a ...interface{}) {
	std.Fatal(format, a...)
}
