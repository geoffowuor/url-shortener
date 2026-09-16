package main

import (
	"log"

	"github.com/geoffowuor/url-shortener/internals/handlers"
	worker "github.com/geoffowuor/url-shortener/internals/workers"

	"github.com/geoffowuor/url-shortener/internals/models"

	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {

	db, err := gorm.Open(sqlite.Open("urls.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	clickQueue := worker.NewClickEventQueue(1000)

	db.AutoMigrate(&models.URL{})
	handler := &handlers.URLHandler{
		DB:         db,
		ClickQueue: clickQueue,
	}

	app := fiber.New()

	app.Post("/api/shorten", handler.CreateShortURL)
	app.Get("/api/stats/:shortCode", handler.GetURLStats)
	app.Get("/:shortCode", handler.RedirectURL)
	app.Get("/api/qrcode/:code", handler.GenerateQR)
	log.Fatal(app.Listen(":3000"))

}
