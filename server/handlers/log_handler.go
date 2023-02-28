package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/layer5io/meshery/server/models"
	"github.com/sirupsen/logrus"
)

type LogRequestBody struct {
	Logs      []Log  `json:"logs,omitempty"`
	SessionId string `json:"sessionId,omitempty"`
}

type Log struct {
	Timestamp       string `json:"timestamp,omitempty"`
	Message         string `json:"message,omitempty"`
	Severity        string `json:"severity,omitempty"`
	CodeAttachments string `json:"codeAttachments,omitempty"`
}

func (h *Handler) LogHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	parsedBody := LogRequestBody{}
	if err := json.NewDecoder(req.Body).Decode(&parsedBody); err != nil {
		http.Error(w, ErrRetrieveData(err).Error(), http.StatusBadRequest)
		return
	}

	logrus.Debugf("Inside Log handler %+v", parsedBody)
	err := createLogFile(user, parsedBody.SessionId, parsedBody.Logs)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	w.Write([]byte("Response provided!!"))
	w.WriteHeader(200)
}

func createLogFile(user *models.User, sessionId string, logs []Log) error {
	filePath, err := getLogFilePath()
	// create set of folders if doesn't exist
	errMkdir := os.MkdirAll(filePath, os.ModePerm)
	if err != nil || errMkdir != nil {
		logrus.Error("Error doing File operation on operating system", err)
		return err
	}

	f, err := os.OpenFile(filepath.Join(filePath, sessionId), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		logrus.Errorf("Failed to open/create/append to logFile with session Id %s", sessionId)
		return err
	}

	defer f.Close()

	byte, err := json.Marshal(logs)
	if err != nil {
		logrus.Error("error marshalling json")
		return err
	}

	if _, err = f.Write(byte); err != nil {
		return err
	}

	return nil
}

func getLogFilePath() (string, error) {
	homDir, e := os.UserHomeDir()
	if e != nil {
		return "", e
	}
	logDir := "/.meshery/logs/"
	return filepath.Join(homDir, logDir), nil
}
