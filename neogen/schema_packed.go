package neogen

import (
	"fmt"
	"reflect"
)

type SchemaType string

const NotSupported SchemaType = "-"

var Kinds = map[reflect.Kind]SchemaType{

	// No need to map all of these types to different stuff (for compatability reasons)
	reflect.Int8:   "int32",
	reflect.Int16:  "int32",
	reflect.Int32:  "int32",
	reflect.Int:    "int32",
	reflect.Uint:   "int32",
	reflect.Uint16: "int32",
	reflect.Uint32: "int32",

	// Big integers should be separate
	reflect.Uint64: "int64",
	reflect.Int64:  "int64",

	// Simple types
	reflect.Array:   "array",
	reflect.Float32: "float32",
	reflect.Float64: "float64",
	reflect.Bool:    "bool",
	reflect.Uint8:   "byte", // This is a special case because this is often used for bytes in Go (therefore should also be that in another languages, it's just a type alias, but the reflect package does not have it)
	reflect.String:  "string",
	reflect.Struct:  "struct",

	// Not supported currently
	reflect.Chan:       NotSupported,
	reflect.Complex128: NotSupported,
	reflect.Complex64:  NotSupported,
	reflect.Func:       NotSupported,
}

type PackedType interface {
	Type() SchemaType
}

type BasicType struct {
	ActualType SchemaType `json:"type"`
}

func (bt BasicType) Type() SchemaType {
	return bt.ActualType
}

type ArrayType struct {
	BasicType

	Element PackedType `json:"element"`
}

type StructType struct {
	BasicType

	Fields map[string]PackedType `json:"fields"`
}

func notSupportedError(st SchemaType) error {
	return fmt.Errorf("the type %s is not supported", string(st))
}

func BuildPackedFor(t reflect.Type) (PackedType, error) {
	kind := t.Kind()
	switch kind {
	case reflect.Struct:
		// Go through all struct fields and build their schemas
		var err error
		fields := map[string]PackedType{}
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			msgTag := field.Tag.Get("msg")
			if msgTag == "" {
				msgTag = field.Name
			}

			fields[msgTag], err = BuildPackedFor(field.Type)
			if err != nil {
				return BasicType{}, err
			}
		}

		return StructType{
			BasicType: BasicType{
				ActualType: Kinds[kind],
			},
			Fields: fields,
		}, nil

	case reflect.Array:
		// Build the type for the array
		arrayElem, err := BuildPackedFor(t.Elem())
		if err != nil {
			return BasicType{}, err
		}

		return ArrayType{
			BasicType: BasicType{
				ActualType: Kinds[kind],
			},
			Element: arrayElem,
		}, nil

	case reflect.Pointer:
		// With msgp pointers just become the regular type
		return BuildPackedFor(t.Elem())

	default:
		st := Kinds[kind]
		if st == NotSupported {
			return BasicType{}, notSupportedError(st)
		}

		return BasicType{
			ActualType: Kinds[kind],
		}, nil
	}
}
