package routers

import (
	"github.com/donnjedarko/paninaro/internal/middleware"
	"github.com/donnjedarko/paninaro/internal/web"
	"github.com/donnjedarko/paninaro/src/handlers"
	"github.com/donnjedarko/paninaro/src/repositories"
	"github.com/donnjedarko/paninaro/src/services"
)

type httpAuthRouter struct {
	Web *web.WebRouter
}

func NewHttpAuthRouter(web *web.WebRouter) web.Router {
	return &httpAuthRouter{
		Web: web,
	}
}

func (r *httpAuthRouter) GetRoutes() {
	userRepo := repositories.NewUserMongoRepository(r.Web.Mongo)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(r.Web.Redis)
	authService := services.NewAuthService(refreshTokenRepo, userRepo)
	handler := handlers.NewHttpAuthHandler(authService)

	g := r.Web.Api().Group("/auth")
	g.Post("/refresh", handler.Refresh)
	g.Post("/token", handler.Signin)
	g.Post("/signup", handler.Signup)
	g.Post("/revoke", middleware.JwtMiddleware, handler.Signout)
	g.Post("/logout", middleware.JwtMiddleware, handler.SignoutAll)

}
