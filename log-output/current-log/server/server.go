package server

import (
	"fmt"
	"io"
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
		resp, err := http.Get("http://ping-pong-svc/ping-count")
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get ping count: %v", err), http.StatusInternalServerError)
			return
		}
		ping_count, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to parse ping-count body %v", err), http.StatusInternalServerError)
		}

		fmt.Fprintln(w, os.Getenv("MESSAGE"))
		w.Write(current_log)
		fmt.Fprintf(w, "Ping / Pongs: %s", ping_count)
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
