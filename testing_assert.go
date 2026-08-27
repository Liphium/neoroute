package neoroute

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tinylib/msgp/msgp"
)

// AssertResponse asserts handler response matches expected response.
//
// Hint: You can pass options from the go-cmp package straight into this function as well. The cmpopts package is useful in that case.
func AssertResponse[RQ any, PQ msgp.UnmarshalPtr[RQ]](t *testing.T, response error, expected RQ, opts ...cmp.Option) {
	t.Helper()
	AssertResponseFunc[RQ, PQ](t, response, expected, func(want, got RQ) bool {
		return cmp.Equal(want, got, opts...)
	})
}

// AssertResponseFunc asserts handler response using compare.
func AssertResponseFunc[RQ any, PQ msgp.UnmarshalPtr[RQ]](t *testing.T, response error, expected RQ, compare func(RQ, RQ) bool) {
	t.Helper()
	actual, userErr, err := GetTestingResponse[RQ, PQ](response)
	if err != nil {
		t.Fatalf("unexpected handler error: %v", err)
	}
	if userErr != "" {
		t.Fatalf("unexpected user error: %s", userErr)
	}
	if compare == nil || !compare(expected, actual) {
		t.Fatalf("response mismatch: expected %#v, got %#v", expected, actual)
	}
}

// AssertResponseOk asserts handler response contains no user or handler error.
func AssertResponseOk(t *testing.T, response error) {
	t.Helper()
	if userErr, err := GetTestingResponseOk(response); err != nil {
		t.Fatalf("unexpected handler error: %v", err)
	} else if userErr != "" {
		t.Fatalf("unexpected user error: %s", userErr)
	}
}

// AssertUserError asserts handler response contains expected user error.
func AssertUserError(t *testing.T, response error, expected string) {
	t.Helper()
	userErr, err := GetTestingResponseOk(response)
	if err != nil {
		t.Fatalf("unexpected handler error: %v", err)
	}
	if userErr != expected {
		t.Fatalf("expected user error %q, got %q", expected, userErr)
	}
}

// AssertEvents asserts adapter contains expected number of events, returning them.
func AssertEvents(t *testing.T, adapter *TestingAdapter, expected int) []event {
	t.Helper()
	events, err := adapter.GetEvents()
	if err != nil {
		t.Fatalf("get events: %v", err)
	}
	if len(events) != expected {
		t.Fatalf("expected %d events, got %d", expected, len(events))
	}
	return events
}

// AssertEventType assert the event at index is actually the type you expect, returning decoded event payload.
func AssertEventType[E any, EP msgp.UnmarshalPtr[E]](t *testing.T, events []event, index int) E {
	t.Helper()
	var zero E
	if index < 0 || index >= len(events) {
		t.Fatalf("event index %d out of range", index)
		return zero
	}
	ev, err := UnmarshalEventTesting[E, EP](events[index].Data)
	if err != nil {
		t.Fatalf("unmarshal event: %v", err)
	}
	return ev
}

// AssertEventFunc asserts the event at index actually has the content you expect, with a custom comparison function.
func AssertEventFunc[E any, EP msgp.UnmarshalPtr[E]](t *testing.T, events []event, index int, expected E, compare func(e1, e2 E) bool) {
	t.Helper()
	actual := AssertEventType[E, EP](t, events, index)
	if !compare(actual, expected) {
		t.Fatalf("event mismatch: expected %#v, got %#v", expected, actual)
	}
}

// AssertEvent asserts the event at index are deep equal to the event you expect.
//
// Hint: You can pass options from the go-cmp package straight into this function as well. The cmpopts package is useful in that case.
func AssertEvent[E any, EP msgp.UnmarshalPtr[E]](t *testing.T, events []event, index int, expected E, opts ...cmp.Option) {
	t.Helper()
	AssertEventFunc[E, EP](t, events, index, expected, func(e1, e2 E) bool {
		return cmp.Equal(e1, e2, opts...)
	})
}
