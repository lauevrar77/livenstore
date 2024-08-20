package bench

func SampleElement[T any](arr []T) T {
	id := RandomInt(uint(len(arr)))
	return arr[id]
}
