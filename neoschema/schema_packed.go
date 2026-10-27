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
	generated, err := buildPackedFor(t, nil)
	if err != nil {
		return nil, err
	}

	generated.CleanRegistries(true)
	return generated, nil
}

// buildPackedFor is an internal recursive function that generates the schema for Golang type. It's a helper function for the public function above.
func buildPackedFor(t reflect.Type, current PackedType) (PackedType, error) {
	var generated PackedType
	kind := t.Kind()
	switch kind {
	case reflect.Struct:

		// If it's a special type (e.g. time.Time), handle accordingly
		switch t.PkgPath() + "." + t.Name() {
		case "time.Time":
			generated = &BasicType{
				ActualType: TypeDate,
			}
		}
		if generated != nil {
			break
		}

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

		// Add struct to the object registry (we do not add it directly to make sure recursion in structs is supported)
		st := &StructType{
			Name: t.Name(),
			BasicType: &BasicType{
				ActualType: Kinds[kind],
			},
			Fields: []StructField{},
			// If the struct is a tuple, mark it as one
			Tuple: t.Implements(reflect.TypeFor[IsTuple]()),
		}
		st.BasicType.Objects = getRegistry(current)
		st.BasicType.Objects[st.Name] = st

		// Go through all struct fields and build their schemas
		st.Fields = []StructField{}
		for field := range t.Fields() {
			// Unexported fields aren't generated (technically they could be with flags, but we're gonna not support that in neogen)
			if !field.IsExported() {
				continue
			}

			msgTag := field.Tag.Get("msg")
			if msgTag == "-" {
				continue
			}
			if msgTag == "" {
				msgTag = field.Name
			}

			packed, err := buildPackedFor(field.Type, st)
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
		arrayElem, err := buildPackedFor(t.Elem(), at)
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
		mapKey, err := buildPackedFor(t.Key(), mt)
		if err != nil {
			return &BasicType{}, err
		}
		mt.Key = mapKey
		mapElem, err := buildPackedFor(t.Elem(), mt)
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
		nullableElem, err := buildPackedFor(t.Elem(), nt)
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
