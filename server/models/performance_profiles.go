package models

import (
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/internal/sql"
	"github.com/lib/pq"

	SMP "github.com/layer5io/service-mesh-performance/spec"
)

// PerformanceProfile represents the performance profile that needs
// to be saved
type PerformanceProfile struct {
	ID *uuid.UUID `json:"id,omitempty"`

	Name              string         `json:"name,omitempty"`
	LastRun           *sql.Time      `json:"last_run,omitempty" gorm:"type:datetime"`
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
	Metadata    	sql.Map `json:"metadata,omitempty"`
	RequestBody    string `json:"request_body,omitempty"`
	ContentType    string `json:"content_type,omitempty"`

	UpdatedAt *sql.Time `json:"updated_at,omitempty"`
	CreatedAt *sql.Time `json:"created_at,omitempty"`
}

type PerformanceTestConfigFile struct {
	Config      *SMP.PerformanceTestConfig `json:"test,omitempty"`
	ServiceMesh *SMP.ServiceMesh           `json:"mesh,omitempty"`
}
