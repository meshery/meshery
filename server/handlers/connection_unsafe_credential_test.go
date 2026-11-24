package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/stretchr/testify/assert"
)

// TestSaveAndUpdateConnection_RejectUnsafeKubeCredential pins issue #24234:
// direct connection write paths (SaveConnection and UpdateConnectionById)
// must reject unsafe Kubernetes credential secrets carrying exec plugins,
// auth-providers, or local file references before persistence.
func TestSaveAndUpdateConnection_RejectUnsafeKubeCredential(t *testing.T) {
	cluster := `"cluster":{"name":"c","cluster":{"server":"https://127.0.0.1:1"}}`
	secrets := map[string]string{
		"exec":          `{"auth":{"name":"u","user":{"exec":{"apiVersion":"client.authentication.k8s.io/v1","command":"leak-cmd","args":["/leak/path"]}}},` + cluster + `}`,
		"auth-provider": `{"auth":{"name":"u","user":{"auth-provider":{"name":"oidc"}}},` + cluster + `}`,
		"tokenFile":     `{"auth":{"name":"u","user":{"tokenFile":"/leak/path"}},` + cluster + `}`,
		"client-key":    `{"auth":{"name":"u","user":{"client-key":"/leak/path"}},` + cluster + `}`,
		"ca path":       `{"auth":{"name":"u","user":{"token":"t"}},"cluster":{"name":"c","cluster":{"server":"https://k","certificate-authority":"/leak/path"}}}`,
	}

	for name, secret := range secrets {
		for _, method := range []string{http.MethodPost, http.MethodPut} {
			t.Run(name+" "+method, func(t *testing.T) {
				h, provider := newDeleteConnectionFixture(t)
				connID := uuid.Must(uuid.NewV4())
				body := fmt.Sprintf(`{"name":"test-conn","kind":"kubernetes","type":"platform","credentialSecret":%s}`, secret)

				var req *http.Request
				if method == http.MethodPost {
					req = httptest.NewRequest(method, "/api/integrations/connections", bytes.NewBufferString(body))
				} else {
					req = httptest.NewRequest(method, "/api/integrations/connections/"+connID.String(), bytes.NewBufferString(body))
					req = mux.SetURLVars(req, map[string]string{"connectionId": connID.String()})
				}
				req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
				rec := httptest.NewRecorder()
				user := &models.User{ID: uuid.Must(uuid.NewV4())}

				if method == http.MethodPost {
					h.SaveConnection(rec, req, nil, user, provider)
				} else {
					h.UpdateConnectionById(rec, req, nil, user, provider)
				}

				assert.Equal(t, http.StatusBadRequest, rec.Code, "status = %d, want 400 (body=%q)", rec.Code, rec.Body.String())
				assert.NotContains(t, rec.Body.String(), "leak-cmd", "response echoes leak-cmd")
				assert.NotContains(t, rec.Body.String(), "/leak/path", "response echoes /leak/path")

				errBody := decodeErrorBody(t, rec.Body.Bytes())
				assert.Equal(t, models.ErrUnsafeKubeCredentialCode, errBody.Code, "error code = %q, want %q", errBody.Code, models.ErrUnsafeKubeCredentialCode)
			})
		}
	}

	t.Run("safe embedded kubeconfig is accepted on SaveConnection", func(t *testing.T) {
		h, provider := newDeleteConnectionFixture(t)
		body := `{"name":"safe-conn","kind":"kubernetes","type":"platform","credentialSecret":{"auth":{"name":"u","user":{"token":"safe-token"}},"cluster":{"name":"c","cluster":{"server":"https://127.0.0.1:1"}}}}`
		req := httptest.NewRequest(http.MethodPost, "/api/integrations/connections", bytes.NewBufferString(body))
		req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
		rec := httptest.NewRecorder()
		user := &models.User{ID: uuid.Must(uuid.NewV4())}

		h.SaveConnection(rec, req, nil, user, provider)

		assert.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
	})
}
