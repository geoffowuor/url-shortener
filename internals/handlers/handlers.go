package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/geoffowuor/url-shortener/internals/models"
	"github.com/geoffowuor/url-shortener/internals/utils"
	"github.com/go-chi/chi/v5"
	"github.com/gofiber/fiber/v3"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type URLHandler struct {
	DB *gorm.DB
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

	fmt.Println("SHORT CODE:", shortCode)

	var url models.URL

	if err := h.DB.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "URL not found",
		})
	}

	fmt.Println("FOUND URL:", url.OriginalURL)

	url.Clicks++

	if err := h.DB.Save(&url).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update clicks",
		})
	}

	fmt.Println("REDIRECTING TO:", url.OriginalURL)

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

	return c.JSON(fiber.Map{
		"original_url": url.OriginalURL,
		"short_code":   url.ShortCode,
		"clicks":       url.Clicks,
	})
}

func (h *URLHandler) GenerateQR(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	shortURL := "http://127.0.0.1:3000/" + code

	png, err := qrcode.Encode(
		shortURL,
		qrcode.Medium,
		512,
	)
	if err != nil {
		http.Error(w, "failed to generate QR code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(png)
}
