package users_postgres

import core_postgres_pool "github.com/eugeniuszglinski/golang-todoapp/internal/core/repository/postgres/pool"

type UsersRepository struct {
	pool core_postgres_pool.Pool
}

func NewUsersRepository(pool core_postgres_pool.Pool) *UsersRepository {
	return &UsersRepository{pool}
}
