package features

import (
	"github.com/ropel12/project-3/app/features/user/repository"
	"github.com/ropel12/project-3/app/features/user/service"
	"go.uber.org/dig"
)

func RegisterRepo(C *dig.Container) error {
	if err := C.Provide(repository.NewUserRepo); err != nil {
		return err
	}
	return nil
}

func RegisterService(C *dig.Container) error {
	if err := C.Provide(service.NewUserService); err != nil {
		return err
	}
	return nil
}
