package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	corev1alpha1 "github.com/meshery/schemas/models/v1alpha1/core"
	"github.com/meshery/schemas/models/v1beta1/workspace"
)

// workspaceMapObject is a DB-safe representation of the generated MapObject type.
type workspaceMapObject map[string]string

func (m workspaceMapObject) GormDataType() string {
	return "text"
}

func (m workspaceMapObject) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}

	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return string(b), nil
}

func (m *workspaceMapObject) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}

	var bytes []byte
	switch typedValue := value.(type) {
	case []byte:
		bytes = typedValue
	case string:
		bytes = []byte(typedValue)
	default:
		return fmt.Errorf("cannot scan type %T into workspaceMapObject", value)
	}

	decoded := make(workspaceMapObject)
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return err
	}

	*m = decoded
	return nil
}

// WorkspaceDBModel keeps the database schema compatible with generated schemas types.
type WorkspaceDBModel struct {
	ID             corev1alpha1.GeneralId `db:"id" json:"id" yaml:"id"`
	CreatedAt      corev1alpha1.CreatedAt `db:"created_at" json:"created_at,omitempty" yaml:"created_at,omitempty"`
	DeletedAt      corev1alpha1.DeletedAt `db:"deleted_at" json:"deleted_at,omitempty" yaml:"deleted_at,omitempty"`
	Description    corev1alpha1.Text      `db:"description" json:"description,omitempty" yaml:"description,omitempty"`
	Metadata       workspaceMapObject     `db:"metadata" json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Name           corev1alpha1.Text      `db:"name" json:"name,omitempty" yaml:"name,omitempty"`
	OrganizationID *corev1alpha1.Uuid     `db:"organization_id" json:"organization_id" yaml:"organization_id"`
	Owner          corev1alpha1.Text      `db:"owner" json:"owner,omitempty" yaml:"owner,omitempty"`
	UpdatedAt      corev1alpha1.UpdatedAt `db:"updated_at" json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
}

func (WorkspaceDBModel) TableName() string {
	return "workspaces"
}

func workspaceDBModelFromSchemaModel(ws *workspace.Workspace) WorkspaceDBModel {
	return WorkspaceDBModel{
		ID:             ws.ID,
		CreatedAt:      ws.CreatedAt,
		DeletedAt:      ws.DeletedAt,
		Description:    ws.Description,
		Metadata:       workspaceMapObject(ws.Metadata),
		Name:           ws.Name,
		OrganizationID: ws.OrganizationID,
		Owner:          ws.Owner,
		UpdatedAt:      ws.UpdatedAt,
	}
}

func workspaceSchemaModelFromDBModel(ws WorkspaceDBModel) workspace.Workspace {
	return workspace.Workspace{
		ID:             ws.ID,
		CreatedAt:      ws.CreatedAt,
		DeletedAt:      ws.DeletedAt,
		Description:    ws.Description,
		Metadata:       corev1alpha1.MapObject(ws.Metadata),
		Name:           ws.Name,
		OrganizationID: ws.OrganizationID,
		Owner:          ws.Owner,
		UpdatedAt:      ws.UpdatedAt,
	}
}
