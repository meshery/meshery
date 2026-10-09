package models

import (
	"encoding/gob"

	"github.com/gofrs/uuid"
	userV1beta2 "github.com/meshery/schemas/models/v1beta2/user"
	userV1beta3 "github.com/meshery/schemas/models/v1beta3/user"
	"github.com/oapi-codegen/runtime/types"
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

// Identity fields of the built-in local provider's single system user. The
// provider is single-user, so several code paths (the /api/user response, the
// persisters' delete-error attribution, the owner joined onto content
// resources) all describe the same person. These constants are the one
// definition they derive from, so renaming the local user cannot silently
// desynchronise those surfaces.
const (
	localProviderUserFirstName = "Meshery"
	localProviderUserLastName  = "Meshery"
	localProviderUserEmail     = "meshery@meshery.local"
)

// User - represents a user in Meshery
type User = userV1beta3.User

// LocalProviderUser returns the built-in local provider's single system user.
func LocalProviderUser() *User {
	avatarURL := ""
	return &User{
		ID:        LocalProviderUserID,
		FirstName: localProviderUserFirstName,
		LastName:  localProviderUserLastName,
		Email:     types.Email(localProviderUserEmail),
		AvatarUrl: &avatarURL,
	}
}

// LocalProviderOwnsContent reports whether the built-in provider's single user
// is the owner of a content resource (design, filter) with the given
// visibility. Published content is the seeded, community-authored catalog that
// SeedContent imports from docs/data/catalog; attributing it to the local user
// would both misattribute third-party work and hand the local user the
// owner-gated edit and visibility affordances over it.
func LocalProviderOwnsContent(visibility string) bool {
	return visibility != Published
}

// LocalProviderContentUser returns LocalProviderUser typed as the schemas
// v1beta2 user.User that design.MesheryPattern.User embeds, so the owner
// profile joined onto content resources is the same identity /api/user reports.
func LocalProviderContentUser() *userV1beta2.User {
	u := LocalProviderUser()
	return &userV1beta2.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		AvatarUrl: u.AvatarUrl,
	}
}

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
