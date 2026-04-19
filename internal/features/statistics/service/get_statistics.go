package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"
	core_errors "github.com/eugeniuszglinski/golang-todoapp/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context, userID *int, from *time.Time, to *time.Time,
) (*domain.Statistics, error) {
	if from != nil && to != nil {
		if from.After(*to) || from.Equal(*to) {
			return nil, fmt.Errorf(
				"'from' date must be before or equal to 'to' date: %w", core_errors.ErrInvalidArgument,
			)
		}
	}

	tasks, err := s.statisticsRepository.GetTasks(ctx, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks from repository: %w", err)
	}

	statistics := calcStatistics(tasks)

	return statistics, nil
}

func calcStatistics(tasks []*domain.Task) *domain.Statistics {
	if len(tasks) == 0 {
		return domain.NewStatistics(0, 0, nil, nil)
	}

	tasksCreated := len(tasks)

	tasksCompleted := 0
	var totalCompletedDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}

		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			totalCompletedDuration += *completionDuration
		}
	}

	tasksCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var tasksAverageCompletionTime *time.Duration
	if tasksCompleted > 0 && totalCompletedDuration != 0 {
		avg := totalCompletedDuration / time.Duration(tasksCompleted)

		tasksAverageCompletionTime = &avg
	}

	return domain.NewStatistics(tasksCreated, tasksCompleted, &tasksCompletedRate, tasksAverageCompletionTime)
}
