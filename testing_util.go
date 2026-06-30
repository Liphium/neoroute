package neoroute

// sameElements checks if two slices contain the same elements, regardless of order.
func sameElements[T comparable](slice1, slice2 []T) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	counts := make(map[T]int)
	for _, item := range slice1 {
		counts[item]++
	}

	for _, item := range slice2 {
		counts[item]--
		if counts[item] < 0 {
			return false
		}
	}

	return true
}
