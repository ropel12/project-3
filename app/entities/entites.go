package entities

import "gorm.io/gorm"

type (
	Event struct {
		gorm.Model
		Name         string  `gorm:"type:varchar(30);not null"`
		StartDate    string  `gorm:"type:timestamp;not null"`
		Duration     float32 `gorm:"type:float;not null"`
		EndDate      string  `gorm:"type:timestamp;not null"`
		Quota        int     `gorm:"not null"`
		Location     string  `gorm:"not null"`
		Detail       string  `gorm:"not null"`
		Image        string  `gorm:"not null"`
		HostedBy     string  `gorm:"not null"`
		UserID       uint    `gorm:"not null"`
		Types        []Type
		Users        []User `gorm:"many2many:participants;"`
		UserComments []UserComments
	}

	Type struct {
		ID      uint   `gorm:"primaryKey;not null;autoIncrement" json:"id,omitempty"`
		Name    string `gorm:"not null" form:"name" json:"name,omitempty" validate:"required" `
		Price   int    `gorm:"not null" form:"price" json:"price,omitempty" validate:"required"`
		EventID uint   `gorm:"not null"`
		Carts   []Carts
	}

	User struct {
		gorm.Model   `json:"-"`
		Name         string `gorm:"not null" json:"name,omitempty"`
		Email        string `gorm:"not null" json:"email,omitempty"`
		Password     string `gorm:"not null" json:"password,omitempty"`
		Address      string `gorm:"not null" json:"address,omitempty"`
		Image        string `gorm:"not null" json:"Image,omitempty"`
		Events       []Event
		UserComments []UserComments
		Carts        []Carts
	}
	UserComments struct {
		UserID  uint   `json:"-"`
		EventID uint   `json:"-"`
		Comment string `json:"comment,omitempty"`
		User    User
	}
	Carts struct {
		UserID    uint
		TypeID    uint
		Qty       int
		Type      Type
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}
)
