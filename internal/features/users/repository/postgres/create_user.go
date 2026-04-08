package users_postgres

import (
	"context"
	"fmt"

	"github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"
)

func (r *UsersRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO todoapp.users (full_name, phone_number)
	VALUES ($1, $2)
	RETURNING id, version, full_name, phone_number;
	`

	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber)

	// row is translated to user model (not a domain.User) to avoid coupling the model with the domain entity
	var userModel UserModel
	err := row.Scan(&userModel.ID, &userModel.Version, &userModel.FullName, &userModel.PhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("scan row error: %w", err)
	}

	userDomain := domain.NewUser(userModel.ID, userModel.Version, userModel.FullName, userModel.PhoneNumber)

	return userDomain, nil
}
