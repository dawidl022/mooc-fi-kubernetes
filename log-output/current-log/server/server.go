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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		current_log, err := os.ReadFile("logs/timestamp_hash.log")
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read log file: %v", err), http.StatusInternalServerError)
			return
		}
		ping_count, err := os.ReadFile("stats/ping_count")
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read ping count file: %v", err), http.StatusInternalServerError)
			return
		}

		w.Write(current_log)
		fmt.Fprintf(w, "Ping / Pongs: %s", ping_count)
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
