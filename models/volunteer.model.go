package models

import (
	"time"

	"github.com/google/uuid"
	pq "github.com/lib/pq"
)

type Volunteer struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Name      string         `gorm:"varchar(255);not null" json:"name,omitempty"`
	Email     string         `gorm:"unique;not null" json:"email,omitempty"`
	Password  string         `gorm:"not null" json:"-"`
	Role      string         `gorm:"not null" json:"role,omitempty"`
	Gender    string         `gorm:"not null" json:"gender,omitempty"`
	Phone     string         `gorm:"varchar(100)" json:"phone,omitempty"`
	Address   string         `gorm:"varchar(100)" json:"address,omitempty"`
	Category  pq.StringArray `gorm:"type:varchar(64)[]" json:"category"`
	NGO       uuid.UUID      `json:"ngo"`
	CreatedAt time.Time      `gorm:"not null" json:"createdAt,omitempty"`
	UpdatedAt time.Time      `gorm:"not null" json:"updatedAt,omitempty"`
}

type CreateVolunteerSchema struct {
	Name     string    `json:"name" validate:"required"`
	Email    string    `json:"email" validate:"required"`
	Password string    `json:"password" validate:"required"`
	Gender   string    `json:"gender" validate:"required"`
	NGO      uuid.UUID `json:"ngo"`
	Category []string  `json:"category,omitempty" validate:"required"`
	Phone    string    `json:"phone,omitempty"`
	Address  string    `json:"address,omitempty"`
}

type UpdateVolunteerSchema struct {
	Name     string    `json:"name,omitempty"`
	Password string    `json:"password,omitempty"`
	Phone    string    `json:"phone,omitempty"`
	NGO      uuid.UUID `json:"ngo"`
}

type LoginVolunteerRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
