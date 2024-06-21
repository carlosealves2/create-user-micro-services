package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type IController interface {
	RegistryController(app *fiber.App, route, method string)
	Controller(ctx *fiber.Ctx) error
}

type BaseController struct {
}

func NewBaseController() *BaseController {
	return &BaseController{}
}

func (c *BaseController) RegistryController(app *fiber.App, route, method string, handler fiber.Handler) {
	switch strings.ToUpper(method) {
	case "GET":
		app.Get(route, handler)
	case "POST":
		app.Post(route, handler)
	}
}
