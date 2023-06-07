package models

import (
	"time"

	"github.com/google/uuid"
)

type Ngo struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Name      string    `gorm:"varchar(255);not null" json:"name,omitempty"`
	Email     string    `gorm:"unique;not null" json:"email,omitempty"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"not null" json:"role,omitempty"`
	Phone     string    `gorm:"varchar(100)" json:"phone,omitempty"`
	Address   string    `gorm:"varchar(100)" json:"address,omitempty"`
	Category  []string  `gorm:"not null" json:"category,omitempty"`
	CreatedAt time.Time `gorm:"not null" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updatedAt,omitempty"`
}

type CreateNgoSchema struct {
	Name     string   `json:"name" validate:"required"`
	Email    string   `json:"email" validate:"required"`
	Password string   `json:"password" validate:"required"`
	Role     string   `json:"role" validate:"required"`
	Address  string   `json:"address,omitempty" validate:"required"`
	Category []string `json:"category,omitempty" validate:"required"`
	Phone    string   `json:"phone,omitempty"`
}

type UpdateNgoSchema struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

type LoginNgoRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
