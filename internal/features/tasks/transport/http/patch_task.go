package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"
	core_logger "github.com/eugeniuszglinski/golang-todoapp/internal/core/logger"
	core_http_request "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"        swaggertype:"string" example:"Cleaning"`
	Description core_http_types.Nullable[string] `json:"description"  swaggertype:"string" example:"null"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"    swaggertype:"boolean"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("'title' can't be null")
		}

		titleLength := len([]rune(*r.Title.Value))
		if titleLength < 1 || titleLength > 100 {
			return fmt.Errorf("'title' must be between 1 and 100")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLength := len([]rune(*r.Description.Value))
			if descriptionLength < 1 || descriptionLength > 1000 {
				return fmt.Errorf("'description' must be between 1 and 1000")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("'completed' can't be null")
		}
	}

	return nil
}

type PatchTaskResponse TaskDtoResponse

// PatchTask    godoc
// @Summary     Patch task
// @Description Change existing task
// @Description ### Logic for patching task (Three-state logic)
// @Description 1. **field not provided**: `description` will be ignored and unchanged in database
// @Description 2. **explicitly provided value**: `"description": "Do cleaning"` will update 'description' in database
// @Description 3. **null value**: `"description": null` will remove 'description' from database (set to NULL)
// @Description Constants: `title` and `completed` cannot be null
// @Tags        tasks
// @Accept      application/json
// @Produce     application/json
// @Param       id path int true "ID of task to patch"
// @Param       request body PatchTaskRequest true "PatchTask request body"
// @Success     200 {object} PatchTaskResponse "Changed task details"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404 {object} core_http_response.ErrorResponse "Task not found"
// @Failure     409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /tasks/{id} [patch]
func (h *TasksHttpHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(logger, w)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task_id path value")

		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")

		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")

		return
	}

	response := PatchTaskResponse(*taskDtoFromDomain(taskDomain))

	responseHandler.JsonResponse(response, http.StatusOK)
}

func taskPatchFromRequest(request PatchTaskRequest) *domain.TaskPatch {
	return domain.NewTaskPatch(request.Title.ToDomain(), request.Description.ToDomain(), request.Completed.ToDomain())
}
