package lib

func InArray[T comparable](val T, array []T) int {
	for i, v := range array {
		if v == val {
			return i
		}
	}
	return -1
}
