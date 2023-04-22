package user

import (
	event "github.com/ropel12/project-3/app/features/event"
	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		Name     string `gorm:"not null"`
		Email    string `gorm:"not null"`
		Password string `gorm:"not null"`
		Address  string `gorm:"not null"`
		Image    string `gorm:"not null"`
		Events   []event.Event
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
	UpdateReq struct {
		Id       int
		Email    string `form:"email"`
		Password string `form:"password"`
		Name     string `form:"name" `
		Address  string `form:"address"`
		Image    string `form:"image"`
	}
)
