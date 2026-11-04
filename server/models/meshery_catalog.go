package models

import (
	"time"

	"github.com/meshery/schemas/models/core"
)

type CatalogRequest struct {
	ID          core.Uuid            `json:"id,omitempty" db:"id"`
	ContentID   core.Uuid            `json:"contentId,omitempty" db:"content_id"`
	ContentName string               `json:"contentName,omitempty" db:"content_name"`
	ContentType CatalogContentType   `json:"contentType,omitempty" db:"content_type"`
	FirstName   string               `json:"firstName,omitempty" db:"first_name"`
	LastName    string               `json:"lastName,omitempty" db:"last_name"`
	Email       string               `json:"email,omitempty" db:"email"`
	Status      CatalogRequestStatus `json:"status,omitempty" db:"status"`
	CreatedAt   time.Time            `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt   time.Time            `json:"updatedAt,omitempty" db:"updated_at"`
}

type CatalogContentType string

const (
	CatalogPattern CatalogContentType = "pattern"
	CatalogFilter  CatalogContentType = "filter"
)

type CatalogRequestStatus string

const (
	CatalogPending  CatalogRequestStatus = "pending"
	CatalogApproved CatalogRequestStatus = "approved"
	CatalogDenied   CatalogRequestStatus = "denied"
)
