package utils

import (
	"encoding/json"
	"fmt"

	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

// ApplyDefaultsToComponentConfiguration mutates component configuration by
// hydrating missing values from JSON schema defaults.
func ApplyDefaultsToComponentConfiguration(comp *component.ComponentDefinition) error {
	if comp == nil {
		return nil
	}

	return ApplyDefaultsToConfiguration(comp.Component.Schema, &comp.Configuration)
}

// ApplyDefaultsToPatternComponents mutates component configurations in a design
// by hydrating missing values from each component schema defaults.
func ApplyDefaultsToPatternComponents(patternFile *pattern.PatternFile) []error {
	if patternFile == nil {
		return nil
	}

	errs := make([]error, 0)
	for _, comp := range patternFile.Components {
		if comp == nil {
			continue
		}

		if err := ApplyDefaultsToComponentConfiguration(comp); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", comp.DisplayName, err))
		}
	}

	return errs
}

// ApplyDefaultsToConfiguration mutates config by hydrating missing values from
// JSON schema defaults. Existing values are never overwritten.
func ApplyDefaultsToConfiguration(schema string, config *map[string]interface{}) error {
	if schema == "" || config == nil {
		return nil
	}

	var schemaNode map[string]interface{}
	if err := json.Unmarshal([]byte(schema), &schemaNode); err != nil {
		return fmt.Errorf("failed to parse component schema for defaults hydration: %w", err)
	}

	if *config == nil {
		*config = map[string]interface{}{}
	}

	hydrated, ok := applySchemaDefaults(schemaNode, *config).(map[string]interface{})
	if !ok {
		// component configuration is expected to be object-like
		*config = map[string]interface{}{}
		return nil
	}

	*config = hydrated
	return nil
}

func applySchemaDefaults(schemaNode map[string]interface{}, value interface{}) interface{} {
	if schemaNode == nil {
		return value
	}

	// Merge defaults from all composed schemas.
	if allOf, ok := schemaNode["allOf"].([]interface{}); ok {
		for _, sub := range allOf {
			subSchema, ok := sub.(map[string]interface{})
			if !ok {
				continue
			}
			value = applySchemaDefaults(subSchema, value)
		}
	}

	if value == nil {
		if def, ok := schemaNode["default"]; ok {
			return deepCopyJSONValue(def)
		}
	}

	schemaType, _ := schemaNode["type"].(string)

	// Object defaults.
	if obj, ok := value.(map[string]interface{}); ok || schemaType == "object" || hasProperties(schemaNode) {
		var objectValue map[string]interface{}
		if ok {
			objectValue = obj
		} else if value == nil {
			objectValue = map[string]interface{}{}
		} else {
			return value
		}

		if props, ok := schemaNode["properties"].(map[string]interface{}); ok {
			for propName, rawPropSchema := range props {
				propSchema, ok := rawPropSchema.(map[string]interface{})
				if !ok {
					continue
				}

				childVal, exists := objectValue[propName]
				if !exists || childVal == nil {
					if propDefault, hasDefault := propSchema["default"]; hasDefault {
						childVal = deepCopyJSONValue(propDefault)
						objectValue[propName] = childVal
						exists = true
					}
				}

				if exists {
					objectValue[propName] = applySchemaDefaults(propSchema, childVal)
				}
			}
		}

		if additionalPropsSchema, ok := schemaNode["additionalProperties"].(map[string]interface{}); ok {
			for key, childVal := range objectValue {
				objectValue[key] = applySchemaDefaults(additionalPropsSchema, childVal)
			}
		}

		value = objectValue
	}

	// Array item defaults.
	if arr, ok := value.([]interface{}); ok || schemaType == "array" {
		var arrayValue []interface{}
		if ok {
			arrayValue = arr
		} else {
			return value
		}

		if itemSchema, ok := schemaNode["items"].(map[string]interface{}); ok {
			for idx := range arrayValue {
				arrayValue[idx] = applySchemaDefaults(itemSchema, arrayValue[idx])
			}
		}

		value = arrayValue
	}

	return value
}

func hasProperties(schemaNode map[string]interface{}) bool {
	props, ok := schemaNode["properties"].(map[string]interface{})
	return ok && len(props) > 0
}

func deepCopyJSONValue(value interface{}) interface{} {
	byt, err := json.Marshal(value)
	if err != nil {
		return value
	}

	var out interface{}
	if err := json.Unmarshal(byt, &out); err != nil {
		return value
	}

	return out
}
