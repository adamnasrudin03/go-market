package dto

import (
	"strings"
	"time"

	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
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

type UserDetailReq struct {
	ID      uint64 `json:"id"`
	NotID   uint64 `json:"not_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Columns string `json:"columns"`
}

func (m *UserDetailReq) Validate() error {
	m.Email = strings.TrimSpace(m.Email)
	m.Name = strings.TrimSpace(m.Name)
	m.Phone = strings.TrimSpace(m.Phone)

	isNotRequired := m.ID == 0 && m.NotID == 0 && m.Email == "" && m.Name == "" && m.Phone == ""
	if isNotRequired {
		return response_mapper.NewError(response_mapper.ErrValidation, response_mapper.NewResponseMultiLang(
			response_mapper.MultiLanguages{
				ID: "Harap masukkan minimal satu parameter yang diperlukan",
				EN: "Please provide at least one required parameter",
			},
		))
	}

	return nil
}
