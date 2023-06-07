package models

import "github.com/google/uuid"

type Blog struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Title       string    `gorm:"varchar(255);not null" json:"title"`
	SubTitle    string    `gorm:"varchar(255);not null" json:"subTitle"`
	Image       string    `gorm:"varchar(255);not null" json:"image"`
	Description string    `gorm:"varchar(255);not null" json:"description"`
	Author      string    `gorm:"varchar(255);not null" json:"author"`
	Date        string    `json:"date"`
}

type CreateBlogSchema struct {
	Title       string `json:"title" validate:"required"`
	SubTitle    string `json:"subTitle" validate:"required"`
	Description string `json:"description" validate:"required"`
	Image       string `json:"image"`
	Author      string `json:"author" validate:"required"`
	Date        string `json:"date" validate:"required"`
}

type UpdateBlogSchema struct {
	Title       string `json:"title"`
	SubTitle    string `json:"subTitle"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Author      string `json:"author"`
	Date        string `json:"date"`
}
