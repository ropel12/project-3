package handler

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"

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
	return c.JSON(http.StatusCreated, CreateWebResponse(http.StatusCreated, "Success operation", map[string]any{"id": id}))
}

func (e *Event) MyEvent(c echo.Context) error {
	uid := helper.GetUid(c.Get("user").(*jwt.Token))
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")
	if page == "" || limit == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Missing limit and page query params", nil))
	}
	newpage, err := strconv.Atoi(page)
	newlimit, err1 := strconv.Atoi(limit)
	if err != nil || err1 != nil {
		e.Dep.Log.Errorf("error handler : %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid query param", nil))
	}
	res, err := e.Service.MyEvent(c.Request().Context(), uid, newlimit, newpage)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success operation", res))
}

func (e *Event) Delete(c echo.Context) error {
	eventid := c.Param("id")
	if eventid == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Missing param id", nil))
	}
	neweventid, err := strconv.Atoi(eventid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid param id", nil))
	}
	if err := e.Service.Delete(c.Request().Context(), neweventid, helper.GetUid(c.Get("user").(*jwt.Token))); err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success operation", nil))
}

func (e *Event) GetAll(c echo.Context) error {
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")
	if page == "" || limit == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Missing limit and page query params", nil))
	}
	newpage, err := strconv.Atoi(page)
	newlimit, err1 := strconv.Atoi(limit)
	if err != nil || err1 != nil {
		e.Dep.Log.Errorf("error handler : %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid query param", nil))
	}
	res, err := e.Service.GetAll(c.Request().Context(), newlimit, newpage)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success operation", res))
}

func (e *Event) Detail(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Missing param id", nil))
	}
	newid, err := strconv.Atoi(id)
	if err != nil {
		e.Dep.Log.Errorf("error handler : %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid query param", nil))
	}
	res, err := e.Service.Detail(c.Request().Context(), newid)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success operation", res))
}

func (e *Event) Update(c echo.Context) error {
	var req entity.ReqUpdate
	if err := c.Bind(&req); err != nil {
		e.Dep.Log.Errorf("Error service: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Request Body", nil))
	}
	var file multipart.File
	fileh, err1 := c.FormFile("image")
	if err1 == nil {
		files, err := fileh.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Cannot Load Image", nil))
		}
		req.Image = fileh.Filename
		file = files
	}
	id, err := e.Service.Update(c.Request().Context(), req, file)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success operation", map[string]any{"id": id}))
}
