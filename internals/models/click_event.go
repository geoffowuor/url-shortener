package models

import (
	"time"
)

type ClickEvent struct {
	ID         uint `gorm:"primaryKey"`
	ShortURLID uint
	IPAddress  string
	UserAgent  string
	Referer    string
	CreatedAt  time.Time
}
