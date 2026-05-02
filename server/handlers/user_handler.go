package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"encoding/json"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
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

	resp, err := provider.GetUsers(token, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"))
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(w, string(resp)); err != nil {
		h.log.Error(err)
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

	// read user preferences from JSON request body
	if err := json.NewDecoder(req.Body).Decode(&prefObj); err != nil {
		h.log.Error(ErrDecoding(err, "user preferences"))
		writeMeshkitError(w, ErrDecoding(err, "user preferences"), http.StatusInternalServerError)
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
