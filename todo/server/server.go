package server

import (
	"log"
	"net/http"
	"os"
	"time"
)

func StartServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	done := make(chan bool)
	go serve(port)
	// wait for any immediate errors
	time.Sleep(time.Second)
	log.Printf("Server started in port %s\n", port)
	<-done
}

func serve(port string) {
	routes()
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
