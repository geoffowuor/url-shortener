package models

import (
	

	"gorm.io/gorm"
)

type URL struct {
	gorm.Model

	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code" gorm:"unique;not null"`
	Clicks      int    `json:"clicks"`
}


