package handlers

import (
	"github.com/geoffowuor/url-shortener/internals/models"
	"github.com/gofiber/fiber/v3"
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
