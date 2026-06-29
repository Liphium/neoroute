// The languages package implements the actual generation of the events and endpoint data types into usable structs, classes or whatever those kinds of objects are called in the target language.
package languages

type SupportedLanguage string

// All supported languages
const (
	LanguageGo SupportedLanguage = "go"
	// TODO: Support more languages (especially JS/TS, Dart)
)

// All conversion functions to generate an object from a packed schema
var conversionFunctions = map[SupportedLanguage]func(packed map[string]any) (string, error){
	LanguageGo: goConversion,
}

// Function to get the acutal return type for a packed schema
var typeFunctions = map[SupportedLanguage]func(packed map[string]any) (string, error){
	LanguageGo: goReturnType,
}

// GenerateObject generates the type defintions for a packed type in any supported language.
func GenerateObject(language SupportedLanguage, packed map[string]any) (string, error) {

	// Approach:
	// 1. Walk the type tree and find lowest hanging types
	// 2. Generate code for lowest hanging types
	// 3. Continue until all types are generated (upper need to embed lower and we don't want duplication)

	return "", nil
}
