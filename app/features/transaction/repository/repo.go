package repository

import (
	entity "github.com/ropel12/project-3/app/entities"
	"github.com/ropel12/project-3/errorr"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	transaction struct {
		log *logrus.Logger
	}

	TransactionRepo interface {
		Create(db *gorm.DB, cart entity.Carts) error
		GetCart(db *gorm.DB, uid int) ([]entity.Carts, error)
	}
)

func NewTransactionRepo(log *logrus.Logger) TransactionRepo {
	return &transaction{log: log}
}

func (t *transaction) Create(db *gorm.DB, cart entity.Carts) error {
	res := db.Where("user_id = ? AND type_id = ? AND deleted_at IS NULL", cart.UserID, cart.TypeID).FirstOrCreate(&cart)
	if res.Error != nil {
		t.log.Errorf("error db : %v", res.Error)
		return errorr.NewInternal("Internal server error")
	}
	if res.RowsAffected == 0 {
		cart.Qty += 1
		if err := db.Model(&cart).Where("user_id = ? AND type_id = ? AND deleted_at IS NULL", cart.UserID, cart.TypeID).Update("qty", cart.Qty).Error; err != nil {
			t.log.Errorf("error db : %v", err)
			return errorr.NewInternal("Internal server error")
		}
	}
	return nil
}

func (t *transaction) GetCart(db *gorm.DB, uid int) ([]entity.Carts, error) {
	carts := []entity.Carts{}
	if err := db.Preload("Type").Where("user_id = ?", uid).Find(&carts).Error; err != nil {
		t.log.Errorf("error db : %v", err)
		return nil, errorr.NewInternal("Internal Server")
	}
	if len(carts) == 0 {
		return nil, errorr.NewBad("Data Not Found")
	}
	return carts, nil
}
