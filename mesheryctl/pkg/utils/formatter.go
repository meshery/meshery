package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/meshery/meshkit/logger"
	"github.com/spf13/viper"
)

// SetupMeshkitLogger initializes and returns a MeshKit Logger instance
func SetupMeshkitLogger(name string, debugLevel bool, output io.Writer) logger.Handler {
	logLevel := viper.GetInt("LOG_LEVEL")

	// Logrus Level constants (Hardcoded to avoid direct import)
	// 4 = InfoLevel
	// 5 = DebugLevel
	const (
		infoLevel  = 4
		debugConst = 5
	)

	// If the debug flag is explicitly set, force the log level to Debug (5)
	if debugLevel {
		logLevel = debugConst
	} else {
		// If debug is off, we check if Viper had a value.
		// If Viper returned 0 (not set), default to Info (4).
		if logLevel == 0 {
			logLevel = infoLevel
		}
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
