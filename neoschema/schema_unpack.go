package neoschema

import (
	"encoding/json"
)

// This makes sure the schema can properly be unpacked into what it was before (MOSTLY AI GENERATED)

var typeRegistry = map[string]func() PackedType{
	string(TypeNotSupported): func() PackedType { return &BasicType{} }, // Catch-all for basic types like "string", "int32", etc.
	string(TypeArray):        func() PackedType { return &ArrayType{} },
	string(TypeStruct):       func() PackedType { return &StructType{} },
	string(TypeReference):    func() PackedType { return &ReferenceType{} },
	string(TypeNullable):     func() PackedType { return &NullableType{} },
}

type rawPackedType struct {
	Type SchemaType `json:"type"`
}

// UnmarshalPackedType uses the registry to instantiate and unmarshal the correct concrete PackedType.
func UnmarshalPackedType(data []byte) (PackedType, error) {
	var probe rawPackedType
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, err
	}

	factory, ok := typeRegistry[string(probe.Type)]
	if !ok {
		factory = typeRegistry[string(TypeNotSupported)]
	}

	obj := factory()

	if err := json.Unmarshal(data, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

// decodePackedMap decodes raw JSON objects to a map of concrete PackedType interface values.
func decodePackedMap(m map[string]json.RawMessage) (map[string]PackedType, error) {
	res := make(map[string]PackedType)

	for k, v := range m {
		obj, err := UnmarshalPackedType(v)
		if err != nil {
			return nil, err
		}
		res[k] = obj
	}

	return res, nil
}

// UnmarshalJSON prevents raw interface unmarshaling errors for BasicType.
func (bt *BasicType) UnmarshalJSON(data []byte) error {
	var aux struct {
		ActualType SchemaType                 `json:"type"`
		Objects    map[string]json.RawMessage `json:"objects,omitempty"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	bt.ActualType = aux.ActualType
	if aux.Objects != nil {
		decoded, err := decodePackedMap(aux.Objects)
		if err != nil {
			return err
		}
		bt.Objects = decoded
	}

	return nil
}

// UnmarshalJSON safely decodes ArrayType and resolves the nested Element interface.
func (at *ArrayType) UnmarshalJSON(data []byte) error {
	var aux struct {
		ActualType SchemaType                 `json:"type"`
		Objects    map[string]json.RawMessage `json:"objects,omitempty"`
		Element    json.RawMessage            `json:"element"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	bt := &BasicType{ActualType: aux.ActualType}
	if aux.Objects != nil {
		decoded, err := decodePackedMap(aux.Objects)
		if err != nil {
			return err
		}
		bt.Objects = decoded
	}
	at.BasicType = bt

	if aux.Element != nil {
		elem, err := UnmarshalPackedType(aux.Element)
		if err != nil {
			return err
		}
		at.Element = elem
	}

	return nil
}

// UnmarshalJSON safely decodes StructType and resolves the nested Fields interface map.
func (st *StructType) UnmarshalJSON(data []byte) error {
	var aux struct {
		ActualType SchemaType                 `json:"type"`
		Objects    map[string]json.RawMessage `json:"objects,omitempty"`
		Name       string                     `json:"name"`
		Fields     map[string]json.RawMessage `json:"fields"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	bt := &BasicType{ActualType: aux.ActualType}
	if aux.Objects != nil {
		decoded, err := decodePackedMap(aux.Objects)
		if err != nil {
			return err
		}
		bt.Objects = decoded
	}
	st.BasicType = bt
	st.Name = aux.Name

	if aux.Fields != nil {
		decoded, err := decodePackedMap(aux.Fields)
		if err != nil {
			return err
		}
		st.Fields = decoded
	}

	return nil
}

// UnmarshalJSON safely decodes ReferenceType.
func (rt *ReferenceType) UnmarshalJSON(data []byte) error {
	var aux struct {
		ActualType SchemaType                 `json:"type"`
		Objects    map[string]json.RawMessage `json:"objects,omitempty"`
		Object     string                     `json:"object"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	bt := &BasicType{ActualType: aux.ActualType}
	if aux.Objects != nil {
		decoded, err := decodePackedMap(aux.Objects)
		if err != nil {
			return err
		}
		bt.Objects = decoded
	}
	rt.BasicType = bt
	rt.Object = aux.Object

	return nil
}

// UnmarshalJSON decodes events and routes on TransporterSchema.
func (t *TransporterSchema) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type   TransporterType            `json:"type"`
		Events map[string]json.RawMessage `json:"events"`
		Routes map[string]RouteSchema     `json:"routes"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	events, err := decodePackedMap(aux.Events)
	if err != nil {
		return err
	}

	t.Type = aux.Type
	t.Events = events
	t.Routes = aux.Routes

	return nil
}

// UnmarshalJSON decodes optional request and response payloads on RouteSchema.
func (r *RouteSchema) UnmarshalJSON(data []byte) error {
	var aux struct {
		HasRequest     bool            `json:"has_request"`
		Request        json.RawMessage `json:"request"`
		HasResponse    bool            `json:"has_response"`
		CanReturnError bool            `json:"can_return_error"`
		Response       json.RawMessage `json:"response"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	r.HasRequest = aux.HasRequest
	r.HasResponse = aux.HasResponse
	r.CanReturnError = aux.CanReturnError

	if aux.Request != nil {
		req, err := UnmarshalPackedType(aux.Request)
		if err != nil {
			return err
		}
		r.Request = req
	}

	if aux.Response != nil {
		res, err := UnmarshalPackedType(aux.Response)
		if err != nil {
			return err
		}
		r.Response = res
	}

	return nil
}

// UnmarshalJSON safely decodes NullableType and resolves the nested Element interface.
func (at *NullableType) UnmarshalJSON(data []byte) error {
	var aux struct {
		ActualType SchemaType                 `json:"type"`
		Objects    map[string]json.RawMessage `json:"objects,omitempty"`
		Element    json.RawMessage            `json:"element"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	bt := &BasicType{ActualType: aux.ActualType}
	if aux.Objects != nil {
		decoded, err := decodePackedMap(aux.Objects)
		if err != nil {
			return err
		}
		bt.Objects = decoded
	}
	at.BasicType = bt

	if aux.Element != nil {
		elem, err := UnmarshalPackedType(aux.Element)
		if err != nil {
			return err
		}
		at.Element = elem
	}

	return nil
}
