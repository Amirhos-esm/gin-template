package mLogger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// ANSI escape codes for colors
const (
    ColorReset   = "\033[0m"
    ColorVerbose = "\033[36m"  // Cyan
    ColorDebug   = "\033[32m"  // Green
    ColorInfo    = "\033[34m"  // Blue
    ColorWarning = "\033[33m"  // Yellow
    ColorError   = "\033[31m"  // Red
    ColorCritical = "\033[35m" // Magenta
)

// LogLevel enum
type LogLevel int

const (
    VERBOSE LogLevel = iota
    DEBUG
    INFO
    WARNING
    ERROR
    CRITICAL
)

type MLogger struct {
	logger   *log.Logger
	enable   bool
	colorful bool
	tag      string
	level LogLevel
}

func New(tag string, enable bool,level LogLevel) *MLogger {
	t := &MLogger{}
	t.logger = log.New(os.Stdout, "43443",  0)
	t.enable = enable
	t.colorful = true
	t.tag = tag
	t.level = level
	return t
}

func createHeader(t *MLogger, logType string, color string) string {
    // Get the current timestamp in the desired format
    timestamp := time.Now().Format("2006-01-02 15:04:05")

    if t.colorful {
        // Return colorful log header with timestamp, tag, and logType
        return fmt.Sprintf("%s[%s] [%s] [%s]%s ",color, timestamp, t.tag, logType, ColorReset)
    } else {
        // Return log header without colors
        return fmt.Sprintf("[%s] [%s] [%s] ", timestamp, t.tag, logType)
    }
}
func (t *MLogger) V(format string, a ...interface{}) {
	if !t.enable || t.level > VERBOSE {
		return
	}
	header := createHeader(t,"VERBOSE",ColorVerbose)
	fmt.Println(header + fmt.Sprintf(format, a...))
}
func (t *MLogger) D(format string, a ...interface{}) {
	if !t.enable || t.level > DEBUG {
		return
	}
	header := createHeader(t,"DEBUG",ColorDebug)
	fmt.Println(header + fmt.Sprintf(format, a...))
}
func (t *MLogger) I(format string, a ...interface{}) {
	if !t.enable || t.level > INFO {
		return
	}
	header := createHeader(t,"INFO",ColorInfo)
	fmt.Println(header + fmt.Sprintf(format, a...))
}
func (t *MLogger) W(format string, a ...interface{}) {
	if !t.enable || t.level > WARNING{
		return
	}
	header := createHeader(t,"WARNING",ColorWarning)
	fmt.Println(header + fmt.Sprintf(format, a...))
}
func (t *MLogger) E(format string, a ...interface{}) {
	if !t.enable || t.level > ERROR{
		return
	}
	header := createHeader(t,"ERROR",ColorError)
	fmt.Println(header + fmt.Sprintf(format, a...))
}
func (t *MLogger) C(format string, a ...interface{}) {
	if !t.enable || t.level > CRITICAL{
		return
	}
	header := createHeader(t,"CRITICAL",ColorCritical)
	fmt.Println(header + fmt.Sprintf(format, a...))
}

func (t *MLogger) Printf(format string, a ...interface{}) {
	if !t.enable || t.level > DEBUG {
		return
	}
	header := createHeader(t,"DEBUG",ColorInfo)
	fmt.Println(header + fmt.Sprintf(format, a...))
}
func (t *MLogger) Print(format string, a ...interface{}) {
	if !t.enable || t.level > DEBUG {
		return
	}
	header := createHeader(t,"DEBUG",ColorInfo)
	fmt.Println(header + fmt.Sprintf(format, a...))
}
func (t *MLogger) Println(format string, a ...interface{}) {
	if !t.enable || t.level > DEBUG {
		return
	}
	header := createHeader(t,"DEBUG",ColorInfo)
	fmt.Println(header + fmt.Sprintf(format, a...))
}