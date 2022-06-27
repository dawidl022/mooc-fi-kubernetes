package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dawidl022/mooc-fi-kubernetes/ping-pong/config"
	"gorm.io/gorm"
)

type server struct {
	inMemoryCounter int
	db              *gorm.DB
}

func StartServer() {
	conf := config.GetConf()

	db, err := initDB(&conf)
	if err != nil {
		log.Println("failed to initialise database")
	}
	s := server{db: db}

	counter := s.getRequestCount(db)
	writeRequestCount(counter)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		counter = s.getRequestCount(db)
		fmt.Fprintf(w, "pong %d", counter)

		writeRequestCount(counter)
		s.incrementRequestCount(db)
	})

	http.HandleFunc("/ping-count", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, strconv.Itoa(counter))
	})

	log.Fatal(http.ListenAndServe(":"+conf.Port, nil))
}

func writeRequestCount(count int) {
	err := os.WriteFile("stats/ping_count", []byte(strconv.Itoa(count)), 0644)
	if err != nil {
		log.Println(err)
	}
}
