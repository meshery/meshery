package handlers

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/meshery/meshery/server/internal/sql"
	"github.com/meshery/meshery/server/models"
)

// contextWithCredentials builds a K8sContext shaped like one the kubeconfig
// endpoints produce: identity fields the clients read, plus the cluster and auth
// blocks copied verbatim out of the uploaded kubeconfig.
func contextWithCredentials() *models.K8sContext {
	return &models.K8sContext{
		ID:           "ctx-id",
		Name:         "prod",
		Server:       "https://1.2.3.4:6443",
		ConnectionID: "conn-id",
		Version:      "v1.29.0",
		Reachable:    true,
		Auth: sql.Map{
			"name": "u1",
			"user": map[string]interface{}{"token": "SECRET-BEARER-TOKEN"},
		},
		Cluster: sql.Map{
			"name": "prod",
			"cluster": map[string]interface{}{
				"server":                     "https://1.2.3.4:6443",
				"certificate-authority-data": "SECRET-INLINED-CA-BYTES",
			},
		},
	}
}

// assertNoCredentials fails if any credential material survived into body, and
// checks that the identity the UI and mesheryctl read did survive.
func assertNoCredentials(t *testing.T, label, body string) {
	t.Helper()

	for _, secret := range []string{"SECRET-BEARER-TOKEN", "SECRET-INLINED-CA-BYTES"} {
		if strings.Contains(body, secret) {
			t.Errorf("%s: credential %q reached the response body:\n%s", label, secret, body)
		}
	}
	for _, key := range []string{`"auth"`, `"cluster"`} {
		if strings.Contains(body, key) {
			t.Errorf("%s: %s block reached the response body:\n%s", label, key, body)
		}
	}
	// ui/components/connections/wizard/types.ts DiscoveredKubeContext reads
	// id/name/server/reachable/connectionId; mesheryctl reads name and
	// connectionId. Stripping must not take those.
	for _, keep := range []string{`"id":"ctx-id"`, `"name":"prod"`, `"server":"https://1.2.3.4:6443"`, `"connectionId":"conn-id"`, `"reachable":true`} {
		if !strings.Contains(body, keep) {
			t.Errorf("%s: identity field %s must survive stripping, body:\n%s", label, keep, body)
		}
	}
}

// TestDiscoveryResponseCarriesNoCredentials pins the shape encoded by
// GetContextsFromK8SConfig (k8sconfig_handler.go, the Encode of strippedContexts).
// The cluster and auth maps hold the uploaded kubeconfig's credentials verbatim,
// including whatever helpers.FlattenMinifyKubeConfig inlined from the server's
// disk, so they must never be echoed back to the uploader.
func TestDiscoveryResponseCarriesNoCredentials(t *testing.T) {
	contexts := []*models.K8sContext{contextWithCredentials()}

	stripped := make([]models.K8sContext, 0, len(contexts))
	for _, ctx := range contexts {
		stripped = append(stripped, models.StripCredentialsForContext(ctx))
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(stripped); err != nil {
		t.Fatalf("encode: %v", err)
	}
	t.Logf("discovery response body:\n%s", body.String())
	assertNoCredentials(t, "discovery", body.String())
}

// TestSaveK8sContextResponseCarriesNoCredentials pins the same contract for the
// registration endpoint, whose four buckets are also encoded back to the
// uploader (addK8SConfig's Encode of saveK8sContextResponse).
func TestSaveK8sContextResponseCarriesNoCredentials(t *testing.T) {
	ctx := contextWithCredentials()

	resp := SaveK8sContextResponse{
		RegisteredContexts: []models.K8sContext{models.StripCredentialsForContext(ctx)},
		ConnectedContexts:  []models.K8sContext{models.StripCredentialsForContext(ctx)},
		IgnoredContexts:    []models.K8sContext{models.StripCredentialsForContext(ctx)},
		ErroredContexts:    []models.K8sContext{models.StripCredentialsForContext(ctx)},
	}

	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(resp); err != nil {
		t.Fatalf("encode: %v", err)
	}
	t.Logf("registration response body:\n%s", body.String())
	assertNoCredentials(t, "registration", body.String())
}
