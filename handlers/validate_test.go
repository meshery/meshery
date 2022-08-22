package handlers

import (
	"testing"
)

var testSchema2 = `{
  "$id": "https://example.com/person.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Person",
  "type": "object",
  "properties": {
    "firstName": {
      "type": "string",
      "description": "The person's first name."
    },
    "lastName": {
      "type": "string",
      "description": "The person's last name."
    },
    "age": {
      "description": "Age in years which must be equal to or greater than zero.",
      "type": "integer",
      "minimum": 0
    }
  }
}`
var testSchema3 = `{
  "$id": "https://example.com/person.schema.json",
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Person",
  "type": "object",
  "properties": {
    "firstName": {
      "type": "string",
      "description": "The person's first name."
    },
    "lastName": {
      "type": "string",
      "description": "The person's last name."
    },
    "age": {
      "description": "Age in years which must be equal to or greater than zero.",
      "type": "integer",
      "minimum": 0
    }
  }, 
  "required": ["age", "lastName"]
}`
var testSchema4 = `{
  "$id": "https://example.com/person.schema.json",
  "title": "Person",
  "type": "object",
  "properties": {
    "firstName": {
      "type": "string",
      "description": "The person's first name."
    },
    "lastName": {
      "type": "string",
      "description": "The person's last name."
    },
    "age": {
      "description": "Age in years which must be equal to or greater than zero.",
      "type": "integer",
      "minimum": 0
    }
  }
}`
var testSchema5 = `
name: string
task: string
properties?: {
  item: "array"
}
`

var value1 = `
firstName: "Meshery"
lastName: "Layer5"
`
var value2 = `
firstName: "Meshery"
lastName: "Layer5"
age: "string"
`
var value3 = `
name: "meshery"
task: "nothing"
`
var value4 = `
firstName: "Meshery"
lastName: "Layer5"
age: 34
`
var value5 = `
firstName: 34
lastName: 54
`

func TestFindSchemaType(t *testing.T) {
	var tests = []struct {
		schema string
		want   validationInputType
	}{
		{testSchema2, jsonschematype},
		{testSchema4, cuetype},
		{testSchema5, cuetype},
	}

	for _, tt := range tests {
		t.Run("schema type", func(t *testing.T) {
			ans := findSchemaType(tt.schema)
			if ans != tt.want {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	var tests = []struct {
		validationItems map[string]validationItem
		want            map[string]bool
	}{
		{
			map[string]validationItem{
				"1": {Schema: testSchema2, Value: value1, ValueType: cuetype},
				"2": {Schema: testSchema2, Value: value2, ValueType: cuetype},
				"3": {Schema: testSchema3, Value: value1, ValueType: cuetype},
				"4": {Schema: testSchema3, Value: value4, ValueType: cuetype},
				"5": {Schema: testSchema5, Value: value3, ValueType: cuetype},
				"6": {Schema: testSchema3, Value: value5, ValueType: cuetype},
			},
			map[string]bool{
				"1": true,
				"2": false,
				"3": false,
				"4": true,
				"5": true,
				"6": false,
			},
		},
	}
	for _, tt := range tests {
		t.Run("meshmodel validate", func(t *testing.T) {
			ans := validate(tt.validationItems)
			for id, vr := range ans {
				if vr.IsValid != tt.want[id] {
					t.Errorf("got %v, want %v for id: %v", ans, tt.want, id)
				}
			}
		})
	}
}
