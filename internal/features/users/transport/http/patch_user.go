package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"
	core_logger "github.com/eugeniuszglinski/golang-todoapp/internal/core/logger"
	core_http_request "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"    swaggertype:"string" example:"Jacek Sutryk"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+48987654321"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("full_name cannot be null")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("full_name must be between 3 and 10")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("phone_number must be between 10 and 15 symbols")
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("phone_number must start with +")
			}
		}
	}
	return nil
}

type PatchUserResponse UserDtoResponse

// PatchUser    godoc
// @Summary     Patch user
// @Description Change existing user
// @Description ### Logic for patching user (Three-state logic)
// @Description 1. **field not provided**: `phone_number` will be ignored and unchanged in database
// @Description 2. **explicitly provided value**: `"phone_number": "+48987654321"` will update 'phone_number' in database
// @Description 3. **null value**: `"phone_number": null` will remove 'phone_number' from database (set to NULL)
// @Description Constants: `full_name` cannot be null
// @Tags        users
// @Accept      application/json
// @Produce     application/json
// @Param       id path int true "ID of user to patch"
// @Param       request body PatchUserRequest true "PatchUser request body"
// @Success     200 {object} PatchUserResponse "Changed user details"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "User not found"
// @Failure     409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /users/{id} [patch]
func (h *UsersHttpHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(logger, w)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID path value")

		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")

		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")

		return
	}

	response := PatchUserResponse(*userDtoFromDomain(userDomain))

	responseHandler.JsonResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) *domain.UserPatch {
	return domain.NewUserPatch(request.FullName.ToDomain(), request.PhoneNumber.ToDomain())
}
