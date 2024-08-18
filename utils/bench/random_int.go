package bench

import "math/rand"

func RandomInt(max uint) int {
	return rand.Intn(int(max))
}
