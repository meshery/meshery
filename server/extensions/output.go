package extensions

import "net/http"

// Router
type Router struct {
	HTTPHandler http.Handler
	Path        string
}

// ExtensionOutput - output for a plugin
type ExtensionOutput struct {
	Router *Router
}
