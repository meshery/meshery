package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/layer5io/meshery/server/models"
	"github.com/sirupsen/logrus"
)

type LogRequestBody struct {
	Logs      []Log  `json:"logs,omitempty"`
	SessionId string `json:"sessionId,omitempty"`
}

type Log struct {
	Timestamp int64  `json:"timestamp,omitempty"` // Timestamp is the Javascript time format, that tells the number of milliseconds elapsed since midnight, January 1, 1970 Universal Coordinated Time (UTC)
	Msg       string `json:"msg,omitempty"`
	Severity  string `json:"severity,omitempty"`
	CodeBlock string `json:"codeBlock,omitempty"`
	Area      string `json:"area,omitempty"`
}

func (h *Handler) LogHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	parsedBody := LogRequestBody{}
	if err := json.NewDecoder(req.Body).Decode(&parsedBody); err != nil {
		http.Error(w, ErrRetrieveData(err).Error(), http.StatusBadRequest)
		return
	}

	err := writeLogToFile(user, parsedBody.SessionId, parsedBody.Logs)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	w.Write([]byte("Logs written on filesystem"))
	w.WriteHeader(200)
}

func writeLogToFile(user *models.User, sessionId string, logs []Log) error {
	f, err := createLogFile(sessionId)
	if err != nil {
		return err
	}
	defer f.Close()

	var stringBytes []byte

	// check for file info and if the file is empty add the log-headers
	fileInfo, _ := f.Stat()
	if fileInfo.Size() == 0 {
		stringBytes = getLogFileHeaders(user)
	}

	// go through logs to write to file as a JSON entry where each line is a log on the file
	for _, log := range logs {
		bytes, err := json.Marshal(log)
		if err != nil {
			logrus.Error(fmt.Sprintf("failed to marshal error log %s", err.Error()))
			continue
		}
		stringBytes = append(stringBytes, bytes...)
		stringBytes = append(stringBytes, []byte("\n")...)
	}

	if _, err = f.Write(stringBytes); err != nil {
		return err
	}

	return nil
}

// createLogFile returns reference to the opened file as a result of request to create or append
func createLogFile(sessionId string) (*os.File, error) {
	filePath, err := getLogFilePath()
	// create set of folders if doesn't exist
	errMkdir := os.MkdirAll(filePath, os.ModePerm)
	if err != nil || errMkdir != nil {
		logrus.Error("Error doing File operation on operating system", err)
		return nil, err
	}

	f, err := os.OpenFile(filepath.Join(filePath, sessionId), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		logrus.Errorf("Failed to open/create/append to logFile with session Id %s [%s]", sessionId, err.Error())
		return nil, err
	}

	return f, nil
}

// getLogFilePath gets the log directory where the log file has to be placed
func getLogFilePath() (string, error) {
	homDir, e := os.UserHomeDir()
	if e != nil {
		return "", e
	}
	logDir := "/.meshery/logs/"
	return filepath.Join(homDir, logDir), nil
}

func getLogFileHeaders(user *models.User) []byte {
	return []byte(fmt.Sprintf(`// Meshery Log File, please share this file to raise issues or problems you report in Meshery
// System Information: 
// Operating System: %s
// CPUs: %d
// Os Architecture: %s
// 
// Time Created: %s
`, runtime.GOOS, runtime.NumCPU(), runtime.GOARCH, time.Now().String()))
}
