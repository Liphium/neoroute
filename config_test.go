package neoroute_test

import (
	"fmt"
	"testing"

	"github.com/Liphium/neoroute"
)

func TestConfig_RunErrorHandler(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cfg  neoroute.Config
		err  error
		want string
	}{
		{
			name: "error without defined error function",
			cfg:  neoroute.Config{},
			err:  fmt.Errorf("some error"),
			want: "Internal Server Error",
		},
		{
			name: "error defined error function",
			cfg: neoroute.Config{
				func(err error) string {
					return fmt.Sprintf("received error: %v", err)
				},
			},
			err:  fmt.Errorf("some error"),
			want: "received error: some error",
		},
		{
			name: "nil error without defined error function",
			cfg:  neoroute.Config{},
			err:  nil,
			want: "Internal Server Error",
		},
		{
			name: "nil error defined error function",
			cfg: neoroute.Config{
				func(err error) string {
					return fmt.Sprintf("received error: %v", err)
				},
			},
			err:  nil,
			want: "received error: <nil>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.RunErrorHandler(tt.err)
			if tt.want != got {
				t.Errorf("RunErrorHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
