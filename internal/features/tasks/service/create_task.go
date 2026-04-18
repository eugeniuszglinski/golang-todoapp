package tasks_service

import (
	"context"
	"fmt"

	"github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"
)

func (s *TasksService) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	if err := task.Validate(); err != nil {
		return nil, fmt.Errorf("task validation failed: %w", err)
	}

	newTask, err := s.tasksRepository.CreateTask(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return newTask, nil
}
