package models

import (
	"github.com/gofrs/uuid"
)

type SmiResult struct {
	ID uuid.UUID `json:"id,omitmepty"`
}
