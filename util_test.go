package neoroute

import (
	"testing"
)

func Test_cleanRoute(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		route string
		want  string
	}{
		{
			name:  "remove leading separator",
			route: string(RouteSeparator) + string(RouteSeparator) + "legal_route",
			want:  "legal_route",
		},
		{
			name:  "remove trailing separator",
			route: "legal_route" + string(RouteSeparator) + string(RouteSeparator),
			want:  "legal_route",
		},
		{
			name:  "remove multiple separators in a row",
			route: "legal" + string(RouteSeparator) + string(RouteSeparator) + "route",
			want:  "legal.route",
		},
		{
			name:  "remove multiple separators in a row with first illegal characters in between",
			route: "legal" + string(RouteSeparator) + "@/$" + string(RouteSeparator) + "route",
			want:  "legal.route",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanRoute(tt.route)
			if tt.want != got {
				t.Errorf("cleanRoute() = %v, want %v", got, tt.want)
			}
		})
	}
}
