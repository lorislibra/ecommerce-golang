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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type httpOrderHandler struct {
	orderService domains.OrderService
}

func NewHttpOrderHandler(orderService domains.OrderService) domains.HttpOrderHandler {
	return &httpOrderHandler{
		orderService: orderService,
	}
}

func (h *httpOrderHandler) Create(c *fiber.Ctx) error {
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(c.UserContext(), cfg.AppTimeout)
	defer cancel()

	// load body from request
	body := new(dtos.OrderCreateBody)
	if err := c.BodyParser(body); err != nil {
		return web.JsonResp(c, fiber.StatusBadRequest, "bad payload", nil)
	}

	// get access token claim
	userClaim := c.Locals(middleware.LocalUserClaimName).(*models.JwtAccessClaim)

	// call all service
	orderResp, err := h.orderService.Create(ctx, body, userClaim.Oid())
	if err != nil {
		return web.JsonRespError(c, err)
	}

	// return order list
	return web.JsonRespData(c, orderResp)
}

func (h *httpOrderHandler) Get(c *fiber.Ctx) error {
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(c.UserContext(), cfg.AppTimeout)
	defer cancel()

	// get access token claim
	userClaim := c.Locals(middleware.LocalUserClaimName).(*models.JwtAccessClaim)

	orderIdParam := c.Params("id")
	orderOid, err := primitive.ObjectIDFromHex(orderIdParam)
	if err != nil {
		return web.JsonResp(c, fiber.StatusBadRequest, "bad id", nil)
	}

	// call all service
	order, err := h.orderService.Get(ctx, orderOid, userClaim.Oid())
	if err != nil {
		return web.JsonRespError(c, err)
	}

	// return order
	return web.JsonRespData(c, order)
}

func (h *httpOrderHandler) GetAll(c *fiber.Ctx) error {
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(c.UserContext(), cfg.AppTimeout)
	defer cancel()

	// call all service
	orders, err := h.orderService.GetAll(ctx)
	if err != nil {
		return web.JsonRespError(c, err)
	}

	// return order
	return web.JsonRespData(c, orders)
}

func (h *httpOrderHandler) GetOwned(c *fiber.Ctx) error {
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(c.UserContext(), cfg.AppTimeout)
	defer cancel()

	// get access token claim
	userClaim := c.Locals(middleware.LocalUserClaimName).(*models.JwtAccessClaim)

	// call all service
	orders, err := h.orderService.GetAllByUser(ctx, userClaim.Oid())
	if err != nil {
		return web.JsonRespError(c, err)
	}

	// return order
	return web.JsonRespData(c, orders)
}

func (h *httpOrderHandler) Cancel(c *fiber.Ctx) error {
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(c.UserContext(), cfg.AppTimeout)
	defer cancel()

	// get access token claim
	userClaim := c.Locals(middleware.LocalUserClaimName).(*models.JwtAccessClaim)

	orderIdParam := c.Params("id")
	orderOid, err := primitive.ObjectIDFromHex(orderIdParam)
	if err != nil {
		return web.JsonResp(c, fiber.StatusBadRequest, "bad id", nil)
	}

	// call all service
	err = h.orderService.Cancel(ctx, orderOid, userClaim.Oid())
	if err != nil {
		return web.JsonRespError(c, err)
	}

	// return order
	return web.JsonRespData(c, "cancelled")

}
