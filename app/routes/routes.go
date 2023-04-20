package routes

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/ropel12/project-3/app/features/user/handler"
	"github.com/ropel12/project-3/config/dependcy"
	"go.uber.org/dig"
)

type Routes struct {
	dig.In
	Depend dependcy.Depend
	User   handler.User
}

func (r *Routes) RegisterRoutes() {
	ro := r.Depend.Echo
	ro.Use(middleware.RemoveTrailingSlash())
	ro.Use(middleware.Logger())
	ro.Use(middleware.Recover())
	ro.POST("/login", r.User.Login)
	ro.POST("/register", r.User.Register)
}
