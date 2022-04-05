package models

import (
	"github.com/donnjedarko/paninaro/src/dtos"
	"github.com/donnjedarko/paninaro/src/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JwtAccessClaim struct {
	jwt.RegisteredClaims
	Role entities.Role
}

func (jac *JwtAccessClaim) Oid() primitive.ObjectID {
	oid, _ := primitive.ObjectIDFromHex(jac.Subject)
	return oid
}

type JwtRefreshClaim struct {
	jwt.RegisteredClaims
}

func (jrc *JwtRefreshClaim) Oid() primitive.ObjectID {
	oid, _ := primitive.ObjectIDFromHex(jrc.Subject)
	return oid
}

type JwtPair struct {
	AccessToken  *dtos.JwtResponse
	RefreshToken *fiber.Cookie
}
