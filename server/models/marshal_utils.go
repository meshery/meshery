package models

import (
	"encoding/json"
	"reflect"
)

// MarshalJSON is a generic function to marshal any struct to JSON bytes
// It returns empty byte slice if marshaling fails
func MarshalJSON(v interface{}) []byte {
	if v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil()) {
		return []byte("{}")
	}

	res, err := json.Marshal(v)
	if err != nil {
		// Log error in production, for now return empty JSON
		return []byte("{}")
	}
	return res
}

// MarshalJSONWithError is a generic function that returns both result and error
func MarshalJSONWithError(v interface{}) ([]byte, error) {
	if v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil()) {
		return []byte("{}"), nil
	}

	return json.Marshal(v)
}
