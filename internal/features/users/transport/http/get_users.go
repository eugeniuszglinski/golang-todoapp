package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/eugeniuszglinski/golang-todoapp/internal/core/logger"
	core_http_response "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/response"
	core_http_utils "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/utils"
)

type GetUsersResponse []*UserDtoResponse

func (h *UsersHttpHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(logger, rw)

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'limit'/'offset' query param")

		return
	}

	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")

		return
	}

	response := GetUsersResponse(usersDtoFromDomains(userDomains))

	responseHandler.JsonResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	limit, err := core_http_utils.GetIntQueryParams(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param error: %w", err)
	}

	offset, err := core_http_utils.GetIntQueryParams(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param error: %w", err)
	}

	return limit, offset, nil
}
