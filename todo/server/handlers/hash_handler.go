package handlers

import (
	"fmt"
	"net/http"

	"github.com/dawidl022/mooc-fi-kubernetes/hashgen"
)

const HASH_LEN = 8

type hashHander struct {
	initialHash string
}

func NewHashHandler() *hashHander {
	return &hashHander{
		initialHash: hashgen.GenRandomHash(HASH_LEN),
	}
}

func (h *hashHander) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s: %s", h.initialHash, hashgen.GenRandomHash(HASH_LEN))
}
