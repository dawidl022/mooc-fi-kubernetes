package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func StartServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	counter := 0

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong %d", counter)
		counter++
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
