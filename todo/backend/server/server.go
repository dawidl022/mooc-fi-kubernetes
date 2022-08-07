package server

import (
	"log"
	"net/http"
	"time"

	"github.com/dawidl022/mooc-fi-kubernetes/todo/config"
	"github.com/nats-io/nats.go"
)

func StartServer() {
	conf := config.GetConf()
	done := make(chan bool)
	go serve(conf)
	// wait for any immediate errors
	time.Sleep(time.Second)
	log.Printf("Server started in port %s\n", conf.Port)
	<-done
}

func serve(conf config.Conf) {
	db, err := initDB(&conf)
	if err != nil {
		log.Fatalf("failed to initialise database: %v", err)
	}

	nc, err := nats.Connect(conf.NatsUrl)
	if err != nil {
		log.Println(err)
	}
	defer nc.Close()

	router := http.NewServeMux()
	routes(router, db, nc)
	configuredRouter := LoggingMiddleware(*log.Default())(router)

	log.Fatal(http.ListenAndServe(":"+conf.Port, configuredRouter))
}
