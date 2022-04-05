package middleware

import (
	"github.com/donnjedarko/paninaro/internal/models"
	"github.com/donnjedarko/paninaro/internal/web"
	"github.com/donnjedarko/paninaro/src/entities"
	"github.com/gofiber/fiber/v2"
)

const (
	LocalUserEntityName = "user"
)

var AdminRoles = []entities.Role{entities.Seller, entities.Admin}

func AdminMiddleware(c *fiber.Ctx) error {
	user := c.Locals(LocalUserClaimName).(*models.JwtAccessClaim)

	if user.Role.In(AdminRoles) {
		return c.Next()
	}

	return web.JsonResp(c, fiber.StatusForbidden, "not allowed", nil)
}
