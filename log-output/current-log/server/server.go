package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const pingPongUrl = "http://ping-pong-svc/"

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
		resp, err := http.Get(pingPongUrl + "ping-count")
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

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get(pingPongUrl + "healthz")
		if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
