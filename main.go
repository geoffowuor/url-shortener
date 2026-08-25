package main

import (
	"log"

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

	db.AutoMigrate(&models.URL{})

	app := fiber.New()

	log.Fatal(app.Listen(":3000"))

}
