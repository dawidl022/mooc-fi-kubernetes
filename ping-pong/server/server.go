package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func StartServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	counter := 0
	writeRequestCount(counter)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong %d", counter)
		counter++
		writeRequestCount(counter)
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func writeRequestCount(count int) {
	err := os.WriteFile("stats/ping_count", []byte(strconv.Itoa(count)), 0644)
	if err != nil {
		log.Println(err)
	}
}
