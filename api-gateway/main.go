package main

import (
	"log"

	"github.com/carlosealves2/create-user-micro-services/api-gateway/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	_ "github.com/carlosealves2/create-user-micro-services/api-gateway/docs"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	createUser := controllers.NewCreateUserController()
	createUser.RegistryController(app, "user/", fiber.MethodPost, createUser.Controller)

	log.Fatalln(app.Listen(":3000"))
}
