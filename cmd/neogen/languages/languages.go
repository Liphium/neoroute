// The languages package implements the actual generation of the events and endpoint data types into usable structs, classes or whatever those kinds of objects are called in the target language.
package languages

import (
	"fmt"
	"maps"
	"slices"

	go_gen "github.com/Liphium/neoroute/cmd/neogen/languages/go"
	"github.com/Liphium/neoroute/neoschema"
)

type SupportedLanguage string

// All supported languages
const (
	LanguageGo SupportedLanguage = "go"
	// TODO: Support more languages (especially JS/TS, Dart)
)

// All object creation functions to generate all the objects in a packed schema
var objectCreationFunctions = map[SupportedLanguage]func(packed map[string]neoschema.PackedType) (string, error){
	LanguageGo: go_gen.CreateObjects,
}

// All functions to generate the transporters (returns a map of file name -> generated transporter)
var generateTransporterFunctions = map[SupportedLanguage]func(packed neoschema.Schema) (map[string]string, error){
	LanguageGo: go_gen.GenerateTransporters,
}

// Validate that no shit is going on with the schema
func ValidateSchema(packed map[string]neoschema.PackedType) error {

	// TODO: Add the following validations:
	// - Make sure no object registries exist except for in the root
	// - Make sure the object registry only has structs as the root (and none in the leafs + there are no other roots than structs)
	// - Make sure there are no struct types outside of the object registry

	return nil
}

// GenerateObjects generates the object type definitions for a packed type in any supported language.
func GenerateObjects(language SupportedLanguage, packed []neoschema.PackedType) (string, error) {

	// Collect all of the objects together
	registry := map[string]neoschema.PackedType{}
	for _, t := range packed {
		if t.ObjectRegistry() != nil {
			maps.Insert(registry, maps.All(t.ObjectRegistry()))
		}
	}

	// Make sure the schema is valid
	if err := ValidateSchema(registry); err != nil {
		return "", fmt.Errorf("invalid schema: %v", err)
	}

	// Try the conversion
	convFunc, ok := objectCreationFunctions[language]
	if !ok {
		return "", fmt.Errorf("unsupported language: %s", language)
	}
	return convFunc(registry)
}

// GenerateTransporters generates the transporters for a packed schema in any supported language.
func GenerateTransporters(language SupportedLanguage, schema neoschema.Schema) (map[string]string, error) {
	genFunc, ok := generateTransporterFunctions[language]
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
	return genFunc(schema)
}

func GetSupported() []SupportedLanguage {
	return slices.Collect(maps.Keys(objectCreationFunctions))
}
