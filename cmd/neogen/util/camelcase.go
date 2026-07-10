package util

import (
	"strings"
	"unicode"
)

// ToCamelCase converts a string into camelCase or PascalCase.
// - firstUp: If true, the first letter of the output is capitalized.
// - Any character that is not a letter or a digit is treated as a separator.
// - Existing camelCase structures and numbers are preserved.
func ToCamelCase(s string, firstUp bool) string {
	if s == "" {
		return ""
	}

	runes := []rune(s)
	var words []string
	var currentWord []rune

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		isLetter := unicode.IsLetter(r)
		isDigit := unicode.IsDigit(r)

		// Treat anything that is not a letter or a number as a separator
		if !isLetter && !isDigit {
			if len(currentWord) > 0 {
				words = append(words, string(currentWord))
				currentWord = nil
			}
			continue
		}

		// Handle boundaries for existing camelCase and transitions with numbers
		if len(currentWord) > 0 {
			prev := currentWord[len(currentWord)-1]
			prevIsLetter := unicode.IsLetter(prev)
			prevIsDigit := unicode.IsDigit(prev)

			switch {
			// Transition: Letter to Digit (e.g., 'r' -> '1' in "user123name")
			case prevIsLetter && isDigit:
				words = append(words, string(currentWord))
				currentWord = nil

			// Transition: Digit to Letter (e.g., '3' -> 'n' in "user123name")
			case prevIsDigit && isLetter:
				words = append(words, string(currentWord))
				currentWord = nil

			// Boundary 1: Lowercase letter followed by an uppercase letter (e.g., 'l' -> 'C' in "camelCase")
			case unicode.IsLower(prev) && unicode.IsUpper(r):
				words = append(words, string(currentWord))
				currentWord = nil

			// Boundary 2: Uppercase acronym transitioning to a new word (e.g., 'N' -> 'P' -> 'a' in "JSONParser")
			case unicode.IsUpper(prev) && unicode.IsUpper(r) && i+1 < len(runes) && unicode.IsLower(runes[i+1]):
				words = append(words, string(currentWord))
				currentWord = nil
			}
		}

		currentWord = append(currentWord, r)
	}

	// Append the final word if there is one
	if len(currentWord) > 0 {
		words = append(words, string(currentWord))
	}

	if len(words) == 0 {
		return ""
	}

	var result strings.Builder
	for i, word := range words {
		if i == 0 {
			if firstUp {
				result.WriteString(capitalize(word))
			} else {
				result.WriteString(strings.ToLower(word))
			}
		} else {
			result.WriteString(capitalize(word))
		}
	}

	return result.String()
}

// capitalize capitalizes the first letter of a word and lowercases the rest, preserving digits.
func capitalize(s string) string {
	runes := []rune(s)
	if len(runes) == 0 {
		return ""
	}
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}
