package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/ropel12/project-3/app/entities"
	"github.com/ropel12/project-3/errorr"
)

type ChargeResponse struct {
	TransactionID          string             `json:"transaction_id"`
	OrderID                string             `json:"order_id"`
	GrossAmount            string             `json:"gross_amount"`
	PaymentType            string             `json:"payment_type"`
	TransactionTime        string             `json:"transaction_time"`
	TransactionStatus      string             `json:"transaction_status"`
	FraudStatus            string             `json:"fraud_status"`
	StatusCode             string             `json:"status_code"`
	Bank                   string             `json:"bank"`
	StatusMessage          string             `json:"status_message"`
	ChannelResponseCode    string             `json:"channel_response_code"`
	ChannelResponseMessage string             `json:"channel_response_message"`
	Currency               string             `json:"currency"`
	ValidationMessages     []string           `json:"validation_messages"`
	PermataVaNumber        string             `json:"permata_va_number"`
	VaNumbers              []coreapi.VANumber `json:"va_numbers"`
	Actions                []coreapi.Action   `json:"actions"`
	PaymentCode            string             `json:"payment_code"`
	QRString               string             `json:"qr_string"`
	Expire                 string             `json:"expiry_time"`
}
type Bank string

const (
	Bni     Bank = "bni"
	Mandiri Bank = "mandiri"
	Cimb    Bank = "cimb"
	Bca     Bank = "bca"
	Bri     Bank = "bri"
)

type Cstore string

const (
	Indomaret Cstore = "indomaret"
	Alafamart Cstore = "alfamart"
)

type Ewallet string

const (
	Gopay Ewallet = "gopay"
	Qris  Ewallet = "qris"
)

type Midtrans struct {
	Midtrans    coreapi.Client
	Req         *coreapi.ChargeReq
	ExpDuration int
	ExpUnit     string
}

func (m *Midtrans) CreateCharge(req entities.ReqCharge) (*ChargeResponse, error) {
	newreq := &coreapi.ChargeReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.Invoice,
			GrossAmt: int64(req.Total),
		},
		Items:           req.ItemsDetails,
		CustomerDetails: req.CustomerDetails,
		CustomExpiry: &coreapi.CustomExpiry{
			ExpiryDuration: m.ExpDuration,
			Unit:           m.ExpUnit,
		},
	}
	m.Req = newreq
	switch req.PaymentType {
	case "bca":
		return m.WithBank(Bca)
	case "mandiri":
		return m.WithBank(Mandiri)
	case "bni":
		return m.WithBank(Bni)
	case "indomaret":
		return m.WithCstore(Indomaret)
	case "alfamart":
		return m.WithCstore(Alafamart)
	case "gopay":
		return m.WithEwallet(Gopay)
	case "qris":
		return m.WithEwallet(Qris)

	}
	return nil, errors.New("payment type not available")
}

func (m *Midtrans) WithBank(bank Bank) (*ChargeResponse, error) {
	m.Req.PaymentType = "bank_transfer"
	m.Req.BankTransfer = &coreapi.BankTransferDetails{
		Bank: midtrans.Bank(bank),
	}
	res, err := m.ChargeCustom(m.Req)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return res, nil
}

func (m *Midtrans) WithCstore(cstore Cstore) (*ChargeResponse, error) {
	m.Req.PaymentType = "cstore"
	m.Req.ConvStore = &coreapi.ConvStoreDetails{
		Store: string(cstore),
	}
	res, err := m.ChargeCustom(m.Req)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return res, nil
}

func (m *Midtrans) WithEwallet(ewallet Ewallet) (*ChargeResponse, error) {
	m.Req.PaymentType = coreapi.CoreapiPaymentType(ewallet)
	res, err := m.ChargeCustom(m.Req)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return res, nil
}

func (m *Midtrans) ChargeCustom(req *coreapi.ChargeReq) (*ChargeResponse, error) {

	resp := &ChargeResponse{}
	jsonReq, _ := json.Marshal(req)
	err := m.Midtrans.HttpClient.Call(http.MethodPost,
		fmt.Sprintf("%s/v2/charge", m.Midtrans.Env.BaseUrl()),
		&m.Midtrans.ServerKey,
		m.Midtrans.Options,
		bytes.NewBuffer(jsonReq),
		resp,
	)
	switch resp.PaymentType {
	case "bank_transfer":
		if resp.PermataVaNumber != "" {

			resp.PaymentCode = resp.PermataVaNumber
		} else {
			resp.PaymentCode = resp.VaNumbers[0].VANumber
		}
	case "gopay", "qris":
		resp.PaymentCode = resp.Actions[0].URL
	}
	if err != nil {

		return resp, errorr.NewBad("Invalid Request Body")
	}
	return resp, nil
}