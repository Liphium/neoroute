package go_gen

import (
	"fmt"
	"strings"

	"github.com/Liphium/neoroute/cmd/neogen/util"
	"github.com/Liphium/neoroute/neoschema"
)

// We don't need to handle
var goTypeMap = map[neoschema.SchemaType]string{
	neoschema.TypeBool:         "bool",
	neoschema.TypeByte:         "byte",
	neoschema.TypeInt32:        "int32",
	neoschema.TypeInt64:        "int64",
	neoschema.TypeFloat32:      "float32",
	neoschema.TypeFloat64:      "float64",
	neoschema.TypeString:       "string",
	neoschema.TypeSerializable: "neoschema.UsableMessage",
}

// Converts any neoschema.PackedType to a Go struct / type that matches it (in case needed)
func CreateObjects(registry map[string]neoschema.PackedType) (string, error) {
	generated := ""
	offset := "	"

	// Parse and generate stuff from the object registry
	if len(registry) > 0 {
		for name, object := range registry {
			name = util.ToCamelCase(name, true)
			if object.Type() != neoschema.TypeStruct {
				return "", fmt.Errorf("object registry has a non-struct type for object %s", name)
			}
			structType := object.(*neoschema.StructType)

			generated += "\ntype " + name + " struct {"

			// Parse the fields of the struct and add them by their types
			for fieldName, packedType := range structType.Fields {
				goType, err := GetType(packedType)
				if err != nil {
					return "", fmt.Errorf("failed to get Go type for field %s in object %s: %v", fieldName, name, err)
				}

				generated += "\n" + offset + util.ToCamelCase(fieldName, true) + " " + goType + " `msg:\"" + fieldName + "\"`"
			}

			generated += "\n}\n"
		}
	}

	return strings.Trim(generated, "\n"), nil
}

// GetType implements getting the basic type for any packed schema (assuming object generation has been completed)
func GetType(packed neoschema.PackedType) (string, error) {
	switch schema := packed.(type) {

	case *neoschema.StructType:
		if schema.Name == "" {
			return "", fmt.Errorf("packed schema does not have a name field for struct type")
		}

		return util.ToCamelCase(schema.Name, true), nil

	case *neoschema.ReferenceType:
		if schema.Object == "" {
			return "", fmt.Errorf("packed schema does not have an object field for reference type")
		}

		return util.ToCamelCase(schema.Object, true), nil

	case *neoschema.ArrayType:
		elemGoType, err := GetType(schema.Element)
		if err != nil {
			return "", err
		}

		return "[]" + elemGoType, nil

	case *neoschema.NullableType:
		elemGoType, err := GetType(schema.Element)
		if err != nil {
			return "", err
		}

		return "*" + elemGoType, nil

	default:
		t, ok := goTypeMap[schema.Type()]
		if !ok {
			return "", fmt.Errorf("packed schema has an unknown type: %s", schema.Type())
		}

		return t, nil
	}
}
