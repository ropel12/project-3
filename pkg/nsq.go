package pkg

import (
	"github.com/nsqio/go-nsq"
	"github.com/ropel12/project-3/config"
	"github.com/ropel12/project-3/errorr"
)

type NSQProducer struct {
	Producer *nsq.Producer
	Env      config.NSQConfig
}

func (np *NSQProducer) Publish(Topic string, message []byte) error {
	switch Topic {
	case "1":
		return np.Producer.Publish(np.Env.Topic, message)
	case "2":
		return np.Producer.Publish(np.Env.Topic2, message)
	case "3":
		return np.Producer.Publish(np.Env.Topic3, message)
	}
	return errorr.NewBad("Topic not available")
}

func (np *NSQProducer) Stop() {
	np.Producer.Stop()
}
