package service

import (
	"context"
	"mime/multipart"

	"github.com/go-playground/validator"
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
	data := entity.Event{
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
