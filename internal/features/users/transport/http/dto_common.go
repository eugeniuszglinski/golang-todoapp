package users_transport_http

import "github.com/eugeniuszglinski/golang-todoapp/internal/core/domain"

type UserDtoResponse struct {
	ID          int     `json:"id"           example:"1"`
	Version     int     `json:"version"      example:"1"`
	FullName    string  `json:"full_name"    example:"John Doe"`
	PhoneNumber *string `json:"phone_number" example:"+48123456789"`
}

func userDtoFromDomain(user *domain.User) *UserDtoResponse {
	return &UserDtoResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func usersDtoFromDomains(users []*domain.User) []*UserDtoResponse {
	usersDto := make([]*UserDtoResponse, len(users))

	for i, user := range users {
		usersDto[i] = userDtoFromDomain(user)
	}

	return usersDto
}
