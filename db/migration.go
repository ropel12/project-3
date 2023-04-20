package db

import (
	entity "github.com/ropel12/project-3/app/features/user"
	"github.com/ropel12/project-3/config"
)

func Migrate(c *config.Config) {
	db, err := config.GetConnection(c)
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(entity.User{}); err != nil {
		panic(err)
	}
}
