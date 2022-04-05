package routers

import (
	"github.com/donnjedarko/paninaro/internal/middleware"
	"github.com/donnjedarko/paninaro/internal/web"
	"github.com/donnjedarko/paninaro/src/handlers"
	"github.com/donnjedarko/paninaro/src/repositories"
	"github.com/donnjedarko/paninaro/src/services"
)

type httpUserRouter struct {
	Web *web.WebRouter
}

func NewHttpUserRouter(web *web.WebRouter) web.Router {
	return &httpUserRouter{
		Web: web,
	}
}

func (r *httpUserRouter) GetRoutes() {
	userRepo := repositories.NewUserMongoRepository(r.Web.Mongo)
	userService := services.NewUserService(userRepo)
	handler := handlers.NewHttpUserHandler(userService)

	g := r.Web.Api().Group("/users", middleware.JwtMiddleware)
	g.Get("/me", handler.Me)

	g = r.Web.Api().Group("/admin/users", middleware.JwtMiddleware, middleware.AdminMiddleware)
	g.Get("/", handler.All)
}
