package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dawidl022/mooc-fi-kubernetes/hashgen"
	"github.com/dawidl022/mooc-fi-kubernetes/log-output/logger/server"
)

func main() {
	s := &server.Status{}
	go logStatus(s)
	server.StartServer(s)
}

func logStatus(s *server.Status) {
	randString := hashgen.GenRandomHash(40)
	for {
		s.String = fmt.Sprintf("%s: %s\n", time.Now().Format(time.RFC3339), randString)
		fmt.Print(s.String)

		err := os.WriteFile("logs/timestamp_hash.log", []byte(s.String), 0644)
		if err != nil {
			fmt.Printf("failed to log to file: %v", err)
		}
		time.Sleep(time.Second * 5)
	}
}
