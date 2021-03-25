package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/lib/pq"
)

// PerformanceProfile represents the performance profile that needs
// to be saved
type PerformanceProfile struct {
	ID *uuid.UUID `json:"id,omitempty"`

	Name              string         `json:"name,omitempty"`
	LastRun           string         `json:"last_run,omitempty"`
	Schedule          *uuid.UUID     `json:"schedule,omitempty"`
	LoadGenerators    pq.StringArray `json:"load_generators,omitempty" gorm:"type:text[]"`
	Endpoints         pq.StringArray `json:"endpoints,omitempty" gorm:"type:text[]"`
	ServiceMesh       string         `json:"service_mesh,omitempty"`
	ConcurrentRequest int            `json:"concurrent_request,omitempty"`
	QPS               int            `json:"qps,omitempty"`
	Duration          string         `json:"duration,omitempty"`
	TotalResults      int            `json:"total_results,omitempty"`

	RequestHeaders string `json:"request_headers,omitempty"`
	RequestCookies string `json:"request_cookies,omitempty"`
	RequestBody    string `json:"request_body,omitempty"`
	ContentType    string `json:"content_type,omitempty"`

	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
