package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/meshery/meshkit/logger"
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

// SetupLogrusFormatter sets up the logrus formatter for improved UX
// This is kept for backwards compatibility
func SetupLogrusFormatter() {
	log.SetFormatter(new(TerminalFormatter))
}

// SetupMeshkitLogger initializes and returns a MeshKit Logger instance
func SetupMeshkitLogger(name string, debugLevel bool, output io.Writer) logger.Handler {
	logLevel := viper.GetInt("LOG_LEVEL")
	if !debugLevel {
		logLevel = int(log.DebugLevel)
	}
	meshkitLogger, err := logger.New(name, logger.Options{
		Format:   logger.TerminalLogFormat,
		LogLevel: logLevel,
		Output:   output,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logger: %v\n", err)
		os.Exit(1)
	}
	return meshkitLogger
}
