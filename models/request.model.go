package models

import (
	"time"

	"github.com/google/uuid"
)

// Request represents a request made by an elderly user
type Request struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	UserID      uuid.UUID `json:"userId"`
	Title       string    `json:"title"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	IsFulfilled bool      `json:"is_fulfilled"`
	NGO         Ngo       `json:"ngo"`
}
type CreateRequestSchema struct {
	UserID      uuid.UUID `json:"userId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	IsFulfilled bool      `json:"is_fulfilled"`
	NGO         Ngo       `json:"ngo"`
}

type UpdateRequestSchema struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	IsFulfilled bool   `json:"is_fulfilled"`
	NGO         Ngo    `json:"ngo"`
}
