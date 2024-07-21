package models

import "time"

type DefaultModel struct {
	CreatedBy uint64    `json:"created_by,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedBy uint64    `json:"updated_by,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
}
