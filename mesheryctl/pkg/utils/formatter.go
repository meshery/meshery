package utils

import (
	"os"

	"github.com/layer5io/meshkit/logger"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

//TerminalFormatter is exported
type TerminalFormatter struct{}

// Format defined the format of output for Logrus logs
// Format is exported
func (f *TerminalFormatter) Format(entry *log.Entry) ([]byte, error) {
	return append([]byte(entry.Message), '\n'), nil
}

// Call this function to setup logrus
func SetupLogrusFormatter() {
	//log formatter for improved UX
	log.SetFormatter(new(TerminalFormatter))
}

// Initialize Meshkit Logger instance
func MeshkitLogger() (logger.Handler, error) {
	log, err := logger.New("mesheryctl", logger.Options{
		Format: logger.SyslogLogFormat,
	})
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	return log, nil
}
