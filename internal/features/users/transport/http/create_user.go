package users_transport_http

import (
	"net/http"

	"github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"
	core_logger "github.com/eugeniuszglinski/golang-todoapp/internal/core/logger"
	core_http_request "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/response"
)

// CreateUserRequest is a dto object for user creation
type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"`
}

// CreateUserResponse is a dto object for user creation response
type CreateUserResponse UserDtoResponse

func (h *UsersHttpHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(logger, rw)

	logger.Debug("handling CreateUser request")

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")

		return
	}

	userDomain := domainFromDto(&request)

	createdUserDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")

		return
	}

	response := CreateUserResponse(*userDtoFromDomain(createdUserDomain))

	responseHandler.JsonResponse(response, http.StatusCreated)
}

func domainFromDto(dto *CreateUserRequest) *domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
