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
	var err error
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
			Fields: map[string]PackedType{},
		}
		if current != nil && current.ObjectRegistry() != nil {
			st.BasicType.Objects = current.ObjectRegistry()
		} else {
			st.BasicType.Objects = map[string]PackedType{}
		}
		st.BasicType.Objects[st.Name] = st

		// Go through all struct fields and build their schemas
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			msgTag := field.Tag.Get("msg")
			if msgTag == "-" {
				continue
			}
			if msgTag == "" {
				msgTag = field.Name
			}

			st.Fields[msgTag], err = buildPackedFor(field.Type, st, t, i)
			if err != nil {
				return &BasicType{}, err
			}
		}

		generated = ReferenceType{
			BasicType: &BasicType{
				ActualType: TypeReference,
				Objects:    st.Objects,
			},
			Object: st.Name,
		}

	case reflect.Array:
		// Build the type for the array
		arrayElem, err := buildPackedFor(t.Elem(), current, nil, 0)
		if err != nil {
			return &BasicType{}, err
		}

		generated = &ArrayType{
			BasicType: &BasicType{
				ActualType: Kinds[kind],
			},
			Element: arrayElem,
		}

	case reflect.Map:
		// Build the type for key and map of the array
		mapKey, err := buildPackedFor(t.Key(), current, nil, 0)
		if err != nil {
			return &BasicType{}, err
		}
		mapElem, err := buildPackedFor(t.Elem(), current, nil, 0)
		if err != nil {
			return &BasicType{}, err
		}

		generated = &MapType{
			BasicType: &BasicType{
				ActualType: TypeMap,
			},
			Key:   mapKey,
			Value: mapElem,
		}

	case reflect.Pointer:
		// Build the type for the nullable
		nullableElem, err := buildPackedFor(t.Elem(), current, nil, 0)
		if err != nil {
			return &BasicType{}, err
		}

		generated = &NullableType{
			BasicType: &BasicType{
				ActualType: TypeNullable,
			},
			Element: nullableElem,
		}

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
