package utils

import (
	"io"
	"os"

	"github.com/layer5io/meshkit/logger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// TerminalFormatter is exported
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
func SetupMeshkitLogger(debugLevel bool, output io.Writer) {
	logLevel := viper.GetInt("LOG_LEVEL")
	if !debugLevel {
		logLevel = int(log.DebugLevel)
	}

	logger, err := logger.New("mesheryctl", logger.Options{
		Format:   logger.TerminalLogFormat,
		LogLevel: logLevel,
		Output:   output,
	})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	Log = logger
}
