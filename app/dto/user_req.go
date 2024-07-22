package dto

import (
	"strings"

	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
)

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
