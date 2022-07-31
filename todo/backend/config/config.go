package config

import (
	"os"
)

type Conf struct {
	DatabaseUrl string
	Port        string
}

func GetConf() Conf {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return Conf{
		DatabaseUrl: os.Getenv("DB_URL"),
		Port:        port,
	}
}
