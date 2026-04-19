package statistics_postgres

import (
	"time"

	"github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"
)

type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func taskDomainFromModel(taskModel *TaskModel) *domain.Task {
	return domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreatedAt,
		taskModel.CompletedAt,
		taskModel.AuthorUserID,
	)
}

func taskDomainsFromModels(tasks []*TaskModel) []*domain.Task {
	taskDomains := make([]*domain.Task, len(tasks))

	for i, taskModel := range tasks {
		if taskModel == nil {
			continue
		}
		taskDomains[i] = taskDomainFromModel(taskModel)
	}

	return taskDomains
}
