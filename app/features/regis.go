package features

import (
	eventrepo "github.com/ropel12/project-3/app/features/event/repository"
	eventserv "github.com/ropel12/project-3/app/features/event/service"
	trxrepo "github.com/ropel12/project-3/app/features/transaction/repository"
	trxserv "github.com/ropel12/project-3/app/features/transaction/service"
	userrepo "github.com/ropel12/project-3/app/features/user/repository"
	userserv "github.com/ropel12/project-3/app/features/user/service"
	"go.uber.org/dig"
)

func RegisterRepo(C *dig.Container) error {
	if err := C.Provide(userrepo.NewUserRepo); err != nil {
		return err
	}
	if err := C.Provide(eventrepo.NewEventRepo); err != nil {
		return err
	}
	if err := C.Provide(trxrepo.NewTransactionRepo); err != nil {
		return err
	}
	return nil
}

func RegisterService(C *dig.Container) error {
	if err := C.Provide(userserv.NewUserService); err != nil {
		return err
	}
	if err := C.Provide(eventserv.NewEventService); err != nil {
		return err
	}
	if err := C.Provide(trxserv.NewTransactionService); err != nil {
		return err
	}

	return nil
}
