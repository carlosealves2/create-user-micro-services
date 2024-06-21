package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type CreateUserController struct {
	*BaseController
}

func NewCreateUserController() *CreateUserController {
	return &CreateUserController{
		BaseController: NewBaseController(),
	}
}

func (c *CreateUserController) Controller(ctx *fiber.Ctx) error {
	ctx.Status(fiber.StatusCreated).WriteString("user created")
	return nil
}
