package models

import (
	"encoding/gob"

	"github.com/gofrs/uuid"
)

func init() {
	gob.Register(&User{})
}

// GlobalTokenForAnonymousResults - stores the global token for anonymous result publishing
var (
	GlobalTokenForAnonymousResults = "dev_token"
)

// User - represents a user in Meshery
type User struct {
	ID        string `json:"id,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Provider  string `json:"provider,omitempty" db:"provider"`
	Email     string `json:"email,omitempty" db:"email"`
	Bio       string `json:"bio,omitempty" db:"bio"`
}

// Role - represents the roles of a user
type Role struct {
	ID        uuid.UUID `json:"id,omitempty" db:"id"`
	RoleName  string    `json:"role_name,omitempty" db:"role_name"`
	CreatedAt string    `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt string    `json:"updated_at,omitempty" db:"updated_at"`
}
