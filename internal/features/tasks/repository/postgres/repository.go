package tasks_postgres

import core_postgres_pool "github.com/eugeniuszglinski/golang-todoapp/internal/core/repository/postgres/pool"

type TasksRepository struct {
	pool core_postgres_pool.Pool
}

func NewTasksRepository(pool core_postgres_pool.Pool) *TasksRepository {
	return &TasksRepository{pool}
}
