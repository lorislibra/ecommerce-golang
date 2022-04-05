package domains

import (
	"context"

	"github.com/donnjedarko/paninaro/src/dtos"
	"github.com/donnjedarko/paninaro/src/entities"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	FindByUserOrEmail(ctx context.Context, username string, email string) (*entities.User, error)
	Find(ctx context.Context, oid primitive.ObjectID) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	FindAll(ctx context.Context) ([]*entities.User, error)
}

type UserService interface {
	Me(ctx context.Context, oid primitive.ObjectID) (*dtos.UserMeResp, error)
	All(ctx context.Context) ([]*dtos.UserMeResp, error)
}

type HttpUserHandler interface {
	Me(c *fiber.Ctx) error
	All(c *fiber.Ctx) error
}
