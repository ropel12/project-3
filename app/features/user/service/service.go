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
		return 0, errorr.NewBad("Password Salah")
	}
	return int(user.ID), nil
}
