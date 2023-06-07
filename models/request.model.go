package models

import (
	"time"

	"github.com/google/uuid"
)

// Request represents a request made by an elderly user
type Request struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	UserID      uuid.UUID `gorm:"not null" json:"userId"`
	Title       string    `gorm:"varchar(255);not null" json:"title"`
	Category    string    `gorm:"varchar(255);not null" json:"category"`
	Description string    `gorm:"varchar(255);not null" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	IsFulfilled bool      `gorm:"not null" json:"isfulfilled"`
	VolunteerID uuid.UUID `gorm:"not null" json:"volId"`
	NGO         uuid.UUID `json:"ngo"`
}
type CreateRequestSchema struct {
	UserID      uuid.UUID `json:"userId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	VolunteerID uuid.UUID `json:"volId"`
	NGO         uuid.UUID `json:"ngo"`
}

type UpdateRequestSchema struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	IsFulfilled bool      `json:"isfulfilled"`
	NGO         uuid.UUID `json:"ngo"`
}
