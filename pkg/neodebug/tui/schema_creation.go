package tui

import (
	"errors"
	"strconv"

	"github.com/Liphium/neoroute/neoschema"
)

func createNode(schema neoschema.PackedType, registry map[string]neoschema.PackedType) SchemaNode {
	switch schema := schema.(type) {
	case *neoschema.StructType:

		// Create the children
		children := make([]StructField, 0, len(schema.Fields))
		for name, field := range schema.Fields {
			children = append(children, StructField{
				Name: name,
				Node: createNode(field, registry),
			})
		}

		return &StructNode{
			name:     schema.Name,
			children: children,
		}

	case *neoschema.NullableType:
		return &NullableNode{
			null:  true,
			other: createNode(schema.Element, registry),
		}

	case *neoschema.ArrayType:
		return &SliceNode{
			element:  schema.Element,
			gap:      -1,
			registry: registry,
			items:    []SchemaNode{},
		}

	case *neoschema.ReferenceType:
		return createNode(registry[schema.Object], registry)

	case *neoschema.BasicType:
		return createBasic(schema)

	default:

	}
	return nil
}

func newValueNode[T any](convert func(string) (T, error)) *ValueNode[T] {
	var zero T
	return &ValueNode[T]{
		convert: convert,
		value:   zero,
	}
}

func createBasic(basic *neoschema.BasicType) SchemaNode {

	switch basic.ActualType {
	case neoschema.TypeInt32:
		return newValueNode(func(s string) (int32, error) {
			i, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				return 0, errors.New("Please enter a valid int32.")
			}
			return int32(i), nil
		})
	case neoschema.TypeInt64:
		return newValueNode(func(s string) (int64, error) {
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return 0, errors.New("Please enter a valid int64.")
			}
			return i, nil
		})
	case neoschema.TypeFloat32:
		return newValueNode(func(s string) (float32, error) {
			f, err := strconv.ParseFloat(s, 32)
			if err != nil {
				return 0, errors.New("Please enter a valid float32.")
			}
			return float32(f), nil
		})
	case neoschema.TypeFloat64:
		return newValueNode(func(s string) (float64, error) {
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return 0, errors.New("Please enter a valid float64.")
			}
			return f, nil
		})
	case neoschema.TypeString:
		// Special because we want prefix and suffix
		return &ValueNode[string]{
			prefix: "\"",
			suffix: "\"",
			convert: func(s string) (string, error) {
				return s, nil
			},
			value: "",
		}
	case neoschema.TypeBool:
		return newValueNode(func(s string) (bool, error) {
			b, err := strconv.ParseBool(s)
			if err != nil {
				return false, errors.New("Please enter true/false.")
			}
			return b, nil
		})
	case neoschema.TypeByte:
		return newValueNode(func(s string) (byte, error) {
			i, err := strconv.ParseUint(s, 10, 8)
			if err != nil {
				return 0, errors.New("Please enter a valid byte (0-255).")
			}
			return byte(i), nil
		})
	}

	return nil
}
