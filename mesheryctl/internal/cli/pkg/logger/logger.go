package mesheryctllogger

import (
	"os"

	"github.com/meshery/meshkit/logger"
	logrus "github.com/sirupsen/logrus"
)

var Log logger.Handler

func NewLogger(level logrus.Level) *logger.Handler {
	if Log != nil {
		return &Log
	}

	options := logger.Options{
		Format:      logger.TerminalLogFormat,
		LogLevel:    int(level),
		Output:      os.Stdout,
		ErrorOutput: os.Stderr,
	}

	log, err := logger.New("mesheryctl", options)
	if err != nil {
		logrus.Errorf("Failed to initialize logger: %v", err)
		os.Exit(1)
	}

	Log = log
	return &Log
}
