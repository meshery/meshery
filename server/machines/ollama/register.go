package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/core"
)

const (
	defaultOllamaTimeout = 10 * time.Second
	ollamaTagsPath       = "/api/tags"
)

// ollamaConn mirrors the non-secret fields in OllamaConnection.json's connectionSchema.
type ollamaConn struct {
	BaseURL      string `json:"baseUrl"`
	DefaultModel string `json:"defaultModel"`
}

// ollamaCred mirrors credentialSchema. APIKey is optional - most local Ollama
// installs have no auth in front of them.
type ollamaCred struct {
	APIKey string `json:"apiKey"`
}

type RegisterAction struct{}

func (ra *RegisterAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (ra *RegisterAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	logLevel := viper.GetInt("LOG_LEVEL")
	if viper.GetBool("DEBUG") {
		logLevel = int(logrus.DebugLevel)
	}
	log, err := logger.New("meshery", logger.Options{
		Format:   logger.SyslogLogFormat,
		LogLevel: logLevel,
	})
	if err != nil {
		// Matches server/machines/prometheus/register.go's existing convention.
		// Inherited from the merged pattern, not introduced here - a repo-wide
		// fix (returning the error instead of os.Exit) is a separate concern
		// from this connection kind.
		logrus.Error(err)
		os.Exit(1)
	}

	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	sysID, _ := ctx.Value(models.SystemIDKey).(*core.Uuid)
	userUUID := user.ID

	eventBuilder := events.NewEvent().
		ActedUpon(userUUID).
		WithCategory("connection").
		WithAction("update").
		FromSystem(*sysID).
		FromOwner(userUUID).
		WithDescription("Failed to interact with the connection.").
		WithSeverity(events.Error)

	connPayload, err := utils.Cast[connections.ConnectionPayload](data)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	metadata, err := utils.Cast[map[string]interface{}](connPayload.MetaData)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	conn, err := utils.MarshalAndUnmarshal[map[string]interface{}, ollamaConn](metadata)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	cred, err := utils.MarshalAndUnmarshal[map[string]interface{}, ollamaCred](connPayload.CredentialSecret)
	if err != nil && !connPayload.SkipCredentialVerification {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	if err := validateOllamaBaseURL(conn.BaseURL); err != nil {
		wrapped := models.ErrOllamaConnectivity(err)
		return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": wrapped}).Build(), wrapped
	}

	if !connPayload.SkipCredentialVerification {
		if err := checkOllamaReachable(ctx, conn.BaseURL, cred.APIKey); err != nil {
			wrapped := models.ErrOllamaConnectivity(err)
			return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": wrapped}).Build(), wrapped
		}
	}

	log.Info("ollama connectivity check passed for " + conn.BaseURL)
	return machines.Exit, nil, nil
}

func (ra *RegisterAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

// checkOllamaReachable hits Ollama's list-models endpoint as the connectivity
// check - lightweight, and doesn't require a model to already be pulled.
func checkOllamaReachable(ctx context.Context, baseURL, apiKey string) error {
	reqCtx, cancel := context.WithTimeout(ctx, defaultOllamaTimeout)
	defer cancel()

	tagsURL := strings.TrimRight(baseURL, "/") + ollamaTagsPath
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, tagsURL, nil)
	if err != nil {
		return fmt.Errorf("building connectivity request: %w", err)
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("reaching %s: %w", baseURL, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama returned status %d from %s", resp.StatusCode, tagsURL)
	}

	var tagsResp struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tagsResp); err != nil {
		return fmt.Errorf("decoding /api/tags response: %w", err)
	}
	return nil
}

// validateOllamaBaseURL blocks the well-known cloud-metadata SSRF traps while
// deliberately still permitting localhost/private/in-cluster addresses, since
// pointing at a self-hosted instance on the local network is the whole point
// of this connection kind. Different mitigation than a host-allowlist
// approach (see PR description for why an allowlist doesn't fit Ollama).
func validateOllamaBaseURL(raw string) error {
	if raw == "" {
		return fmt.Errorf("baseUrl is required")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("not a valid URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("scheme must be http or https, got %q", u.Scheme)
	}

	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("missing host")
	}

	blockedHosts := map[string]bool{
		"169.254.169.254":          true, // AWS/GCP/Azure IMDS
		"metadata.google.internal": true,
	}
	if blockedHosts[strings.ToLower(host)] {
		return fmt.Errorf("baseUrl host %q is not permitted", host)
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsLinkLocalUnicast() {
		return fmt.Errorf("baseUrl host %q is a link-local address and is not permitted", host)
	}

	return nil
}
