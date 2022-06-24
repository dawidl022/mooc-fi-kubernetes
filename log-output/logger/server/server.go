package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Status struct {
	String string
}

func StartServer(s *Status) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, s.String)
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
