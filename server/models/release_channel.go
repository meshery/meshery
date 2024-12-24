package models

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshkit/encoding"
	"github.com/layer5io/meshkit/logger"
	"github.com/pkg/errors"
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

type AnonymousFlowResponse struct {
	AccessToken  string             `json:"access_token"`
	Capabilities ProviderProperties `json:"capability,omitempty"`
	UserID       uuid.UUID          `json:"user_id,omitempty"`
}

// This function is triggered when the token is invalid, the channel is in Kanvas mode, or when a token is present but the capabilities lack the required extension.
// Process flow:
// 1. Retrieve the capabilities of the anonymous user.
// 2. Construct the connection payload.
// 3. Generate the user using the special token.
// 4. If the navigator is missing, prompt the user to download the navigator extension.

// The appearance of the default /error page on the Kanvas site could be due to one of the following reasons:
// 1. The remote provider is down, causing an error while fetching capabilities.
// 2. The provided special token is either empty or incorrect, resulting in an unauthorized error.
// 3. An issue occurred during the download of the extension package.

// For exact diagnostics, check the server logs in case of any of these errors.

func (k *Kanvas) Intercept(req *http.Request, res http.ResponseWriter) {
	providerProperties := k.Provider.GetProviderProperties()
	providerURL, _ := url.Parse(k.Provider.GetProviderURL())
	errorUI := "/error"
	ep, exists := providerProperties.Capabilities.GetEndpointForFeature(PersistAnonymousUser)
	if !exists {
		err := ErrInvalidCapability("PersistAnonymousUser", k.Provider.Name())
		k.log.Error(err)
		http.Redirect(res, req, errorUI, http.StatusFound)
		return
	}

	credential := make(map[string]interface{}, 0)
	connectionPayload := connections.BuildMesheryConnectionPayload(req.Context().Value(MesheryServerURL).(string), credential)

	providerURL = providerURL.JoinPath(ep)

	buf, _ := encoding.Marshal(connectionPayload)
	data := bytes.NewReader(buf)

	client := &http.Client{}
	newReq, _ := http.NewRequest("POST", providerURL.String(), data)

	newReq.Header.Set("X-API-Key", GlobalTokenForAnonymousResults)

	resp, err := client.Do(newReq)
	if err != nil {
		err = ErrUnreachableRemoteProvider(err)
		k.log.Error(err)
		http.Redirect(res, req, errorUI, http.StatusFound)
		return
	}
	defer resp.Body.Close()

	flowResponse := AnonymousFlowResponse{}
	err = json.NewDecoder(resp.Body).Decode(&flowResponse)
	if err != nil {
		err = ErrUnmarshal(err, "user flow response")
		k.log.Error(err)
		http.Redirect(res, req, errorUI, http.StatusFound)
		return
	}

	k.Provider.SetJWTCookie(res, flowResponse.AccessToken)
	flowResponse.Capabilities.ProviderURL = k.Provider.GetProviderURL()

	err = k.Provider.WriteCapabilitiesForUser(flowResponse.UserID.String(), &flowResponse.Capabilities)
	if err != nil {
		err = ErrDBPut(errors.Wrapf(err, "failed to write capabilities for the user %s", flowResponse.UserID.String()))
		k.log.Error(err)
		http.Redirect(res, req, errorUI, http.StatusFound)

		return
	}
	k.Provider.SetProviderProperties(flowResponse.Capabilities)
	// Download the package for the user only if they have extension capability
	// The download is skipped if package already exists.
	if len(flowResponse.Capabilities.Extensions.Navigator) > 0 {
		flowResponse.Capabilities.DownloadProviderExtensionPackage(k.log)
	}
	redirectURL := GetRedirectURLForNavigatorExtension(&providerProperties, k.log)
	if redirectURL == "/" {
		k.log.Info("No navigator extension found, redirecting to /error")
		redirectURL = errorUI
	}
	http.Redirect(res, req, redirectURL, http.StatusFound)
}
