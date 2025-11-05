package pipeline

import (
	"github.com/sirupsen/logrus"
)

// UpdateHandler defines the function signature for handling updates
type UpdateHandler func(data interface{}) error

// HandleUpdate safely calls the update handler if it's not nil
func HandleUpdate(handler UpdateHandler, data interface{}) error {
	if handler == nil {
		logrus.Warn("UpdateHandler is nil; skipping update handling to prevent crash")
		return nil
	}

	return handler(data)
}
