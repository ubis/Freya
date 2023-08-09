package log

import (
	"os"
	"strings"

	"github.com/op/go-logging"
	"github.com/ubis/Freya/share/directory"
)

var log = logging.MustGetLogger("example")

// Init logging system which will create a log file or append one
func Init(name string) {
	b2 := logging.NewLogBackend(os.Stderr, "", 0)
	path := directory.Root() + "/log/"
	fname := path + strings.ToLower(name) + ".log"
	format := logging.MustStringFormatter(
		`%{time:2006-01-02 15:04:05.000} [%{level}] %{message}`)
	logging.SetBackend(logging.NewBackendFormatter(b2, format))

	log.Infof("Opening %s file...", fname)

	// create directory, if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0775)
	}

	// open log file for writing
	var f, err = os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Error(err)
		return
	}

	var b1 = logging.NewLogBackend(f, "", 0)
	logging.SetBackend(logging.NewBackendFormatter(b1, format),
		logging.NewBackendFormatter(b2, format))
	log.Info(name + " init")
}

// Instance returns a Logger instance
func Instance() *logging.Logger {
	return log
}

// Fatal is similar to Critical followed by a call to os.Exit(1)
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Fatalf is similar to Criticalf followed by a call to os.Exit(1)
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Panic is similar to Critical followed by a call to panic()
func Panic(args ...interface{}) {
	log.Panic(args...)
}

// Panicf is similar to Criticalf followed by a call to panic()
func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

// Critical logs a message using CRITICAL as log level
func Critical(args ...interface{}) {
	log.Critical(args...)
}

// Criticalf logs a message using CRITICAL as log level
func Criticalf(format string, args ...interface{}) {
	log.Criticalf(format, args...)
}

// Error logs a message using ERROR as log level
func Error(args ...interface{}) {
	log.Error(args...)
}

// Errorf logs a message using ERROR as log level
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Warning logs a message using WARNING as log level
func Warning(args ...interface{}) {
	log.Warning(args...)
}

// Warningf logs a message using WARNING as log level
func Warningf(format string, args ...interface{}) {
	log.Warningf(format, args...)
}

// Notice logs a message using NOTICE as log level
func Notice(args ...interface{}) {
	log.Notice(args...)
}

// Noticef logs a message using NOTICE as log level
func Noticef(format string, args ...interface{}) {
	log.Noticef(format, args...)
}

// Info logs a message using INFO as log level
func Info(args ...interface{}) {
	log.Info(args...)
}

// Infof logs a message using INFO as log level
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Debug logs a message using DEBUG as log level
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Debugf logs a message using DEBUG as log level
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}
