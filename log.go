package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
)

var (
	gLogger *Logger
	eLogger *Logger
	isInit  bool
)

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
	out          io.Writer
	baseLogger   *log.Logger
	baseFile     *os.File
	showFuncName bool
}

func init() {
	gLogger = new(Logger)
	gLogger.out = os.Stdout
	gLogger.level = debugLevel
	gLogger.flag = log.LstdFlags | log.Lmicroseconds | log.Lshortfile
	gLogger.baseLogger = log.New(gLogger.out, "", gLogger.flag)

	eLogger = new(Logger)
	eLogger.out = os.Stderr
	eLogger.level = warnLevel
	eLogger.flag = gLogger.flag
	eLogger.baseLogger = log.New(eLogger.out, "", eLogger.flag)
}

// GetLogger StdLog and ErrLog
func GetLogger() (*Logger, *Logger) {
	return gLogger, eLogger
}

// GetOutput Stdout and Stderr
func GetOutput() (io.Writer, io.Writer) {
	return gLogger.out, eLogger.out
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
	gLogger.Debug(format, a...)
}

// Info Info
func Info(format string, a ...interface{}) {
	gLogger.Info(format, a...)
}

// Warn Warn
func Warn(format string, a ...interface{}) {
	gLogger.Warn(format, a...)
}

// Error Error
func Error(format string, a ...interface{}) {
	eLogger.Error(format, a...)
}

// Error Error
func Panic(format string, a ...interface{}) {
	eLogger.Panic(format, a...)
}

// Fatal Fatal
func Fatal(format string, a ...interface{}) {
	eLogger.Fatal(format, a...)
}

// ReloadLogger 重新加载日志级别配置
func ReloadLogger(level string) {
	gLogger.SetLevel(getLevel(level))
}

// SetShowFuncName 设置日志是否显示函数名
func SetShowFuncName(isShow bool) {
	gLogger.showFuncName = isShow
	eLogger.showFuncName = isShow
}

type Options func(g *Logger, e *Logger)

// WithShowFuncName 设置日志是否显示函数名
func WithShowFuncName() Options {
	return func(g *Logger, e *Logger) {
		g.showFuncName = true
		e.showFuncName = true
	}
}

// WithPath 设置日志输出路径
func WithPath(logDir string) Options {
	return func(g *Logger, e *Logger) {
		g.out = g.combineFileWithStdWriter(filepath.Join(logDir, globalFileName), os.Stdout)
		e.out = e.combineFileWithStdWriter(filepath.Join(logDir, errFileName), os.Stderr)
	}
}

// Set log
func Set(level string, flag int, options ...Options) {
	if isInit {
		return
	}
	isInit = true
	gLogger.level, eLogger.level = getLevel(level), getLevel(level)
	gLogger.flag, eLogger.flag = flag, flag
	for _, op := range options {
		op(gLogger, eLogger)
	}
	gLogger.baseLogger = log.New(gLogger.out, "", gLogger.flag)
	eLogger.baseLogger = log.New(eLogger.out, "", eLogger.flag)
}

// CloseLogger 关闭日志
func CloseLogger() {
	gLogger.Close()
	eLogger.Close()
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
