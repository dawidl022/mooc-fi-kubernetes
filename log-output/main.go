package main

import (
	"fmt"
	"time"

	"github.com/dawidl022/mooc-fi-kubernetes/hashgen"
	"github.com/dawidl022/mooc-fi-kubernetes/log-output/server"
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
		time.Sleep(time.Second * 5)
	}
}
