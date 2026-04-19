package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"
	core_logger "github.com/eugeniuszglinski/golang-todoapp/internal/core/logger"
	core_http_request "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/eugeniuszglinski/golang-todoapp/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"`
	TasksCompleted             int      `json:"tasks_completed"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time"`
}

func (h *StatisticsHttpHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(logger, w)

	userId, from, to, err := getUserIdFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user_id/from/to query params")

		return
	}

	statistics, err := h.statisticsService.GetStatistics(ctx, userId, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics")

		return
	}

	response := toDtoFromDomain(statistics)

	responseHandler.JsonResponse(response, http.StatusOK)
}

func toDtoFromDomain(statistics *domain.Statistics) *GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String()
		avgTime = &duration
	}
	return &GetStatisticsResponse{
		statistics.TasksCreated,
		statistics.TasksCompleted,
		statistics.TasksCompletedRate,
		avgTime,
	}
}

func getUserIdFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIdQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)

	userID, err := core_http_request.GetIntQueryParams(r, userIdQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	from, err := core_http_request.GetDateQueryParams(r, fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' query param: %w", err)
	}

	to, err := core_http_request.GetDateQueryParams(r, toQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' query param: %w", err)
	}

	return userID, from, to, nil
}
