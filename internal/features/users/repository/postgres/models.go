package users_postgres

import "github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"

type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func userDomainsFromModels(users []*UserModel) []*domain.User {
	userDomains := make([]*domain.User, len(users))

	for i, userModel := range users {
		if userModel == nil {
			continue
		}
		userDomains[i] = domain.NewUser(userModel.ID, userModel.Version, userModel.FullName, userModel.PhoneNumber)
	}

	return userDomains
}
