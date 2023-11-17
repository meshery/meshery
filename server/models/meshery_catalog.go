package models

import (
	"github.com/gofrs/uuid"
	"time"
)

type CatalogRequest struct {
	ID          uuid.UUID            `json:"id,omitempty" db:"id"`
	ContentID   uuid.UUID            `json:"content_id,omitempty" db:"content_id"`
	ContentName string               `json:"content_name,omitempty" db:"content_name"`
	ContentType CatalogContentType   `json:"content_type,omitempty" db:"content_type"`
	FirstName   string               `json:"first_name,omitempty" db:"first_name"`
	LastName    string               `json:"last_name,omitempty" db:"last_name"`
	Email       string               `json:"email,omitempty" db:"email"`
	Status      CatalogRequestStatus `json:"status,omitempty" db:"status"`
	CreatedAt   time.Time            `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at,omitempty" db:"updated_at"`
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
