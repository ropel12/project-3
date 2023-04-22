package routes

import (
	"github.com/labstack/echo/v4/middleware"
	eventhand "github.com/ropel12/project-3/app/features/event/handler"
	userhand "github.com/ropel12/project-3/app/features/user/handler"
	"github.com/ropel12/project-3/config/dependcy"
	"go.uber.org/dig"
)

type Routes struct {
	dig.In
	Depend dependcy.Depend
	User   userhand.User
	Event  eventhand.Event
}

func (r *Routes) RegisterRoutes() {
	ro := r.Depend.Echo
	ro.Use(middleware.RemoveTrailingSlash())
	ro.Use(middleware.Logger())
	ro.Use(middleware.Recover())
	//No Auth
	ro.POST("/login", r.User.Login)
	ro.POST("/register", r.User.Register)
	//Auth Area
	rauth := ro.Group("", middleware.JWT([]byte(r.Depend.Config.JwtSecret)))
	/// Users
	rauth.PUT("/users", r.User.Update)
	rauth.DELETE("/users", r.User.Delete)
	rauth.GET("/users", r.User.GetProfile)

	///Events
	rauth.POST("/events", r.Event.Create)

}
