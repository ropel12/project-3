package container

import (
	"context"
	"os"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
	feat "github.com/ropel12/project-3/app/features"
	"github.com/ropel12/project-3/config"
	"github.com/ropel12/project-3/pkg"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

var (
	Container = dig.New()
)

func RunAll() {
	Container := Container
	if err := Container.Provide(config.InitConfiguration); err != nil {
		panic(err)
	}
	if err := Container.Provide(config.GetConnection); err != nil {
		panic(err)
	}
	if err := Container.Provide(config.NewRedis); err != nil {
		panic(err)
	}
	if err := Container.Provide(echo.New); err != nil {
		panic(err)
	}
	if err := Container.Provide(NewLog); err != nil {
		panic(err)
	}
	if err := Container.Provide(NewStorage); err != nil {
		panic(err)
	}
	if err := feat.RegisterRepo(Container); err != nil {
		panic(err)
	}
	if err := feat.RegisterService(Container); err != nil {
		panic(err)
	}

}

func NewStorage(cfg *config.Config) (*pkg.StorageGCP, error) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", cfg.GCP.Credential)
	client, err := storage.NewClient(context.Background())
	if err != nil {
		return nil, err
	}
	return &pkg.StorageGCP{
		ClG:        client,
		ProjectID:  cfg.GCP.PRJID,
		BucketName: cfg.GCP.BCKNM,
		Path:       cfg.GCP.Path,
	}, nil
}

func NewLog() (*log.Logger, error) {
	var logger = log.New()
	file, _ := os.OpenFile("output.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	logger.SetOutput(file)
	logger.SetFormatter(&log.JSONFormatter{})
	return logger, nil
}
