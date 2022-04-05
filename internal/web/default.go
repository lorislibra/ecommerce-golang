package web

import "github.com/gofiber/fiber/v2"

func NotImplementedHandler(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}
