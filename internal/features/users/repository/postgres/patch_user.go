package users_postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"
	core_errors "github.com/eugeniuszglinski/golang-todoapp/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) PatchUser(ctx context.Context, userID int, user *domain.User) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE todoapp.users
	SET full_name=$1, phone_number=$2, version=version+1
	WHERE id=$3 AND version = $4
	RETURNING id, version, full_name, phone_number
	`

	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber, userID, user.Version)

	var userModel UserModel
	if err := row.Scan(&userModel.ID, &userModel.Version, &userModel.FullName, &userModel.PhoneNumber); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(
				"user with id='%d' concurrently accessed: %w", userID, core_errors.ErrConflict,
			)
		}

		return nil, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(userModel.ID, userModel.Version, userModel.FullName, userModel.PhoneNumber)

	return userDomain, nil
}
