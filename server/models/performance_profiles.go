package models

import (
	"github.com/lib/pq"
	"github.com/meshery/meshery/server/internal/sql"
	"github.com/meshery/schemas/models/core"
	perfprofile "github.com/meshery/schemas/models/v1beta3/performance_profile"
)

// PerformanceProfile represents the performance profile that needs
// to be saved
type PerformanceProfile struct {
	ID *core.Uuid `json:"id,omitempty"`

	Name              string         `json:"name,omitempty"`
	LastRun           *sql.Time      `json:"lastRun,omitempty" gorm:"type:datetime"`
	Schedule          *core.Uuid     `json:"schedule,omitempty"`
	LoadGenerators    pq.StringArray `json:"loadGenerators,omitempty" gorm:"type:text[]"`
	Endpoints         pq.StringArray `json:"endpoints,omitempty" gorm:"type:text[]"`
	ServiceMesh       string         `json:"serviceMesh,omitempty"`
	ConcurrentRequest int            `json:"concurrentRequest,omitempty"`
	QPS               int            `json:"qps,omitempty"`
	Duration          string         `json:"duration,omitempty"`
	TotalResults      int            `json:"totalResults,omitempty"`

	RequestHeaders string  `json:"requestHeaders,omitempty"`
	RequestCookies string  `json:"requestCookies,omitempty"`
	Metadata       sql.Map `json:"metadata,omitempty"`
	RequestBody    string  `json:"requestBody,omitempty"`
	ContentType    string  `json:"contentType,omitempty"`

	UpdatedAt *sql.Time `json:"updatedAt,omitempty"`
	CreatedAt *sql.Time `json:"createdAt,omitempty"`
}

// PerformanceTestConfigFile is the uploadable envelope that wraps a runtime
// PerformanceTestConfig together with the technology under test. It replaces
// the legacy SMP {test, mesh} file format; "mesh" is now a free-form Meshery
// Registry model name rather than the constraining SMP ServiceMesh enum.
type PerformanceTestConfigFile struct {
	Config      *perfprofile.PerformanceTestConfig `json:"test,omitempty"`
	ServiceMesh string                             `json:"mesh,omitempty"`
}
