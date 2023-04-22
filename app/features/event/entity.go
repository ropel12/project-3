package event

import (
	"gorm.io/gorm"
)

type (
	Event struct {
		gorm.Model
		Name      string  `gorm:"type:varchar(30);not null"`
		StartDate string  `gorm:"type:timestamp;not null"`
		Duration  float32 `gorm:"type:float;not null"`
		EndDate   string  `gorm:"type:timestamp;not null"`
		Quota     int     `gorm:"not null"`
		Location  string  `gorm:"not null"`
		Detail    string  `gorm:"not null"`
		Image     string  `gorm:"not null"`
		HostedBy  string  `gorm:"not null"`
		UserID    uint    `gorm:"not null"`
		Types     []Type
	}

	Type struct {
		Id      uint   `gorm:"primaryKey;not null;autoIncrement"`
		Name    string `gorm:"not null" form:"name" json:"name" validate:"required" `
		Price   int    `gorm:"not null" form:"price" json:"price" validate:"required"`
		EventID uint   `gorm:"not null"`
	}

	ReqCreate struct {
		Name      string  `form:"name" validate:"required"`
		StartDate string  `form:"date" validate:"required"`
		Duration  float32 `form:"duration" validate:"required"`
		Details   string  `form:"details" validate:"required"`
		Quota     int     `form:"quota" validate:"required"`
		HostedBy  string  `form:"hosted_by" validate:"required"`
		Location  string  `form:"location" validate:"required"`
		Rtype     string  `form:"type" json:"type" validate:"required"`
		Types     []Type
		Image     string
		Uid       int
	}
)
