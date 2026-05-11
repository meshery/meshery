package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"encoding/json"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	userSchema "github.com/meshery/schemas/models/v1beta2/user"
	pkgerrors "github.com/pkg/errors"

	"github.com/meshery/meshery/server/models"
)

// UserHandler returns info about the logged in user
func (h *Handler) UserHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, user *models.User, _ models.Provider) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		obj := "user data"
		h.log.Error(models.ErrEncoding(err, obj))
		writeMeshkitError(w, models.ErrEncoding(err, obj), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	userID := mux.Vars(r)["id"]
	_, err := uuid.FromString(userID)
	if err != nil {
		writeMeshkitError(w, ErrInvalidUUID(err), http.StatusBadRequest)
		return
	}
	resp, err := provider.GetUserByID(r, userID)
	if err != nil {
		if errors.Is(err, models.ErrUserIsSystemInstance) {
			// The requested ID is this Meshery instance's own system UUID; it
			// isn't a real user record. Return a 204 so callers can render a
			// "system" placeholder without treating it as a fetch failure.
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusNotFound)
		return
	}

	if resp == nil {
		writeMeshkitError(w, ErrUserNotFound(userID), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(w, string(resp)); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) GetUsers(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		writeMeshkitError(w, ErrFetchToken(fmt.Errorf("token not found in request context")), http.StatusInternalServerError)
		return
	}

	q := req.URL.Query()
	pageSize := q.Get("pageSize")
	if pageSize == "" {
		pageSize = q.Get("pagesize")
	}

	resp, err := provider.GetUsers(token, q.Get("page"), pageSize, q.Get("search"), q.Get("order"), q.Get("filter"))
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		obj := "users data"
		h.log.Error(models.ErrEncoding(err, obj))
		writeMeshkitError(w, models.ErrEncoding(err, obj), http.StatusInternalServerError)
		return
	}
}

// UserPrefsHandler updates anonymous stats for user or for persisting load test preferences
func (h *Handler) UserPrefsHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if req.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(prefObj); err != nil {
			obj := "user preference object"
			h.log.Error(models.ErrEncoding(err, obj))
			writeMeshkitError(w, models.ErrEncoding(err, obj), http.StatusInternalServerError)
		}
		return
	}

	defer func() {
		_ = req.Body.Close()
	}()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrDecoding(err, "user preferences"))
		writeMeshkitError(w, ErrDecoding(err, "user preferences"), http.StatusInternalServerError)
		return
	}

	if req.URL.Path == "/api/identity/users/preferences" {
		mergedPref, err := preferenceFromSchemaPayload(body, prefObj)
		if err != nil {
			h.log.Error(ErrDecoding(err, "user preferences"))
			writeMeshkitError(w, ErrDecoding(err, "user preferences"), http.StatusBadRequest)
			return
		}
		prefObj = mergedPref
	} else if err := json.Unmarshal(body, &prefObj); err != nil {
		h.log.Error(ErrDecoding(err, "user preferences"))
		writeMeshkitError(w, ErrDecoding(err, "user preferences"), http.StatusBadRequest)
		return
	}

	// only validate load test data when LoadTestPreferences is send
	if prefObj.LoadTestPreferences != nil {
		// validate load test data
		qps := prefObj.LoadTestPreferences.QueriesPerSecond
		if qps < 0 {
			err := fmt.Errorf("QPS value less than 0")
			h.log.Error(ErrSavingUserPreference(err))
			writeMeshkitError(w, ErrSavingUserPreference(err), http.StatusBadRequest)
			return
		}

		dur := prefObj.LoadTestPreferences.Duration
		if _, err := time.ParseDuration(dur); err != nil {
			err = pkgerrors.Wrap(err, "unable to parse test duration")
			h.log.Error(ErrSavingUserPreference(err))
			writeMeshkitError(w, ErrSavingUserPreference(err), http.StatusBadRequest)
			return
		}

		c := prefObj.LoadTestPreferences.ConcurrentRequests
		if c < 0 {
			err := fmt.Errorf("number of concurrent requests less than 0")
			h.log.Error(ErrSavingUserPreference(err))
			writeMeshkitError(w, ErrSavingUserPreference(err), http.StatusBadRequest)
			return
		}

		loadGen := prefObj.LoadTestPreferences.LoadGenerator
		loadGenSupported := false
		for _, lg := range []models.LoadGenerator{models.FortioLG, models.Wrk2LG, models.NighthawkLG} {
			if lg.Name() == loadGen {
				loadGenSupported = true
			}
		}
		if !loadGenSupported {
			err := fmt.Errorf("invalid load generator: %s", loadGen)
			h.log.Error(ErrSavingUserPreference(err))
			writeMeshkitError(w, ErrSavingUserPreference(err), http.StatusBadRequest)
			return
		}
	}

	if err := provider.RecordPreferences(req, user.UserId, prefObj); err != nil {
		err := fmt.Errorf("unable to save user preferences: %v", err)
		h.log.Error(ErrSavingUserPreference(err))
		writeMeshkitError(w, ErrSavingUserPreference(err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(prefObj); err != nil {
		obj := "user preferences"
		h.log.Error(models.ErrEncoding(err, obj))
		writeMeshkitError(w, models.ErrEncoding(err, obj), http.StatusInternalServerError)
		return
	}
}

func preferenceFromSchemaPayload(body []byte, existing *models.Preference) (*models.Preference, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	selectedOrganizationID := ""
	for _, key := range []string{"selectedOrganizationId", "selectedOrganizationID"} {
		if selectedOrg, ok := raw[key]; ok {
			if err := json.Unmarshal(selectedOrg, &selectedOrganizationID); err != nil {
				return nil, err
			}
			if selectedOrganizationID == "all" {
				delete(raw, "selectedOrganizationId")
				delete(raw, "selectedOrganizationID")
			}
			break
		}
	}

	normalizedBody, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}

	prefUpdate := &userSchema.PreferencePayload{}
	if err := json.Unmarshal(normalizedBody, prefUpdate); err != nil {
		return nil, err
	}

	merged := *defaultPreferenceForPatch()
	if existing != nil {
		merged = *existing
	}

	if prefUpdate.MeshAdapters != nil {
		if err := convertSchemaPreferenceValue(prefUpdate.MeshAdapters, &merged.MeshAdapters); err != nil {
			return nil, err
		}
	}
	if prefUpdate.Grafana != nil {
		grafana, err := applyGrafanaPreferencePayload(merged.Grafana, prefUpdate.Grafana)
		if err != nil {
			return nil, err
		}
		merged.Grafana = grafana
	}
	if prefUpdate.Prometheus != nil {
		prometheus, err := applyPrometheusPreferencePayload(merged.Prometheus, prefUpdate.Prometheus)
		if err != nil {
			return nil, err
		}
		merged.Prometheus = prometheus
	}
	if prefUpdate.LoadTestPrefs != nil {
		merged.LoadTestPreferences = applyLoadTestPreferencesPayload(merged.LoadTestPreferences, prefUpdate.LoadTestPrefs)
	}
	if prefUpdate.AnonymousUsageStats != nil {
		merged.AnonymousUsageStats = *prefUpdate.AnonymousUsageStats
	}
	if prefUpdate.AnonymousPerfResults != nil {
		merged.AnonymousPerfResults = *prefUpdate.AnonymousPerfResults
	}
	if prefUpdate.DashboardPreferences != nil {
		merged.DashboardPreferences = *prefUpdate.DashboardPreferences
	}
	if prefUpdate.SelectedOrganizationId != nil {
		merged.SelectedOrganizationID = prefUpdate.SelectedOrganizationId.String()
	}
	if selectedOrganizationID != "" {
		merged.SelectedOrganizationID = selectedOrganizationID
	}
	if prefUpdate.SelectedWorkspaceForOrganizations != nil {
		merged.SelectedWorkspaceForOrganizations = *prefUpdate.SelectedWorkspaceForOrganizations
	}
	if prefUpdate.UsersExtensionPreferences != nil {
		merged.UsersExtensionPreferences = *prefUpdate.UsersExtensionPreferences
	}
	if prefUpdate.RemoteProviderPreferences != nil {
		merged.RemoteProviderPreferences = *prefUpdate.RemoteProviderPreferences
	}

	return &merged, nil
}

func defaultPreferenceForPatch() *models.Preference {
	return &models.Preference{
		AnonymousUsageStats:               true,
		AnonymousPerfResults:              true,
		DashboardPreferences:              map[string]interface{}{},
		SelectedWorkspaceForOrganizations: map[string]string{},
		UsersExtensionPreferences:         map[string]interface{}{},
		RemoteProviderPreferences:         map[string]interface{}{},
	}
}

func convertSchemaPreferenceValue(from interface{}, to interface{}) error {
	data, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, to)
}

func applyGrafanaPreferencePayload(existing *models.Grafana, payload *userSchema.Grafana) (*models.Grafana, error) {
	merged := models.Grafana{}
	if existing != nil {
		merged = *existing
	}
	if payload.GrafanaUrl != nil {
		merged.GrafanaURL = *payload.GrafanaUrl
	}
	if payload.GrafanaApiKey != nil {
		merged.GrafanaAPIKey = *payload.GrafanaApiKey
	}
	if payload.SelectedBoardsConfigs != nil {
		if err := convertSchemaPreferenceValue(payload.SelectedBoardsConfigs, &merged.GrafanaBoards); err != nil {
			return nil, err
		}
	}
	return &merged, nil
}

func applyPrometheusPreferencePayload(existing *models.Prometheus, payload *userSchema.Prometheus) (*models.Prometheus, error) {
	merged := models.Prometheus{}
	if existing != nil {
		merged = *existing
	}
	if payload.PrometheusUrl != nil {
		merged.PrometheusURL = *payload.PrometheusUrl
	}
	if payload.SelectedPrometheusBoardsConfigs != nil {
		if err := convertSchemaPreferenceValue(payload.SelectedPrometheusBoardsConfigs, &merged.SelectedPrometheusBoardsConfigs); err != nil {
			return nil, err
		}
	}
	return &merged, nil
}

func applyLoadTestPreferencesPayload(existing *models.LoadTestPreferences, payload *userSchema.LoadTestPreferences) *models.LoadTestPreferences {
	merged := models.LoadTestPreferences{}
	if existing != nil {
		merged = *existing
	}
	if payload.C != nil {
		merged.ConcurrentRequests = *payload.C
	}
	if payload.Qps != nil {
		merged.QueriesPerSecond = *payload.Qps
	}
	if payload.T != nil {
		merged.Duration = *payload.T
	}
	if payload.Gen != nil {
		merged.LoadGenerator = *payload.Gen
	}
	return &merged
}

func (h *Handler) ShareDesignHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	statusCode, err := provider.ShareDesign(r)
	if err != nil {
		writeMeshkitError(w, ErrShareDesign(err), statusCode)
		return
	}

	writeJSONMessage(w, map[string]string{"message": "Design shared"}, http.StatusOK)
}

func (h *Handler) ShareFilterHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	statusCode, err := provider.ShareFilter(r)
	if err != nil {
		writeMeshkitError(w, ErrShareFilter(err), statusCode)
		return
	}

	writeJSONMessage(w, map[string]string{"message": "Filter shared"}, http.StatusOK)
}
