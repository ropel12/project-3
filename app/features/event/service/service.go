package service

import (
	"context"
	"math"
	"mime/multipart"

	"github.com/go-playground/validator"
	entity2 "github.com/ropel12/project-3/app/entities"
	entity "github.com/ropel12/project-3/app/features/event"
	"github.com/ropel12/project-3/app/features/event/repository"
	"github.com/ropel12/project-3/config/dependcy"
	"github.com/ropel12/project-3/errorr"
	"github.com/ropel12/project-3/helper"
)

type (
	event struct {
		repo      repository.EventRepo
		validator *validator.Validate
		dep       dependcy.Depend
	}

	EventService interface {
		Create(ctx context.Context, req entity.ReqCreate, file multipart.File) (int, error)
		MyEvent(ctx context.Context, uid, limit, page int) (*entity.Response, error)
		Delete(ctx context.Context, id int, uid int) error
	}
)

func NewEventService(repo repository.EventRepo, dep dependcy.Depend) EventService {
	return &event{
		repo:      repo,
		dep:       dep,
		validator: validator.New(),
	}
}

func (e *event) Create(ctx context.Context, req entity.ReqCreate, file multipart.File) (int, error) {
	if err := e.validator.Struct(req); err != nil {
		e.dep.Log.Errorf("Error Service: %v", err)
		return 0, errorr.NewBad("Invalid and missing request body")
	}
	if err := e.dep.Gcp.UploadFile(file, req.Image); err != nil {
		return 0, errorr.NewBad(err.Error())
	}
	data := entity2.Event{
		Name:      req.Name,
		StartDate: req.StartDate,
		Duration:  req.Duration,
		EndDate:   helper.GenerateEndTime(req.StartDate, req.Duration),
		Quota:     req.Quota,
		Location:  req.Location,
		Detail:    req.Details,
		Image:     req.Image,
		HostedBy:  req.HostedBy,
		UserID:    uint(req.Uid),
		Types:     req.Types,
	}
	id, err := e.repo.Create(e.dep.Db.WithContext(ctx), data)
	if err != nil {
		return 0, err
	}

	return *id, nil
}

func (e *event) MyEvent(ctx context.Context, uid, limit, page int) (*entity.Response, error) {
	offset := (page - 1) * limit
	data, total, err := e.repo.GetByUid(e.dep.Db.WithContext(ctx), e.dep.Rds, uid, limit, offset)
	if err != nil {
		return nil, err
	}
	res := new(entity.Response)
	res.Page = page
	res.Limit = limit
	res.TotalPage = int(math.Ceil(float64(total) / float64(limit)))
	res.TotalData = total
	var datas []*entity.ResponseEvent
	for _, val := range data {
		newdata := new(entity.ResponseEvent)
		newdata.Id = int(val.ID)
		newdata.Name = val.Name
		newdata.Date = val.StartDate
		newdata.EndDate = val.EndDate
		newdata.Location = val.Location
		newdata.HostedBy = val.HostedBy
		newdata.Image = val.Image
		newdata.Participants = len(val.Users)
		datas = append(datas, newdata)
	}
	res.Data = datas
	return res, nil
}

func (e *event) Delete(ctx context.Context, id int, uid int) error {
	if err := e.repo.Delete(e.dep.Db.WithContext(ctx), id, uid); err != nil {
		return err
	}
	return nil
}
