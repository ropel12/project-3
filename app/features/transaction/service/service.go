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
		GetCart(ctx context.Context, uid int) (*entity.Response, error)
		CreateTransaction(ctx context.Context, req entity.ReqCheckout) (*entity.Transaction, error)
		UpdateStatus(ctx context.Context, status, invoice string) error
		GetDetail(ctx context.Context, invoice string, uid int) (*entity.Response, error)
		GetHistoryByuid(ctx context.Context, uid int) (*entity.Response, error)
		GetByStatus(ctx context.Context, uid int, status string) (*entity.Response, error)
		GetTickets(ctx context.Context, invoice string, uid int) (*entity.Response, error)
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

func (t *transaction) GetCart(ctx context.Context, uid int) (*entity.Response, error) {
	carts := []entity.Cart{}
	data, err := t.repo.GetCart(t.dep.Db.WithContext(ctx), uid)
	if err != nil {
		return nil, err
	}
	total := 0
	for _, val := range data {
		subtotal := val.Qty * val.Type.Price
		cart := entity.Cart{
			EventId:   int(val.Type.EventID),
			TypeID:    int(val.TypeID),
			TypeName:  val.Type.Name,
			TypePrice: val.Type.Price,
			Qty:       val.Qty,
			Subtotal:  subtotal,
		}
		total += subtotal
		carts = append(carts, cart)
	}
	res := entity.Response{
		Total: total,
		Data:  carts,
	}
	return &res, nil
}

func (t *transaction) CreateTransaction(ctx context.Context, req entity.ReqCheckout) (*entity.Transaction, error) {
	if err := t.validator.Struct(req); err != nil {
		return nil, errorr.NewBad("Invalid or missing request body")
	}
	var total, totalqty int
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
			totalqty += val.Qty
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
	wg.Wait()
	if err := t.repo.CheckQuota(t.dep.Db.WithContext(ctx), req.EventId, totalqty); err != nil {
		return nil, err
	}
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
			return nil, errorr.NewBad(err.Error())
		}
		if res.Expire == "" {
			res.Expire = helper.GenerateExpiretime(res.TransactionTime, t.dep.Mds.ExpDuration)
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
		return nil, err
	}
	if total > 0 {
		encodeddata, _ := json.Marshal(map[string]any{"invoice": invoice, "total": total, "name": userdetail.Name, "email": userdetail.Email, "payment_code": paymentcode, "payment_method": req.PaymentType, "expire": expire})
		err := t.dep.Nsq.Publish("1", encodeddata)
		if err != nil {
			t.dep.Log.Errorf("Failed to publish to NSQ: %v", err)
		}
	} else {
		if err := t.repo.UpdateQuotaEvent(t.dep.Db.WithContext(ctx), req.EventId, totalqty); err != nil {
			return nil, err
		}
	}
	res := entity.Transaction{
		Total:         int64(total),
		Expire:        expire,
		PaymentMethod: req.PaymentType,
		PaymentCode:   paymentcode,
		Invoice:       invoice,
	}
	return &res, nil
}

func (t *transaction) UpdateStatus(ctx context.Context, status, invoice string) error {

	user := t.repo.GetDetailUserByInvoice(t.dep.Db.WithContext(ctx), invoice)
	encodeddata, _ := json.Marshal(map[string]any{"invoice": invoice, "email": user.User.Email, "name": user.User.Name})
	switch status {
	case "paid":
		eventId, qty, err := t.repo.GetQtyByInvoice(t.dep.Db.WithContext(ctx), invoice)
		if err != nil {
			return err
		}
		err = t.repo.CheckQuota(t.dep.Db.WithContext(ctx), eventId, qty)
		if err != nil {
			return err
		}
		if err := t.repo.UpdateStatusTrasansaction(t.dep.Db.WithContext(ctx), invoice, status); err != nil {
			return err
		}
		if err := t.repo.UpdateQuotaEvent(t.dep.Db.WithContext(ctx), eventId, qty); err != nil {
			return err
		}
		go func() {
			if err := t.dep.Nsq.Publish("2", encodeddata); err != nil {
				t.dep.Log.Errorf("Failed to publish to NSQ: %v", err)
			}
		}()
		if err := t.dep.Pusher.Publish(map[string]string{"invoice": invoice, "status": "success"}); err != nil {
			t.dep.Log.Errorf("Failed to publish to PusherJs: %v", err)
		}
	case "cancel":
		if err := t.repo.UpdateStatusTrasansaction(t.dep.Db.WithContext(ctx), invoice, status); err != nil {
			return err
		}
		go func() {
			if err := t.dep.Nsq.Publish("3", encodeddata); err != nil {
				t.dep.Log.Errorf("Failed to publish to NSQ: %v", err)
			}
		}()
		if err := t.dep.Pusher.Publish(map[string]string{"invoice": invoice, "status": "canceled"}); err != nil {
			t.dep.Log.Errorf("Failed to publish to PusherJs: %v", err)
		}
	case "refund":
		if err := t.repo.UpdateStatusTrasansaction(t.dep.Db.WithContext(ctx), invoice, status); err != nil {
			return err
		}
		go func() {
			if err := t.dep.Nsq.Publish("4", encodeddata); err != nil {
				t.dep.Log.Errorf("Failed to publish to NSQ: %v", err)
			}
		}()
		if err := t.dep.Pusher.Publish(map[string]string{"invoice": invoice, "status": "refund", "reason": "out of quota"}); err != nil {
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
		Invoice:       invoice,
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

func (t *transaction) GetByStatus(ctx context.Context, uid int, status string) (*entity.Response, error) {
	if status != "paid" && status != "pending" && status != "" {
		return nil, errorr.NewBad("Data not found")
	}
	data, err := t.repo.GetByStatus(t.dep.Db.WithContext(ctx), uid, status)
	if err != nil {
		return nil, err
	}
	restrx := []entity.Transaction{}
	trxchan := make(chan entity.Transaction)
	wg := &sync.WaitGroup{}
	go func(wg *sync.WaitGroup, datachan <-chan entity.Transaction) {
		for val := range datachan {
			restrx = append(restrx, val)
			wg.Done()
		}
	}(wg, trxchan)
	wg.Add(len(data))
	for _, val := range data {
		res := entity.Transaction{
			Invoice:   val.Invoice,
			EventName: val.Event.Name,
		}
		trxchan <- res
	}
	wg.Wait()
	close(trxchan)
	response := entity.Response{
		Data: restrx,
	}
	return &response, nil
}

func (t *transaction) GetTickets(ctx context.Context, invoice string, uid int) (*entity.Response, error) {
	data, err := t.repo.GetTicketByInvoice(t.dep.Db.WithContext(ctx), invoice, uid)
	if err != nil {
		return nil, err
	}
	var tickets []entity.TicketTransaction
	ticketchan := make(chan entity.TicketTransaction)
	wg := &sync.WaitGroup{}
	go func(wg *sync.WaitGroup, ticketchan <-chan entity.TicketTransaction) {
		for val := range ticketchan {
			for i := 0; i < val.Qty; i++ {
				tickets = append(tickets, val)
			}
			wg.Done()
		}
	}(wg, ticketchan)
	wg.Add(len(data.TransactionItems))
	for _, val := range data.TransactionItems {
		ticket := entity.TicketTransaction{
			TicketType: val.Type.Name,
			EventName:  data.Event.Name,
			Location:   data.Event.Location,
			Date:       data.Event.StartDate,
			HostedBy:   data.Event.HostedBy,
			Qty:        val.Qty,
		}
		ticketchan <- ticket
	}
	wg.Wait()
	close(ticketchan)
	res := entity.Response{
		Data: tickets,
	}
	return &res, nil
}
