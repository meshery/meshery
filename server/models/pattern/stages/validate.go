package stages

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshery/server/models/pattern/jsonschema"
	"github.com/layer5io/meshery/server/models/pattern/patterns/k8s"
	"github.com/layer5io/meshery/server/models/pattern/resource/selector"
)

func Validator(prov ServiceInfoProvider, act ServiceActionProvider) ChainStageFunction {
	s := selector.New(prov)

	return func(data *Data, err error, next ChainStageNextFunction) {
		if err != nil {
			act.Terminate(err)
			return
		}

		data.PatternSvcWorkloadCapabilities = map[string]core.WorkloadCapability{}
		data.PatternSvcTraitCapabilities = map[string][]core.TraitCapability{}

		for svcName, svc := range data.Pattern.Services {
			wc, ok := s.Workload(svc.Type)
			if !ok {
				act.Terminate(fmt.Errorf("invalid workload of type: %s", svc.Type))
				return
			}

			var svcSettings map[string]interface{}
			//deep copy settings for validation
			if k8s.Format {
				svcSettings = k8s.Format.Prettify(svc.Settings, true)
			} else {
				svcSettings = svc.Settings
			}

			//Validate workload definition
			if err := validateWorkload(svcSettings, wc); err != nil {
				act.Terminate(fmt.Errorf("invalid workload definition: %s", err))
				return
			}
			if k8s.Format {
				svc.Settings = k8s.Format.DePrettify(svc.Settings, false)
			}
			// Store the workload capability in the metadata
			data.PatternSvcWorkloadCapabilities[svcName] = wc

			data.PatternSvcTraitCapabilities[svcName] = []core.TraitCapability{}

			// Validate traits applied to this workload
			for trName, tr := range svc.Traits {
				tc, ok := s.Trait(trName)
				if !ok {
					act.Terminate(fmt.Errorf("invalid trait of type: %s", svc.Type))
					return
				}

				if err := validateTrait(tr, tc, svc.Type); err != nil {
					act.Terminate(err)
					return
				}

				// Store the trait capability in the metadata
				data.PatternSvcTraitCapabilities[svcName] = append(data.PatternSvcTraitCapabilities[svcName], tc)
			}
		}

		if next != nil {
			next(data, nil)
		}
	}
}

func validateWorkload(comp map[string]interface{}, wc core.WorkloadCapability) error {
	// Create schema validator from the schema
	rs := jsonschema.GlobalJSONSchema()
	if err := json.Unmarshal([]byte(wc.OAMRefSchema), rs); err != nil {
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

func validateTrait(trait interface{}, tc core.TraitCapability, compType string) error {
	// Create schema validator from the schema
	rs := jsonschema.GlobalJSONSchema()
	if err := json.Unmarshal([]byte(tc.OAMRefSchema), rs); err != nil {
		return fmt.Errorf("failed to create schema: %s", err)
	}

	// Check if the trait applied to the component is legal or not
	isLegal := isLegalTrait(tc, compType)
	if !isLegal {
		return fmt.Errorf(
			"%s trait is not applicable to %s",
			tc.OAMDefinition.Name,
			compType,
		)
	}

	// Create json of the trait's properties for validation
	jsonCompTraitProp, err := json.Marshal(trait)
	if err != nil {
		return fmt.Errorf(
			"failed to generate schema from the PatternFile's service type %s trait: %s",
			compType,
			err,
		)
	}

	// Validate the json against the schema
	errs, err := rs.ValidateBytes(context.TODO(), jsonCompTraitProp)
	if err != nil {
		return fmt.Errorf("error occurred during schema validation: %s", err)
	}
	if len(errs) > 0 {
		return fmt.Errorf("invalid traits: %s", errs)
	}

	return nil
}

func isLegalTrait(tc core.TraitCapability, compType string) bool {
	if len(tc.OAMDefinition.Spec.AppliesToWorkloads) == 0 {
		return true
	}

	for _, wl := range tc.OAMDefinition.Spec.AppliesToWorkloads {
		if wl == compType {
			return true
		}
	}

	return false
}
