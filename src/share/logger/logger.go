package logger

import (
    "os"
    "share/directory"
    "github.com/op/go-logging"
)

var logPath = directory.Root() + "/log/"
var log     = logging.MustGetLogger("example")
var format  = logging.MustStringFormatter(
    `%{time:2006-01-02 15:04:05.000} [%{level}] %{message}`,
)

// Initializes logging system, creates log file for writing logs
// and returns Logger instance
func Init(name string) *logging.Logger {
    var backend2 = logging.NewLogBackend(os.Stderr, "", 0)
    logging.SetBackend(logging.NewBackendFormatter(backend2, format))

    log.Infof("Opening %s file...", logPath + name + ".log")

    // create directory, if doesn't exist
    if _, err := os.Stat(logPath); os.IsNotExist(err) {
        os.Mkdir(logPath, 0775)
    }

    // open log file for writing
    var file, err = os.OpenFile(
        logPath + name + ".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
    if err != nil {
        log.Error(err)
    } else {
        var backend1 = logging.NewLogBackend(file, "", 0)
        logging.SetBackend(
            logging.NewBackendFormatter(backend1, format),
            logging.NewBackendFormatter(backend2, format),
        )
    }

    return log
}

// Returns Logger instance
func Instance() *logging.Logger {
    return log
}