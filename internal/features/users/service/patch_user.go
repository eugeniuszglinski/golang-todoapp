package users_service

import (
	"context"
	"fmt"

	"github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"
)

func (s *UsersService) PatchUser(ctx context.Context, userID int, userPatch *domain.UserPatch) (*domain.User, error) {
	userDomain, err := s.usersRepository.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err := userDomain.ApplyPatch(userPatch); err != nil {
		return nil, fmt.Errorf("failed to apply patch: %w", err)
	}

	userDomain, err = s.usersRepository.PatchUser(ctx, userID, userDomain)
	if err != nil {
		return nil, fmt.Errorf("failed to patch user: %w", err)
	}

	return userDomain, nil
}
