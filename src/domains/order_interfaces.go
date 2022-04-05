package domains

import (
	"context"

	"github.com/donnjedarko/paninaro/src/dtos"
	"github.com/donnjedarko/paninaro/src/entities"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderRepository interface {
	Find(ctx context.Context, orderOid primitive.ObjectID, userOid primitive.ObjectID) (*entities.Order, error)
	FindFull(ctx context.Context, orderOid primitive.ObjectID, userOid primitive.ObjectID) (*entities.Order, error)
	FindAll(ctx context.Context) ([]*entities.Order, error)
	FindAllByUser(ctx context.Context, userOid primitive.ObjectID) ([]*entities.Order, error)
	Create(ctx context.Context, order *entities.Order) error
	EditStatus(ctx context.Context, orderOid primitive.ObjectID, userOid primitive.ObjectID, status string) (bool, error)
	Cancel(ctx context.Context, orderOid primitive.ObjectID, userOid primitive.ObjectID) error
}

type OrderService interface {
	Cancel(ctx context.Context, orderOid primitive.ObjectID, userOid primitive.ObjectID) error
	Create(ctx context.Context, order *dtos.OrderCreateBody, userOid primitive.ObjectID) (*dtos.OrderResp, error)
	Get(ctx context.Context, orderOid primitive.ObjectID, userOid primitive.ObjectID) (*dtos.OrderResp, error)
	GetAllByUser(ctx context.Context, orderOid primitive.ObjectID) ([]*dtos.OrderResp, error)
	GetAll(ctx context.Context) ([]*dtos.OrderResp, error)
}

type HttpOrderHandler interface {
	GetAll(c *fiber.Ctx) error
	GetOwned(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Get(c *fiber.Ctx) error
	Cancel(c *fiber.Ctx) error
}
