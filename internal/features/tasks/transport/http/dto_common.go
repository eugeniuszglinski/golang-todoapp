package tasks_transport_http

import (
	"time"

	"github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"
)

type TaskDtoResponse struct {
	ID           int        `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	Completed    bool       `json:"completed"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	AuthorUserID int        `json:"author_user_id"`
}

func taskDtoFromDomain(task *domain.Task) *TaskDtoResponse {
	return &TaskDtoResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func tasksDtoFromDomains(tasks []*domain.Task) []*TaskDtoResponse {
	tasksDto := make([]*TaskDtoResponse, len(tasks))

	for i, task := range tasks {
		tasksDto[i] = taskDtoFromDomain(task)
	}

	return tasksDto
}
