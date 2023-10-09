package handlers

import (
	"github.com/google/uuid"
	"github.com/layer5io/meshkit/database"
)

type MesheryHostsPersistor struct {
	DB *database.Handler
}

type MesheryHostsDisplay struct {
	ID       uuid.UUID           `json:"id"`
	Hostname string              `json:"hostname"`
	Port     int                 `json:"port"`
	Summary  hostIndividualCount `json:"summary"`
}
type MesheryHostsContextPage struct {
	Page        int                   `json:"page"`
	PageSize    int                   `json:"page_size"`
	Count       int64                 `json:"total_count"`
	Registrants []MesheryHostsDisplay `json:"registrants"`
}
type hostIndividualCount struct {
	Models        int64 `json:"models"`
	Components    int64 `json:"components"`
	Relationships int64 `json:"relationships"`
	Policies      int64 `json:"policies"`
}
