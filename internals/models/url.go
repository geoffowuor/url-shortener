package models

import (
	"time"

	"gorm.io/gorm"
)

type URL struct {
	gorm.Model

	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code" gorm:"unique;not null"`
	Clicks      int    `json:"clicks"`
}

type ClickEvent struct {
	ShortURLID uint
	IPAddress  string
	UserAgent  string
	Referer    string
	CreatedAt  time.Time
}
