package core_http_middleware

import "net/http"

// Middleware is a function that takes an http.Handler and returns a new http.Handler that wraps the original one.
type Middleware func(http.Handler) http.Handler

func ChainMiddleware(h http.Handler, middlewares ...Middleware) http.Handler {
	if len(middlewares) == 0 {
		return h
	}

	// Wrap the handler with the middlewares in reverse order
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}
