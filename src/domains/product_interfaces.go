package domains

import (
	"context"

	"github.com/donnjedarko/paninaro/src/dtos"
	"github.com/donnjedarko/paninaro/src/entities"
	"github.com/gofiber/fiber/v2"
)

type ProductRepository interface {
	Find(ctx context.Context, sku string) (*entities.Product, error)
	FindAll(ctx context.Context) ([]*entities.Product, error)
	Create(ctx context.Context, product *entities.Product) error
	Edit(ctx context.Context, sku string, product *dtos.ProductUpdateBody) (bool, error)
	SetHidden(ctx context.Context, sku string, value bool) (bool, error)
	FindMany(ctx context.Context, skus []string) ([]*entities.Product, error)
}

type ProductService interface {
	All(ctx context.Context) ([]*dtos.ProductResp, error)
	Get(ctx context.Context, sku string) (*dtos.ProductResp, error)
	Create(ctx context.Context, product *dtos.ProductCreateBody) (*dtos.ProductResp, error)
	Edit(ctx context.Context, sku string, product *dtos.ProductUpdateBody) error
	SetHidden(ctx context.Context, sku string, value bool) error
}

type HttpProductHandler interface {
	All(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Edit(c *fiber.Ctx) error
	Hide(c *fiber.Ctx) error
	UnHide(c *fiber.Ctx) error
}
