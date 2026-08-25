package models

type URL struct {
	ID          string `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
	Clicks      int    `json:"clicks"`
}
