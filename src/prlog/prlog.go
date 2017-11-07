package prlog

import (
	"os"

	"github.com/Sirupsen/logrus"
)

//logger is the logger used in the whole application
var logger = logrus.New()

//New initializes the logger used within the whole application
func New() {
	logger.Formatter = &logrus.JSONFormatter{}

	logger.Out = os.Stdout

	file, err := os.OpenFile("poprep.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		logger.Out = file
	} else {
		logger.Info("Failed to log to file, using default stderr")
	}
}

//Get the logger instance
func Get() *logrus.Logger {
	return logger
}
