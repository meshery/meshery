package utils

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// OnlyStringFormatterForLogrus is a custom logrus formatter that returns only string of message as output
type OnlyStringFormatterForLogrus struct{}

// Format defined the format of output for Logrus logs
func (obj OnlyStringFormatterForLogrus) Format(entry *logrus.Entry) ([]byte, error) {
	if entry != nil {
		return []byte(entry.Message + "\n"), nil
	}
	return []byte(""), errors.New("no entry received")
}
