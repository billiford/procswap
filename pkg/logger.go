package procswap

import (
	"fmt"
	"os"
	"time"

	"github.com/shiena/ansicolor"
)

type logLevel string

const (
	// green   = "\033[97;42m"
	white  = "\033[90;47m"
	yellow = "\033[90;43m"
	red    = "\033[97;41m"
	// blue    = "\033[97;44m"
	// magenta = "\033[97;45m"
	cyan  = "\033[97;46m"
	reset = "\033[0m"

	logLevelDebug = logLevel("DEBUG")
	logLevelInfo  = logLevel("INFO")
	logLevelWarn  = logLevel("WARN")
	logLevelError = logLevel("ERROR")
	logLevelFatal = logLevel("FATAL")
)

// logWithLevel logs a given message in a nice format.
func logWithLevel(level logLevel, message string, newline ...bool) {
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

func logDebug(message string, newline ...bool) {
	logWithLevel(logLevelDebug, message, newline...)
}

func logInfo(message string, newline ...bool) {
	logWithLevel(logLevelInfo, message, newline...)
}

func logWarn(message string, newline ...bool) {
	logWithLevel(logLevelWarn, message, newline...)
}

func logError(message string, newline ...bool) {
	logWithLevel(logLevelError, message, newline...)
}

// Log and exit.
func logFatal(message string, newline ...bool) {
	logWithLevel(logLevelFatal, message, newline...)
	os.Exit(1)
}
