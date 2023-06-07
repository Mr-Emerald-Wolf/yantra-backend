package models

import "github.com/google/uuid"

type Event struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Title       *string   `gorm:"varchar(255);not null" json:"title"`
	SubTitle    *string   `gorm:"varchar(255);not null" json:"subTitle"`
	Description *string   `gorm:"varchar(255);not null" json:"description"`
	Location    *string   `gorm:"varchar(255);not null" json:"location"`
	Date        *string   `json:"date"`
	StartTime   *string   `json:"startTime"`
	EndTime     *string   `json:"endTime"`
	NGO         uuid.UUID `json:"ngo"`
}

type CreateEventSchema struct {
	Title       *string   `json:"title" validate:"required"`
	SubTitle    *string   `json:"subTitle" validate:"required"`
	Description *string   `json:"description" validate:"required"`
	Location    *string   `json:"location" validate:"required"`
	Date        *string   `json:"date" validate:"required"`
	StartTime   *string   `json:"startTime" validate:"required"`
	EndTime     *string   `json:"endTime" validate:"required"`
	NGO         uuid.UUID `json:"ngo" validate:"required"`
}

type UpdateEventSchema struct {
	Title       *string `json:"title"`
	SubTitle    *string `json:"subTitle"`
	Description *string `json:"description"`
	Location    *string `json:"location"`
	Date        *string `json:"date"`
	StartTime   *string `json:"startTime"`
	EndTime     *string `json:"endTime"`
}
