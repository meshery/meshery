package stages

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshery/server/models/pattern/jsonschema"
	"github.com/layer5io/meshery/server/models/pattern/patterns/k8s"
	"github.com/layer5io/meshery/server/models/pattern/resource/selector"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
)

func Validator(prov ServiceInfoProvider, act ServiceActionProvider) ChainStageFunction {
	s := selector.New(act.GetRegistry(), prov)

	return func(data *Data, err error, next ChainStageNextFunction) {
		if err != nil {
			act.Terminate(err)
			return
		}

		data.PatternSvcWorkloadCapabilities = map[string]meshmodel.ComponentDefinition{}
		data.PatternSvcTraitCapabilities = map[string][]core.TraitCapability{}

		for svcName, svc := range data.Pattern.Services {
			wc, ok := s.Workload(svc.Type, svc.Version, svc.Model)
			if !ok {
				act.Terminate(fmt.Errorf("invalid workload of type: %s", svc.Type))
				return
			}

			if k8s.Format {
				svc.Settings = k8s.Format.DePrettify(svc.Settings, false)
			}
			//Validate workload definition
			if err := validateWorkload(svc.Settings, wc); err != nil {
				act.Terminate(fmt.Errorf("invalid workload definition: %s", err))
				return
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

func validateWorkload(comp map[string]interface{}, wc meshmodel.ComponentDefinition) error {
	// Create schema validator from the schema
	rs := jsonschema.GlobalJSONSchema()
	if err := json.Unmarshal([]byte(wc.Schema), rs); err != nil {
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
