package config

import (
	"os"

	"github.com/nats-io/nats.go"
)

type Conf struct {
	DatabaseUrl string
	Port        string
	NatsUrl     string
}

func GetConf() Conf {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	natsUrl := os.Getenv("NATS_URL")
	if natsUrl == "" {
		natsUrl = nats.DefaultURL
	}

	return Conf{
		DatabaseUrl: os.Getenv("DB_URL"),
		Port:        port,
		NatsUrl:     natsUrl,
	}
}
