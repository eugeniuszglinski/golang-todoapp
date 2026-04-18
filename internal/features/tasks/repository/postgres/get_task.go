package tasks_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"
	core_errors "github.com/eugeniuszglinski/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/eugeniuszglinski/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) GetTask(ctx context.Context, id int) (*domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
	FROM todoapp.tasks
	WHERE id = $1;
	`
	row := r.pool.QueryRow(ctx, query, id)

	var taskModel TaskModel

	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserID,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return nil, fmt.Errorf("task with id='%d': %w", id, core_errors.ErrNotFound)
		}

		return nil, fmt.Errorf("failed to scan task from row: %w", err)
	}

	taskDomain := taskDomainFromModel(&taskModel)

	return taskDomain, nil
}
