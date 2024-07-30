package stages

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshery/server/helpers"
	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshery/server/models/pattern/planner"
	"github.com/layer5io/meshkit/logger"
	models "github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"

	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/layer5io/meshkit/utils"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

type CompConfigPair struct {
	Component model.ComponentDefinition
	Hosts     []connection.Connection
}

const ProvisionSuffixKey = ".isProvisioned"

func Provision(prov ServiceInfoProvider, act ServiceActionProvider, log logger.Handler) ChainStageFunction {
	return func(data *Data, err error, next ChainStageNextFunction) {
		fmt.Println("line 29 reached", err)

		if err != nil {
			act.Terminate(err)
			return
		}

		processAnnotations(data.Pattern)

		// Create provision plan
		plan, err := planner.CreatePlan(*data.Pattern, prov.IsDelete())
		fmt.Println("plan-----------: 40 ",data.Pattern.Components, err)
		if err != nil {
			act.Terminate(err)
			return
		}

		// Check feasibility of the generated plan
		if !plan.IsFeasible() {
			act.Terminate(fmt.Errorf("infeasible execution: detected cycle in the plan"))
			return
		}

		// config, err := data.Pattern.GenerateApplicationConfiguration()
		// if err != nil {
		// 	act.Terminate(fmt.Errorf("failed to generate application configuration: %s", err))
		// 	return
		// }
		errs := []error{}

		// Execute the plan
		_ = plan.Execute(func(name string, component model.ComponentDefinition) bool {
			ccp := CompConfigPair{}

			core.AssignAdditionalLabels(&component)

			// Create component definition
			// // Create application component
			// comp, err := data.Pattern.GetApplicationComponent(name)
			// if err != nil {
			// 	return false
			// }

			// Generate hosts list
			ccp.Hosts = generateHosts(
				component,
				act.GetRegistry(),
			)

			var annotations map[string]string

			fmt.Println("ccp.Hosts-----------: ", ccp.Hosts)
			_annotations, ok := component.Configuration["annotations"]
			if !ok {
				_annotations = map[string]string{}
			} else {
				annotations, err = utils.Cast[map[string]string](_annotations)
				if err != nil {
					errs = append(errs, err)
					return false
				}
			}


			// Get annotations for the component and merge with existing, if any
			component.Configuration["annotations"] = helpers.MergeStringMaps(
				annotations,
				getAdditionalAnnotations(data.Pattern),
			)

			fmt.Println("ccp.Component-----------: ", ccp.Component)
			ccp.Component = component
			// // Add configuration only if traits are applied to the component
			// if len(svc.Traits) > 0 {
			// 	ccp.Configuration = config
			// }

			msg, err := act.Provision(ccp)
			if err != nil {
				errs = append(errs, err)
				return false
			}
			data.Lock.Lock()
			// Store that this service was provisioned successfully
			data.Other[fmt.Sprintf("%s%s", name, ProvisionSuffixKey)] = msg
			data.Lock.Unlock()

			return true
		}, log)

		if next != nil {
			next(data, mergeErrors(errs))
		}
	}
}

func processAnnotations(pattern *pattern.PatternFile) {
	components := []*model.ComponentDefinition{}
	for _, component := range pattern.Components {
		if !component.Metadata.IsAnnotation {
			components = append(components, component)
		}
	}
	pattern.Components = components
}

func generateHosts(cd model.ComponentDefinition, reg *meshmodel.RegistryManager) []connection.Connection {
	// res := map[connection.Connection]bool{}
	_connection := reg.GetRegistrant(&cd)
	fmt.Println("id :::::::::::", cd.Id.String())
	// tcs, _ := reg.GetRegistrants(&cd)
	// for _, tc := range tcs {
	// 	res[tc.Host] = true
	// }

	return []connection.Connection{_connection}
}

func mergeErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	var errMsg []string
	for _, err := range errs {
		errMsg = append(errMsg, err.Error())
	}

	return fmt.Errorf(strings.Join(errMsg, "\n"))
}

// move into meshkit and change annotations prefix name

func getAdditionalAnnotations(pattern *pattern.PatternFile) map[string]string {
	annotations := make(map[string]string, 2)
	annotations[fmt.Sprintf("%s.name", models.MesheryAnnotationPrefix)] = pattern.Name
	annotations[fmt.Sprintf("%s.id", models.MesheryAnnotationPrefix)] = pattern.Id.String()
	return annotations
}
