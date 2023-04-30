package handler

import (
	"log"
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

func (u *Transaction) CreateTransaction(c echo.Context) error {
	req := entity.ReqCheckout{}
	if err := c.Bind(&req); err != nil {
		u.Dep.Log.Errorf("Error service: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Request Body", nil))
	}
	req.UserId = helper.GetUid(c.Get("user").(*jwt.Token))
	invoice, err := u.Service.CreateTransaction(c.Request().Context(), req)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusCreated, CreateWebResponse(http.StatusCreated, "Success Operation", map[string]any{"data": invoice}))
}

func (u *Transaction) MidtransNotification(c echo.Context) error {
	midres := MidtransNotifResponse{}
	if err := c.Bind(&midres); err != nil {
		log.Printf("[ERROR] When Binding Midtrans Reponse : %v", err)
	}

	switch midres.TransactionStatus {
	case "settlement":
		if err := u.Service.UpdateStatus(c.Request().Context(), "Success", midres.OrderID); err != nil {
			log.Printf("[ERROR]When update status: %v", err)
		}
	case "expire":
		if err := u.Service.UpdateStatus(c.Request().Context(), "Cancel", midres.OrderID); err != nil {
			log.Printf("[ERROR]When update status: %v", err)
		}

	}
	return nil
}
func (u *Transaction) GetDetail(c echo.Context) error {
	invoice := c.Param("invoice")
	uid := helper.GetUid(c.Get("user").(*jwt.Token))
	res, err := u.Service.GetDetail(c.Request().Context(), invoice, uid)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success operation", res))
}

func (u *Transaction) MyHistory(c echo.Context) error {
	uid := helper.GetUid(c.Get("user").(*jwt.Token))
	res, err := u.Service.GetHistoryByuid(c.Request().Context(), uid)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success operation", res))
}

func (u *Transaction) GetByStatus(c echo.Context) error {
	uid := helper.GetUid(c.Get("user").(*jwt.Token))
	status := c.QueryParam("status")
	if status == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "query param status is missing", nil))
	}
	res, err := u.Service.GetByStatus(c.Request().Context(), uid, status)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success operation", res))
}
