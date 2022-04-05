package routers

import (
	"github.com/donnjedarko/paninaro/internal/middleware"
	"github.com/donnjedarko/paninaro/internal/web"
	"github.com/donnjedarko/paninaro/src/handlers"
	"github.com/donnjedarko/paninaro/src/repositories"
	"github.com/donnjedarko/paninaro/src/services"
)

type httpProductRouter struct {
	Web *web.WebRouter
}

func NewHttpProductRouter(web *web.WebRouter) web.Router {
	return &httpProductRouter{
		Web: web,
	}
}

func (r *httpProductRouter) GetRoutes() {
	productRepo := repositories.NewProductRepository(r.Web.Mongo)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewHttpProductHandler(productService)

	g := r.Web.Api().Group("/products", middleware.JwtMiddleware)
	g.Get("/", productHandler.All)
	g.Get("/:sku", productHandler.Get)

	g = r.Web.Api().Group("/admin/products", middleware.JwtMiddleware, middleware.AdminMiddleware)
	g.Post("/", productHandler.Create)
	g.Post("/:sku/edit", productHandler.Edit)
	g.Post("/:sku/hide", productHandler.Hide)
	g.Post("/:sku/unhide", productHandler.UnHide)
}
