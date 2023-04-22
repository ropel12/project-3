package handler

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	entity "github.com/ropel12/project-3/app/features/event"
	"github.com/ropel12/project-3/app/features/event/service"
	"github.com/ropel12/project-3/config/dependcy"
	"github.com/ropel12/project-3/helper"
	"go.uber.org/dig"
)

type Event struct {
	dig.In
	Service service.EventService
	Dep     dependcy.Depend
}

func (e *Event) Create(c echo.Context) error {
	var req entity.ReqCreate
	if err := c.Bind(&req); err != nil {
		e.Dep.Log.Errorf("Error handler: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Request Body", nil))
	}
	json.Unmarshal([]byte(req.Rtype), &req.Types)
	filehead, err := c.FormFile("image")
	if err != nil {
		e.Dep.Log.Errorf("Error handler: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Missing image in request body", nil))
	}
	file, err := filehead.Open()
	if err != nil {
		e.Dep.Log.Errorf("Error handler: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Error when opening the image", nil))
	}
	req.Image = filehead.Filename
	req.Uid = helper.GetUid(c.Get("user").(*jwt.Token))
	id, err := e.Service.Create(c.Request().Context(), req, file)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "OK", map[string]any{"id": id}))
}
