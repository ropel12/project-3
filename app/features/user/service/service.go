package service

import (
	"context"

	"github.com/go-playground/validator"
	entity "github.com/ropel12/project-3/app/features/user"
	"github.com/ropel12/project-3/app/features/user/repository"
	"github.com/ropel12/project-3/config/dependcy"
	"github.com/ropel12/project-3/errorr"
	"github.com/ropel12/project-3/helper"
)

type (
	user struct {
		repo      repository.UserRepo
		validator *validator.Validate
		dep       dependcy.Depend
	}
	UserService interface {
		Login(ctx context.Context, req entity.LoginReq) (int, error)
		Register(ctx context.Context, req entity.RegisterReq) error
	}
)

func NewUserService(repo repository.UserRepo, dep dependcy.Depend) UserService {
	return &user{repo: repo, dep: dep, validator: validator.New()}
}

func (u *user) Login(ctx context.Context, req entity.LoginReq) (int, error) {
	if err := u.validator.Struct(req); err != nil {
		u.dep.Log.Errorf("Error Service: %v", err)
		return 0, err
	}
	user, err := u.repo.FindByEmail(u.dep.Db.WithContext(ctx), req.Email)
	if err != nil {
		return 0, err
	}
	if err := helper.VerifyPassword(user.Password, req.Password); err != nil {
		u.dep.Log.Errorf("Error Service : %v", err)
		return 0, errorr.NewBad("Wrong password")
	}
	return int(user.ID), nil
}

func (u *user) Register(ctx context.Context, req entity.RegisterReq) error {
	if err := u.validator.Struct(req); err != nil {
		u.dep.Log.Errorf("Error service: %v", err)
		return errorr.NewBad("Request body not valid")
	}
	_, err := u.repo.FindByEmail(u.dep.Db.WithContext(ctx), req.Email)
	if err == nil {
		return errorr.NewBad("Email already registered")
	}
	passhash, err := helper.HashPassword(req.Password)
	if err != nil {
		u.dep.Log.Errorf("Erorr service: %v", err)
		return errorr.NewBad("Register failed")
	}
	data := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Address:  req.Address,
		Password: passhash,
		Image:    "default.jpg",
	}
	err = u.repo.Create(u.dep.Db.WithContext(ctx), data)
	if err != nil {
		return errorr.NewInternal("Failed to create account")
	}
	return nil
}
