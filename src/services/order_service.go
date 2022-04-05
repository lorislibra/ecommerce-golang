package services

import (
	"context"
	"time"

	"github.com/donnjedarko/paninaro/internal/web"
	"github.com/donnjedarko/paninaro/src/domains"
	"github.com/donnjedarko/paninaro/src/dtos"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type orderService struct {
	orderRepo   domains.OrderRepository
	userRepo    domains.UserRepository
	productRepo domains.ProductRepository
}

func NewOrderService(
	orderRepo domains.OrderRepository,
	prodRepo domains.ProductRepository,
	userRepo domains.UserRepository,
) domains.OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		productRepo: prodRepo,
		userRepo:    userRepo,
	}
}

func (s *orderService) Create(
	ctx context.Context,
	order *dtos.OrderCreateBody,
	userOid primitive.ObjectID,
) (*dtos.OrderResp, error) {
	skus := order.Skus()

	products, err := s.productRepo.FindMany(ctx, skus)
	if err != nil {
		return nil, web.ErrInternal
	}

	if len(products) == 0 {
		return nil, web.NewWebErrMessage(fiber.StatusBadRequest, "no items in the order")
	}

	orderE := order.ToEntity(products)
	orderE.UserOid = userOid
	orderE.CreatedAt = time.Now()

	err = s.orderRepo.Create(ctx, orderE)
	if err != nil {
		return nil, web.ErrInternal
	}

	return dtos.OrderRespFromEntity(orderE), nil
}

func (s *orderService) Cancel(ctx context.Context, orderOid primitive.ObjectID, userOid primitive.ObjectID) error {
	err := s.orderRepo.Cancel(ctx, orderOid, userOid)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return web.ErrNotFound
		}
		return web.ErrInternal
	}
	return nil
}

func (s *orderService) Get(
	ctx context.Context,
	orderOid primitive.ObjectID,
	userOid primitive.ObjectID,
) (*dtos.OrderResp, error) {
	order, err := s.orderRepo.Find(ctx, orderOid, userOid)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, web.ErrNotFound
		}
		return nil, web.ErrInternal
	}

	return dtos.OrderRespFromEntity(order), nil
}

func (s *orderService) GetAll(ctx context.Context) ([]*dtos.OrderResp, error) {
	orders, err := s.orderRepo.FindAll(ctx)
	if err != nil {
		return nil, web.ErrInternal
	}

	orderDtos := make([]*dtos.OrderResp, 0, len(orders))
	for _, order := range orders {
		orderDtos = append(orderDtos, dtos.OrderRespFromEntity(order))
	}

	return orderDtos, nil
}

func (s *orderService) GetAllByUser(ctx context.Context, userOid primitive.ObjectID) ([]*dtos.OrderResp, error) {
	orders, err := s.orderRepo.FindAllByUser(ctx, userOid)
	if err != nil {
		return nil, web.ErrInternal
	}

	orderDtos := make([]*dtos.OrderResp, 0, len(orders))
	for _, order := range orders {
		orderDtos = append(orderDtos, dtos.OrderRespFromEntity(order))
	}

	return orderDtos, nil
}
