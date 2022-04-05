package handlers

import (
	"context"

	"github.com/donnjedarko/paninaro/config"
	"github.com/donnjedarko/paninaro/internal/middleware"
	"github.com/donnjedarko/paninaro/internal/models"
	"github.com/donnjedarko/paninaro/internal/web"
	"github.com/donnjedarko/paninaro/src/domains"
	"github.com/donnjedarko/paninaro/src/dtos"
	"github.com/gofiber/fiber/v2"
)

type httpAuthHandler struct {
	authService domains.AuthService
}

func NewHttpAuthHandler(authService domains.AuthService) domains.HttpAuthHandler {
	return &httpAuthHandler{
		authService: authService,
	}
}

func (h *httpAuthHandler) Signin(c *fiber.Ctx) error {
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(c.UserContext(), cfg.AppTimeout)
	defer cancel()

	// load body from request
	body := new(dtos.UserSigninBody)
	if err := c.BodyParser(body); err != nil {
		return web.JsonResp(c, fiber.StatusBadRequest, "bad payload", nil)
	}

	// call signin service
	jwtPair, err := h.authService.Signin(ctx, body)
	if err != nil {
		return web.JsonRespError(c, err)
	}

	// set the refresh token cookie if login ok
	c.Cookie(jwtPair.RefreshToken)

	// return access token
	return web.JsonResp(c, fiber.StatusOK, "ok", jwtPair.AccessToken)
}

func (h *httpAuthHandler) Signup(c *fiber.Ctx) error {
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(c.UserContext(), cfg.AppTimeout)
	defer cancel()

	// load body from request
	body := new(dtos.UserSignupBody)
	if err := c.BodyParser(body); err != nil {
		return web.JsonResp(c, fiber.StatusBadRequest, "bad payload", nil)
	}

	// call signup service
	jwtPair, err := h.authService.Signup(ctx, body)
	if err != nil {
		return web.JsonRespError(c, err)
	}

	// set refresh token cookie if signup ok
	c.Cookie(jwtPair.RefreshToken)

	// return access token
	return web.JsonResp(c, fiber.StatusOK, "ok", jwtPair.AccessToken)
}

func (h *httpAuthHandler) Signout(c *fiber.Ctx) error {
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(c.UserContext(), cfg.AppTimeout)
	defer cancel()

	// user info from access token middleware
	userClaim := c.Locals(middleware.LocalUserClaimName).(*models.JwtAccessClaim)

	// call signout service
	err := h.authService.Signout(ctx, userClaim.Oid(), userClaim.ID)
	if err != nil {
		return web.JsonRespError(c, err)
	}

	// expire the refresh token
	c.ClearCookie(cfg.JwtRefreshTokenCookieName)

	return web.JsonRespData(c, "signout done")
}

func (h *httpAuthHandler) SignoutAll(c *fiber.Ctx) error {
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(c.UserContext(), cfg.AppTimeout)
	defer cancel()

	// user info from access token middleware
	userClaim := c.Locals(middleware.LocalUserClaimName).(*models.JwtAccessClaim)

	// call the signout all service
	err := h.authService.Signout(ctx, userClaim.Oid(), "")
	if err != nil {
		return web.JsonRespError(c, err)
	}

	// expire the refresh token
	c.ClearCookie(cfg.JwtRefreshTokenCookieName)
	return web.JsonRespData(c, "signout done")
}

func (h *httpAuthHandler) Refresh(c *fiber.Ctx) error {
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(c.UserContext(), cfg.AppTimeout)
	defer cancel()

	// get refresh token value
	refreshToken := c.Cookies(cfg.JwtRefreshTokenCookieName)

	// call refresh service
	jwtPair, err := h.authService.Refresh(ctx, refreshToken)
	if err != nil {
		return web.JsonRespError(c, err)
	}

	// set the new refresh token
	c.Cookie(jwtPair.RefreshToken)

	// return the new access token
	return web.JsonRespData(c, jwtPair.AccessToken)
}
