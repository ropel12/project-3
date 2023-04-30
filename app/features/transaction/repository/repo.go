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
		CreateTransaction(db *gorm.DB, data entity.Transaction) error
		GetDetailUserById(db *gorm.DB, uid int) *entity.User
		GetDetailUserByInvoice(db *gorm.DB, invoice string) *entity.Transaction
		UpdateStatusTrasansaction(db *gorm.DB, invoice string, status string) error
		GetByInvoice(db *gorm.DB, invoice string, uid int) (*entity.Transaction, error)
		GetHistory(db *gorm.DB, uid int) ([]entity.Transaction, error)
		GetByStatus(db *gorm.DB, uid int, status string) ([]entity.Transaction, error)
		CheckQuota(db *gorm.DB, eventid, qty int) error
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

func (t *transaction) DeleteCart(db *gorm.DB, uid int) error {

	if rowaffc := db.Where("user_id = ?", uid).Delete(&entity.Carts{}); rowaffc.RowsAffected == 0 {
		t.log.Errorf("Error Db : %s", "Cart data does not exist")
		return errorr.NewBad("Cannot Create transaction")
	}
	return nil
}

func (t *transaction) CreateTransaction(db *gorm.DB, data entity.Transaction) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&data).Error; err != nil {
			t.log.Errorf("Error Db : %v", err)
			return errorr.NewInternal("Internal server error")
		}
		if rowaffc := db.Where("user_id = ?", data.UserID).Delete(&entity.Carts{}); rowaffc.RowsAffected == 0 {
			t.log.Errorf("Error Db : %s", "Cart data does not exist")
			return errorr.NewBad("Cannot Create transaction")
		}
		return nil
	})
}

func (t *transaction) GetDetailUserById(db *gorm.DB, uid int) *entity.User {
	userdata := entity.User{}
	if err := db.Find(&userdata, uid).Error; err != nil {
		return nil
	}
	return &userdata
}

func (t *transaction) GetDetailUserByInvoice(db *gorm.DB, invoice string) *entity.Transaction {
	userdata := entity.Transaction{}
	userdata.Invoice = invoice
	if err := db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,name,email")
	}).Find(&userdata).Error; err != nil {
		return nil
	}
	return &userdata
}

func (t *transaction) UpdateStatusTrasansaction(db *gorm.DB, invoice string, status string) error {
	if err := db.Model(&entity.Transaction{}).Where("invoice = ?", invoice).Update("status", status).Error; err != nil {
		t.log.Errorf("Error db : %v", err)
		return errorr.NewInternal("Internal server error")
	}
	return nil
}

func (t *transaction) GetByInvoice(db *gorm.DB, invoice string, uid int) (*entity.Transaction, error) {
	res := entity.Transaction{}
	if err := db.Preload("TransactionItems").Where("invoice = ? AND user_id = ?", invoice, uid).Find(&res).Error; err != nil {
		t.log.Errorf("Error db : %v", err)
		return nil, errorr.NewInternal("Internal server error")
	}
	if res.Invoice == "" {
		return nil, errorr.NewBad("Data not found")
	}
	return &res, nil
}
func (t *transaction) GetHistory(db *gorm.DB, uid int) ([]entity.Transaction, error) {
	res := []entity.Transaction{}
	if err := db.Preload("Event", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Users", func(db *gorm.DB) *gorm.DB {
			return db.Select("id")
		}).Select("id,start_date,end_date,name,hosted_by,image,location")
	}).Where("transactions.user_id = ? AND transactions.status='paid' AND e.start_date < NOW()", uid).Joins("join events e on e.id = transactions.event_id").Find(&res).Error; err != nil {
		t.log.Errorf("error db : %v", err)
		return nil, errorr.NewInternal("Internal server error")
	}
	if len(res) == 0 {
		return nil, errorr.NewBad("data not found")
	}
	return res, nil
}

func (t *transaction) CheckQuota(db *gorm.DB, eventid, qty int) error {
	var event entity.Event
	if err := db.First(&event, eventid).Error; err != nil {
		t.log.Printf("Error Db: %v", err)
		return errorr.NewInternal("Internal server error")
	}
	if qty > event.Quota {
		return errorr.NewBad("You have exceeded the quota")
	}
	return nil
}

func (t *transaction) GetByStatus(db *gorm.DB, uid int, status string) ([]entity.Transaction, error) {
	res := []entity.Transaction{}
	if err := db.Preload("Event", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,name")
	}).Where("user_id = ? AND status=?", uid, status).Find(&res).Error; err != nil {
		t.log.Errorf("error db : %v", err)
		return nil, errorr.NewInternal("Internal server error")
	}
	if len(res) == 0 {
		return nil, errorr.NewBad("data not found")
	}
	return res, nil
}
