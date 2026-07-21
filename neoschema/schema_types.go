package neoschema

import (
	"reflect"

	"github.com/tinylib/msgp/msgp"
)

// Interface that bundles all of the interfaces for a proper message that can be sent by neoroute.
type UsableMessage interface {
	msgp.Marshaler
	msgp.Unmarshaler
	msgp.Decodable
	msgp.Encodable
	msgp.Sizer
}

type SchemaType string

// All schema types
const (
	TypeInt32        SchemaType = "int32"
	TypeInt64        SchemaType = "int64"
	TypeFloat32      SchemaType = "float32"
	TypeFloat64      SchemaType = "float64" // Does not exist in MessagePack, but kept here for compatability with other languages
	TypeBool         SchemaType = "bool"
	TypeByte         SchemaType = "byte" // This does not exist in MessagePack, but helps us let message pack identify byte arrays
	TypeString       SchemaType = "string"
	TypeArray        SchemaType = "array"
	TypeMap          SchemaType = "map"
	TypeStruct       SchemaType = "struct"
	TypeSerializable SchemaType = "serializable"
	TypeNullable     SchemaType = "nullable"
	TypeReference    SchemaType = "reference"
	TypeNotSupported SchemaType = "-"
)

var Kinds = map[reflect.Kind]SchemaType{

	// No need to map all of these types to different stuff (for compatability reasons)
	reflect.Int8:   TypeInt32,
	reflect.Int16:  TypeInt32,
	reflect.Int32:  TypeInt32,
	reflect.Int:    TypeInt64, // At least 32 bits in size, meaning it could also be 64 bit, so we're doing this for safety
	reflect.Uint:   TypeInt32,
	reflect.Uint16: TypeInt32,
	reflect.Uint32: TypeInt32,

	// Big integers should be separate
	reflect.Uint64: TypeInt64,
	reflect.Int64:  TypeInt64,

	// Simple types
	reflect.Array:   TypeArray,
	reflect.Slice:   TypeArray,
	reflect.Float32: TypeFloat32,
	reflect.Float64: TypeFloat64,
	reflect.Bool:    TypeBool,
	reflect.Uint8:   TypeByte, // This is a special case because this is often used for bytes in Go (therefore should also be that in another languages, it's just a type alias, but the reflect package does not have it)
	reflect.String:  TypeString,
	reflect.Struct:  TypeStruct,
	reflect.Map:     TypeMap,

	// For interfaces, we can't do anything special, but we can just give them straight to the thing anyway
	reflect.Interface: TypeSerializable,

	// Not supported currently
	reflect.Chan:       TypeNotSupported,
	reflect.Complex128: TypeNotSupported,
	reflect.Complex64:  TypeNotSupported,
	reflect.Func:       TypeNotSupported,
}

type PackedType interface {
	Type() SchemaType
	ObjectRegistry() map[string]PackedType
	SetRegistry(registry map[string]PackedType)

	// Removes registries from all children types, except from the root one
	CleanRegistries(root bool)
}

type BasicType struct {
	ActualType SchemaType            `json:"type"`
	Objects    map[string]PackedType `json:"objects,omitempty"`
}

func (bt *BasicType) Type() SchemaType {
	return bt.ActualType
}

func (bt *BasicType) ObjectRegistry() map[string]PackedType {
	return bt.Objects
}

func (bt *BasicType) SetRegistry(registry map[string]PackedType) {
	bt.Objects = registry
}

func (bt *BasicType) CleanRegistries(root bool) {
	if bt == nil {
		return
	}

	if root {
		// Keep this root level's Object registry, but clean all in the object registry.
		for _, obj := range bt.Objects {
			if obj != nil {
				obj.CleanRegistries(false)
			}
		}
		return
	}

	// Clear the registry reference on nested children
	bt.Objects = nil
}

type ArrayType struct {
	*BasicType

	Element PackedType `json:"element"`
}

func (at *ArrayType) CleanRegistries(root bool) {
	at.Element.CleanRegistries(false)
	at.BasicType.CleanRegistries(root)
}

type StructType struct {
	*BasicType

	Name   string                `json:"name"`
	Fields map[string]PackedType `json:"fields"`
}

func (st *StructType) CleanRegistries(root bool) {
	for _, v := range st.Fields {
		v.CleanRegistries(false)
	}
	st.BasicType.CleanRegistries(root)
}

type ReferenceType struct {
	*BasicType

	Object string `json:"object"`
}

type NullableType struct {
	*BasicType

	Element PackedType `json:"element"`
}

func (at *NullableType) CleanRegistries(root bool) {
	at.Element.CleanRegistries(false)
	at.BasicType.CleanRegistries(root)
}

type MapType struct {
	*BasicType

	Key   PackedType `json:"key"`
	Value PackedType `json:"value"`
}

func (at *MapType) CleanRegistries(root bool) {
	at.Key.CleanRegistries(false)
	at.Value.CleanRegistries(false)
	at.BasicType.CleanRegistries(root)
}
