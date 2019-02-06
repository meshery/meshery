package handlers

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if !h.config.ByPassAuth {
			isValid := h.validateAuth(req)
			logrus.Debugf("validate auth: %t", isValid)
			if !isValid {
				http.Redirect(w, req, "/play/login", http.StatusFound)
				return
			}
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

func (h *Handler) validateAuth(req *http.Request) bool {
	sess, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err == nil {
		logrus.Debugf("session: %v", sess)
		return !sess.IsNew
	}
	logrus.Errorf("session invalid, error: %v", err)
	return false
}
