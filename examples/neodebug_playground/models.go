package main

import "time"

// ScalarRequest covers every primitive editor, including date and file inputs.
type ScalarRequest struct {
	Int32   int32     `msg:"int32"`
	Int64   int64     `msg:"int64"`
	Float32 float32   `msg:"float32"`
	Float64 float64   `msg:"float64"`
	Bool    bool      `msg:"bool"`
	Byte    byte      `msg:"byte"`
	String  string    `msg:"string"`
	Date    time.Time `msg:"date"`
	File    []byte    `msg:"file"`
}

// Detail is reused to exercise struct references in different containers.
type Detail struct {
	Label   string `msg:"label"`
	Enabled bool   `msg:"enabled"`
}

// CollectionRequest exercises adding, removing and navigating nested values.
type CollectionRequest struct {
	Strings  []string            `msg:"strings"`
	Numbers  []int32             `msg:"numbers"`
	Matrix   [][]float64         `msg:"matrix"`
	Details  []Detail            `msg:"details"`
	Labels   map[string]string   `msg:"labels"`
	ByNumber map[int32]Detail    `msg:"by_number"`
	Groups   map[string][]Detail `msg:"groups"`
	Nested   Detail              `msg:"nested"`
	Another  Detail              `msg:"another"`
	Files    [][]byte            `msg:"files"`
	Dates    []time.Time         `msg:"dates"`
}

// NullableRequest covers toggling null on primitives and composite values.
type NullableRequest struct {
	String         *string            `msg:"string"`
	Int32          *int32             `msg:"int32"`
	Int64          *int64             `msg:"int64"`
	Float32        *float32           `msg:"float32"`
	Float64        *float64           `msg:"float64"`
	Bool           *bool              `msg:"bool"`
	Byte           *byte              `msg:"byte"`
	Date           *time.Time         `msg:"date"`
	File           *[]byte            `msg:"file"`
	Detail         *Detail            `msg:"detail"`
	Strings        *[]string          `msg:"strings"`
	OptionalItems  []*Detail          `msg:"optional_items"`
	OptionalValues map[string]*Detail `msg:"optional_values"`
}

// AllTypesRequest combines the focused fixtures in one large request editor.
type AllTypesRequest struct {
	Scalars     ScalarRequest     `msg:"scalars"`
	Collections CollectionRequest `msg:"collections"`
	Nullable    NullableRequest   `msg:"nullable"`
}
