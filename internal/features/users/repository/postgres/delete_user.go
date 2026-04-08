package users_postgres

import (
	"context"
	"fmt"

	core_errors "github.com/eugeniuszglinski/golang-todoapp/internal/core/errors"
)

func (r *UsersRepository) DeleteUser(ctx context.Context, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE FROM todoapp.users WHERE id = $1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to execute delete user command: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with ID='%d': %w", userID, core_errors.ErrNotFound)
	}

	return nil
}
