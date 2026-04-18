package tasks_postgres

import (
	"context"
	"fmt"

	core_errors "github.com/eugeniuszglinski/golang-todoapp/internal/core/errors"
)

func (r *TasksRepository) DeleteTask(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `DELETE FROM todoapp.tasks WHERE id = $1;`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to execute delete task command: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("task with ID='%d': %w", id, core_errors.ErrNotFound)
	}

	return nil
}
