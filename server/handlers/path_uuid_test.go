package handlers

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/v1beta3/environment"
	"github.com/meshery/schemas/models/v1beta3/workspace"
)

const (
	pathUUIDWorkspace  = "11111111-1111-1111-1111-111111111111"
	pathUUIDSecond     = "22222222-2222-2222-2222-222222222222"
	pathUUIDThird      = "33333333-3333-3333-3333-333333333333"
	pathUUIDMalformed  = "not-a-uuid"
	pathUUIDOrgID      = "44444444-4444-4444-4444-444444444444"
	pathUUIDUpdateBody = `{"name":"updated","organizationId":"44444444-4444-4444-4444-444444444444"}`
)

// pathUUIDSpyProvider embeds DefaultLocalProvider and overrides every provider
// method used by the workspace and environment handlers that take an ID path
// variable. Each override records the method name and the ID arguments it was
// called with, so a test can tell whether (and with what) the provider was
// reached.
type pathUUIDSpyProvider struct {
	*models.DefaultLocalProvider
	calls []string
	ids   []string
}

func newPathUUIDSpyProvider() *pathUUIDSpyProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &pathUUIDSpyProvider{DefaultLocalProvider: base}
}

func (p *pathUUIDSpyProvider) record(method string, ids ...string) {
	p.calls = append(p.calls, method)
	p.ids = ids
}

func (p *pathUUIDSpyProvider) GetWorkspaceByID(_ *http.Request, workspaceID, _ string) ([]byte, error) {
	p.record("GetWorkspaceByID", workspaceID)
	return []byte(`{}`), nil
}

func (p *pathUUIDSpyProvider) DeleteWorkspace(_ *http.Request, workspaceID string) ([]byte, error) {
	p.record("DeleteWorkspace", workspaceID)
	return []byte(`{}`), nil
}

func (p *pathUUIDSpyProvider) UpdateWorkspace(_ *http.Request, _ *workspace.WorkspaceUpdatePayload, workspaceID string) (*workspace.Workspace, error) {
	p.record("UpdateWorkspace", workspaceID)
	return &workspace.Workspace{}, nil
}

func (p *pathUUIDSpyProvider) GetEnvironmentsOfWorkspace(_ *http.Request, workspaceID, _, _, _, _, _ string) ([]byte, error) {
	p.record("GetEnvironmentsOfWorkspace", workspaceID)
	return []byte(`{}`), nil
}

func (p *pathUUIDSpyProvider) GetDesignsOfWorkspace(_ *http.Request, workspaceID, _, _, _, _, _ string, _ []string) ([]byte, error) {
	p.record("GetDesignsOfWorkspace", workspaceID)
	return []byte(`{}`), nil
}

func (p *pathUUIDSpyProvider) GetViewsOfWorkspace(_ *http.Request, workspaceID, _, _, _, _, _ string) ([]byte, error) {
	p.record("GetViewsOfWorkspace", workspaceID)
	return []byte(`{}`), nil
}

func (p *pathUUIDSpyProvider) GetTeamsOfWorkspace(_ *http.Request, workspaceID, _, _, _, _, _ string) ([]byte, error) {
	p.record("GetTeamsOfWorkspace", workspaceID)
	return []byte(`{}`), nil
}

func (p *pathUUIDSpyProvider) AddEnvironmentToWorkspace(_ *http.Request, workspaceID, environmentID string) ([]byte, error) {
	p.record("AddEnvironmentToWorkspace", workspaceID, environmentID)
	return []byte(`{}`), nil
}

func (p *pathUUIDSpyProvider) RemoveEnvironmentFromWorkspace(_ *http.Request, workspaceID, environmentID string) ([]byte, error) {
	p.record("RemoveEnvironmentFromWorkspace", workspaceID, environmentID)
	return []byte(`{}`), nil
}

func (p *pathUUIDSpyProvider) AddDesignToWorkspace(_ *http.Request, workspaceID, designID string) ([]byte, error) {
	p.record("AddDesignToWorkspace", workspaceID, designID)
	return []byte(`{}`), nil
}

func (p *pathUUIDSpyProvider) RemoveDesignFromWorkspace(_ *http.Request, workspaceID, designID string) ([]byte, error) {
	p.record("RemoveDesignFromWorkspace", workspaceID, designID)
	return []byte(`{}`), nil
}

func (p *pathUUIDSpyProvider) AddViewToWorkspace(_ *http.Request, workspaceID, viewID string) ([]byte, error) {
	p.record("AddViewToWorkspace", workspaceID, viewID)
	return []byte(`{}`), nil
}

func (p *pathUUIDSpyProvider) RemoveViewFromWorkspace(_ *http.Request, workspaceID, viewID string) ([]byte, error) {
	p.record("RemoveViewFromWorkspace", workspaceID, viewID)
	return []byte(`{}`), nil
}

func (p *pathUUIDSpyProvider) AddTeamToWorkspace(_ *http.Request, workspaceID, teamID string) ([]byte, error) {
	p.record("AddTeamToWorkspace", workspaceID, teamID)
	return []byte(`{}`), nil
}

func (p *pathUUIDSpyProvider) RemoveTeamFromWorkspace(_ *http.Request, workspaceID, teamID string) ([]byte, error) {
	p.record("RemoveTeamFromWorkspace", workspaceID, teamID)
	return []byte(`{}`), nil
}

func (p *pathUUIDSpyProvider) GetEnvironmentByID(_ *http.Request, environmentID, _ string) ([]byte, error) {
	p.record("GetEnvironmentByID", environmentID)
	return []byte(`{}`), nil
}

func (p *pathUUIDSpyProvider) DeleteEnvironment(_ *http.Request, environmentID string) ([]byte, error) {
	p.record("DeleteEnvironment", environmentID)
	return []byte(`{}`), nil
}

func (p *pathUUIDSpyProvider) UpdateEnvironment(_ *http.Request, _ *environment.EnvironmentPayload, environmentID string) (*environment.Environment, error) {
	p.record("UpdateEnvironment", environmentID)
	return &environment.Environment{}, nil
}

func (p *pathUUIDSpyProvider) GetConnectionsOfEnvironment(_ *http.Request, environmentID, _, _, _, _, _ string) ([]byte, error) {
	p.record("GetConnectionsOfEnvironment", environmentID)
	return []byte(`{}`), nil
}

func (p *pathUUIDSpyProvider) AddConnectionToEnvironment(_ *http.Request, environmentID, connectionID string) ([]byte, error) {
	p.record("AddConnectionToEnvironment", environmentID, connectionID)
	return []byte(`{}`), nil
}

func (p *pathUUIDSpyProvider) RemoveConnectionFromEnvironment(_ *http.Request, environmentID, connectionID string) ([]byte, error) {
	p.record("RemoveConnectionFromEnvironment", environmentID, connectionID)
	return []byte(`{}`), nil
}

type pathUUIDHandlerFunc func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider)

// pathUUIDCase describes one handler: the path variables it reads (in the
// order the provider receives them), the method it must reach, and the query
// string and body needed to get past its other validation.
type pathUUIDCase struct {
	name       string
	httpMethod string
	vars       []string
	query      string
	body       string
	provider   string
	call       pathUUIDHandlerFunc
}

func pathUUIDCases() []pathUUIDCase {
	type hf = pathUUIDHandlerFunc
	var (
		wsGet hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.GetWorkspaceByIdHandler(w, r, nil, nil, p)
		}
		wsDelete hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.DeleteWorkspaceHandler(w, r, nil, nil, p)
		}
		wsUpdate hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.UpdateWorkspaceHandler(w, r, nil, nil, p)
		}
		wsEnvs hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.GetEnvironmentsOfWorkspaceHandler(w, r, nil, nil, p)
		}
		wsDesigns hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.GetDesignsOfWorkspaceHandler(w, r, nil, nil, p)
		}
		wsViews hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.GetViewsOfWorkspaceHandler(w, r, nil, nil, p)
		}
		wsTeams hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.GetTeamsOfWorkspaceHandler(w, r, nil, nil, p)
		}
		wsAddEnv hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.AddEnvironmentToWorkspaceHandler(w, r, nil, nil, p)
		}
		wsRemEnv hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.RemoveEnvironmentFromWorkspaceHandler(w, r, nil, nil, p)
		}
		wsAddDesign hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.AddDesignToWorkspaceHandler(w, r, nil, nil, p)
		}
		wsRemDesign hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.RemoveDesignFromWorkspaceHandler(w, r, nil, nil, p)
		}
		wsAddView hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.AddViewToWorkspaceHandler(w, r, nil, nil, p)
		}
		wsRemView hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.RemoveViewFromWorkspaceHandler(w, r, nil, nil, p)
		}
		wsAddTeam hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.AddTeamToWorkspaceHandler(w, r, nil, nil, p)
		}
		wsRemTeam hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.RemoveTeamFromWorkspaceHandler(w, r, nil, nil, p)
		}
		envGet hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.GetEnvironmentByIDHandler(w, r, nil, nil, p)
		}
		envDelete hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.DeleteEnvironmentHandler(w, r, nil, nil, p)
		}
		envUpdate hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.UpdateEnvironmentHandler(w, r, nil, nil, p)
		}
		envConns hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.GetConnectionsOfEnvironmentHandler(w, r, nil, nil, p)
		}
		envAddConn hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.AddConnectionToEnvironmentHandler(w, r, nil, nil, p)
		}
		envRemConn hf = func(h *Handler, w http.ResponseWriter, r *http.Request, p models.Provider) {
			h.RemoveConnectionFromEnvironmentHandler(w, r, nil, nil, p)
		}
	)
	return []pathUUIDCase{
		{name: "GetWorkspaceByIdHandler", httpMethod: http.MethodGet, vars: []string{"id"}, query: "orgId=" + pathUUIDOrgID, provider: "GetWorkspaceByID", call: wsGet},
		{name: "DeleteWorkspaceHandler", httpMethod: http.MethodDelete, vars: []string{"id"}, provider: "DeleteWorkspace", call: wsDelete},
		{name: "UpdateWorkspaceHandler", httpMethod: http.MethodPut, vars: []string{"id"}, body: pathUUIDUpdateBody, provider: "UpdateWorkspace", call: wsUpdate},
		{name: "GetEnvironmentsOfWorkspaceHandler", httpMethod: http.MethodGet, vars: []string{"id"}, provider: "GetEnvironmentsOfWorkspace", call: wsEnvs},
		{name: "GetDesignsOfWorkspaceHandler", httpMethod: http.MethodGet, vars: []string{"id"}, provider: "GetDesignsOfWorkspace", call: wsDesigns},
		{name: "GetViewsOfWorkspaceHandler", httpMethod: http.MethodGet, vars: []string{"id"}, provider: "GetViewsOfWorkspace", call: wsViews},
		{name: "GetTeamsOfWorkspaceHandler", httpMethod: http.MethodGet, vars: []string{"id"}, provider: "GetTeamsOfWorkspace", call: wsTeams},
		{name: "AddEnvironmentToWorkspaceHandler", httpMethod: http.MethodPost, vars: []string{"id", "environmentID"}, provider: "AddEnvironmentToWorkspace", call: wsAddEnv},
		{name: "RemoveEnvironmentFromWorkspaceHandler", httpMethod: http.MethodDelete, vars: []string{"id", "environmentID"}, provider: "RemoveEnvironmentFromWorkspace", call: wsRemEnv},
		{name: "AddDesignToWorkspaceHandler", httpMethod: http.MethodPost, vars: []string{"id", "designID"}, provider: "AddDesignToWorkspace", call: wsAddDesign},
		{name: "RemoveDesignFromWorkspaceHandler", httpMethod: http.MethodDelete, vars: []string{"id", "designID"}, provider: "RemoveDesignFromWorkspace", call: wsRemDesign},
		{name: "AddViewToWorkspaceHandler", httpMethod: http.MethodPost, vars: []string{"id", "viewID"}, provider: "AddViewToWorkspace", call: wsAddView},
		{name: "RemoveViewFromWorkspaceHandler", httpMethod: http.MethodDelete, vars: []string{"id", "viewID"}, provider: "RemoveViewFromWorkspace", call: wsRemView},
		{name: "AddTeamToWorkspaceHandler", httpMethod: http.MethodPost, vars: []string{"id", "teamID"}, provider: "AddTeamToWorkspace", call: wsAddTeam},
		{name: "RemoveTeamFromWorkspaceHandler", httpMethod: http.MethodDelete, vars: []string{"id", "teamID"}, provider: "RemoveTeamFromWorkspace", call: wsRemTeam},
		{name: "GetEnvironmentByIDHandler", httpMethod: http.MethodGet, vars: []string{"id"}, query: "orgId=" + pathUUIDOrgID, provider: "GetEnvironmentByID", call: envGet},
		{name: "DeleteEnvironmentHandler", httpMethod: http.MethodDelete, vars: []string{"id"}, provider: "DeleteEnvironment", call: envDelete},
		{name: "UpdateEnvironmentHandler", httpMethod: http.MethodPut, vars: []string{"id"}, body: pathUUIDUpdateBody, provider: "UpdateEnvironment", call: envUpdate},
		{name: "GetConnectionsOfEnvironmentHandler", httpMethod: http.MethodGet, vars: []string{"environmentID"}, provider: "GetConnectionsOfEnvironment", call: envConns},
		{name: "AddConnectionToEnvironmentHandler", httpMethod: http.MethodPost, vars: []string{"environmentID", "connectionID"}, provider: "AddConnectionToEnvironment", call: envAddConn},
		{name: "RemoveConnectionFromEnvironmentHandler", httpMethod: http.MethodDelete, vars: []string{"environmentID", "connectionID"}, provider: "RemoveConnectionFromEnvironment", call: envRemConn},
	}
}

// validPathUUIDs returns distinct valid UUIDs for each path variable of c.
func (c pathUUIDCase) validPathUUIDs() map[string]string {
	pool := []string{pathUUIDWorkspace, pathUUIDSecond, pathUUIDThird}
	vars := make(map[string]string, len(c.vars))
	for i, name := range c.vars {
		vars[name] = pool[i]
	}
	return vars
}

func (c pathUUIDCase) run(t *testing.T, vars map[string]string, p models.Provider) *httptest.ResponseRecorder {
	t.Helper()
	target := "/api/test"
	if c.query != "" {
		target += "?" + c.query
	}
	req := httptest.NewRequest(c.httpMethod, target, strings.NewReader(c.body))
	req = mux.SetURLVars(req, vars)
	rec := httptest.NewRecorder()
	h := newTestHandler(t, map[string]models.Provider{}, "")
	c.call(h, rec, req, p)
	return rec
}

func TestRequirePathUUID(t *testing.T) {
	const valid = "0f8fad5b-d9cb-469f-a165-70867728950e"
	cases := []struct {
		name       string
		value      string
		wantOK     bool
		wantStatus int
	}{
		{name: "valid lowercase UUID", value: valid, wantOK: true, wantStatus: http.StatusOK},
		{name: "valid uppercase UUID", value: strings.ToUpper(valid), wantOK: true, wantStatus: http.StatusOK},
		{name: "not a UUID", value: "not-a-uuid", wantStatus: http.StatusBadRequest},
		{name: "empty string", value: "", wantStatus: http.StatusBadRequest},
		{name: "UUID with one character removed", value: valid[:len(valid)-1], wantStatus: http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tc.value})
			rec := httptest.NewRecorder()

			id, ok := h.requirePathUUID(rec, req, "id")

			if ok != tc.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tc.wantOK)
			}
			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, tc.wantStatus, rec.Body.String())
			}
			if tc.wantOK {
				if id != tc.value {
					t.Errorf("id = %q, want %q unchanged", id, tc.value)
				}
				return
			}
			if id != "" {
				t.Errorf("id = %q, want empty on failure", id)
			}
			if decoded := decodeErrorBody(t, rec.Body.Bytes()); decoded.Code != models.ErrInvalidUUIDValueCode {
				t.Errorf("code = %q, want %q", decoded.Code, models.ErrInvalidUUIDValueCode)
			}
		})
	}
}

func TestWorkspaceAndEnvironmentHandlers_MalformedPathUUID_ReturnsBadRequest(t *testing.T) {
	for _, hc := range pathUUIDCases() {
		for _, bad := range hc.vars {
			t.Run(hc.name+"/malformed_"+bad, func(t *testing.T) {
				vars := hc.validPathUUIDs()
				vars[bad] = pathUUIDMalformed
				provider := newPathUUIDSpyProvider()

				rec := hc.run(t, vars, provider)

				if rec.Code != http.StatusBadRequest {
					t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusBadRequest, rec.Body.String())
				}
				if decoded := decodeErrorBody(t, rec.Body.Bytes()); decoded.Code != models.ErrInvalidUUIDValueCode {
					t.Errorf("code = %q, want %q", decoded.Code, models.ErrInvalidUUIDValueCode)
				}
				if len(provider.calls) != 0 {
					t.Fatalf("provider must not be called, got calls %v", provider.calls)
				}
			})
		}
	}
}

func TestWorkspaceAndEnvironmentHandlers_ValidPathUUID_ReachesProvider(t *testing.T) {
	for _, hc := range pathUUIDCases() {
		t.Run(hc.name, func(t *testing.T) {
			vars := hc.validPathUUIDs()
			provider := newPathUUIDSpyProvider()

			rec := hc.run(t, vars, provider)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusOK, rec.Body.String())
			}
			if !reflect.DeepEqual(provider.calls, []string{hc.provider}) {
				t.Fatalf("provider calls = %v, want exactly [%s]", provider.calls, hc.provider)
			}
			want := make([]string, 0, len(hc.vars))
			for _, name := range hc.vars {
				want = append(want, vars[name])
			}
			if !reflect.DeepEqual(provider.ids, want) {
				t.Errorf("provider received IDs %v, want %v", provider.ids, want)
			}
		})
	}
}
