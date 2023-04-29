package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

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
		GetDetail(ctx context.Context, invoice string, uid int) (*entity.Response, error)
		GetHistoryByuid(ctx context.Context, uid int) (*entity.Response, error)
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
	var status, expire, date, paymentcode string
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
	if total > 0 {
		res, err := t.dep.Mds.CreateCharge(reqcharge)
		if err != nil {
			return "", errorr.NewBad(err.Error())
		}
		if res.Expire == "" {
			res.Expire = helper.GenerateExpiretime(res.TransactionTime, 1)
		}
		status = "pending"
		date = res.TransactionTime
		expire = res.Expire
		paymentcode = res.PaymentCode
	} else {
		status = "paid"
		date = time.Now().Format("2006-01-02 15:04:05")
		expire = "0000-00-00 00:00:00"
		paymentcode = "-"
	}

	trxdata := entity2.Transaction{
		Invoice:          invoice,
		PaymentMethod:    req.PaymentType,
		Status:           status,
		Date:             date,
		Total:            total,
		UserID:           uint(req.UserId),
		EventID:          uint(req.EventId),
		Expire:           expire,
		PaymentCode:      paymentcode,
		TransactionItems: transactionitems,
	}

	if err := t.repo.CreateTransaction(t.dep.Db.WithContext(ctx), trxdata); err != nil {
		return "", err
	}
	if total > 0 {
		encodeddata, _ := json.Marshal(map[string]any{"invoice": invoice, "total": total, "name": userdetail.Name, "email": userdetail.Email, "payment_code": paymentcode, "payment_method": req.PaymentType, "expire": expire})
		err := t.dep.Nsq.Publish("1", encodeddata)
		if err != nil {
			t.dep.Log.Errorf("Failed to publish to NSQ: %v", err)
		}
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

func (t *transaction) GetDetail(ctx context.Context, invoice string, uid int) (*entity.Response, error) {
	data, err := t.repo.GetByInvoice(t.dep.Db.WithContext(ctx), invoice, uid)
	if err != nil {
		return nil, err
	}
	itemsdetail := []entity.ItemDetails{}
	trasactiondata := entity.Transaction{
		Total:         int64(data.Total),
		Date:          data.Date,
		Expire:        data.Expire,
		PaymentMethod: data.PaymentMethod,
		Status:        data.Status,
		PaymentCode:   data.PaymentCode,
	}
	itemschan := make(chan entity.ItemDetails)
	wg := sync.WaitGroup{}
	go func(wg *sync.WaitGroup, itemschan <-chan entity.ItemDetails) {
		for val := range itemschan {
			itemsdetail = append(itemsdetail, val)
			wg.Done()
		}
	}(&wg, itemschan)
	wg.Add(len(data.TransactionItems))
	for _, val := range data.TransactionItems {
		itemdetail := entity.ItemDetails{
			Name:     val.Type.Name,
			Price:    val.Price,
			Qty:      val.Qty,
			SubTotal: val.Qty * val.Price,
		}
		itemschan <- itemdetail
	}
	wg.Wait()
	close(itemschan)
	trasactiondata.ItemDetails = itemsdetail
	res := entity.Response{
		Data: trasactiondata,
	}
	return &res, nil
}
func (t *transaction) GetHistoryByuid(ctx context.Context, uid int) (*entity.Response, error) {
	restrx := []entity.EventTransaction{}
	data, err := t.repo.GetHistory(t.dep.Db.WithContext(ctx), uid)
	if err != nil {
		return nil, err
	}
	trxchan := make(chan entity.EventTransaction)
	wg := &sync.WaitGroup{}
	go func(wg *sync.WaitGroup, data <-chan entity.EventTransaction) {
		for val := range data {
			restrx = append(restrx, val)
			wg.Done()
		}

	}(wg, trxchan)
	wg.Add(len(data))
	for _, val := range data {
		eventrx := entity.EventTransaction{
			Id:           int(val.EventID),
			Date:         val.Event.StartDate,
			Location:     val.Event.Location,
			EndDate:      val.Event.EndDate,
			HostedBy:     val.Event.HostedBy,
			Image:        val.Event.Image,
			Participants: len(val.Event.Users),
		}
		trxchan <- eventrx
	}
	wg.Wait()
	close(trxchan)
	res := entity.Response{
		Data: restrx,
	}
	return &res, nil
}
