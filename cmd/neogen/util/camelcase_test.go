package util_test

import (
	"testing"

	"github.com/Liphium/neoroute/cmd/neogen/util"
)

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		firstUp bool
		want    string
	}{
		{
			name:    "already camelCase (firstUp=false)",
			input:   "camelCase",
			firstUp: false,
			want:    "camelCase",
		},
		{
			name:    "already camelCase (firstUp=true)",
			input:   "camelCase",
			firstUp: true,
			want:    "CamelCase",
		},
		{
			name:    "snake_case (firstUp=false)",
			input:   "already_camel_case",
			firstUp: false,
			want:    "alreadyCamelCase",
		},
		{
			name:    "snake_case (firstUp=true)",
			input:   "already_camel_case",
			firstUp: true,
			want:    "AlreadyCamelCase",
		},
		{
			name:    "multiple mixed separators",
			input:   "multiple---separators___here",
			firstUp: false,
			want:    "multipleSeparatorsHere",
		},
		{
			name:    "handling acronyms",
			input:   "JSONParser",
			firstUp: false,
			want:    "jsonParser",
		},
		{
			name:    "numbers preserved inside string",
			input:   "user123name",
			firstUp: false,
			want:    "user123Name",
		},
		{
			name:    "numbers preserved inside string with firstUp",
			input:   "user123name",
			firstUp: true,
			want:    "User123Name",
		},
		{
			name:    "numbers separated by special characters",
			input:   "V3___release_4_final",
			firstUp: true,
			want:    "V3Release4Final",
		},
		{
			name:    "unicode support",
			input:   "café_au_lait",
			firstUp: false,
			want:    "caféAuLait",
		},
		{
			name:    "empty string",
			input:   "",
			firstUp: false,
			want:    "",
		},
		{
			name:    "only non-letters and non-digits",
			input:   "---___!!!",
			firstUp: false,
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.ToCamelCase(tt.input, tt.firstUp)
			if got != tt.want {
				t.Errorf("ToCamelCase(%q, %t) = %q; want %q", tt.input, tt.firstUp, got, tt.want)
			}
		})
	}
}
