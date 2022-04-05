package dtos

import (
	"time"

	"github.com/donnjedarko/paninaro/src/entities"
)

type UserMeResp struct {
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstname"`
	LastName  string    `json:"lastname"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func UserMeRespFromEntity(e *entities.User) *UserMeResp {
	return &UserMeResp{
		Username:  e.Username,
		FirstName: e.FirstName,
		LastName:  e.LastName,
		Email:     e.Email,
		Role:      e.Role.String(),
		CreatedAt: e.CreatedAt,
	}
}
