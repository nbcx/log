package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync/atomic"
)

// Logger Logger
type Log struct {
	level  int32
	opt    *Option
	format Formatter
	writer io.Writer
	extra  [fatalLevel + 1][]io.Writer
}

// Close It's dangerous to call the method on logging
func (logger *Log) Close() {
	// if logger.baseFile != nil {
	// 	_ = logger.baseFile.Close()
	// }

	// // logger.handle = nil
	// logger.baseFile = nil
}

func formatPattern(f interface{}, v []interface{}) string {
	var msg string
	switch f := f.(type) {
	case string:
		msg = f
		if len(v) == 0 {
			return msg
		}
		if !strings.Contains(msg, "%") {
			// do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return msg
}

func (l *Log) print(level int32, format any, a ...interface{}) {
	msg := formatPattern(format, a)
	if level < atomic.LoadInt32(&l.level) {
		return
	}
	if l.format == nil {
		panic("logger closed case format is nil")
	}

	s := l.format.Format(l.opt, level, msg, a...)
	_, _ = l.writer.Write(s)
	// 额外的日志输出通道
	for _, v := range l.extra[level] {
		_, _ = v.Write(s)
	}

	if level == panicLevel {
		panic(fmt.Sprintf(msg, a...))
	}

	if level == fatalLevel {
		os.Exit(1)
	}
}

// SetLevel SetLevel
func (l *Log) SetLevel(level int32) {
	atomic.StoreInt32(&l.level, level)
}

// Debug Debug
func (l *Log) Debug(format any, a ...interface{}) {
	l.print(debugLevel, format, a...)
}

// Info Info
func (l *Log) Info(format any, a ...interface{}) {
	l.print(infoLevel, format, a...)
}

// Warn Warn
func (l *Log) Warn(format any, a ...interface{}) {
	l.print(warnLevel, format, a...)
}

// Error Error
func (l *Log) Error(format any, a ...interface{}) {
	l.print(errorLevel, format, a...)
}

// Panic Panic
func (l *Log) Panic(format any, a ...interface{}) {
	l.print(panicLevel, format, a...)
}

// Fatal Fatal
func (l *Log) Fatal(format any, a ...interface{}) {
	l.print(fatalLevel, format, a...)
}

// WithPath 设置日志输出路径
func (l *Log) WithLevel(level string) {
	l.SetLevel(getLevel(level))
}

// 设置日志格式器
func (l *Log) SetFormatter(level string) {
	l.SetLevel(getLevel(level))
}

// 设置日志输出
func (l *Log) SetWriter(w io.Writer) {
	l.writer = w
}

// 为指定等级的日志设置额外的输出
// 通常用于需要特别关注的紧急日志
func (l *Log) SetLevelWriter(level string, w ...io.Writer) {
	lv := getLevel(level)
	l.extra[lv] = append(l.extra[lv], w...)
}

// 设置日志显示格式
// 需要注意，对于一些自定义的formatter，它并不是绝对生效的
func (l *Log) SetFlag(flag int) {
	l.opt.SetFlags(flag)
}

// 通常我们需要把日志传递给第三方模块使用，并想要标记是第三方模块
// 那么可以使用此函数获取一个带有指定标记的日志实例
func (l *Log) SetPrefix(flag int) *Log {
	l.opt.SetFlags(flag)

	// todo: wait do
	return l
}
