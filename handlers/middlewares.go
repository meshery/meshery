package handlers

import (
	"net/http"
	"os"
)

var byPassAuth = (os.Getenv("BYPASS_AUTH") == "true")

func authMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if !byPassAuth && !validateAuth(req) {
			http.Redirect(w, req, "/play/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}
