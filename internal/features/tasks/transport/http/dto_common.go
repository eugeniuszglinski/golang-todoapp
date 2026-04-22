package tasks_transport_http

import (
	"time"

	"github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"
)

type TaskDtoResponse struct {
	ID           int        `json:"id"             example:"1"`
	Version      int        `json:"version"        example:"1"`
	Title        string     `json:"title"          example:"Homework"`
	Description  *string    `json:"description"    example:"Do it before Friday"`
	Completed    bool       `json:"completed"      example:"false"`
	CreatedAt    time.Time  `json:"created_at"     example:"2026-01-27T11:30:00Z"`
	CompletedAt  *time.Time `json:"completed_at"   example:"null"`
	AuthorUserID int        `json:"author_user_id" example:"1"`
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
