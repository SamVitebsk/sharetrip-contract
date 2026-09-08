package main

import (
	"log"

	contractapi "github.com/SamVitebsk/sharetrip-contract/internal/api"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: contractapi.ErrorHandler,
	})
	app.Use(swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: "./api/contract.yaml",
		Path:     "docs",
		Title:    "ShareTrip Contract Service API documentation",
	}))

	contractapi.RegisterRoutes(app)

	log.Fatal(app.Listen(":9190"))
}
