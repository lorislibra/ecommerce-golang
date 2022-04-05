package web

import (
	"github.com/donnjedarko/paninaro/infrastructures/db"
	"github.com/gofiber/fiber/v2"
)

type Router interface {
	GetRoutes()
}

type WebRouter struct {
	App   *fiber.App
	Mongo db.MongoInstance
	Redis db.RedisInstance
}

func NewWebRouter(app *fiber.App, mongo db.MongoInstance, redis db.RedisInstance) *WebRouter {
	return &WebRouter{
		App:   app,
		Mongo: mongo,
		Redis: redis,
	}
}

func (wr *WebRouter) Api() fiber.Router {
	return wr.App.Group("/api/")
}
