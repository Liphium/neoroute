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

// allOrderedSubsets returns all ordered subsets of the given slice.
func allOrderedSubsets[T any](items []T) [][]T {
	var result [][]T
	var current []T
	visited := make([]bool, len(items))

	var backtrack func()
	backtrack = func() {
		if len(current) > 0 {
			combination := make([]T, len(current))
			copy(combination, current)
			result = append(result, combination)
		}

		for i := 0; i < len(items); i++ {
			if !visited[i] {
				visited[i] = true
				current = append(current, items[i])

				backtrack()

				current = current[:len(current)-1]
				visited[i] = false
			}
		}
	}

	backtrack()
	return result
}
