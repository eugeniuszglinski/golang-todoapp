package tasks_service

import (
	"context"
	"fmt"

	"github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"
)

func (s *TasksService) GetTask(ctx context.Context, id int) (*domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task from repository: %w", err)
	}

	return task, nil
}
