package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	entity "github.com/ropel12/project-3/app/entities"

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
		GetByUid(db *gorm.DB, rds *redis.Client, uid int, limit int, offset int) ([]*entity.Event, int, error)
		Delete(db *gorm.DB, id int, uid int) error
		GetAll(db *gorm.DB, rds *redis.Client, limit int, offset int) ([]*entity.Event, int, error)
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

func (e *event) GetByUid(db *gorm.DB, rds *redis.Client, uid int, limit int, offset int) ([]*entity.Event, int, error) {
	var (
		dbres    []*entity.Event
		redisres []byte
	)
	var count int64
	db.Model(&entity.Event{}).Where("user_id = ?", uid).Count(&count)
	ctx, redisCancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer redisCancel()
	key := fmt.Sprintf("eventuser:uid:%d:limit:%d:offset:%d", uid, limit, offset)
	err := rds.Get(ctx, key).Scan(&redisres)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			e.log.Errorf("Error redis : %v", err)
			return nil, 0, errorr.NewInternal("Server internal Error")
		}
		if err := db.Preload("Users").Order("id DESC").Where("user_id = ?", uid).Limit(limit).Offset(offset).Find(&dbres).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, 0, errorr.NewBad("Data not found")
			}
			e.log.Errorf("Error db : %v", err)
			return nil, 0, errorr.NewInternal("Server internal Error")
		}
		dbval, _ := json.Marshal(dbres)
		if op := rds.Set(ctx, key, dbval, time.Duration(2)*time.Hour); op.Err() != nil {
			e.log.Errorf("error set redis val : %v", err)
			return nil, 0, errorr.NewInternal("Server internal Error")
		}
		return dbres, int(count), nil
	}
	if err := json.Unmarshal([]byte(redisres), &dbres); err != nil {
		e.log.Errorf("error unmarshal redis val: %v", err)
		return nil, 0, errorr.NewInternal("Server internal Error")
	}
	return dbres, int(count), nil
}

func (e *event) Delete(db *gorm.DB, id int, uid int) error {
	eventdata := entity.Event{}
	db.Find(&eventdata, id)
	if eventdata.Name == "" {
		return errorr.NewBad("Id not found")
	}
	if eventdata.UserID != uint(uid) {
		return errorr.NewBad("Cannot delete event")
	}
	if err := db.Delete(&entity.Event{}, id).Error; err != nil {
		e.log.Errorf("error database : %v", err)
		return errorr.NewInternal("Internal Server Error")
	}
	return nil
}

func (e *event) GetAll(db *gorm.DB, rds *redis.Client, limit int, offset int) ([]*entity.Event, int, error) {
	var (
		dbres    []*entity.Event
		redisres []byte
	)
	var count int64
	db.Model(&entity.Event{}).Count(&count)
	ctx, redisCancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer redisCancel()
	key := fmt.Sprintf("eventall:limit:%d:offset:%d", limit, offset)
	err := rds.Get(ctx, key).Scan(&redisres)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			e.log.Errorf("Error redis : %v", err)
			return nil, 0, errorr.NewInternal("Server internal Error")
		}
		if err := db.Preload("Users").Order("id DESC").Limit(limit).Offset(offset).Find(&dbres).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, 0, errorr.NewBad("Data not found")
			}
			e.log.Errorf("Error db : %v", err)
			return nil, 0, errorr.NewInternal("Server internal Error")
		}
		dbval, _ := json.Marshal(dbres)
		if op := rds.Set(ctx, key, dbval, time.Duration(2)*time.Hour); op.Err() != nil {
			e.log.Errorf("error set redis val : %v", err)
			return nil, 0, errorr.NewInternal("Server internal Error")
		}
		return dbres, int(count), nil
	}
	if err := json.Unmarshal([]byte(redisres), &dbres); err != nil {
		e.log.Errorf("error unmarshal redis val: %v", err)
		return nil, 0, errorr.NewInternal("Server internal Error")
	}
	return dbres, int(count), nil
}
