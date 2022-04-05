package handlers

import (
	"context"

	"github.com/donnjedarko/paninaro/config"
	"github.com/donnjedarko/paninaro/internal/middleware"
	"github.com/donnjedarko/paninaro/internal/models"
	"github.com/donnjedarko/paninaro/internal/web"
	"github.com/donnjedarko/paninaro/src/domains"
	"github.com/gofiber/fiber/v2"
)

type httpUserHandler struct {
	userService domains.UserService
}

func NewHttpUserHandler(userService domains.UserService) domains.HttpUserHandler {
	return &httpUserHandler{
		userService: userService,
	}
}

func (h *httpUserHandler) Me(c *fiber.Ctx) error {
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(c.UserContext(), cfg.AppTimeout)
	defer cancel()

	// get access token claim
	userClaim := c.Locals(middleware.LocalUserClaimName).(*models.JwtAccessClaim)

	// call me service
	resp, err := h.userService.Me(ctx, userClaim.Oid())
	if err != nil {
		return web.JsonRespError(c, err)
	}

	return web.JsonResp(c, fiber.StatusOK, "ok", resp)
}

func (h *httpUserHandler) All(c *fiber.Ctx) error {
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(c.UserContext(), cfg.AppTimeout)
	defer cancel()

	// call all service
	users, err := h.userService.All(ctx)
	if err != nil {
		return web.JsonRespError(c, err)
	}

	// return user list
	return web.JsonRespData(c, users)
}
