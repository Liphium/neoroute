package languages

import (
	"fmt"

	"github.com/Liphium/neoroute/neoschema"
)

// We don't need to handle
var goTypeMap = map[neoschema.SchemaType]string{
	neoschema.TypeBool:    "bool",
	neoschema.TypeByte:    "byte",
	neoschema.TypeInt32:   "int32",
	neoschema.TypeInt64:   "int64",
	neoschema.TypeFloat32: "float32",
	neoschema.TypeFloat64: "float64",
	neoschema.TypeString:  "string",
}

// Converts any neoschema.PackedType to a Go struct / type that matches it (in case needed)
func goConversion(packed map[string]any) (string, error) {
	t, ok := packed["type"].(string)
	if !ok {
		return "", fmt.Errorf("packed schema does not have a type field")
	}

	switch neoschema.SchemaType(t) {
	case neoschema.TypeArray:
		// Get the element type
		elemType, ok := packed["element"].(map[string]any)
		if !ok {
			return "", fmt.Errorf("packed schema does not have an element field for array type")
		}

		elemGoType, err := goConversion(elemType)
		if err != nil {
			return "", err
		}

		return "[]" + elemGoType, nil
	}

	return "", nil
}

func goReturnType(packed map[string]any) (string, error) {
	return "", nil
}
