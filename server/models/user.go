package models

import (
	"encoding/gob"

	userV1beta2 "github.com/meshery/schemas/models/v1beta2/user"
)

func init() {
	gob.Register(&User{})
}

// GlobalTokenForAnonymousResults - stores the global token for anonymous result publishing
var (
	GlobalTokenForAnonymousResults = "dev_token"
)

// User - represents a user in Meshery
type User = userV1beta2.User

type AllUsers struct {
	Page       int     `json:"page"`
	PageSize   int     `json:"pageSize"`
	Data       []*User `json:"data"`
	TotalCount int     `json:"totalCount"`
}

type UserKeys struct {
	ID          string `json:"id,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Function    string `json:"function,omitempty"`
	Category    string `json:"category,omitempty"`
	Description string `json:"description,omitempty"`
	Subcategory string `json:"subcategory,omitempty"`
	CreatedAt  string `json:"createdAt,omitempty"`
	UpdatedAt  string `json:"updatedAt,omitempty"`
	DeletedAt  string `json:"deletedAt,omitempty"`
}
