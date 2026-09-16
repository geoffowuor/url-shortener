package models

import (
	"time"
	
)

type ClickEvent struct {
	ShortURLID uint
	IPAddress  string
	UserAgent  string
	Referer    string
	CreatedAt  time.Time
}
