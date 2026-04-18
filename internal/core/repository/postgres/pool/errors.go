package core_postgres_pool

import "errors"

var (
	ErrNoRows             = errors.New("no rows returned")
	ErrViolatesForeignKey = errors.New("violated foreign key")
	ErrUnknown            = errors.New("unknown error")
)
