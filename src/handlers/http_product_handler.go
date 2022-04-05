package handlers

import (
	"context"

	"github.com/donnjedarko/paninaro/config"
	"github.com/donnjedarko/paninaro/internal/web"
	"github.com/donnjedarko/paninaro/src/domains"
	"github.com/donnjedarko/paninaro/src/dtos"
	"github.com/gofiber/fiber/v2"
)

type httpProductHandler struct {
	productService domains.ProductService
}

func NewHttpProductHandler(productService domains.ProductService) domains.HttpProductHandler {
	return &httpProductHandler{
		productService: productService,
	}
}

func (h *httpProductHandler) All(c *fiber.Ctx) error {
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(c.UserContext(), cfg.AppTimeout)
	defer cancel()

	// call all service
	products, err := h.productService.All(ctx)
	if err != nil {
		return web.JsonRespError(c, err)
	}

	// return product list
	return web.JsonRespData(c, products)
}

func (h *httpProductHandler) Get(c *fiber.Ctx) error {
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(c.UserContext(), cfg.AppTimeout)
	defer cancel()

	// call all service
	sku := c.Params("sku")
	product, err := h.productService.Get(ctx, sku)
	if err != nil {
		return web.JsonRespError(c, err)
	}

	// return product list
	return web.JsonRespData(c, product)
}

func (h *httpProductHandler) Create(c *fiber.Ctx) error {
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(c.UserContext(), cfg.AppTimeout)
	defer cancel()

	// load body from request
	body := new(dtos.ProductCreateBody)
	if err := c.BodyParser(body); err != nil {
		return web.JsonResp(c, fiber.StatusBadRequest, "bad payload", nil)
	}

	// call all service
	productResp, err := h.productService.Create(ctx, body)
	if err != nil {
		return web.JsonRespError(c, err)
	}

	// return product list
	return web.JsonRespData(c, productResp)
}

func (h *httpProductHandler) Edit(c *fiber.Ctx) error {
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(c.UserContext(), cfg.AppTimeout)
	defer cancel()

	// load body from request
	body := new(dtos.ProductUpdateBody)
	if err := c.BodyParser(body); err != nil {
		return web.JsonResp(c, fiber.StatusBadRequest, "bad payload", nil)
	}

	// call edit service
	sku := c.Params("sku")
	err := h.productService.Edit(ctx, sku, body)
	if err != nil {
		return web.JsonRespError(c, err)
	}

	return web.JsonRespData(c, "product updated")
}

func (h *httpProductHandler) Hide(c *fiber.Ctx) error {
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(c.UserContext(), cfg.AppTimeout)
	defer cancel()

	sku := c.Params("sku")
	err := h.productService.SetHidden(ctx, sku, true)
	if err != nil {
		return web.JsonRespError(c, err)
	}

	return web.JsonRespData(c, "product visibility updated")
}

func (h *httpProductHandler) UnHide(c *fiber.Ctx) error {
	cfg := config.Get()
	ctx, cancel := context.WithTimeout(c.UserContext(), cfg.AppTimeout)
	defer cancel()

	sku := c.Params("sku")
	err := h.productService.SetHidden(ctx, sku, false)
	if err != nil {
		return web.JsonRespError(c, err)
	}

	return web.JsonRespData(c, "product visibility updated")
}
