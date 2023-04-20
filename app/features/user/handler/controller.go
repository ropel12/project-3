package handler

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	entity "github.com/ropel12/project-3/app/features/user"
	"github.com/ropel12/project-3/app/features/user/service"
	"github.com/ropel12/project-3/config/dependcy"
	"github.com/ropel12/project-3/errorr"
	"github.com/ropel12/project-3/helper"
	"go.uber.org/dig"
)

type User struct {
	dig.In
	Service service.UserService
	Dep     dependcy.Depend
}

func (u *User) Login(c echo.Context) error {
	var req entity.LoginReq
	var token string
	if err := c.Bind(&req); err != nil {
		u.Dep.Log.Errorf("Error handler: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Request Body", nil))
	}
	uid, err := u.Service.Login(c.Request().Context(), req)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		token = helper.GenerateJWT(uid, u.Dep)
	}()
	if err != nil {
		if err, ok := err.(errorr.BadRequest); ok {
			return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, err.Error(), nil))
		} else {
			return c.JSON(http.StatusInternalServerError, CreateWebResponse(http.StatusInternalServerError, err.Error(), nil))
		}
	}
	wg.Wait()
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", map[string]any{"token": token}))
}
