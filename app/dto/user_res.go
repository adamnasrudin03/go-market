package dto

import (
	"time"

	"github.com/adamnasrudin03/go-market/app/models"
)

type UserRes struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *UserRes) ConvertFromModel(input models.User) *UserRes {
	m = &UserRes{
		ID:        input.ID,
		Name:      input.Name,
		Phone:     input.Phone,
		Email:     input.Email,
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
	}

	return m
}
