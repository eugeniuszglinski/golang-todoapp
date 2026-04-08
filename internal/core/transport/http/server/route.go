package core_http_server

import "net/http"

// Route is defined in the core package and serves as a contract that allows
// each feature to declare its own set of HTTP routes — specifying the method,
// path, and handler it exposes. These routes are consumed by HttpServer, which
// registers them with the multiplexer, enabling each feature to independently
// tell the server how to handle incoming HTTP requests that belong to it.
type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

func NewRoute(method, path string, handler http.HandlerFunc) *Route {
	return &Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}
