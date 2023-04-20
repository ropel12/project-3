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
		Image    string
	}
	LoginReq struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	RegisterReq struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
		Name     string `json:"name" validate:"required"`
		Address  string `json:"address" validate:"required"`
	}
)
