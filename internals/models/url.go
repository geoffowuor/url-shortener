package models

import (
	"gorm.io/gorm"
)

type URL struct {
	gorm.Model
	ID          string `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
	Clicks      int    `json:"clicks"`
}
