package models

import (
	"encoding/json"
	"time"

	"database/sql"

	isql "github.com/meshery/meshery/server/internal/sql"
	"github.com/meshery/meshkit/models/catalog/v1alpha1"
	"github.com/meshery/schemas/models/core"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
)

type DesignType string

type DesignTypeResponse struct {
	Type                DesignType `json:"designType"`
	SupportedExtensions []string   `json:"supported_extensions"`
}

func GetDesignsTypes() (r []DesignTypeResponse) {
	r = append(r, DesignTypeResponse{
		Type:                HelmChart,
		SupportedExtensions: []string{".tgz"},
	},
		DesignTypeResponse{
			Type:                DockerCompose,
			SupportedExtensions: []string{".yaml", ".yml"},
		},
		DesignTypeResponse{
			Type:                K8sManifest,
			SupportedExtensions: []string{".yaml", ".yml"},
		},
		DesignTypeResponse{
			Type:                Design,
			SupportedExtensions: []string{".yaml", ".yml"},
		})
	return
}

const (
	HelmChart     DesignType = "Helm Chart"
	DockerCompose DesignType = "Docker Compose"
	K8sManifest   DesignType = "Kubernetes Manifest"
	Design        DesignType = "Design"
)

// reason for adding this constucts is because these has been removed in latest client-go
// https://github.com/kubernetes/client-go/commit/0f17f43973be78f6dcaf6d9a8614fcb35be40d5c#diff-b49fe30cb74d2c3c9c0ca1438056432985f3cad978fd6440f91b695e16195ded
type ListMetaApplyConfiguration struct {
	SelfLink           *string `json:"selfLink,omitempty"`
	ResourceVersion    *string `json:"resourceVersion,omitempty"`
	Continue           *string `json:"continue,omitempty"`
	RemainingItemCount *int64  `json:"remainingItemCount,omitempty"`
}

type StatusCauseApplyConfiguration struct {
	Type    *metav1.CauseType `json:"reason,omitempty"`
	Message *string           `json:"message,omitempty"`
	Field   *string           `json:"field,omitempty"`
}

type StatusDetailsApplyConfiguration struct {
	Name              *string                         `json:"name,omitempty"`
	Group             *string                         `json:"group,omitempty"`
	Kind              *string                         `json:"kind,omitempty"`
	UID               *types.UID                      `json:"uid,omitempty"`
	Causes            []StatusCauseApplyConfiguration `json:"causes,omitempty"`
	RetryAfterSeconds *int32                          `json:"retryAfterSeconds,omitempty"`
}

type StatusApplyConfiguration struct {
	// TypeMetaApplyConfiguration  `json:",inline"`
	*ListMetaApplyConfiguration `json:"metadata,omitempty"`
	Status                      *string                          `json:"status,omitempty"`
	Message                     *string                          `json:"message,omitempty"`
	Reason                      *metav1.StatusReason             `json:"reason,omitempty"`
	Details                     *StatusDetailsApplyConfiguration `json:"details,omitempty"`
	Code                        *int32                           `json:"code,omitempty"`
}

// MesheryPattern represents the patterns that needs to be saved
type MesheryPattern struct {
	ID *core.Uuid `json:"id,omitempty"`

	Name        string `json:"name,omitempty"`
	PatternFile string `json:"patternFile"`
	// Meshery doesn't have the user id fields
	// but the remote provider is allowed to provide one
	UserID *string `json:"user_id"`

	Location      isql.Map             `json:"location"`
	Visibility    string               `json:"visibility"`
	CatalogData   v1alpha1.CatalogData `json:"catalogData,omitempty" gorm:"type:bytes;serializer:json"`
	Type          sql.NullString       `json:"type"`
	SourceContent []byte               `json:"source_content"`

	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`

	ViewCount       int       `json:"viewCount" db:"view_count"`
	ShareCount      int       `json:"shareCount" db:"share_count"`
	DownloadCount   int       `json:"downloadCount" db:"download_count"`
	CloneCount      int       `json:"cloneCount" db:"clone_count"`
	DeploymentCount int       `json:"deploymentCount" db:"deployment_count"`
	WorkspaceID     core.Uuid `json:"workspace_id,omitempty" db:"-"`
	OrgID           core.Uuid `json:"orgId,omitempty" db:"-"`
}

// UnmarshalJSON dual-accepts the canonical camelCase count keys
// (`viewCount`, `shareCount`, `downloadCount`, `cloneCount`,
// `deploymentCount`) and the legacy snake_case spellings
// (`view_count`, `share_count`, `download_count`, `clone_count`,
// `deployment_count`) on the MesheryPattern wire. This is the Phase 2.K
// cascade shim: meshery/schemas v1.2.0 flipped the canonical property
// names to camelCase, and inbound responses from remote providers may
// still carry the legacy snake_case form during the deprecation window.
// Canonical wins when both spellings are present for a given field.
//
// Other fields unmarshal via stdlib default rules through the
// embedded-alias pattern; the five count fields are explicitly reset
// before the precedence switch so a reused receiver does not carry
// stale counts when the next payload omits both spellings.
//
// Remove once every known upstream producer (meshery-cloud remote
// provider, Kanvas catalog API) has migrated off the snake_case
// spellings.
func (m *MesheryPattern) UnmarshalJSON(data []byte) error {
	type alias MesheryPattern
	aux := &struct {
		*alias
		ViewCountCanonical       *int `json:"viewCount,omitempty"`
		ViewCountLegacy          *int `json:"view_count,omitempty"`
		ShareCountCanonical      *int `json:"shareCount,omitempty"`
		ShareCountLegacy         *int `json:"share_count,omitempty"`
		DownloadCountCanonical   *int `json:"downloadCount,omitempty"`
		DownloadCountLegacy      *int `json:"download_count,omitempty"`
		CloneCountCanonical      *int `json:"cloneCount,omitempty"`
		CloneCountLegacy         *int `json:"clone_count,omitempty"`
		DeploymentCountCanonical *int `json:"deploymentCount,omitempty"`
		DeploymentCountLegacy    *int `json:"deployment_count,omitempty"`
	}{alias: (*alias)(m)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	m.ViewCount = 0
	switch {
	case aux.ViewCountCanonical != nil:
		m.ViewCount = *aux.ViewCountCanonical
	case aux.ViewCountLegacy != nil:
		m.ViewCount = *aux.ViewCountLegacy
	}
	m.ShareCount = 0
	switch {
	case aux.ShareCountCanonical != nil:
		m.ShareCount = *aux.ShareCountCanonical
	case aux.ShareCountLegacy != nil:
		m.ShareCount = *aux.ShareCountLegacy
	}
	m.DownloadCount = 0
	switch {
	case aux.DownloadCountCanonical != nil:
		m.DownloadCount = *aux.DownloadCountCanonical
	case aux.DownloadCountLegacy != nil:
		m.DownloadCount = *aux.DownloadCountLegacy
	}
	m.CloneCount = 0
	switch {
	case aux.CloneCountCanonical != nil:
		m.CloneCount = *aux.CloneCountCanonical
	case aux.CloneCountLegacy != nil:
		m.CloneCount = *aux.CloneCountLegacy
	}
	m.DeploymentCount = 0
	switch {
	case aux.DeploymentCountCanonical != nil:
		m.DeploymentCount = *aux.DeploymentCountCanonical
	case aux.DeploymentCountLegacy != nil:
		m.DeploymentCount = *aux.DeploymentCountLegacy
	}
	return nil
}

// MesheryCatalogPatternRequestBody refers to the type of request body
// that PublishCatalogPattern would receive
type MesheryCatalogPatternRequestBody struct {
	ID          core.Uuid `json:"id,omitempty"`
	CatalogData isql.Map  `json:"catalog_data,omitempty"`
}

// MesheryCatalogPatternRequestBody refers to the type of request body
// that CloneMesheryPatternHandler would receive
type MesheryClonePatternRequestBody struct {
	Name string `json:"name,omitempty"`
}

// GetPatternName takes in a stringified patternfile and extracts the name from it
// patternfile can be in yaml and json format
func GetPatternName(stringifiedFile string) (string, error) {
	out := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(stringifiedFile), &out); err != nil {
		if err := json.Unmarshal([]byte(stringifiedFile), &out); err != nil {
			return "", err
		}
	}

	// Get Name from the file
	name, ok := out["name"].(string)
	if !ok {
		return "", ErrPatternFileName
	}

	return name, nil
}

type MesheryPatternFileDeployPayload struct {
	PatternFile string    `json:"pattern_file"`
	PatternID   core.Uuid `json:"pattern_id"`
}
