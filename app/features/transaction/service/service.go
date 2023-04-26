package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-playground/validator"
	"github.com/midtrans/midtrans-go"
	entity2 "github.com/ropel12/project-3/app/entities"
	entity "github.com/ropel12/project-3/app/features/transaction"
	"github.com/ropel12/project-3/app/features/transaction/repository"
	"github.com/ropel12/project-3/config/dependcy"
	"github.com/ropel12/project-3/errorr"
	"github.com/ropel12/project-3/helper"
)

type (
	transaction struct {
		repo      repository.TransactionRepo
		validator *validator.Validate
		dep       dependcy.Depend
	}
	TransactionService interface {
		CreateCart(ctx context.Context, req entity.ReqCart) error
		GetCart(ctx context.Context, uid int) ([]entity.Cart, error)
		CreateTransaction(ctx context.Context, req entity.ReqCheckout) (string, error)
		UpdateStatus(ctx context.Context, status, invoice string) error
	}
)

func NewTransactionService(repo repository.TransactionRepo, dep dependcy.Depend) TransactionService {
	return &transaction{
		repo:      repo,
		dep:       dep,
		validator: validator.New(),
	}
}

func (t *transaction) CreateCart(ctx context.Context, req entity.ReqCart) error {
	if err := t.validator.Struct(req); err != nil {
		return errorr.NewBad("Missing or Invalid Request Body")
	}
	data := entity2.Carts{
		UserID: uint(req.UID),
		TypeID: uint(req.TypeID),
		Qty:    1,
	}
	if err := t.repo.Create(t.dep.Db.WithContext(ctx), data); err != nil {
		return err
	}
	return nil
}

func (t *transaction) GetCart(ctx context.Context, uid int) ([]entity.Cart, error) {
	carts := []entity.Cart{}
	data, err := t.repo.GetCart(t.dep.Db.WithContext(ctx), uid)
	if err != nil {
		return nil, err
	}
	for _, val := range data {
		cart := entity.Cart{
			EventId:   int(val.Type.EventID),
			TypeID:    int(val.TypeID),
			TypeName:  val.Type.Name,
			TypePrice: val.Type.Price,
			Qty:       val.Qty,
			Subtotal:  val.Qty * val.Type.Price,
		}
		carts = append(carts, cart)
	}
	return carts, nil
}

func (t *transaction) CreateTransaction(ctx context.Context, req entity.ReqCheckout) (string, error) {
	if err := t.validator.Struct(req); err != nil {
		return "", errorr.NewBad("Invalid and missing request body")
	}
	var total int
	itemdetails := []midtrans.ItemDetails{}
	transactionitems := []entity2.TransactionItems{}
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for id, val := range req.ItemDetails {
			Item := midtrans.ItemDetails{
				ID:    fmt.Sprintf("%d", id+1),
				Name:  val.Name,
				Price: int64(val.Price),
				Qty:   int32(val.Qty),
			}
			itemdetails = append(itemdetails, Item)
			total += val.SubTotal
		}
	}()
	go func() {
		defer wg.Done()
		for _, val := range req.ItemDetails {
			item := entity2.TransactionItems{
				TypeID: uint(val.TypeId),
				Qty:    val.Qty,
				Price:  val.Price,
			}
			transactionitems = append(transactionitems, item)
		}
	}()

	userdetail := t.repo.GetDetailUserById(t.dep.Db.WithContext(ctx), req.UserId)
	invoice := helper.GenerateInvoice(req.EventId, req.UserId)
	customerdetails := &midtrans.CustomerDetails{
		FName: userdetail.Name,
		Email: userdetail.Email,
		BillAddr: &midtrans.CustomerAddress{
			FName:   userdetail.Name,
			Address: userdetail.Address,
		},
	}
	wg.Wait()
	reqcharge := entity2.ReqCharge{
		PaymentType:     req.PaymentType,
		Invoice:         invoice,
		Total:           total,
		ItemsDetails:    &itemdetails,
		CustomerDetails: customerdetails,
	}
	res, err := t.dep.Mds.CreateCharge(reqcharge)
	if err != nil {
		return "", errorr.NewBad(err.Error())
	}
	if res.Expire == "" {
		res.Expire = helper.GenerateExpiretime(res.TransactionTime, 1)
	}
	trxdata := entity2.Transaction{
		Invoice:          invoice,
		PaymentMethod:    req.PaymentType,
		Status:           res.TransactionStatus,
		Date:             res.TransactionTime,
		Total:            total,
		UserID:           uint(req.UserId),
		EventID:          uint(req.EventId),
		Expire:           res.Expire,
		PaymentCode:      res.PaymentCode,
		TransactionItems: transactionitems,
	}

	if err := t.repo.CreateTransaction(t.dep.Db.WithContext(ctx), trxdata); err != nil {
		return "", err
	}
	encodeddata, _ := json.Marshal(map[string]any{"invoice": invoice, "total": total, "name": userdetail.Name, "email": userdetail.Email, "payment_code": res.PaymentCode, "payment_method": res.PaymentType, "expire": res.Expire})
	err = t.dep.Nsq.Publish("1", encodeddata)
	if err != nil {
		t.dep.Log.Errorf("Failed to publish to NSQ: %v", err)
	}
	return invoice, nil
}

func (t *transaction) UpdateStatus(ctx context.Context, status, invoice string) error {

	if err := t.repo.UpdateStatusTrasansaction(t.dep.Db.WithContext(ctx), invoice, status); err != nil {
		return err
	}
	user := t.repo.GetDetailUserByInvoice(t.dep.Db.WithContext(ctx), invoice)
	encodeddata, _ := json.Marshal(map[string]any{"invoice": invoice, "email": user.User.Email, "name": user.User.Name})
	switch status {
	case "Success":
		go func() {
			if err := t.dep.Nsq.Publish("2", encodeddata); err != nil {
				t.dep.Log.Errorf("Failed to publish to NSQ: %v", err)
			}
		}()
		if err := t.dep.Pusher.Publish(map[string]string{"invoice": invoice, "status": "success"}); err != nil {
			t.dep.Log.Errorf("Failed to publish to PusherJs: %v", err)
		}
	case "Cancel":
		go func() {
			if err := t.dep.Nsq.Publish("3", encodeddata); err != nil {
				t.dep.Log.Errorf("Failed to publish to NSQ: %v", err)
			}
		}()
		if err := t.dep.Pusher.Publish(map[string]string{"invoice": invoice, "status": "canceled"}); err != nil {
			t.dep.Log.Errorf("Failed to publish to PusherJs: %v", err)
		}
	}
	return nil
}
