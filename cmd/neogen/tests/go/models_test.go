package generated

import (
	"reflect"
	"testing"
	"time"

	client "github.com/Liphium/neoroute/client/sdk/v2"
	"github.com/Liphium/neoroute/v2/neoschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type field struct {
	name string
	typ  reflect.Type
	tag  string
}

func fieldOf[T any](name, tag string) field { return field{name, reflect.TypeFor[T](), tag} }

func TestModelShapes(t *testing.T) {
	cases := []struct {
		typ    reflect.Type
		fields []field
	}{
		{reflect.TypeFor[Address](), []field{fieldOf[string]("City", "city"), fieldOf[*int64]("Zip", "zip")}},
		{reflect.TypeFor[Inner](), []field{fieldOf[string]("Note", "note")}},
		{reflect.TypeFor[User](), []field{fieldOf[string]("Name", "name"), fieldOf[*Address]("Addr", "addr"), fieldOf[[]string]("Tags", "tags"), fieldOf[Inner]("Inner", "inner"), fieldOf[[]User]("Friends", "friends")}},
		{reflect.TypeFor[StandaloneInner](), []field{fieldOf[string]("Info", "info")}},
		{reflect.TypeFor[Standalone](), []field{fieldOf[int64]("Code", "code"), fieldOf[*string]("Label", "label"), fieldOf[[]StandaloneInner]("Items", "items"), fieldOf[time.Time]("SomeStamp", "some_stamp")}},
		{reflect.TypeFor[Coordinates](), []field{fieldOf[float64]("X", "x"), fieldOf[float64]("Y", "y")}},
		{reflect.TypeFor[Location](), []field{fieldOf[string]("Name", "name"), fieldOf[Coordinates]("Coordinates", "coordinates")}},
		{reflect.TypeFor[Point](), []field{fieldOf[float64]("X", "x"), fieldOf[string]("Label", "label")}},
		{reflect.TypeFor[Envelope](), []field{
			fieldOf[Point]("Point", "point"), fieldOf[[]*Point]("Items", "items"), fieldOf[[][]Point]("Matrix", "matrix"),
			fieldOf[map[string]*Point]("Lookup", "lookup"),
			fieldOf[*Envelope]("Next", "next"), fieldOf[[]byte]("Bytes", "bytes"), fieldOf[time.Time]("Date", "date"), fieldOf[neoschema.UsableMessage]("Opaque", "opaque"),
		}},
	}
	for _, tc := range cases {
		t.Run(tc.typ.Name(), func(t *testing.T) {

			require.Equal(t, len(tc.fields), tc.typ.NumField(), "field count")
			for i, want := range tc.fields {
				got := tc.typ.Field(i)
				assert.Equal(t, want.name, got.Name, "field %d name", i)
				assert.True(t, want.typ == got.Type, "field %s: want %v, got %v", want.name, want.typ, got.Type)
				assert.Equal(t, want.tag, got.Tag.Get("msg"), "field %s tag", want.name)
			}
		})
	}
}

func TestTupleEncoding(t *testing.T) {

	// Test with a point
	point := Point{X: 4, Label: "test"}
	bytes, err := point.MarshalMsg(nil)
	require.NoError(t, err)
	rest := []byte{}
	var decoded Point
	rest, err = decoded.UnmarshalMsg(bytes)
	require.NoError(t, err)
	assert.Empty(t, rest)
	assert.Equal(t, point, decoded)

	// Test with coordinates
	location := Location{Name: "home", Coordinates: Coordinates{X: 1, Y: 2}}
	bytes, err = location.MarshalMsg(nil)
	require.NoError(t, err)
	var decodedLocation Location
	rest, err = decodedLocation.UnmarshalMsg(bytes)
	require.NoError(t, err)
	assert.Empty(t, rest)
	assert.Equal(t, location, decodedLocation)
}

func TestConnectorSignatures(t *testing.T) {
	methods := map[string]reflect.Type{
		"SendGetUser":  reflect.TypeFor[func(User, ...func(*client.RequestOptions)) (User, error)](),
		"SendExchange": reflect.TypeFor[func(Envelope, ...func(*client.RequestOptions)) (Envelope, error)](),
		"SendGet":      reflect.TypeFor[func(...func(*client.RequestOptions)) (Envelope, error)](),
		"SendOk":       reflect.TypeFor[func(Envelope, ...func(*client.RequestOptions)) error](),
		"SendPush":     reflect.TypeFor[func(Envelope, ...func(*client.RequestOptions)) error](),
		"SendSignal":   reflect.TypeFor[func(...func(*client.RequestOptions)) error](),
		"SendPing":     reflect.TypeFor[func(...func(*client.RequestOptions)) error](),
	}
	for _, connector := range []any{&HttpConnector{}, &WsConnector{}} {
		for name, want := range methods {
			got := reflect.ValueOf(connector).MethodByName(name)

			require.True(t, got.IsValid(), "%T missing %s", connector, name)
			assert.True(t, want == got.Type(), "%T.%s: want %v, got %v", connector, name, want, got.Type())
		}
	}
	got := reflect.ValueOf(&HttpConnector{}).MethodByName("ReceiveUpdate")
	want := reflect.TypeFor[func(func(Envelope))]()
	require.True(t, got.IsValid(), "missing ReceiveUpdate")
	assert.True(t, want == got.Type(), "ReceiveUpdate: want %v, got %v", want, got.Type())
}
