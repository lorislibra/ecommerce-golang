package web

import (
	"github.com/gofiber/fiber/v2"
)

type JsonResponseBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JsonResp(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	resp := JsonResponseBody{
		Code:    statusCode,
		Message: message,
		Data:    data,
	}
	return c.Status(statusCode).JSON(resp)
}

func JsonRespError(c *fiber.Ctx, err error) error {
	webErr, ok := err.(*WebError)
	if ok {
		return JsonResp(c, webErr.Status, webErr.Error(), nil)
	}
	return JsonResp(c, fiber.StatusNotFound, err.Error(), nil)
}

func JsonRespData(c *fiber.Ctx, data interface{}) error {
	return JsonResp(c, fiber.StatusOK, "ok", data)
}
