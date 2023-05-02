package container

import (
	"context"
	"encoding/csv"
	"io"
	"os"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/nsqio/go-nsq"
	"github.com/pusher/pusher-http-go/v5"
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
	if err := Container.Provide(NewMidtrans); err != nil {
		panic(err)
	}
	if err := Container.Provide(NewNSQ); err != nil {
		panic(err)
	}
	if err := Container.Provide(NewPusher); err != nil {
		panic(err)
	}
	if err := Container.Provide(NewValidation); err != nil {
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

func NewMidtrans(cfg *config.Config) *pkg.Midtrans {
	return &pkg.Midtrans{
		Midtrans: coreapi.Client{
			ServerKey:  cfg.Midtrans.ServerKey,
			ClientKey:  cfg.Midtrans.ClientKey,
			Env:        midtrans.EnvironmentType(cfg.Midtrans.Env),
			HttpClient: midtrans.GetHttpClient(midtrans.EnvironmentType(cfg.Midtrans.Env)),
			Options: &midtrans.ConfigOptions{
				PaymentOverrideNotification: &cfg.Midtrans.URLHandler,
				PaymentAppendNotification:   &cfg.Midtrans.URLHandler,
			},
		},
		ExpDuration: cfg.Midtrans.ExpiryDuration,
		ExpUnit:     cfg.Midtrans.Unit,
	}
}

func NewNSQ(conf *config.Config) (np *pkg.NSQProducer, err error) {
	np = &pkg.NSQProducer{}
	np.Env = conf.NSQ
	nsqConfig := nsq.NewConfig()
	np.Producer, err = nsq.NewProducer(np.Env.Host+":"+np.Env.Port, nsqConfig)
	if err != nil {
		return nil, err
	}

	return np, nil
}
func NewPusher(conf *config.Config) (ps *pkg.Pusher) {
	ps = &pkg.Pusher{}
	ps.Env = conf.Pusher
	ps.Client = &pusher.Client{
		AppID:   ps.Env.AppId,
		Key:     ps.Env.Key,
		Secret:  ps.Env.Secret,
		Cluster: ps.Env.Cluster,
		Secure:  ps.Env.Secure,
	}
	return ps
}

func NewValidation() (*pkg.Validation, error) {
	badwords := make(map[string]struct{})
	wd, _ := os.Getwd()
	file, err := os.Open(wd + "/pkg/badword.csv")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		badwords[record[0]] = struct{}{}
	}
	return &pkg.Validation{Badwords: badwords}, nil
}
