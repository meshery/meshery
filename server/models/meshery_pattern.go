package models

import (
	"time"

	"database/sql"
	"github.com/gofrs/uuid"
	isql "github.com/layer5io/meshery/server/internal/sql"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
)

type DesignType string

type DesignTypeResponse struct {
	Type                DesignType `json:"design_type"`
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
	ID *uuid.UUID `json:"id,omitempty"`

	Name        string `json:"name,omitempty"`
	PatternFile string `json:"pattern_file"`
	// Meshery doesn't have the user id fields
	// but the remote provider is allowed to provide one
	UserID *string `json:"user_id"`

	Location      isql.Map       `json:"location"`
	Visibility    string         `json:"visibility"`
	CatalogData   isql.Map       `json:"catalog_data,omitempty"`
	Type          sql.NullString `json:"type"`
	SourceContent []byte         `json:"source_content"`

	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

// MesheryCatalogPatternRequestBody refers to the type of request body
// that PublishCatalogPattern would receive
type MesheryCatalogPatternRequestBody struct {
	ID          uuid.UUID `json:"id,omitempty"`
	CatalogData isql.Map  `json:"catalog_data,omitempty"`
}

// MesheryCatalogPatternRequestBody refers to the type of request body
// that CloneMesheryPatternHandler would receive
type MesheryClonePatternRequestBody struct {
	Name string `json:"name,omitempty"`
}

// GetPatternName takes in a stringified patternfile and extracts the name from it
func GetPatternName(stringifiedFile string) (string, error) {
	out := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(stringifiedFile), &out); err != nil {
		return "", err
	}

	// Get Name from the file
	name, ok := out["name"].(string)
	if !ok {
		return "", ErrPatternFileName
	}

	return name, nil
}
