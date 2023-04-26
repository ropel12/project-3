package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ropel12/project-3/errorr"
)

type (
	WebResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    any    `json:"data,omitempty"`
	}

	MidtransNotifResponse struct {
		StatusCode        string `json:"status_code"`
		OrderID           string `json:"order_id"`
		TransactionTime   string `json:"transaction_time"`
		TransactionStatus string `json:"transaction_status"`
		FraudStatus       string `json:"fraud_status"`
	}
)

func CreateWebResponse(code int, message string, data any) any {
	return WebResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func CreateErrorResponse(err error, c echo.Context) error {
	if err, ok := err.(errorr.BadRequest); ok {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, err.Error(), nil))
	}
	return c.JSON(http.StatusInternalServerError, CreateWebResponse(http.StatusInternalServerError, err.Error(), nil))
}
