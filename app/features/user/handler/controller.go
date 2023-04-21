package handler

import (
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	entity "github.com/ropel12/project-3/app/features/user"
	"github.com/ropel12/project-3/app/features/user/service"
	"github.com/ropel12/project-3/config/dependcy"
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
		return CreateErrorResponse(err, c)
	}
	wg.Wait()
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", map[string]any{"token": token}))
}

func (u *User) Register(c echo.Context) error {
	var req entity.RegisterReq
	if err := c.Bind(&req); err != nil {
		u.Dep.Log.Errorf("Error service: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Request Body", nil))
	}
	if err := u.Service.Register(c.Request().Context(), req); err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", nil))
}

func (u *User) Update(c echo.Context) error {
	var req entity.UpdateReq
	if err := c.Bind(&req); err != nil {
		u.Dep.Log.Errorf("Error service: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Request Body", nil))
	}
	req.Id = helper.GetUid(c.Get("user").(*jwt.Token))
	var filee multipart.File
	file, err1 := c.FormFile("image")
	if err1 == nil {
		files, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Cannot Load Image", nil))
		}
		req.Image = file.Filename
		filee = files
	}
	data, err := u.Service.Update(c.Request().Context(), req, filee)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	res := map[string]any{
		"id":       data.ID,
		"image":    data.Image,
		"password": data.Password,
		"email":    data.Email,
		"name":     data.Name,
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", res))
}

func (u *User) Delete(c echo.Context) error {
	if err := u.Service.Delete(c.Request().Context(), helper.GetUid(c.Get("user").(*jwt.Token))); err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", nil))
}
