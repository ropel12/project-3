package repository

import (
	entity "github.com/ropel12/project-3/app/features/event"
	"github.com/ropel12/project-3/errorr"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	event struct {
		log *logrus.Logger
	}
	EventRepo interface {
		Create(db *gorm.DB, data entity.Event) (*int, error)
	}
)

func NewEventRepo(log *logrus.Logger) EventRepo {
	return &event{log}
}
func (e *event) Create(db *gorm.DB, data entity.Event) (*int, error) {
	if err := db.Create(&data).Error; err != nil {
		return nil, errorr.NewInternal("cannot create event")
	}
	id := int(data.ID)
	return &id, nil
}
