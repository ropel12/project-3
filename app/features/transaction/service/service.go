package service

import (
	"context"

	"github.com/go-playground/validator"
	entity2 "github.com/ropel12/project-3/app/entities"
	entity "github.com/ropel12/project-3/app/features/transaction"
	"github.com/ropel12/project-3/app/features/transaction/repository"
	"github.com/ropel12/project-3/config/dependcy"
	"github.com/ropel12/project-3/errorr"
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
