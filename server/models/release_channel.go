package models

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/layer5io/meshkit/logger"
)

type ReleaseChannel interface {
	Intercept(req *http.Request, res http.ResponseWriter)
}

type Default struct{}

func (Default) Intercept(req *http.Request, res http.ResponseWriter) {}

func NewReleaseChannelInterceptor(channel string, provider Provider, log logger.Handler) ReleaseChannel {
	stable := &Stable{}
	edge := &Edge{}
	kanvas := &Kanvas{Provider: provider, log: log}

	switch channel {
	case kanvas.String():
		return kanvas
	case stable.String():
		fallthrough
	case edge.String():
		fallthrough
	default:
		return nil
	}
}

type Stable struct{ Default }

func (Stable) String() string {
	return "stable"
}

type Edge struct{ Default }

func (*Edge) String() string {
	return "edge"
}

type Kanvas struct {
	Provider Provider
	log      logger.Handler
}

func (Kanvas) String() string {
	return "kanvas"
}

func (k *Kanvas) Intercept(req *http.Request, res http.ResponseWriter) {
	providerProperties := k.Provider.GetProviderProperties()
	providerURL, _ := url.Parse(k.Provider.GetProviderURL())

	ep, exists := providerProperties.Capabilities.GetEndpointForFeature(PersistAnonymousUser)
	if !exists {
		k.log.Warn(ErrInvalidCapability("PersistAnonymousUser", k.Provider.Name()))
		return
	}

	providerURL = providerURL.JoinPath(ep)

	client := &http.Client{}
	newReq, _ := http.NewRequest("GET", providerURL.String(), nil)

	newReq.Header.Set("X-API-Key", GlobalTokenForAnonymousResults)

	resp, err := client.Do(newReq)
	if err != nil {
		k.log.Error(ErrDoRequest(err, "GET", providerURL.String()))
		return
	}
	defer resp.Body.Close()

	var token string
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == TokenCookieName {
			token = cookie.Value
			break
		}
	}

	if token == "" {
		k.log.Error(ErrGetToken(fmt.Errorf("token not found in the response")))
		return
	}

	k.Provider.SetJWTCookie(res, token)
	redirectURL := getRedirectURLForNavigatorExtension(&providerProperties)

	http.Redirect(res, req, redirectURL, http.StatusFound)
}
