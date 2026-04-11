package users_postgres

import (
	"context"
	"fmt"

	"github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"
)

func (r *UsersRepository) GetUsers(ctx context.Context, limit *int, offset *int) ([]*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, full_name, phone_number
	FROM todoapp.users
	ORDER BY id ASC
	LIMIT $1 OFFSET $2;
	`
	// limit and offset will not be considered by pgx library when they are nil
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query error: %w: ", err)
	}
	defer rows.Close()

	var userModels []*UserModel
	for rows.Next() {
		var userModel UserModel

		err := rows.Scan(&userModel.ID, &userModel.Version, &userModel.FullName, &userModel.PhoneNumber)
		if err != nil {
			return nil, fmt.Errorf("scan error: %w: ", err)
		}

		userModels = append(userModels, &userModel)
	}
	// pgx can optimize SELECT query and do not download all data from DB to RAM at once,
	// but retain portions of data step by step while rows.Next(). Possible error can be accessed with rows.Err()
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows error: %w: ", err)
	}

	return userDomainsFromModels(userModels), nil
}
