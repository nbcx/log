package log

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// brush is a color join function
type brush func(string) string

// newBrush returns a fix color Brush
func newBrush(color string) brush {
	pre := "\033["
	reset := "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = []brush{
	newBrush("1;37"), // Emergency          white

	newBrush("1;34"), // Informational      blue
	newBrush("1;33"), // Warning            yellow
	newBrush("1;31"), // Error              red
	newBrush("1;35"), // Critical           magenta
	newBrush("1;36"), // Alert              cyan

	newBrush("1;44"), // Debug              Background blue
	newBrush("1;32"), // Notice             green
}

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})

	w32Green   = string([]byte{27, 91, 52, 50, 109})
	w32White   = string([]byte{27, 91, 52, 55, 109})
	w32Yellow  = string([]byte{27, 91, 52, 51, 109})
	w32Red     = string([]byte{27, 91, 52, 49, 109})
	w32Blue    = string([]byte{27, 91, 52, 52, 109})
	w32Magenta = string([]byte{27, 91, 52, 53, 109})
	w32Cyan    = string([]byte{27, 91, 52, 54, 109})

	reset = string([]byte{27, 91, 48, 109})

	levelConsolePrefix = [fatalLevel + 1]string{"[D]", "[I]", "[W]", "[E]", "[P]", "[F]"}
)

var (
	once     sync.Once
	colorMap map[string]string
)

func initColor() {
	if runtime.GOOS == "windows" {
		green = w32Green
		white = w32White
		yellow = w32Yellow
		red = w32Red
		blue = w32Blue
		magenta = w32Magenta
		cyan = w32Cyan
	}
	colorMap = map[string]string{
		// by color
		"green":  green,
		"white":  white,
		"yellow": yellow,
		"red":    red,
		// by method
		"GET":     blue,
		"POST":    cyan,
		"PUT":     yellow,
		"DELETE":  red,
		"PATCH":   green,
		"HEAD":    magenta,
		"OPTIONS": white,
	}
}

// ColorByStatus return color by http code
// 2xx return Green
// 3xx return White
// 4xx return Yellow
// 5xx return Red
func ColorByStatus(code int) string {
	once.Do(initColor)
	switch {
	case code >= 200 && code < 300:
		return colorMap["green"]
	case code >= 300 && code < 400:
		return colorMap["white"]
	case code >= 400 && code < 500:
		return colorMap["yellow"]
	default:
		return colorMap["red"]
	}
}

// ColorByMethod return color by http code
func ColorByMethod(method string) string {
	once.Do(initColor)
	if c := colorMap[method]; c != "" {
		return c
	}
	return reset
}

// ResetColor return reset color
func ResetColor() string {
	return reset
}

// consoleWriter implements LoggerInterface and writes messages to terminal.
type consoleWriter struct{}

// NewConsole creates ConsoleWriter returning as LoggerInterface.
func NewConsole() *consoleWriter {
	return &consoleWriter{}
}

func (c *consoleWriter) Format(lm *Option, level int32, format string, args ...interface{}) []byte {
	buf := getBuffer()
	defer putBuffer(buf)

	// time
	lm.Time(buf, time.Now())
	// level
	*buf = append(*buf, colors[level](levelConsolePrefix[level])+" "...)
	// msg
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	*buf = append(*buf, msg...)

	if len(*buf) == 0 || (*buf)[len(*buf)-1] != '\n' {
		*buf = append(*buf, '\n')
	}
	return *buf
}
