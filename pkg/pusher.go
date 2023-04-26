package pkg

import (
	"github.com/pusher/pusher-http-go/v5"
	"github.com/ropel12/project-3/config"
)

type Pusher struct {
	Client *pusher.Client
	Env    config.PusherConfig
}

func (p *Pusher) Publish(data any) error {
	return p.Client.Trigger(p.Env.Channel, p.Env.Event, data)
}
