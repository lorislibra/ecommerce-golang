package routes

import (
	"github.com/donnjedarko/paninaro/internal/web"
	"github.com/donnjedarko/paninaro/src/routers"
)

type MainRouter struct {
	Web *web.WebRouter
}

func NewMainRouter(web *web.WebRouter) *MainRouter {
	return &MainRouter{
		Web: web,
	}
}

func (r *MainRouter) GetRoutes() {
	routers.NewHttpUserRouter(r.Web).GetRoutes()
	routers.NewHttpAuthRouter(r.Web).GetRoutes()
	routers.NewHttpProductRouter(r.Web).GetRoutes()
	routers.NewHttpOrderRouter(r.Web).GetRoutes()
}
