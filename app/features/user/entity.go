package user

import (
	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		Name     string
		Email    string
		Password string
		Address  string
	}
	LoginReq struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
)
