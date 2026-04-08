package core_http_response

import "net/http"

var (
	statusCodeUninitialized = -1
)

// ResponseWriter is a wrapper around http.ResponseWriter that captures the status code of the response.
type ResponseWriter struct {
	http.ResponseWriter

	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, statusCode: statusCodeUninitialized}
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (w *ResponseWriter) GetStatusCodeOrPanic() int {
	if w.statusCode == statusCodeUninitialized {
		panic("uninitialized status code")
	}

	return w.statusCode
}
