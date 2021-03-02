package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// PerformanceProfile represents the performance profile that needs
// to be saved
type PerformanceProfile struct {
	ID *uuid.UUID `json:"id,omitempty"`

	Name              string     `json:"name,omitempty"`
	LastRun           *time.Time `json:"last_run,omitempty"`
	Schedule          *uuid.UUID `json:"schedule,omitempty"`
	LoadGenerators    []string   `json:"load_generators,omitempty"`
	Endpoints         []string   `json:"endpoints,omitempty"`
	ServiceMesh       string     `json:"service_mesh,omitempty"`
	ConcurrentRequest int        `json:"concurrent_request,omitempty"`
	QPS               int        `json:"qps,omitempty"`
	Duration          string     `json:"duration,omitempty"`

	RequestHeaders string `json:"request_headers,omitempty"`
	RequestCookies string `json:"request_cookies,omitempty"`
	RequestBody    string `json:"request_body,omitempty"`
	ContentType    string `json:"content_type,omitempty"`
}
