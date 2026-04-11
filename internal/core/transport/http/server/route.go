package core_http_server

import (
	"net/http"

	core_http_middleware "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/middleware"
)

// Route is defined in the core package and serves as a contract that allows
// each feature to declare its own set of HTTP routes — specifying the method,
// path, and handler it exposes. These routes are consumed by HttpServer, which
// registers them with the multiplexer, enabling each feature to independently
// tell the server how to handle incoming HTTP requests that belong to it.
type Route struct {
	Method      string
	Path        string
	Handler     http.HandlerFunc
	Middlewares []core_http_middleware.Middleware
}

func NewRoute(method, path string, handler http.HandlerFunc, middlewares ...core_http_middleware.Middleware) *Route {
	return &Route{
		Method:      method,
		Path:        path,
		Handler:     handler,
		Middlewares: middlewares,
	}
}

func (r *Route) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(r.Handler, r.Middlewares...)
}
