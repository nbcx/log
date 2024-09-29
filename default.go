package log

import (
	"io"
	"sync/atomic"
	"time"
)

// A Logger represents an active logging object that generates lines of
// output to an [io.Writer]. Each logging operation makes a single call to
// the Writer's Write method. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type DefaultFormat struct {
	// outMu sync.Mutex
	// out   io.Writer // destination for output

	prefix atomic.Pointer[string] // prefix on each line to identify the logger (but see Lmsgprefix)
	// flag   atomic.Int32           // properties
	// isDiscard atomic.Bool
}

// New creates a new [Logger]. The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line, or
// after the log header if the [Lmsgprefix] flag is provided.
// The flag argument defines the logging properties.
func New(out io.Writer, prefix string, flag int) *DefaultFormat {
	l := new(DefaultFormat)
	return l
}

// Default returns the standard logger used by the package-level output functions.
func Default() *Log { return std }

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// Format implements log.Formatter.
func (l *DefaultFormat) Format(lm *Option, level int32, s string, args ...interface{}) []byte {

	now := time.Now() // get this early.

	buf := getBuffer()
	defer putBuffer(buf)
	// prefix = lm.PrintLevel(level)
	// fmt.Println("formatHeader", "now:", now, "prefix:", "flag:", flag)
	l.formatHeader(lm, buf, now, lm.PrintLevel(level), lm.Flag)
	*buf = append(*buf, s...)
	// *buf = appendOutput(*buf)
	if len(*buf) == 0 || (*buf)[len(*buf)-1] != '\n' {
		*buf = append(*buf, '\n')
	}

	return *buf
}

// formatHeader writes log header to buf in following order:
//   - l.prefix (if it's not blank and Lmsgprefix is unset),
//   - date and/or time (if corresponding flags are provided),
//   - file and line number (if corresponding flags are provided),
//   - l.prefix (if it's not blank and Lmsgprefix is set).
func (l *DefaultFormat) formatHeader(lm *Option, buf *[]byte, t time.Time, level string, flag int) {
	if flag&Lmsgprefix == 0 {
		*buf = append(*buf, l.Prefix()...)
	}
	lm.Time(buf, t)
	*buf = append(*buf, level...)
	lm.File(buf)
	if flag&Lmsgprefix != 0 {
		*buf = append(*buf, l.Prefix()...)
	}

	*buf = append(*buf, " "...)
}

// Prefix returns the output prefix for the logger.
func (l *DefaultFormat) Prefix() string {
	if p := l.prefix.Load(); p != nil {
		return *p
	}
	return ""
}
