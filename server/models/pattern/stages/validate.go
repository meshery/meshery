package stages

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshery/server/models/pattern/jsonschema"
	"github.com/layer5io/meshery/server/models/pattern/resource/selector"
	"github.com/meshery/schemas/models/v1beta1/component"

	"gopkg.in/yaml.v2"
)

var specialComps = map[string][]string{
	"EnvoyFilter": {"configPatches", "patch", "value"},
	"WasmPlugin":  {"pluginConfig"},
}

func hydrateComponentWithOriginalType(compType string, spec interface{}) error {
	if spec == nil {
		return fmt.Errorf("empty spec provided for component: %s", compType)
	}

	specValue := spec.(map[string]interface{})

	switch compType {
	case "EnvoyFilter":
		configPatches, ok := specValue["configPatches"].([]interface{})
		if !ok {
			return fmt.Errorf("empty configpatch provided for component: %s", compType)
		}
		for _, cp := range configPatches {
			patch, ok := cp.(map[string]interface{})
			if ok {
				value, ok := patch["patch"].(map[string]interface{})
				if !ok {
					return fmt.Errorf("cannot hydrate config for comp: %s", compType)
				}
				err := formatValue("value", value)
				if err != nil {
					return fmt.Errorf("cannot hydrate config for comp: %s", compType)
				}
			}
		}
	case "WasmPlugin":
		err := formatValue("pluginConfig", specValue)
		if err != nil {
			return fmt.Errorf("cannot hydrate config for comp: %s", compType)
		}
	}
	return nil
}

func formatValue(path string, val map[string]interface{}) error {
	updatedValue := make(map[string]interface{})
	byt, _ := val[path].(string)
	_ = yaml.Unmarshal([]byte(byt), &updatedValue)
	val[path] = updatedValue
	return nil
}

func Validator(prov ServiceInfoProvider, act ServiceActionProvider, validate bool) ChainStageFunction {
	s := selector.New(act.GetRegistry())

	return func(data *Data, err error, next ChainStageNextFunction) {
		if err != nil {
			act.Terminate(err)
			return
		}

		for _, component := range data.Pattern.Components {
			wc, err := s.GetDefinition(component.Component.Kind, component.Model.Model.Version, component.Model.Name, component.Component.Version, false)
			if err != nil {
				act.Terminate(err)
				return
			}
			act.Log(fmt.Sprintf("%s version for %s: %s", component.Model.Name, component.DisplayName, wc.Model.Model.Version)) //Eg: kubernetes version for Namespace: v1.25.0
			if core.Format {
				component.Configuration = core.Format.DePrettify(component.Configuration, false)
			}
			//Validate component definition
			if validate {
				if err := validateWorkload(component.Configuration, wc); err != nil {
					act.Terminate(fmt.Errorf("invalid component configuration for %s: %s", component.DisplayName, err.Error()))
					return
				}
			}

			if _, ok := specialComps[component.Component.Kind]; ok {
				err := hydrateComponentWithOriginalType(component.Component.Kind, component.Configuration["spec"])
				if err != nil {
					act.Terminate(err)
					return
				}
			}

			// Store the corresponding definition
			data.DeclartionToDefinitionMapping[component.Id] = wc
		}

		if next != nil {
			next(data, nil)
		}
	}
}

func validateWorkload(comp map[string]interface{}, wc component.ComponentDefinition) error {
	// skip the validation if the component does not have a schema and has isAnnotation set to true.
	isAnnotation := wc.Metadata.IsAnnotation
	if wc.Component.Schema == "" && isAnnotation {
		return nil
	}
	schemaByt := []byte(wc.Component.Schema)
	// Create schema validator from the schema
	rs := jsonschema.GlobalJSONSchema()
	if err := json.Unmarshal(schemaByt, rs); err != nil {
		return fmt.Errorf("failed to create schema: %s", err)
	}

	// Create json settings
	jsonSettings, err := json.Marshal(comp)
	if err != nil {
		return fmt.Errorf("failed to generate schema from the PatternFile settings: %s", err)
	}

	// Validate the json against the schema
	errs, err := rs.ValidateBytes(context.TODO(), jsonSettings)
	if err != nil {
		return fmt.Errorf("error occurred during schema validation: %s", err)
	}
	if len(errs) > 0 {
		return fmt.Errorf("invalid settings: %s", errs)
	}

	return nil
}
