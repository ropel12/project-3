package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
	entity "github.com/ropel12/project-3/app/entities"
	"github.com/ropel12/project-3/helper"

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
		GetAll(db *gorm.DB, rds *redis.Client, limit int, offset int, search string) ([]*entity.Event, int, error)
		GetById(db *gorm.DB, id int) (*entity.Event, error)
		Update(db *gorm.DB, event entity.Event) (*entity.Event, error)
		CreateComment(db *gorm.DB, comment entity.UserComments) (*entity.UserComments, error)
		CreateTicket(db *gorm.DB, comment entity.Type) (*entity.Type, error)
		DeleteTicket(db *gorm.DB, id int) (*entity.Type, error)
		JoinEvent(db *gorm.DB, participant entity.Participants) (*entity.Participants, error)
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
		if op := rds.Set(ctx, key, dbval, time.Duration(2)*time.Minute); op.Err() != nil {
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
	var eventtrx int64
	db.Model(&entity.Transaction{}).Where("event_id=?", id).Count(&eventtrx)
	if eventtrx > 0 {
		return errorr.NewBad("Cannot delete event")
	}
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

func (e *event) GetAll(db *gorm.DB, rds *redis.Client, limit int, offset int, search string) ([]*entity.Event, int, error) {
	var (
		dbres    []*entity.Event
		redisres []byte
	)
	var count int64
	search = "%" + search + "%"
	db.Model(&entity.Event{}).Where("start_date > NOW() AND quota > 0 AND (name like ? or start_date like ? or duration like ? or location like ? or hosted_by like ?)", search, search, search, search, search).Count(&count)
	ctx, redisCancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer redisCancel()
	key := fmt.Sprintf("eventall:limit:%d:offset:%d:search:%s", limit, offset, search)
	err := rds.Get(ctx, key).Scan(&redisres)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			e.log.Errorf("Error redis : %v", err)
			return nil, 0, errorr.NewInternal("Server internal Error")
		}
		if err := db.Preload("Users").Where("start_date > NOW() AND quota > 0 AND (name like ? or start_date like ? or duration like ? or location like ? or hosted_by like ?)", search, search, search, search, search).Order("id DESC").Limit(limit).Offset(offset).Find(&dbres).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, 0, errorr.NewBad("Data not found")
			}
			e.log.Errorf("Error db : %v", err)
			return nil, 0, errorr.NewInternal("Server internal Error")
		}
		dbval, _ := json.Marshal(dbres)
		if op := rds.Set(ctx, key, dbval, time.Duration(2)*time.Minute); op.Err() != nil {
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
func (e *event) GetById(db *gorm.DB, id int) (*entity.Event, error) {
	var (
		dbres entity.Event
	)
	if err := db.Preload("UserComments", func(db *gorm.DB) *gorm.DB {
		return db.Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id,name,image")
		})
	}).Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Select("id,name,image")
	}).Preload("Types").First(&dbres, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorr.NewBad("Data not found")
		}
		e.log.Errorf("Error db : %v", err)
		return nil, errorr.NewInternal("Server internal Error")
	}

	return &dbres, nil
}

func (e *event) Update(db *gorm.DB, newdata entity.Event) (*entity.Event, error) {
	data := entity.Event{}
	if err := db.First(&data, newdata.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorr.NewBad("Data not found")
		}
		e.log.Errorf("error db :%v", err)
		return nil, errorr.NewInternal("Internal Server Error")
	}
	t, err := time.Parse(time.RFC3339, data.StartDate)
	t2, err2 := time.Parse(time.RFC3339, data.EndDate)
	if err != nil || err2 != nil {
		e.log.Error("[ERROR]When parse time, %v", err)
	}
	data.StartDate = t.Format("2006-01-02 15:04:05")
	data.EndDate = t2.Format("2006-01-02 15:04:05")
	if newdata.Duration != 0 {
		data.EndDate = helper.GenerateEndTime(data.StartDate, newdata.Duration)
	}
	v := reflect.ValueOf(newdata)
	n := reflect.ValueOf(&data).Elem()
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Interface().(type) {
		case string:
			val := v.Field(i).Interface().(string)
			if val != "" {
				n.Field(i).SetString(val)
			}
		case float32:
			val := v.Field(i).Interface().(float32)
			if val != 0 {
				n.Field(i).SetFloat(float64(val))
			}
		case int:
			val := v.Field(i).Interface().(int)
			if val != 0 {
				n.Field(i).SetInt(int64(val))
			}
		}
	}
	err = db.Transaction(func(db *gorm.DB) error {
		if err := db.Save(&data).Error; err != nil {
			e.log.Errorf("error db :%v", err)
			return errorr.NewInternal("Internal Server Error")
		}
		if len(newdata.Types) > 0 {
			for _, val := range newdata.Types {
				if err := db.Model(&entity.Type{}).Where("id=?", val.ID).Updates(val).Error; err != nil {
					e.log.Errorf("error db :%v", err)
					return errorr.NewInternal("Internal Server Error")
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &newdata, nil
}

func (e *event) CreateComment(db *gorm.DB, comment entity.UserComments) (*entity.UserComments, error) {

	if err := db.Create(&comment).Error; err != nil {
		return nil, errorr.NewInternal("Internal Server error")
	}
	return &comment, nil
}

func (e *event) CreateTicket(db *gorm.DB, comment entity.Type) (*entity.Type, error) {

	if err := db.Create(&comment).Error; err != nil {
		return nil, errorr.NewInternal("Internal Server error")
	}
	return &comment, nil
}

func (e *event) DeleteTicket(db *gorm.DB, id int) (*entity.Type, error) {
	res := entity.Type{}
	err := db.Transaction(func(db *gorm.DB) error {
		if err := db.Find(&res, id).Error; err != nil {
			e.log.Errorf("Error db: %v", err)
			return errorr.NewInternal("Internal server error")
		}
		if res.ID == 0 {
			return errorr.NewBad("TicketId doesn't exist")
		}
		if err := db.Delete(&res).Error; err != nil {
			e.log.Errorf("Error db: %v", err)
			return errorr.NewInternal("Internal server error")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (e *event) JoinEvent(db *gorm.DB, participant entity.Participants) (*entity.Participants, error) {
	var count int64
	db.Model(&entity.Participants{}).Where("user_id=? AND event_id=?", participant.UserID, participant.EventID).Count(&count)
	if count > 0 {
		return nil, errorr.NewBad("you have already joined the event")
	}
	var eventid int
	if err := db.Model(&entity.Transaction{}).Where("user_id=? AND event_id=? AND status='paid'", participant.UserID, participant.EventID).Select("event_id").Scan(&eventid).Error; err != nil {
		e.log.Errorf("error db: %v", err)
		return nil, errorr.NewInternal("Internal server error")
	}
	if eventid == 0 {
		return nil, errorr.NewBad("Cannot join the event")
	}
	if err := db.Create(&participant).Error; err != nil {
		return nil, errorr.NewInternal("Internal server error")
	}
	return &participant, nil
}
