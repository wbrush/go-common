package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var errorLog = logrus.New()
var infoLog = logrus.New()
var debugLog = logrus.New()

func init() {
	errorLog.Out = os.Stderr
	infoLog.Out = os.Stdout
	debugLog.Out = os.Stdout
}

func Debug() *logrus.Logger {
	return debugLog
}

func Info() *logrus.Logger {
	return infoLog
}

func Error() *logrus.Logger {
	return errorLog
}
