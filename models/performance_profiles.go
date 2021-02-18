package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// PerformanceProfile represents the performance profile that needs
// to be saved
type PerformanceProfile struct {
	ID *uuid.UUID `json:"id,omitempty"`

	Name           string     `json:"name,omitempty"`
	LastRun        *time.Time `json:"last_run,omitempty"`
	Schedule       *uuid.UUID `json:"schedule,omitempty"`
	LoadGenerators []string   `json:"load_generators,omitempty"`
	Endpoints      []string   `json:"endpoints,omitempty"`
}
