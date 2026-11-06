package neoschema_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Liphium/neoroute/v2/neoschema"
	"github.com/stretchr/testify/require"
)

func TestMapKeysMustBeStrings(t *testing.T) {
	for _, typ := range []reflect.Type{
		reflect.TypeFor[map[int32]string](),
		reflect.TypeFor[map[uint64]string](),
		reflect.TypeFor[map[float64]string](),
		reflect.TypeFor[map[bool]string](),
		reflect.TypeFor[struct{ Values []map[int]string }](),
		reflect.TypeFor[map[string]*map[int]string](),
	} {
		t.Run(typ.String(), func(t *testing.T) {
			_, err := neoschema.BuildPackedFor(typ)
			require.ErrorIs(t, err, neoschema.ErrMapKeyNotString)
		})
	}
}

func TestStringMapKeys(t *testing.T) {
	type Key string
	for _, typ := range []reflect.Type{
		reflect.TypeFor[map[string]string](),
		reflect.TypeFor[map[Key][]map[string]int](),
	} {
		t.Run(typ.String(), func(t *testing.T) {
			packed, err := neoschema.BuildPackedFor(typ)
			require.NoError(t, err)
			data, err := json.Marshal(packed)
			require.NoError(t, err)
			_, err = neoschema.UnmarshalPackedType(data)
			require.NoError(t, err)
		})
	}
}

func TestLoadingSchemaRejectsNonStringMapKeys(t *testing.T) {
	for _, key := range []neoschema.SchemaType{
		neoschema.TypeInt32, neoschema.TypeInt64, neoschema.TypeFloat64, neoschema.TypeBool,
	} {
		t.Run(string(key), func(t *testing.T) {
			packed := &neoschema.MapType{
				BasicType: &neoschema.BasicType{ActualType: neoschema.TypeMap},
				Key:       &neoschema.BasicType{ActualType: key},
				Value:     &neoschema.BasicType{ActualType: neoschema.TypeString},
			}
			data, err := json.Marshal(packed)
			require.NoError(t, err)
			_, err = neoschema.UnmarshalPackedType(data)
			require.ErrorIs(t, err, neoschema.ErrMapKeyNotString)
		})
	}
}
