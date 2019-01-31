package handlers

import (
	"encoding/base64"
	"html/template"
	"os"

	"github.com/dghubble/sessions"
	"github.com/layer5io/meshery/meshes"
)

const (
	sessionUserKey       = "twitterID"
	sessionUserName      = "twitterUserName"
	sessionTwitterToken  = "token"
	sessionTwitterSecret = "secret"
	cookieSuffix         = "_referrer"
)

var (
	sessionName   = os.Getenv("EVENT") + "_pl"
	sessionSecret = base64.StdEncoding.EncodeToString([]byte(sessionName))
	sessionStore  = sessions.NewCookieStore([]byte(sessionSecret), nil)
	cookieName    = os.Getenv("EVENT") + cookieSuffix + "_pl"

	dashTempl       = template.Must(template.ParseFiles("../public/dashboard.html"))
	getAOTokenTempl = template.Must(template.ParseFiles("../public/get-ao-token.html"))
)

// ServerConfig - config type to hold the needed instances/properties
type ServerConfig struct {
	MeshClient meshes.MeshClient
}
