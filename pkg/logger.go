package procswap

import (
	"fmt"
	"os"
	"time"

	"github.com/shiena/ansicolor"
)

type LogLevel string

const (
	// green   = "\033[97;42m"
	white  = "\033[90;47m"
	yellow = "\033[90;43m"
	red    = "\033[97;41m"
	// blue    = "\033[97;44m"
	// magenta = "\033[97;45m"
	cyan  = "\033[97;46m"
	reset = "\033[0m"

	logLevelDebug = LogLevel("DEBUG")
	logLevelInfo  = LogLevel("INFO")
	logLevelWarn  = LogLevel("WARN")
	logLevelError = LogLevel("ERROR")
	logLevelFatal = LogLevel("FATAL")
)

func LogDebug(message string, newline ...bool) {
	l(logLevelDebug, message, newline...)
}

func LogInfo(message string, newline ...bool) {
	l(logLevelInfo, message, newline...)
}

func LogWarn(message string, newline ...bool) {
	l(logLevelWarn, message, newline...)
}

func LogError(message string, newline ...bool) {
	l(logLevelError, message, newline...)
}

func LogFatal(message string, newline ...bool) {
	l(logLevelFatal, message, newline...)
}

// l logs a given message in a nice format.
//
// A single letter function name? I can't just name it "log"
// since that overwrites the package imported that we need.
//
// Certainly there's a convention for this... I'll get back
// to it later.
func l(level LogLevel, message string, newline ...bool) {
	nl := true

	// Allow us to define if a newline is appended or not.
	if len(newline) > 0 {
		nl = newline[0]
	}

	var logColor string

	switch level {
	case logLevelDebug:
		logColor = white
	case logLevelInfo:
		logColor = cyan
	case logLevelWarn:
		logColor = yellow
	case logLevelError:
		logColor = red
	case logLevelFatal:
		logColor = red
	}

	messageFormat := "[PROCSWAP] %v |%s %-5s %s| %s"
	if nl {
		messageFormat += "\n"
	}

	w := ansicolor.NewAnsiColorWriter(os.Stdout)

	// the log package always adds a newline even if one is not present, so
	// just use fmt for this.
	fmt.Fprintf(w, messageFormat,
		time.Now().Format("2006/01/02 - 15:04:05"),
		logColor, level, reset,
		message,
	)
}
