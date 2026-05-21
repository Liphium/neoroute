package neoroute

import "strings"

func containsSeparator(route string) bool {
	return strings.ContainsRune(route, RouteSeparator)
}

func cleanRoute(route string) string {

	// Remove leading and trailing separators
	route = strings.Trim(route, string(RouteSeparator))

	unfilteredRunes := []rune(route)

	var runes []rune
	for _, r := range unfilteredRunes {
		if _, ok := allowedRouteRunes[r]; ok {
			runes = append(runes, r)
		}
	}

	if len(runes) == 0 {
		return ""
	}

	var result strings.Builder

	result.WriteRune(runes[0])

	for i := 1; i < len(runes); i++ {

		// Remove if multiple separators in a row
		if runes[i] == runes[i-1] && runes[i] == RouteSeparator {
			continue
		}

		result.WriteRune(runes[i])
	}

	return result.String()
}
