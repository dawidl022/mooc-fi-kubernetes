package main

import (
	"fmt"
	"time"

	"github.com/dawidl022/mooc-fi-kubernetes/hashgen"
)

func main() {
	randString := hashgen.GenRandomHash(40)
	for {
		fmt.Printf("%s: %s\n", time.Now().Format(time.RFC3339), randString)
		time.Sleep(time.Second * 5)
	}
}
