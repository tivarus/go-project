package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

func New() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger
}
