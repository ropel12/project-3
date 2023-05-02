package routes

import (
	"github.com/labstack/echo/v4/middleware"
	eventhand "github.com/ropel12/project-3/app/features/event/handler"
	Trx "github.com/ropel12/project-3/app/features/transaction/handler"
	userhand "github.com/ropel12/project-3/app/features/user/handler"
	"github.com/ropel12/project-3/config/dependcy"
	"go.uber.org/dig"
)

type Routes struct {
	dig.In
	Depend dependcy.Depend
	User   userhand.User
	Event  eventhand.Event
	Trx    Trx.Transaction
}

func (r *Routes) RegisterRoutes() {
	ro := r.Depend.Echo
	ro.Use(middleware.RemoveTrailingSlash())
	ro.Use(middleware.Logger())
	ro.Use(middleware.Recover())
	ro.Use(middleware.CORS())
	//No Auth
	ro.POST("/login", r.User.Login)
	ro.POST("/register", r.User.Register)
	ro.GET("/events", r.Event.GetAll)
	ro.GET("/events/:id", r.Event.Detail)

	///Third-Party Payment Notification
	ro.POST("/notif", r.Trx.MidtransNotification)

	//Auth Area
	rauth := ro.Group("", middleware.JWT([]byte(r.Depend.Config.JwtSecret)))
	/// Users
	rauth.PUT("/users", r.User.Update)
	rauth.DELETE("/users", r.User.Delete)
	rauth.GET("/users", r.User.GetProfile)
	rauth.GET("/users/events", r.Event.MyEvent)
	rauth.GET("/users/history", r.Trx.MyHistory)
	rauth.GET("/users/transactions", r.Trx.GetByStatus)

	///Events
	rauth.POST("/events", r.Event.Create)
	rauth.PUT("/events", r.Event.Update)
	rauth.DELETE("/events/:id", r.Event.Delete)

	/// Trasanction
	rauth.POST("/transactions/cart", r.Trx.CreateCart)
	rauth.GET("/transactions/cart", r.Trx.GetCart)
	rauth.POST("/transactions/checkout", r.Trx.CreateTransaction)
	rauth.GET("/transactions/:invoice", r.Trx.GetDetail)

	// Tickets
	rauth.GET("/tickets/:invoice", r.Trx.GetTickets)
	rauth.POST("/tickets", r.Event.CreateTicket)
	rauth.DELETE("/tickets/:id", r.Event.DeleteTicket)

	//Comments
	rauth.POST("/comments", r.Event.CreateComment)

}
