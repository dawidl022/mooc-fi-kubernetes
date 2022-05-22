package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func genRandomString(length int) string {
	rand.Seed(time.Now().Unix())
	builder := strings.Builder{}
	for i := 0; i < length; i++ {
		builder.WriteRune('a' + rune(rand.Intn(26)))
	}
	return builder.String()
}

func main() {
	randString := genRandomString(36)
	for {
		fmt.Printf("%s: %s\n", time.Now().Format(time.RFC3339), randString)
		time.Sleep(time.Second * 5)
	}
}
