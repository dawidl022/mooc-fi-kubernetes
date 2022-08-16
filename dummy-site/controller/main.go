package main

import (
	"log"

	"github.com/dawidl022/mooc-fi-kubernetes/dummy-site/controller/gateway"
)

func main() {
	ctrl, err := gateway.NewController()
	if err != nil {
		log.Fatalf("failed to initialise controller: %v", err)
	}
	log.Fatal(ctrl.Start())
}
