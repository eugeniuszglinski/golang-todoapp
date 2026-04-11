package users_transport_http

import (
	"net/http"

	core_logger "github.com/eugeniuszglinski/golang-todoapp/internal/core/logger"
	core_http_request "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/response"
)

type DeleteUserResponse UserDtoResponse

func (h *UsersHttpHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(logger, w)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID path value")

		return
	}

	if err := h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete user")

		return
	}

	responseHandler.NoContentResponse()
}
