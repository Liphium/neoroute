package neoroute

import (
	"fmt"
	"testing"
)

type testSessionData struct{}

func TestConfig_RunErrorHandler(t *testing.T) {
	r := NewNeoRouter(Config[testSessionData]{})
	session := NewSession("test-id", testSessionData{}, SessionTransporterCallbacks[testSessionData]{})
	ctx := NewTestingCtx(r, "test.route", session)

	tests := []struct {
		name string
		cfg  Config[testSessionData]
		err  error
		want string
	}{
		{
			name: "error without defined error function",
			cfg:  Config[testSessionData]{},
			err:  fmt.Errorf("some error"),
			want: "Internal Server Error",
		},
		{
			name: "error defined error function",
			cfg: Config[testSessionData]{
				ErrorHandler: func(err error, c *Ctx[testSessionData]) string {
					return fmt.Sprintf("received error: %v", err)
				},
			},
			err:  fmt.Errorf("some error"),
			want: "received error: some error",
		},
		{
			name: "nil error without defined error function",
			cfg:  Config[testSessionData]{},
			err:  nil,
			want: "Internal Server Error",
		},
		{
			name: "nil error defined error function",
			cfg: Config[testSessionData]{
				ErrorHandler: func(err error, c *Ctx[testSessionData]) string {
					return fmt.Sprintf("received error: %v", err)
				},
			},
			err:  nil,
			want: "received error: <nil>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.RunErrorHandler(tt.err, ctx)
			if tt.want != got {
				t.Errorf("RunErrorHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
