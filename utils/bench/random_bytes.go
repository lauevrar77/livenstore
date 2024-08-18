package bench

import "math/rand"

func RandomBytes(s uint) []byte {
	b := make([]byte, s)
	for i := range b {
		b[i] = byte(rand.Intn(256))
	}
	return b
}
