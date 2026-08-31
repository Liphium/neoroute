package neoschema

import (
	"fmt"
	"reflect"

)

func notSupportedError(kind reflect.Kind) error {
	return fmt.Errorf("the type %s is not supported", kind.String())
}

// BuildPackedFor generates a schema from a Golang type using the reflect package.
func BuildPackedFor(t reflect.Type) (PackedType, error) {
	generated, err := buildPackedFor(t, nil, nil, 0)
	if err != nil {
		return nil, err
	}

	generated.CleanRegistries(true)
	return generated, nil
}

// buildPackedFor is the internal recursive function.
func buildPackedFor(t reflect.Type, current PackedType, parent reflect.Type, fieldIndex int) (PackedType, error) {
	var generated PackedType
	kind := t.Kind()
	switch kind {
	case reflect.Struct:

		// If the struct is already in the registry, use that instead
		if current != nil && current.ObjectRegistry()[t.Name()] != nil {
			generated = ReferenceType{
				BasicType: &BasicType{
					ActualType: TypeReference,
					Objects:    current.ObjectRegistry(),
				},
				Object: t.Name(),
			}
			break
		}

		st := &StructType{
			Name: t.Name(),
			BasicType: &BasicType{
				ActualType: Kinds[kind],
			},
			Fields: []StructField{},
		}
		st.BasicType.Objects = getRegistry(current)
		st.BasicType.Objects[st.Name] = st

		// Go through all struct fields and build their schemas
		st.Fields = []StructField{}
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			msgTag := field.Tag.Get("msg")
			if msgTag == "-" {
				continue
			}
			if msgTag == "" {
				msgTag = field.Name
			}

			packed, err := buildPackedFor(field.Type, st, t, i)
			if err != nil {
				return &BasicType{}, err
			}

			st.Fields = append(st.Fields, StructField{
				Name: msgTag,
				Type: packed,
			})
		}

		generated = ReferenceType{
			BasicType: &BasicType{
				ActualType: TypeReference,
				Objects:    st.Objects,
			},
			Object: st.Name,
		}

	case reflect.Array, reflect.Slice:

		at := &ArrayType{
			BasicType: &BasicType{
				ActualType: Kinds[kind],
			},
		}
		at.BasicType.Objects = getRegistry(current)

		// Build the type for the array
		arrayElem, err := buildPackedFor(t.Elem(), at, nil, 0)
		if err != nil {
			return &BasicType{}, err
		}
		at.Element = arrayElem

		generated = at

	case reflect.Map:

		mt := &MapType{
			BasicType: &BasicType{
				ActualType: TypeMap,
			},
		}
		mt.BasicType.Objects = getRegistry(current)

		// Build the type for key and map of the array
		mapKey, err := buildPackedFor(t.Key(), mt, nil, 0)
		if err != nil {
			return &BasicType{}, err
		}
		mt.Key = mapKey
		mapElem, err := buildPackedFor(t.Elem(), mt, nil, 0)
		if err != nil {
			return &BasicType{}, err
		}
		mt.Value = mapElem

		generated = mt

	case reflect.Pointer:

		nt := &NullableType{
			BasicType: &BasicType{
				ActualType: TypeNullable,
			},
		}
		nt.BasicType.Objects = getRegistry(current)

		// Build the type for the nullable
		nullableElem, err := buildPackedFor(t.Elem(), nt, nil, 0)
		if err != nil {
			return &BasicType{}, err
		}
		nt.Element = nullableElem

		generated = nt

	default:
		st := Kinds[kind]
		if st == TypeNotSupported {
			return &BasicType{}, notSupportedError(kind)
		} else if st == "" {
			generated = &BasicType{
				ActualType: TypeSerializable,
			}
			break
		}

		generated = &BasicType{
			ActualType: Kinds[kind],
		}
	}

	// Fix registry
	if generated.ObjectRegistry() == nil {
		if current != nil && current.ObjectRegistry() != nil {
			generated.SetRegistry(current.ObjectRegistry())
		} else {
			generated.SetRegistry(map[string]PackedType{})
		}
	}

	// Remove all registries from children
	return generated, nil
}

func getRegistry(current PackedType) map[string]PackedType {
	if current != nil && current.ObjectRegistry() != nil {
		return current.ObjectRegistry()
	}
	return map[string]PackedType{}
}
