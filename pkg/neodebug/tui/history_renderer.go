package tui

import (
	"fmt"
	"strings"
)

// render renders a any value with the least lines possible for it to still be readable.
func render(value any, width int) string {

	switch v := value.(type) {
	case map[string]any:

		totalLength := 0
		fields := make([]string, len(v))

		i := 0
		for n, e := range v {
			if totalLength > width {
				break
			}

			const nameLength = 6 // "name": value,
			rendered := render(e, max(width-totalLength-len(n)-nameLength, 0))
			totalLength += len(rendered) + len(n) + nameLength
			fields[i] = fmt.Sprintf("\"%s\": %s", n, rendered)
			i++
		}

		if totalLength <= width {
			return "{" + strings.Join(fields, ", ") + "}"
		}

		// Otherwise render again, but this time with proper padding (this will allow maps / slices that before could not be rendered in one line to potentially be (-> more width))
		i = 0
		for n, e := range v {
			fields[i] = fmt.Sprintf("\"%s\": %s", n, render(e, width-structPadding-1 /* comma */))
			i++
		}

		return "{\n" + structChildStyle.Render(strings.Join(fields, ",\n")+",") + "\n}"

	case []any:

		totalLength := 0
		fields := make([]string, len(v))

		for i, e := range v {
			rendered := render(e, max(width-totalLength, 0))
			fields[i] = rendered
			totalLength += len(rendered) + 2 /* space + comma between the thing and the next value */
		}

		if totalLength /* we don't need to -2 here cause brackets */ <= width {
			return "[" + strings.Join(fields, ", ") + "]"
		}

		// Otherwise render again, but this time with proper padding (this will allow maps / slices that before could not be rendered in one line to potentially be (-> more width))
		for i, e := range v {
			fields[i] = render(e, width-structPadding-1 /* comma */)
		}

		return "[\n" + structChildStyle.Render(strings.Join(fields, ",\n")+",") + "\n]"

	case string:
		return "\"" + v + "\""

	case float64, float32:
		return fmt.Sprintf("%2.f", v)

	default:
		// Just render everything else as a string, this is fine for basically everything
		return fmt.Sprintf("%v", v)
	}
}
