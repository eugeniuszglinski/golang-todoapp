package users_service

import (
	"context"
	"fmt"

	"github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"
)

func (s *UsersService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("user validation failed: %w", err)
	}

	user, err := s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user failed: %w", err)
	}

	return user, nil
}
