package entities

import "gorm.io/gorm"

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
		Users     []User `gorm:"many2many:participants;"`
	}

	Type struct {
		Id      uint   `gorm:"primaryKey;not null;autoIncrement"`
		Name    string `gorm:"not null" form:"name" json:"name" validate:"required" `
		Price   int    `gorm:"not null" form:"price" json:"price" validate:"required"`
		EventID uint   `gorm:"not null"`
	}

	User struct {
		gorm.Model
		Name     string `gorm:"not null"`
		Email    string `gorm:"not null"`
		Password string `gorm:"not null"`
		Address  string `gorm:"not null"`
		Image    string `gorm:"not null"`
		Events   []Event
	}
)
