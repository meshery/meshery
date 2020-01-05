package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/layer5io/meshery/models"
	"github.com/sirupsen/logrus"
)

// AuthMiddleware is a middleware to validate if a user is authenticated
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		isValid := h.validateAuth(req)
		// logrus.Debugf("validate auth: %t", isValid)
		if !isValid {
			// if h.GetProviderType() == models.CloudProviderType {
			// 	http.Redirect(w, req, "/login", http.StatusFound)
			// } else { // Local Provider
			// 	h.LoginHandler(w, req)
			// }
			// return
			if h.GetProviderType() == models.CloudProviderType {
				http.Redirect(w, req, "/login", http.StatusFound)
				return
			}
			// Local Provider
			h.LoginHandler(w, req, true)
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

func (h *Handler) validateAuth(req *http.Request) bool {
	sess, err := h.config.Provider.GetSession(req)
	if err == nil {
		// logrus.Debugf("session: %v", sess)
		return !sess.IsNew
	}
	logrus.Errorf("session invalid, error: %v", err)
	return false
}

// SessionInjectorMiddleware - is a middleware which injects user and session object
func (h *Handler) SessionInjectorMiddleware(next func(http.ResponseWriter, *http.Request, *sessions.Session, *models.Preference, *models.User)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// ensuring session is intact before running load test
		session, err := h.config.Provider.GetSession(req)
		if err != nil {
			logrus.Errorf("Error: unable to get session: %v", err)
			http.Error(w, "unable to get session", http.StatusUnauthorized)
			return
		}

		user, _ := h.config.Provider.GetUserDetails(req)

		prefObj, err := h.config.Provider.ReadFromPersister(user.UserID)
		if err != nil {
			logrus.Warn("unable to read session from the session persister, starting with a new one")
		}

		if prefObj == nil {
			prefObj = &models.Preference{
				AnonymousUsageStats:  true,
				AnonymousPerfResults: true,
			}
		}

		next(w, req, session, prefObj, user)
	})
}
