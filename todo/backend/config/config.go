package config

import (
	"fmt"
	"net/url"
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
		DatabaseUrl: getDbUrl(),
		Port:        port,
	}
}

func getDbUrl() string {
	host := url.PathEscape(os.Getenv("DB_HOST"))
	name := url.PathEscape(os.Getenv("DB_NAME"))
	user := url.PathEscape(os.Getenv("DB_USER"))
	pass := url.PathEscape(os.Getenv("DB_PASSWORD"))

	fmt.Println(os.Getenv("DB_HOST"))
	fmt.Println(host)

	return fmt.Sprintf("postgresql://%s:%s@%s:5432/%s?sslmode=disable", user, pass, host, name)
}
