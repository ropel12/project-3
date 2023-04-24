package handler

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	entity "github.com/ropel12/project-3/app/features/transaction"
	"github.com/ropel12/project-3/app/features/transaction/service"
	"github.com/ropel12/project-3/config/dependcy"
	"github.com/ropel12/project-3/helper"
	"go.uber.org/dig"
)

type Transaction struct {
	dig.In
	Service service.TransactionService
	Dep     dependcy.Depend
}

func (u *Transaction) CreateCart(c echo.Context) error {
	var req entity.ReqCart
	if err := c.Bind(&req); err != nil {
		u.Dep.Log.Errorf("Error service: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Request Body", nil))
	}
	req.UID = helper.GetUid(c.Get("user").(*jwt.Token))
	if err := u.Service.CreateCart(c.Request().Context(), req); err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", nil))
}

func (u *Transaction) GetCart(c echo.Context) error {
	uid := helper.GetUid(c.Get("user").(*jwt.Token))
	res, err := u.Service.GetCart(c.Request().Context(), uid)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", res))
}
