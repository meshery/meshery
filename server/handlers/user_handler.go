package handlers

import (
	"fmt"
	"net/http"
	"time"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/layer5io/meshery/server/models"
)

// UserHandler returns info about the logged in user
func (h *Handler) UserHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, user *models.User, _ models.Provider) {
	if err := json.NewEncoder(w).Encode(user); err != nil {
		obj := "user data"
		h.log.Error(ErrEncoding(err, obj))
		http.Error(w, ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/user/profile/{id} UserAPI idGetUserByIDHandler
// Handle GET for User info by ID
//
// Returns User info
// responses:
// 	200: userInfo

func (h *Handler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	userID := mux.Vars(r)["id"]
	resp, err := provider.GetUserByID(r, userID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}

// swagger:route GET /api/identity/users UserAPI idGetAllUsersHandler
// Handles GET for all Users
//
// # Users can be further filtered through query parameters
//
// ```?order={field}``` orders on the passed field
//
// ```?page={page-number}``` Default page number is 0
//
// ```?pagesize={pagesize}``` Default pagesize is 20
//
// ```?search={username|email|first_name|last_name}``` If search is non empty then a greedy search is performed
//
// ```?filter={condition}```
// responses:
// 	200: users

func (h *Handler) GetUsers(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}

	q := req.URL.Query()

	resp, err := provider.GetUsers(token, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"))
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}

// swagger:route GET /api/identity/users/keys UserKeysAPI idGetAllUsersKeysHandler
// Handles GET for all Keys for users
//
// Returns all keys for user
// responses:
// 	200: keys

func (h *Handler) GetUsersKeys(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}

	q := req.URL.Query()
	resp, err := provider.GetUsersKeys(token, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"))
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}

// swagger:route GET /api/user/prefs UserAPI idGetUserTestPrefs
// Handle GET Requests for User Load Test Preferences
//
// Returns User Load Test Preferences
// responses:
// 	200: userLoadTestPrefsRespWrapper

// swagger:route POST /api/user/prefs UserAPI idPostUserTestPrefs
// Handle POST Requests for User Load Test Preferences
//
// Updates User Load Test Preferences
// responses:
// 	200: userLoadTestPrefsRespWrapper

// UserPrefsHandler updates anonymous stats for user or for persisting load test preferences
func (h *Handler) UserPrefsHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if req.Method == http.MethodGet {
		if err := json.NewEncoder(w).Encode(prefObj); err != nil {
			obj := "user preference object"
			h.log.Error(ErrEncoding(err, obj))
			http.Error(w, ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
		}
		return
	}

	defer func() {
		_ = req.Body.Close()
	}()

	// read user preferences from JSON request body
	if err := json.NewDecoder(req.Body).Decode(&prefObj); err != nil {
		h.log.Error(ErrDecoding(err, "user preferences"))
		http.Error(w, ErrDecoding(err, "user preferences").Error(), http.StatusInternalServerError)
		return
	}

	// only validate load test data when LoadTestPreferences is send
	if prefObj.LoadTestPreferences != nil {
		// validate load test data
		qps := prefObj.LoadTestPreferences.QueriesPerSecond
		if qps < 0 {
			w.WriteHeader(http.StatusBadRequest)
			err := fmt.Errorf("QPS value less than 0")
			h.log.Error(ErrSavingUserPreference(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dur := prefObj.LoadTestPreferences.Duration
		if _, err := time.ParseDuration(dur); err != nil {
			err = errors.Wrap(err, "unable to parse test duration")
			h.log.Error(ErrSavingUserPreference(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		c := prefObj.LoadTestPreferences.ConcurrentRequests
		if c < 0 {
			err := fmt.Errorf("number of concurrent requests less than 0")
			h.log.Error(ErrSavingUserPreference(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := provider.RecordPreferences(req, user.UserID, prefObj); err != nil {
		err := fmt.Errorf("unable to save user preferences: %v", err)
		h.log.Error(ErrSavingUserPreference(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(prefObj); err != nil {
		obj := "user preferences"
		h.log.Error(ErrEncoding(err, obj))
		http.Error(w, ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route POST /api/content/design/share ShareContent idPostShareContent
// Handle POST request for Sharing content
//
// Used to share designs with others
// responses:
// 	200:
//  403:
//  500:

func (h *Handler) ShareDesignHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	statusCode, err := provider.ShareDesign(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err.Error()), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "Design shared")
}

// swagger:route POST /api/content/filter/share ShareContent idPostShareContent
// Handle POST request for Sharing content
//
// Used to share filters with others
// responses:
// 	200:
//  403:
//  500:

func (h *Handler) ShareFilterHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	statusCode, err := provider.ShareDesign(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err.Error()), statusCode)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "Filter shared")
}
