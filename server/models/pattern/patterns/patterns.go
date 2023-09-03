package patterns

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/gofrs/uuid"

	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/pattern/patterns/application"
	"github.com/layer5io/meshery/server/models/pattern/patterns/k8s"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	"github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/spf13/viper"
)

func ProcessOAM(kconfigs []string, oamComps []string, oamConfig string, isDel bool, patternName string, ec *models.Signal, userID string, provider models.Provider, hostname registry.IHost, skipCrdAndOperator bool) (string, error) {
	var comps []v1alpha1.Component
	var config v1alpha1.Configuration
	mesheryInstanceID, _ := viper.Get("INSTANCE_ID").(*uuid.UUID)
	userUUID, _ := uuid.FromString(userID)
	action := "deploy"
	if isDel {
		action = "undeploy"
	}
	for _, oamComp := range oamComps {
		var comp v1alpha1.Component
		if err := json.Unmarshal([]byte(oamComp), &comp); err != nil {
			return "", err
		}

		comps = append(comps, comp)
	}

	if err := json.Unmarshal([]byte(oamConfig), &config); err != nil {
		return "", err
	}

	var msgs []string
	var errs []error
	var kclis []*kubernetes.Client
	for _, config := range kconfigs {
		cli, err := kubernetes.New([]byte(config))
		if err != nil {
			errs = append(errs, err)
			continue
		}
		kclis = append(kclis, cli)
	}
	
	var wg sync.WaitGroup
	for _, kcli := range kclis {
		wg.Add(1)
		go func(kcli *kubernetes.Client) {
			defer wg.Done()

			for _, comp := range comps {
				if comp.Spec.Model == "core" {
					if err := application.Deploy(kcli, comp, config, isDel); err != nil {
						var description string
						if isDel {
							description = fmt.Sprintf("Error undeploying %s/%s", patternName, comp.Name)
						} else {
							description = fmt.Sprintf("Error deploying application %s", comp.Name)
						}
						errs = append(errs, err)
						
						// Format bove ProbableCause, SuggestedRemediation,..... as meshkit er and add to metadata
						event := events.NewEvent().FromSystem(*mesheryInstanceID).WithSeverity(events.Error).WithCategory("pattern").WithAction(action).WithDescription(description).FromUser(userUUID).Build()
						err := provider.PersistEvent(event)
						if err != nil {
							// When unable to persist event notify the user, not inside notification center, but have a status symbol in the center to denote whether events are being persisted/subscription is active/.. such event will have category event itself handle them especially.
							evt := events.NewEvent().FromSystem(*mesheryInstanceID).WithSeverity(events.Alert).WithCategory("event").WithAction("persist").WithDescription("Failed persisting events").FromUser(userUUID).Build()
							go ec.Publish(userUUID, evt)
						}
						
						continue
					}
					var description string
					if !isDel {
						description = fmt.Sprintf("Deployed %s/%s", patternName, comp.Name)
						msgs = append(msgs, "Deployed application: "+comp.Name)
					} else {
						description = fmt.Sprintf("Undeployed %s/%s", patternName, comp.Name)
						msgs = append(msgs, "Deleted application: "+comp.Name)
					}
					event := events.NewEvent().FromSystem(*mesheryInstanceID).WithSeverity(events.Informational).WithCategory("pattern").WithAction(action).WithDescription(description).FromUser(userUUID).WithDescription(description).Build()
					err := provider.PersistEvent(event)
					if err != nil {
						evt := events.NewEvent().FromSystem(*mesheryInstanceID).WithSeverity(events.Alert).WithCategory("event").WithAction("persist").WithDescription("Failed persisting events").FromUser(userUUID).Build()
						go ec.Publish(userUUID, evt)
					}
					continue
				}
				if !skipCrdAndOperator && hostname != nil && comp.Spec.Model != (registry.Kubernetes{}).String() {
					var description string
					severity := events.Informational
					if !isDel {
						description = fmt.Sprintf("Detected dependency for %s/%s, deploying dependent model %s.", patternName, comp.Name, comp.Spec.Model)
						msgs = append(msgs, fmt.Sprintf("Deployed %s: %s", comp.Spec.Type, comp.Name))
					} else {
						description = fmt.Sprintf("Detected dependency for %s/%s, undeploying dependent model %s.", patternName, comp.Name, comp.Spec.Model)
						msgs = append(msgs, fmt.Sprintf("Deleted %s: %s", comp.Spec.Type, comp.Name))
					}
					fmt.Println("host: ", hostname, " comp: ", comp.Name, "type: ", comp.Spec.Type, comp)
					// Deploys resources that are required inside cluster for successful deployment of the design.
					result, err := hostname.HandleDependents(comp, kcli, !isDel)
					// If dependencies were not resolved fail forward, there can be case that dependency already exist in the cluster.
					
					eventMetadata := map[string]interface{}{
						"summary": result,
					}

					if err != nil {
						eventMetadata["error"] = err
						severity = events.Error
						errs = append(errs, err)
					}
	
					event := events.NewEvent().FromSystem(*mesheryInstanceID).WithSeverity(severity).WithCategory("pattern").WithAction(action).WithDescription(description).FromUser(userUUID).WithMetadata(eventMetadata).Build()
					err = provider.PersistEvent(event)
					if err != nil {
						evt := events.NewEvent().FromSystem(*mesheryInstanceID).WithSeverity(events.Alert).WithCategory("event").WithAction("persist").WithDescription("Failed persisting events").FromUser(userUUID).Build()
						go ec.Publish(userUUID, evt)		
					}
				}
				//All other components will be handled directly by Kubernetes
				//TODO: Add a Mapper utility function which carries the logic for X hosts can handle Y components under Z circumstances.

				severity := events.Informational
				eventMetadata := make(map[string]interface{})
				description := fmt.Sprintf("Deployed %s/%s.", patternName, comp.Name)
				if err := k8s.Deploy(kcli, comp, config, isDel); err != nil {
					errs = append(errs, err)
					severity = events.Error
					eventMetadata["error"] = err
					var description string
					if isDel {
						description = fmt.Sprintf("Error undeploying %s/%s", patternName, comp.Name)
					} else {
						description = fmt.Sprintf("Error deploying %s/%s", patternName, comp.Name)
					}
	
					event := events.NewEvent().FromSystem(*mesheryInstanceID).WithSeverity(severity).WithCategory("pattern").WithAction(action).WithDescription(description).FromUser(userUUID).WithMetadata(eventMetadata).Build()
					err := provider.PersistEvent(event)
					if err != nil {
						evt := events.NewEvent().FromSystem(*mesheryInstanceID).WithSeverity(severity).WithCategory("event").WithAction("persist").WithDescription("Failed persisting events").FromUser(userUUID).Build()
						go ec.Publish(userUUID, evt)		
					}
					continue
				}
				if !isDel {
					msgs = append(msgs, fmt.Sprintf("Deployed %s: %s", comp.Spec.Type, comp.Name))
				} else {
					description = fmt.Sprintf("Undeployed %s/%s.", patternName, comp.Name)
	
					msgs = append(msgs, fmt.Sprintf("Deleted %s: %s", comp.Spec.Type, comp.Name))
				}
				event := events.NewEvent().FromSystem(*mesheryInstanceID).WithSeverity(severity).WithCategory("pattern").WithAction(action).WithDescription(description).FromUser(userUUID).Build()
				err := provider.PersistEvent(event)
				if err != nil {
					evt := events.NewEvent().FromSystem(*mesheryInstanceID).WithSeverity(events.Alert).WithCategory("event").WithAction("persist").WithDescription("Failed persisting events").FromUser(userUUID).Build()
					go ec.Publish(userUUID, evt)		
				}
			}
		}(kcli)
	}
	wg.Wait()
	return strings.Join(msgs, "\n"), mergeErrors(errs)
}

func mergeErrors(errs []error) error {
	var msgs []string

	for _, err := range errs {
		msgs = append(msgs, err.Error())
	}

	if len(msgs) == 0 {
		return nil
	}

	return fmt.Errorf(strings.Join(msgs, "\n"))
}
