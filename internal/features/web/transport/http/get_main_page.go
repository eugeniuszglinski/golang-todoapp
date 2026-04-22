package web_transport_http

import (
	"net/http"

	core_logger "github.com/eugeniuszglinski/golang-todoapp/internal/core/logger"
	core_http_response "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/response"
)

func (h *WebHttpHandler) GetMainPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(logger, w)

	html, err := h.webService.GetMainPage()
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get index.html for main page")

		return
	}

	responseHandler.HtmlResponse(html)
}
