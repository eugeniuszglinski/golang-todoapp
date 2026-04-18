package tasks_service

import (
	"context"
	"fmt"

	"github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"
)

func (s *TasksService) PatchTask(ctx context.Context, taskID int, taskPatch *domain.TaskPatch) (*domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	if err := task.ApplyPatch(taskPatch); err != nil {
		return nil, fmt.Errorf("failed to apply patch: %w", err)
	}

	task, err = s.tasksRepository.PatchTask(ctx, taskID, task)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}
	
	return task, nil
}
