package handlers

import (
	"errors"
	"fmt"
	"time"

	"github.com/geoffowuor/url-shortener/internals/models"
	"github.com/geoffowuor/url-shortener/internals/utils"
	worker "github.com/geoffowuor/url-shortener/internals/workers"
	"github.com/gofiber/fiber/v3"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type URLHandler struct {
	DB         *gorm.DB
	ClickQueue *worker.ClickEventQueue
}

func NewURLHandler(db *gorm.DB) *URLHandler {
	return &URLHandler{
		DB: db,
	}
}

func (h *URLHandler) CreateShortURL(c fiber.Ctx) error {
	var url models.URL

	if err := c.Bind().Body(&url); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	for {
		url.ShortCode = utils.GenerateShortCode(6)

		var existing models.URL

		result := h.DB.Where("short_code = ?", url.ShortCode).First(&existing)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			break
		}

		if result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate short URL",
			})
		}
	}

	if err := h.DB.Create(&url).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create short URL",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(url)
}
func (h *URLHandler) RedirectURL(c fiber.Ctx) error {
	shortCode := c.Params("shortCode")

	var url models.URL

	if err := h.DB.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "URL not found",
		})
	}

	url.Clicks++

	if err := h.DB.Save(&url).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update clicks",
		})
	}

	// Analytics event
	event := models.ClickEvent{
		ShortURLID: url.ID,
		IPAddress:  c.IP(),
		UserAgent:  c.Get("User-Agent"),
		Referer:    c.Get("Referer"),
		CreatedAt:  time.Now(),
	}

	fmt.Println("QUEUEING EVENT:", event)

	select {
	case h.ClickQueue.Events <- event:
		fmt.Println("EVENT QUEUED")
	default:
		fmt.Println("EVENT QUEUE FULL")
	}

	return c.Redirect().Status(fiber.StatusFound).To(url.OriginalURL)
}
func (h *URLHandler) GetURLStats(c fiber.Ctx) error {
	shortCode := c.Params("shortCode")

	var url models.URL

	if err := h.DB.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "URL not found",
		})
	}

	var events []models.ClickEvent

	if err := h.DB.
		Where("short_url_id = ?", url.ID).
		Order("created_at DESC").
		Find(&events).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch click analytics",
		})
	}

	uniqueIPs := make(map[string]struct{})

	clicksByDate := make(map[string]int)

	referrers := make(map[string]int)

	for _, event := range events {
		uniqueIPs[event.IPAddress] = struct{}{}

		date := event.CreatedAt.Format("2006-01-02")
		clicksByDate[date]++

		referer := event.Referer
		if referer == "" {
			referer = "direct"
		}

		referrers[referer]++
	}

	return c.JSON(fiber.Map{
		"short_code":   url.ShortCode,
		"original_url": url.OriginalURL,

		"summary": fiber.Map{
			"total_clicks":    len(events),
			"unique_visitors": len(uniqueIPs),
		},

		"clicks_by_date": clicksByDate,

		"referrers": referrers,

		"recent_clicks": events,
	})
}

func (h *URLHandler) GenerateQR(c fiber.Ctx) error {
	shortCode := c.Params("code")

	shortURL := "http://127.0.0.1:3000/" + shortCode

	png, err := qrcode.Encode(
		shortURL,
		qrcode.Medium,
		512,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate QR code",
		})
	}

	c.Set("Content-Type", "image/png")
	c.Status(fiber.StatusOK)
	c.Write(png)
	return nil
}
