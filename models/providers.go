package models

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type ProviderType string

const (
	LocalProviderType ProviderType = "local"
	CloudProviderType ProviderType = "cloud"
)

type Provider interface {
	SessionPersister
	// Returns ProviderType
	GetProviderType() ProviderType
	// InitiateLogin - does the needed check, returns a true to indicate "return" or false to continue
	InitiateLogin(http.ResponseWriter, *http.Request)
	GetSession(req *http.Request) (*sessions.Session, error)
	GetUserDetails(*http.Request) (*User, error)
	GetProviderToken(req *http.Request) (string, error)
	Logout(http.ResponseWriter, *http.Request)
	FetchResults(req *http.Request, page, pageSize, search, order string) ([]byte, error)
	PublishResults(req *http.Request, data []byte) (string, error)
	PublishMetrics(tokenVal string, data []byte) error

	RecordPreferences(req *http.Request, userID string, data *Session) error
}
