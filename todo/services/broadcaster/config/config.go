package config

import (
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
)

type Config struct {
	NatsUrl             string
	SendMessageUrl      string
	SendMessageTemplate string
}

func GetConf() (*Config, error) {
	natsUrl := os.Getenv("NATS_URL")
	msgUrl := os.Getenv("SEND_MESSAGE_URL")
	tmpl := os.Getenv("SEND_MESSAGE_TEMPLATE")

	if natsUrl == "" {
		natsUrl = nats.DefaultURL
	}

	for _, ev := range []string{"SEND_MESSAGE_URL", "SEND_MESSAGE_TEMPLATE"} {
		if ev == "" {
			return nil, fmt.Errorf("environment variable %s not set", ev)
		}
	}
	return &Config{
		NatsUrl:             natsUrl,
		SendMessageUrl:      msgUrl,
		SendMessageTemplate: tmpl,
	}, nil
}
