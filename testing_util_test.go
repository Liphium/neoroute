package neoroute

import "testing"

func Test_sameElements(t *testing.T) {
	// Test case for ints
	t.Run("integers", func(t *testing.T) {
		tests := []struct {
			name   string
			slice1 []int
			slice2 []int
			want   bool
		}{
			{"identical", []int{1, 2, 3}, []int{1, 2, 3}, true},
			{"shuffled", []int{1, 2, 3}, []int{3, 2, 1}, true},
			{"different lengths", []int{1, 2}, []int{1, 2, 3}, false},
			{"mismatched values", []int{1, 2, 3}, []int{1, 2, 4}, false},
			{"duplicates mismatch", []int{1, 2, 2}, []int{1, 1, 2}, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := sameElements(tt.slice1, tt.slice2); got != tt.want {
					t.Errorf("sameElements() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	// Test case for strings
	t.Run("strings", func(t *testing.T) {
		a := []string{"apple", "banana"}
		b := []string{"banana", "apple"}
		if !sameElements(a, b) {
			t.Errorf("Expected true for matching string slices")
		}
	})
}
