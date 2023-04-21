package repository

import (
	"reflect"

	entity "github.com/ropel12/project-3/app/features/user"
	"github.com/ropel12/project-3/errorr"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	user struct {
		log *logrus.Logger
	}
	UserRepo interface {
		Create(db *gorm.DB, user entity.User) error
		FindByEmail(db *gorm.DB, email string) (*entity.User, error)
		Delete(db *gorm.DB, user entity.User) error
		Update(db *gorm.DB, user entity.User) (*entity.User, error)
	}
)

func NewUserRepo(log *logrus.Logger) UserRepo {
	return &user{log}
}

func (u *user) Create(db *gorm.DB, user entity.User) error {
	if err := db.Create(&user).Error; err != nil {
		u.log.Errorf("error Db: %v ", err)
		return errorr.NewInternal(err.Error())
	}
	return nil
}

func (u *user) FindByEmail(db *gorm.DB, email string) (*entity.User, error) {
	res := entity.User{}
	err := db.Where("email = ?", email).Find(&res).Error
	if res.Email == "" {
		return nil, errorr.NewBad("Email not registered")
	}
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			u.log.Errorf("error Db: %v", err)
			return nil, errorr.NewInternal(err.Error())
		} else {
			u.log.Errorf("error Db: %v", err)
			return nil, errorr.NewBad(err.Error())
		}
	}
	return &res, nil
}
func (u *user) Delete(db *gorm.DB, user entity.User) error {
	if err := db.Delete(&user).Error; err != nil {
		u.log.Errorf("error Db: %v", err)
		return errorr.NewInternal(err.Error())
	}
	return nil
}
func (u *user) Update(db *gorm.DB, user entity.User) (*entity.User, error) {
	newdata := entity.User{}
	if err := db.First(&newdata, user.ID).Error; err == gorm.ErrRecordNotFound {
		u.log.Errorf("error Db: %v", err)
		return nil, errorr.NewBad("Id Not Found")
	}
	v := reflect.ValueOf(user)
	n := reflect.ValueOf(&newdata).Elem()
	for i := 0; i < v.NumField(); i++ {
		if val, ok := v.Field(i).Interface().(string); ok {
			if val != "" {
				n.Field(i).SetString(val)
			}
		}
	}
	if err := db.Save(&newdata).Error; err != nil {
		u.log.Errorf("error Db: %v")
		return nil, errorr.NewInternal("error update user")
	}
	return &newdata, nil
}

func (u *user) GetById(db *gorm.DB, id int) (*entity.User, error) {
	if err := db.Delete(&user).Error; err != nil {
		u.log.Errorf("error Db: %v", err)
		return nil, errorr.NewInternal(err.Error())
	}
	res := entity.User{}
	err := db.Find(&res).Error
	if res.Email == "" {
		return nil, errorr.NewBad("Email not registered")
	}
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			u.log.Errorf("error Db: %v", err)
			return nil, errorr.NewInternal(err.Error())
		} else {
			u.log.Errorf("error Db: %v", err)
			return nil, errorr.NewBad(err.Error())
		}
	}
	return &res, nil

}
