package mesheryctllogger

import (
	"os"
	"sync"

	"github.com/meshery/meshkit/logger"
	logrus "github.com/sirupsen/logrus"
)

var (
	Log  logger.Handler
	once sync.Once
)

func GetMeshkitLogger(level logrus.Level) logger.Handler {

	once.Do(func() {
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
	})

	return Log
}
