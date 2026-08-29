package neoroute

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp -unexported
type reqStructTesting struct {
	SomeField int `msgp:"some_field"`
}

type respStructTesting struct {
	SomeField string `msgp:"some_field"`
}

func TestRouteTest(t *testing.T) {

	type Response[RS any, PS msgp.UnmarshalPtr[RS]] struct {
		Resp      RS
		UserError string
		Error     error
	}

	tests := []struct {
		name        string
		routingFunc func(r *Router[bool])
		req         reqStructTesting
		want        Response[respStructTesting, *respStructTesting]
		sessionData bool
		route       string
	}{
		{
			name: "simple route",
			routingFunc: func(r *Router[bool]) {
				r.Route("some-route", func(c *ResCtx[bool, respStructTesting], req reqStructTesting) error {
					return c.Respond(respStructTesting{
						SomeField: fmt.Sprintf("%d", req.SomeField),
					})
				})
			},
			route: "some-route",
			want: Response[respStructTesting, *respStructTesting]{
				Resp:      respStructTesting{SomeField: "42"},
				UserError: "",
				Error:     nil,
			},
			req: reqStructTesting{SomeField: 42},
		},
		{
			name: "middleware applies",
			routingFunc: func(r *Router[bool]) {
				r.Use("some-route", func(c *Ctx[bool]) error {
					c.session.UpdateData(func(data *bool) { *data = true })
					return nil
				})
				r.Route("some-route", func(c *ResCtx[bool, respStructTesting], req reqStructTesting) error {
					return c.Respond(respStructTesting{
						SomeField: fmt.Sprintf("%d", req.SomeField),
					})
				})
			},
			route: "some-route",
			want: Response[respStructTesting, *respStructTesting]{
				Resp:      respStructTesting{SomeField: "42"},
				UserError: "",
				Error:     nil,
			},
			req:         reqStructTesting{SomeField: 42},
			sessionData: true,
		},
		{
			name: "user error on route",
			routingFunc: func(r *Router[bool]) {
				r.Route("some-route", func(c *ResCtx[bool, respStructTesting], req reqStructTesting) error {
					return NewError("user error")
				})
			},
			route: "some-route",
			want: Response[respStructTesting, *respStructTesting]{
				Resp:      respStructTesting{},
				UserError: "user error",
				Error:     nil,
			},
			req: reqStructTesting{},
		},
		{
			name: "middleware fails",
			routingFunc: func(r *Router[bool]) {
				r.Use("some-route", func(c *Ctx[bool]) error {
					return NewError("middleware error")
				})
				r.Route("some-route", func(c *ResCtx[bool, respStructTesting], req reqStructTesting) error {
					c.Session().SetData(true)
					return NewError("user error")
				})
			},
			route: "some-route",
			want: Response[respStructTesting, *respStructTesting]{
				Resp:      respStructTesting{},
				UserError: "middleware error",
				Error:     nil,
			},
			req: reqStructTesting{},
		},
		{
			name: "first middleware fails, second middleware passes",
			routingFunc: func(r *Router[bool]) {
				r.Use("some-route", func(c *Ctx[bool]) error {
					return NewError("middleware error")
				})
				r.Use("some-route", func(c *Ctx[bool]) error {
					c.Session().SetData(true)
					return nil
				})
				r.Route("some-route", func(c *ResCtx[bool, respStructTesting], req reqStructTesting) error {
					return NewError("user error")
				})
			},
			route: "some-route",
			want: Response[respStructTesting, *respStructTesting]{
				Resp:      respStructTesting{},
				UserError: "middleware error",
				Error:     nil,
			},
			req: reqStructTesting{},
		},
		{
			name: "run after func gets executed",
			routingFunc: func(r *Router[bool]) {
				r.Route("some-route", func(c *ResCtx[bool, respStructTesting], req reqStructTesting) error {
					c.RunAfter(func() {
						c.Session().SetData(true)
					})
					return NewError("user error")
				})
			},
			route: "some-route",
			want: Response[respStructTesting, *respStructTesting]{
				Resp:      respStructTesting{},
				UserError: "user error",
				Error:     nil,
			},
			req:         reqStructTesting{},
			sessionData: true,
		},
		{
			name: "run after gets executed after middleware fails",
			routingFunc: func(r *Router[bool]) {
				r.Use("some-route", func(c *Ctx[bool]) error {
					c.RunAfter(func() {
						c.Session().SetData(true)
					})
					return NewError("middleware error")
				})
				r.Route("some-route", func(c *ResCtx[bool, respStructTesting], req reqStructTesting) error {
					return NewError("user error")
				})
			},
			route: "some-route",
			want: Response[respStructTesting, *respStructTesting]{
				Resp:      respStructTesting{},
				UserError: "middleware error",
				Error:     nil,
			},
			req:         reqStructTesting{},
			sessionData: true,
		},
		{
			name: "run afters stack",
			routingFunc: func(r *Router[bool]) {
				r.Use("some-route", func(c *Ctx[bool]) error {
					c.RunAfter(func() {
						c.Session().SetData(false)
					}, func() {
						c.Session().SetData(false)
					})
					c.RunAfter(func() {
						c.Session().SetData(true)
					})
					return NewError("middleware error")
				})
				r.Route("some-route", func(c *ResCtx[bool, respStructTesting], req reqStructTesting) error {
					return NewError("user error")
				})
			},
			route: "some-route",
			want: Response[respStructTesting, *respStructTesting]{
				Resp:      respStructTesting{},
				UserError: "middleware error",
				Error:     nil,
			},
			req:         reqStructTesting{},
			sessionData: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter(Config[bool]{})
			tt.routingFunc(r)

			session := NewTestingSession(false, "some-session")

			resp, userError, err := RouteTest[bool, reqStructTesting, respStructTesting](r, session, tt.route, tt.req)

			got := Response[respStructTesting, *respStructTesting]{
				Resp:      resp,
				UserError: userError,
				Error:     err,
			}

			if tt.sessionData != session.Data() {
				t.Errorf("session data mismatch: got %v, want %v", session.Data(), tt.sessionData)
			}

			if !cmp.Equal(got, tt.want, cmpopts.EquateErrors()) {
				t.Errorf("RouteTest() mismatch (-got +want):\n%s",
					cmp.Diff(got, tt.want, cmpopts.EquateErrors()),
				)
			}
		})
	}
}
