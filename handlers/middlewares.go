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
			http.Redirect(w, req, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

func (h *Handler) validateAuth(req *http.Request) bool {
	sess, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err == nil {
		// logrus.Debugf("session: %v", sess)
		return !sess.IsNew
	}
	logrus.Errorf("session invalid, error: %v", err)
	return false
}

func (h *Handler) SessionInjectorMiddleware(next func(http.ResponseWriter, *http.Request, *sessions.Session, *models.User)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// ensuring session is intact before running load test
		session, err := h.config.SessionStore.Get(req, h.config.SessionName)
		if err != nil {
			logrus.Errorf("Error: unable to get session: %v", err)
			http.Error(w, "unable to get session", http.StatusUnauthorized)
			return
		}

		user, _ := session.Values["user"].(*models.User)

		next(w, req, session, user)
	})
}
