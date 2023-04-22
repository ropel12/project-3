package db

import (
	event "github.com/ropel12/project-3/app/features/event"
	user "github.com/ropel12/project-3/app/features/user"
	"github.com/ropel12/project-3/config"
)

func Migrate(c *config.Config) {
	db, err := config.GetConnection(c)
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(user.User{}, event.Event{}, event.Type{}); err != nil {
		panic(err)
	}
}
