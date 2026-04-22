package web_transport_http

import (
	core_http_server "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/server"
)

type WebHttpHandler struct {
	webService WebService
}

type WebService interface {
	GetMainPage() ([]byte, error)
}

func NewWebHttpHandler(webService WebService) *WebHttpHandler {
	return &WebHttpHandler{webService}
}

func (h *WebHttpHandler) Routes() []*core_http_server.Route {
	return []*core_http_server.Route{
		&core_http_server.Route{Path: "/", Handler: h.GetMainPage},
	}
}
