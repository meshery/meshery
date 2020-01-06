package models

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// ProviderType - for representing provider types
type ProviderType string

const (
	// LocalProviderType - represents local providers
	LocalProviderType ProviderType = "local"

	// CloudProviderType - represents cloud providers
	CloudProviderType ProviderType = "cloud"
)

// Provider - interface for providers
type Provider interface {
	PreferencePersister
	// Returns ProviderType
	GetProviderType() ProviderType
	// InitiateLogin - does the needed check, returns a true to indicate "return" or false to continue
	InitiateLogin(http.ResponseWriter, *http.Request, bool)
	GetSession(req *http.Request) (*sessions.Session, error)
	GetUserDetails(*http.Request) (*User, error)
	GetProviderToken(req *http.Request) (string, error)
	Logout(http.ResponseWriter, *http.Request)
	FetchResults(req *http.Request, page, pageSize, search, order string) ([]byte, error)
	PublishResults(req *http.Request, data []byte) (string, error)
	PublishMetrics(tokenVal string, data []byte) error

	RecordPreferences(req *http.Request, userID string, data *Preference) error
}
