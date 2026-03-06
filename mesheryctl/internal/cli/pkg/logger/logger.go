package mesheryctllogger

import (
	"os"
	"sync"

	"github.com/meshery/meshkit/logger"
	logrus "github.com/sirupsen/logrus"
)

var lock = &sync.Mutex{}

var Log logger.Handler

func GetLogger(level logrus.Level) *logger.Handler {
	if Log == nil {
		lock.Lock()
		defer lock.Unlock()
		if Log == nil {

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
		} else {
			Log.SetLevel(level)
		}
	}

	return &Log
}
