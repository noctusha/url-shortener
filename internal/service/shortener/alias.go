package shortener

import (
	"math/rand"
	"strings"
	"time"
)

func Random(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	chars := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	b := strings.Builder{}
	b.Grow(size)
	for i := 0; i < size; i++ {
		b.WriteByte(chars[rnd.Intn(len(chars))])
	}
	return b.String()
}
