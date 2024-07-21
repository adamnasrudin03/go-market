package models

import (
	"log"

	help "github.com/adamnasrudin03/go-helpers"
	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
	"gorm.io/gorm"
)

// User represents the model for an user
type User struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"not null"`
	Phone    string `json:"phone" gorm:"not null;uniqueIndex"`
	Email    string `json:"email" gorm:"not null;uniqueIndex"`
	Password string `json:"password,omitempty" gorm:"not null"`
	DefaultModel
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPass, err := help.HashPassword(u.Password)
	if err != nil {
		log.Printf("failed hash password: %v ", err)
		return response_mapper.ErrHashPasswordFailed()
	}

	u.Password = hashedPass
	return nil
}
