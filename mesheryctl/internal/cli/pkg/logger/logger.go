package mesheryctllogger

import (
	"io"
	"os"
	"sync"

	"github.com/meshery/meshkit/logger"
	logrus "github.com/sirupsen/logrus"
)

// MesheryctlLogger is a wrapper around the meshkit logger.Handler to provide a consistent logging interface into mesheryctl.
type MesheryctlLogger struct {
	Log logger.Handler
}

var (
	Log  *MesheryctlLogger
	once sync.Once
)

func GetMeshkitLogger(level logrus.Level) *MesheryctlLogger {

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

		Log = &MesheryctlLogger{Log: log}
	})

	return Log
}

func (ml *MesheryctlLogger) SetLevel(level logrus.Level) {
	ml.Log.SetLevel(level)
}

func (ml *MesheryctlLogger) GetLevel() logrus.Level {
	return logrus.Level(ml.Log.GetLevel())
}

func (ml *MesheryctlLogger) Fatal(err error) {
	ml.Log.Fatal(err)
}

func (ml *MesheryctlLogger) Fatalf(format string, args ...interface{}) {
	ml.Log.Fatalf(format, args...)
}

func (ml *MesheryctlLogger) Error(err error) {
	ml.Log.Error(err)
}

func (ml *MesheryctlLogger) Errorf(format string, args ...interface{}) {
	ml.Log.Errorf(format, args...)
}

func (ml *MesheryctlLogger) Info(description ...interface{}) {
	ml.Log.Info(description...)
}

func (ml *MesheryctlLogger) Infof(format string, args ...interface{}) {
	ml.Log.Infof(format, args...)
}

func (ml *MesheryctlLogger) Debug(description ...interface{}) {
	ml.Log.Debug(description...)
}

func (ml *MesheryctlLogger) Debugf(format string, args ...interface{}) {
	ml.Log.Debugf(format, args...)
}

func (ml *MesheryctlLogger) Warn(err error) {
	ml.Log.Warn(err)
}

func (ml *MesheryctlLogger) Warnf(format string, args ...interface{}) {
	ml.Log.Warnf(format, args...)
}

func (ml *MesheryctlLogger) UpdateLogOutput(w io.Writer) {
	ml.Log.UpdateLogOutput(w)
}

func (ml *MesheryctlLogger) ControllerLogger() logger.Handler {
	return ml.Log
}

func (ml *MesheryctlLogger) DatabaseLogger() logger.Handler {
	return ml.Log
}
