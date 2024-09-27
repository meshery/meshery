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

func (k *Kanvas) Intercept(req *http.Request, res http.ResponseWriter) {
	providerProperties := k.Provider.GetProviderProperties()
	providerURL, _ := url.Parse(k.Provider.GetProviderURL())
	errorUI := "/error"
	ep, exists := providerProperties.Capabilities.GetEndpointForFeature(PersistAnonymousUser)
	if !exists {
		k.log.Error(ErrInvalidCapability("PersistAnonymousUser", k.Provider.Name()))
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
		k.log.Error(ErrUnreachableRemoteProvider(err))
		http.Redirect(res, req, errorUI, http.StatusFound)
		return
	}
	defer resp.Body.Close()

	flowResponse := AnonymousFlowResponse{}
	err = json.NewDecoder(resp.Body).Decode(&flowResponse)
	if err != nil {
		k.log.Error(ErrUnmarshal(err, "user flow response"))
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
	redirectURL := getRedirectURLForNavigatorExtension(&providerProperties)

	http.Redirect(res, req, redirectURL, http.StatusFound)
}
