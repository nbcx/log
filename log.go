package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/nbcx/log/internal"
)

var std *Logger

// levels
const (
	debugLevel = 0
	infoLevel  = 1
	warnLevel  = 2
	errorLevel = 3
	panicLevel = 4
	fatalLevel = 5
)

// log file name
const (
	globalFileName = "service.log"
	errFileName    = "service.err.log"
)

const (
	printDebugLevel = "[debug]"
	printInfoLevel  = "[info ]"
	printWarnLevel  = "[warn ]"
	printErrorLevel = "[error]"
	printPanicLevel = "[panic]"
	printFatalLevel = "[fatal]"
)

// Logger Logger
type Logger struct {
	level        int32
	flag         int // properties
	baseLogger   *internal.Logger
	baseFile     *os.File
	showFuncName bool
}

func init() {
	std = new(Logger)
	// gLogger.out = os.Stdout
	std.level = debugLevel
	std.flag = internal.LstdFlags | internal.Lmicroseconds | internal.Lshortfile
	std.baseLogger = internal.Default() // internal.New(os.Stdout, "", std.flag) // todo: 这里默认可以使用internal的std实例

	// eLogger = new(Logger)
	// // eLogger.out = os.Stderr
	// eLogger.level = warnLevel
	// eLogger.flag = gLogger.flag
	// eLogger.baseLogger = internal.New(os.Stderr, "", eLogger.flag)
}

// GetLogger StdLog and ErrLog
func GetLogger() *Logger {
	return std
}

// GetOutput Stdout and Stderr
func GetOutput() io.Writer {
	return std.baseLogger.Writer()
}

func getLevel(level string) int32 {
	switch strings.ToLower(level) {
	case "debug":
		return debugLevel
	case "info":
		return infoLevel
	case "warn":
		return warnLevel
	case "error":
		return errorLevel
	case "fatal":
		return fatalLevel
	default:
		return debugLevel
	}
}

func (logger *Logger) combineFileWithStdWriter(path string, stdWriter io.Writer) io.Writer {
	_ = CreateDirIfNotExists(filepath.Dir(path))
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	logger.baseFile = file

	return io.MultiWriter(stdWriter, logger.baseFile)
}

// Close It's dangerous to call the method on logging
func (logger *Logger) Close() {
	if logger.baseFile != nil {
		_ = logger.baseFile.Close()
	}

	logger.baseLogger = nil
	logger.baseFile = nil
}

func (logger *Logger) doPrintf(level int32, printLevel string, format string, a ...interface{}) {
	if level < atomic.LoadInt32(&logger.level) {
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	if logger.showFuncName {
		format = fmt.Sprintf("%s[%s] %s", printLevel, runFuncName(5), format)
	} else {
		format = fmt.Sprintf("%s %s", printLevel, format)
	}

	if level == panicLevel {
		logger.baseLogger.Panicf(format, a...)
	} else {
		_ = logger.baseLogger.Output(4, fmt.Sprintf(format, a...))
	}
	if level == fatalLevel {
		os.Exit(1)
	}
}

// SetLevel SetLevel
func (logger *Logger) SetLevel(level int32) {
	atomic.StoreInt32(&logger.level, level)
}

// Debug Debug
func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

// Info Info
func (logger *Logger) Info(format string, a ...interface{}) {
	logger.doPrintf(infoLevel, printInfoLevel, format, a...)
}

// Warn Warn
func (logger *Logger) Warn(format string, a ...interface{}) {
	logger.doPrintf(warnLevel, printWarnLevel, format, a...)
}

// Error Error
func (logger *Logger) Error(format string, a ...interface{}) {
	logger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

// Panic Panic
func (logger *Logger) Panic(format string, a ...interface{}) {
	logger.doPrintf(panicLevel, printPanicLevel, format, a...)
}

// Fatal Fatal
func (logger *Logger) Fatal(format string, a ...interface{}) {
	logger.doPrintf(fatalLevel, printFatalLevel, format, a...)
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

// ReloadLogger 重新加载日志级别配置
func ReloadLogger(level string) {
	std.SetLevel(getLevel(level))
}

// CloseLogger 关闭日志
func CloseLogger() {
	std.Close()
}

// 获取正在运行的函数名
func runFuncName(depth int) string {
	pc := make([]uintptr, 1)
	runtime.Callers(depth, pc)
	f := runtime.FuncForPC(pc[0])
	// 分割字符串 并返回最后一个元素
	arr := strings.Split(f.Name(), ".")
	if len(arr) < 1 {
		return f.Name()
	}

	return arr[len(arr)-1]
}

// CreateDirIfNotExists 目录不存在时创建目录
func CreateDirIfNotExists(path string) error {
	if !Exists(path) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

// Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
