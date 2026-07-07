package neoroute

import (
	"reflect"
	"slices"
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
			want:  "legal/route",
		},
		{
			name:  "remove multiple separators in a row with first illegal characters in between",
			route: "legal" + string(RouteSeparator) + "@$" + string(RouteSeparator) + "route",
			want:  "legal/route",
		},
		{
			name:  "empty route",
			route: "",
			want:  "",
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

func Test_buildSubroutes(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		route string
		want  []string
	}{
		{
			name:  "3 level route",
			route: "group1/group2/route",
			want:  []string{"group1", "group1/group2", "group1/group2/route"},
		},
		{
			name:  "route with leading and trailing separators",
			route: string(RouteSeparator) + "group1/group2/route" + string(RouteSeparator),
			want:  []string{"group1", "group1/group2", "group1/group2/route"},
		},
		{
			name:  "route with multiple separators in a row",
			route: "group_1" + string(RouteSeparator) + string(RouteSeparator) + "group2" + string(RouteSeparator) + "route",
			want:  []string{"group_1", "group_1/group2", "group_1/group2/route"},
		},
		{
			name:  "empty route",
			route: "",
			want:  []string{""},
		},
		{
			name:  "single route",
			route: "route_1",
			want:  []string{"route_1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildSubroutes(tt.route)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildSubroutes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitRoute(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		route string
		want  []string
	}{
		{
			name:  "split empty route",
			route: "",
			want:  []string{""},
		},
		{
			name:  "split single route",
			route: "single_route",
			want:  []string{"single_route"},
		},
		{
			name:  "split multiple routes",
			route: "first_route" + string(RouteSeparator) + "second_route" + string(RouteSeparator) + "third_route",
			want:  []string{"first_route", "second_route", "third_route"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitRoute(tt.route)
			if !slices.Equal(got, tt.want) {
				t.Errorf("splitRoute() = %v, want %v", got, tt.want)
			}
		})
	}
}
