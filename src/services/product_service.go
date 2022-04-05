package services

import (
	"context"
	"time"

	"github.com/donnjedarko/paninaro/internal/web"
	"github.com/donnjedarko/paninaro/src/domains"
	"github.com/donnjedarko/paninaro/src/dtos"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type productService struct {
	productRepo domains.ProductRepository
}

func NewProductService(productRepo domains.ProductRepository) domains.ProductService {
	return &productService{
		productRepo: productRepo,
	}
}

func (s *productService) All(ctx context.Context) ([]*dtos.ProductResp, error) {
	products, err := s.productRepo.FindAll(ctx)
	if err != nil {
		return nil, web.ErrInternal
	}

	productDtos := make([]*dtos.ProductResp, 0, len(products))
	for _, product := range products {
		if !product.Hidden {
			productDtos = append(productDtos, dtos.ProductRespFromEntity(product))
		}
	}

	return productDtos, nil
}

func (s *productService) Get(ctx context.Context, sku string) (*dtos.ProductResp, error) {
	product, err := s.productRepo.Find(ctx, sku)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, web.ErrNotFound
		}
		return nil, web.ErrInternal
	}

	if product.Hidden {
		return nil, web.ErrNotFound
	}

	return dtos.ProductRespFromEntity(product), nil
}

func (s *productService) Create(ctx context.Context, product *dtos.ProductCreateBody) (*dtos.ProductResp, error) {
	productE := product.ToEntity()
	productE.CreatedAt = time.Now()

	err := s.productRepo.Create(ctx, productE)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, web.NewWebErrMessage(fiber.StatusBadRequest, "sku already exist")
		}
		return nil, web.ErrInternal
	}

	return dtos.ProductRespFromEntity(productE), nil
}

func (s *productService) Edit(ctx context.Context, sku string, product *dtos.ProductUpdateBody) error {
	updated, err := s.productRepo.Edit(ctx, sku, product)
	if err != nil {
		return web.ErrInternal
	}

	if !updated {
		return web.NewWebErrMessage(fiber.StatusBadRequest, "product not updated")
	}

	return nil
}

func (s *productService) SetHidden(ctx context.Context, sku string, value bool) error {
	updated, err := s.productRepo.SetHidden(ctx, sku, value)
	if err != nil {
		return web.ErrInternal
	}

	if !updated {
		return web.NewWebErrMessage(fiber.StatusBadRequest, "product not updated")
	}

	return nil
}
