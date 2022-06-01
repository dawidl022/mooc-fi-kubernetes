package hashgen

import (
	"math/rand"
	"strings"
	"time"
)

func GenRandomHash(length int) string {
	rand.Seed(time.Now().Unix())
	builder := strings.Builder{}
	for i := 0; i < length; i++ {
		r := rand.Intn(36)
		var c rune
		if r >= 26 {
			c = '0' + rune(r-26)
		} else {
			c = 'a' + rune(r)
		}
		builder.WriteRune(c)
	}
	return builder.String()
}
