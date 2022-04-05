package domains

import (
	"context"
	"time"

	"github.com/donnjedarko/paninaro/internal/models"
	"github.com/donnjedarko/paninaro/src/dtos"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshTokenRepository interface {
	Save(ctx context.Context, uid string, tokenId string, refreshToken string, expire time.Duration) error
	SaveAndDelete(ctx context.Context, uid string, tokenId string, refreshToken string, expire time.Duration, oldTokenId string) error
	Delete(ctx context.Context, uid string, tokenId string) (bool, error)
	DeleteAll(ctx context.Context, uid string) (bool, error)
	Exist(ctx context.Context, uid string, tokenId string) (bool, error)
}

type AuthService interface {
	Signin(ctx context.Context, body *dtos.UserSigninBody) (*models.JwtPair, error)
	Signup(ctx context.Context, body *dtos.UserSignupBody) (*models.JwtPair, error)
	Signout(ctx context.Context, oid primitive.ObjectID, refreshTokenId string) error
	Refresh(ctx context.Context, refreshToken string) (*models.JwtPair, error)
}

type HttpAuthHandler interface {
	Signin(c *fiber.Ctx) error
	Signup(c *fiber.Ctx) error
	Signout(c *fiber.Ctx) error
	SignoutAll(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
}
