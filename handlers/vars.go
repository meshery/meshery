package handlers

import (
	"html/template"
)

const (
	// sessionUserKey       = "twitterID"
	// sessionUserName      = "twitterUserName"
	// sessionTwitterToken  = "token"
	// sessionTwitterSecret = "secret"
	cookieSuffix = "_referrer"
	// saasTokenName        = "meshery_saas"
)

var (
	// sessionName   = os.Getenv("EVENT") + "_pl"
	// sessionSecret = base64.StdEncoding.EncodeToString([]byte(sessionName))
	// sessionStore  = sessions.NewCookieStore([]byte(sessionSecret), nil)
	// cookieName    = os.Getenv("EVENT") + cookieSuffix + "_pl"

	dashTempl         = template.Must(template.ParseFiles("../public/dashboard.html"))
	getK8SConfigTempl = template.Must(template.ParseFiles("../public/get-k8s-config.html"))
)
