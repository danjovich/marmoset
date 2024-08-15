package utils

// safely gets arr[i], with bound checks and sentinel boolean
func At[T any](arr []T, i int) (T, bool) {
	if i >= 0 && i < len(arr) {
		return arr[i], true
	}

	var defaultT T
	return defaultT, false
}

// returns a slice of the keys of a map
// func Keys[K comparable, V any](aMap map[K]V) []K {
// 	keys := make([]K, len(aMap))

// 	i := 0
// 	for k := range aMap {
// 		keys[i] = k
// 		i++
// 	}

// 	return keys
// }
