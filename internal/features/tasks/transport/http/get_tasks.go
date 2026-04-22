package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/eugeniuszglinski/golang-todoapp/internal/core/logger"
	core_http_request "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/response"
)

type GetTasksResponse []*TaskDtoResponse

// GetTasks     godoc
// @Summary     Get tasks list
// @Description Get tasks list with optional pagination and/or filtering by author ID
// @Tags        tasks
// @Produce     application/json
// @Param       user_id query int false "Filtering by author ID"
// @Param       limit query int false "Limit of tasks to return"
// @Param       offset query int false "Offset of returned tasks"
// @Success     200 {object} GetTasksResponse "List of tasks"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "Task not found"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /tasks [get]
func (h *TasksHttpHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(logger, w)

	userID, limit, offset, err := getUserIdLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user_id/limit/offset query params")

		return
	}

	tasksDomains, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tasks")

		return
	}

	response := GetTasksResponse(tasksDtoFromDomains(tasksDomains))

	responseHandler.JsonResponse(response, http.StatusOK)
}

func getUserIdLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		userIdQueryParamKey = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userID, err := core_http_request.GetIntQueryParams(r, userIdQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param error: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParams(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param error: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParams(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param error: %w", err)
	}

	return userID, limit, offset, nil
}
