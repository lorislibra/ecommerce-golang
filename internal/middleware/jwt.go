package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/donnjedarko/paninaro/config"
	"github.com/donnjedarko/paninaro/internal/models"
	"github.com/donnjedarko/paninaro/internal/web"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var JwtSignMethod = jwt.SigningMethodHS256

const (
	LocalJwtTokenName  = "jwt_token"
	LocalUserClaimName = "user_claim"
)

func GetJwtKey(t *jwt.Token) (interface{}, error) {
	return config.Get().PrivateKey, nil
}

func jwtFromHeader(header string, authScheme string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		auth := c.Get(header)
		l := len(authScheme)
		if len(auth) > l+1 && strings.EqualFold(auth[:l], authScheme) {
			return auth[l+1:], nil
		}
		return "", errors.New("missing or malformed JWT")
	}
}

func JwtMiddleware(c *fiber.Ctx) error {
	jwtExtractor := jwtFromHeader("Authorization", "Bearer")
	token, err := jwtExtractor(c)
	if err != nil {
		return web.JsonResp(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	claim := new(models.JwtAccessClaim)
	jwtToken, err := jwt.ParseWithClaims(token, claim, GetJwtKey)
	if err != nil {
		return web.JsonResp(c, fiber.StatusUnauthorized, "invalid or expired jwt", nil)
	}

	c.Locals(LocalJwtTokenName, jwtToken)
	c.Locals(LocalUserClaimName, claim)
	return c.Next()
}

func JwtMiddlewareExpired(c *fiber.Ctx) error {
	jwtExtractor := jwtFromHeader("Authorization", "Bearer")
	token, err := jwtExtractor(c)
	if err != nil {
		return web.JsonResp(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	claim := new(models.JwtAccessClaim)
	jwtToken, err := jwt.ParseWithClaims(token, claim, GetJwtKey)
	if err != nil {
		v, ok := err.(*jwt.ValidationError)
		if !ok || v.Errors != jwt.ValidationErrorExpired || !claim.VerifyExpiresAt(time.Now(), true) {
			return web.JsonResp(c, fiber.StatusUnauthorized, "invalid jwt", nil)
		}
	}

	c.Locals(LocalJwtTokenName, jwtToken)
	c.Locals(LocalUserClaimName, claim)
	return c.Next()
}
