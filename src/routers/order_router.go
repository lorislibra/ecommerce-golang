package routers

import (
	"github.com/donnjedarko/paninaro/internal/middleware"
	"github.com/donnjedarko/paninaro/internal/web"
	"github.com/donnjedarko/paninaro/src/handlers"
	"github.com/donnjedarko/paninaro/src/repositories"
	"github.com/donnjedarko/paninaro/src/services"
)

type httpOrderRouter struct {
	Web *web.WebRouter
}

func NewHttpOrderRouter(web *web.WebRouter) web.Router {
	return &httpOrderRouter{
		Web: web,
	}
}

func (r *httpOrderRouter) GetRoutes() {
	orderRepo := repositories.NewOrderRepository(r.Web.Mongo)
	productRepo := repositories.NewProductRepository(r.Web.Mongo)
	userRepo := repositories.NewUserMongoRepository(r.Web.Mongo)

	orderService := services.NewOrderService(orderRepo, productRepo, userRepo)
	orderHandler := handlers.NewHttpOrderHandler(orderService)

	g := r.Web.Api().Group("/orders", middleware.JwtMiddleware)
	g.Post("/", orderHandler.Create)
	g.Get("/:id", orderHandler.Get)
	g.Get("/", orderHandler.GetOwned)
	g.Post("/:id/cancel", orderHandler.Cancel)

	g = r.Web.Api().Group("/admin/orders", middleware.JwtMiddleware, middleware.AdminMiddleware)
	g.Get("/", orderHandler.GetAll)
}
