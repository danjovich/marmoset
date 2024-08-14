package utils

func At[T any](arr []T, i int) (T, bool) {
	if i >= 0 && i < len(arr) {
		return arr[i], true
	}

	var defaultT T
	return defaultT, false
}
