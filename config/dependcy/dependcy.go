package dependcy

import (
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/ropel12/project-3/config"
	"github.com/ropel12/project-3/pkg"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type Depend struct {
	dig.In
	Db     *gorm.DB
	Config *config.Config
	Echo   *echo.Echo
	Log    *logrus.Logger
	Gcp    *pkg.StorageGCP
	Rds    *redis.Client
	Mds    *pkg.Midtrans
	Nsq    *pkg.NSQProducer
}
