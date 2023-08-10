package stages

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshery/server/models/pattern/jsonschema"
	"github.com/layer5io/meshery/server/models/pattern/resource/selector"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
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

func Validator(prov ServiceInfoProvider, act ServiceActionProvider, skipValidation bool) ChainStageFunction {
	s := selector.New(act.GetRegistry(), prov)

	return func(data *Data, err error, next ChainStageNextFunction) {
		if err != nil {
			act.Terminate(err)
			return
		}

		data.PatternSvcWorkloadCapabilities = map[string]meshmodel.ComponentDefinition{}
		data.PatternSvcTraitCapabilities = map[string][]core.TraitCapability{}

		for svcName, svc := range data.Pattern.Services {
			wc, err := s.Workload(svc.Type, svc.Version, svc.Model, svc.APIVersion)
			if err != nil {
				act.Terminate(err)
				return
			}
			act.Log(fmt.Sprintf("%s version for %s: %s", svc.Model, svc.Name, wc.Model.Version)) //Eg: kubernetes version for Namespace: v1.25.0
			if core.Format {
				svc.Settings = core.Format.DePrettify(svc.Settings, false)
			}
			//Validate component definition
			if !skipValidation {
				if err := validateWorkload(svc.Settings, wc); err != nil {
					act.Terminate(fmt.Errorf("invalid component configuration for %s: %s", svc.Name, err.Error()))
					return
				}
			}

			if _, ok := specialComps[svc.Type]; ok {
				err := hydrateComponentWithOriginalType(svc.Type, svc.Settings["spec"])
				if err != nil {
					act.Terminate(err)
					return
				}
			}

			// Store the workload capability in the metadata
			data.PatternSvcWorkloadCapabilities[svcName] = wc

			data.PatternSvcTraitCapabilities[svcName] = []core.TraitCapability{}

			//DEPRECATED: `traits` will be no-op for pattern engine
			// Validate traits applied to this workload
			// for trName, tr := range svc.Traits {
			// 	tc, ok := s.Trait(trName)
			// 	if !ok {
			// 		act.Terminate(fmt.Errorf("invalid trait of type: %s", svc.Type))
			// 		return
			// 	}

			// 	if err := validateTrait(tr, tc, svc.Type); err != nil {
			// 		act.Terminate(err)
			// 		return
			// 	}

			// 	// Store the trait capability in the metadata
			// 	data.PatternSvcTraitCapabilities[svcName] = append(data.PatternSvcTraitCapabilities[svcName], tc)
			// }
		}

		if next != nil {
			next(data, nil)
		}
	}
}

func validateWorkload(comp map[string]interface{}, wc meshmodel.ComponentDefinition) error {
	schemaByt := []byte(wc.Schema)
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
