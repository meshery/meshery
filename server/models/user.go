package models

import (
	"encoding/gob"

	userV1beta "github.com/meshery/schemas/models/v1beta1/user"
)

func init() {
	gob.Register(&User{})
}

// GlobalTokenForAnonymousResults - stores the global token for anonymous result publishing
var (
	GlobalTokenForAnonymousResults = "dev_token"
)

// User - represents a user in Meshery
type User = userV1beta.User

type AllUsers struct {
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	Data       []*User `json:"data"`
	TotalCount int     `json:"total_count"`
}

type UserKeys struct {
	ID          string `json:"id,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Function    string `json:"function,omitempty"`
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
	Subcategory string `json:"subcategory,omitempty"`
	Created_at  string `json:"created_at,omitempty"`
	Updated_at  string `json:"updated_at,omitempty"`
	Deleted_at  string `json:"deleted_at,omitempty"`
}
