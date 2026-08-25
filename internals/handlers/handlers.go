package handlers

import (
	"github.com/geoffowuor/url-shortener/internals/models"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"math/rand"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type URLHandler struct {
	DB *gorm.DB
}

func NewURLHandler(db *gorm.DB) *URLHandler {
	return &URLHandler{
		DB: db,
	}
}

func GenerateShortCode(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, length)

	for i := range code {
		code[i] = letters[r.Intn(len(letters))]
	}

	return string(code)
}

func (h *URLHandler) CreateShortURL(c fiber.Ctx) error {
	var url models.URL
	if err := c.Bind().Body(&url); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Internal Server error",
		})
	}

	if err := h.DB.Create(&url).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server error",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(url)
}
