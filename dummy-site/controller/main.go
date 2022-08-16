package main

import (
	"log"

	"github.com/dawidl022/mooc-fi-kubernetes/dummy-site/controller/processor"
)

func main() {
	a, err := processor.NewApplier()
	if err != nil {
		log.Fatal(err)
	}
	url := "https://www.wikipedia.org/"
	a.ApplyUntilDestroyed(processor.UrlToWebsite(url), url)
}
