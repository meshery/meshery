package models

import (
	"encoding/gob"

	"github.com/gofrs/uuid"
	userV1beta3 "github.com/meshery/schemas/models/v1beta3/user"
)

func init() {
	gob.Register(&User{})
}

// GlobalTokenForAnonymousResults - stores the global token for anonymous result publishing
var (
	GlobalTokenForAnonymousResults = "dev_token"
)

// LocalProviderUserID is the stable synthetic id the built-in local provider
// uses to key preferences and persister-scoped data for its single "meshery"
// system user. schemas v1beta3 dropped the string User.UserId ("meshery"); this
// deterministic (namespaced) uuid replaces it so the local persister key stays
// stable across restarts. It is a non-zero UUID on purpose - the zero value
// (uuid.Nil) is treated as "unset" on the remote path (see the
// `user.ID == uuid.Nil` guards in this package).
var LocalProviderUserID = uuid.NewV5(uuid.NamespaceDNS, "meshery-local-provider-user")

// User - represents a user in Meshery
type User = userV1beta3.User

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
	Created_at  string `json:"createdAt,omitempty"`
	Updated_at  string `json:"updatedAt,omitempty"`
	Deleted_at  string `json:"deletedAt,omitempty"`
}
