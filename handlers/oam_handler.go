package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/meshes"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshery/models/oam"
	OAM "github.com/layer5io/meshery/models/oam"
	"github.com/layer5io/meshery/store"
	"github.com/sirupsen/logrus"
)

// policies are hardcoded here BUT they should be
// fetched from a "service registry" as soon as we have
// one
var policies = [][2]string{}

// PatternFileHandler handles the requested related to pattern files
func (h *Handler) PatternFileHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {

	// Read the PatternFile
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "failed to read request body: %s", err)
		return
	}

	isDel := r.Method == http.MethodDelete

	// Generate the pattern file object
	patternFile, err := oam.NewPatternFile(body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "failed to parse to PatternFile: %s", err)
		return
	}

	// Get execution plan
	plan, err := oam.CreatePlan(patternFile, policies)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "failed to create an execution plan: %s", err)
		return
	}

	logrus.Debugf("Services: %+#v\n", plan.Data.Services)

	// Adpater to be used
	adapter := ""

	// Check for feasibility
	if feasible := plan.IsFeasible(); !feasible {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, "invalid Pattern, execution is infeasible")
		return
	}

	// Get sorted application components
	var comps []string
	var internalErr error
	plan.Execute(func(svcName string, _ oam.Service) bool {
		// Convert the current service into application component
		comp, err := patternFile.GetApplicationComponent(svcName)
		if err != nil {
			internalErr = err
			return false
		}

		// Get the schema for the component
		key := fmt.Sprintf(
			"/meshery/registry/definition/%s/%s/%s",
			"core.oam.dev/v1alpha1",
			"WorkloadDefinition",
			comp.Spec.Type,
		)
		workload, ok := store.Get(key)
		if !ok {
			internalErr = fmt.Errorf("invalid Pattern, service type %s does not exist", comp.Spec.Type)
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, internalErr.Error())
			return false
		}

		workloadCap, err := oam.ValidateWorkload(workload, comp)
		if err != nil {
			internalErr = fmt.Errorf("invalid Pattern: %s", err)
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, internalErr.Error())
			return false
		}

		if adapter == "" {
			adapter = workloadCap.Host
		} else if adapter != workloadCap.Host {
			// Cross adapter config detection
			internalErr = fmt.Errorf("cross adapter components detected" +
				": run \"$mesheryctl pattern get components\" to get list of components and there adapters")
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, internalErr.Error())
			return false
		}

		// Convert component into json
		jsonComp, err := json.Marshal(comp)
		if err != nil {
			internalErr = fmt.Errorf("something went wrong: %s", err)
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, internalErr.Error())
			return false
		}

		comps = append(comps, string(jsonComp))
		return true
	})

	// Check if there were some errors recorded during traversing the components
	if internalErr != nil {
		return
	}

	// Check if the adapter is still not set, fail immediately
	if adapter == "" {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "something went wrong")
		return
	}

	// Get the application configuration
	aConfig, err := patternFile.GenerateApplicationConfiguration()
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "invalid Pattern, application configuration couldn't be parsed: %s", err)
		return
	}

	for _, comps := range aConfig.Spec.Components {
		for _, tr := range comps.Traits {
			// Get the schema for the component
			key := fmt.Sprintf(
				"/meshery/registry/definition/%s/%s/%s",
				"core.oam.dev/v1alpha1",
				"TraitDefinition",
				tr.Name,
			)
			traitDef, ok := store.Get(key)
			if !ok {
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "invalid Pattern, trait %s does not exist", tr.Name)
				return
			}

			trait, err := oam.ValidateTrait(traitDef, comps, patternFile)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "invalid Pattern: %s", err)
				return
			}

			if adapter != trait.Host {
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "cross adapter trait detected"+
					": run \"$mesheryctl pattern get traits\" to get list of traits and there adapters")

				return
			}
		}
	}

	// Convert configuration into json
	jsonConfig, err := json.Marshal(aConfig)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "something went wrong: %s", err)
		return
	}

	// Execute the action
	if err = executeAction(r.Context(), prefObj, adapter, user.UserID, isDel, comps, string(jsonConfig)); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "%s", err)
	}
}

// OAMRegisterHandler handles OAM registry related operations
//
// These operations can be:
// 1. Adding a workload/trait/scope
// 2. Getting list of workloads/traits/scopes
func (h *Handler) OAMRegisterHandler(rw http.ResponseWriter, r *http.Request) {
	typ := mux.Vars(r)["type"]

	if !(typ == "workload" || typ == "trait" || typ == "scope") {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	method := r.Method
	if method == "POST" {
		if err := POSTOAMRegisterHandler(typ, r); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			logrus.Debug(err)
			rw.Write([]byte(err.Error()))
			return
		}
	}
}

// POSTOAMRegisterHandler handles registering OMA objects
func POSTOAMRegisterHandler(typ string, r *http.Request) error {
	// Get the body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if typ == "workload" {
		return OAM.RegisterWorkload(body)
	}
	if typ == "trait" {
		return OAM.RegisterTrait(body)
	}
	if typ == "scope" {
		return OAM.RegisterScope(body)
	}

	return nil
}

// fmt.Println("COMPONENT:", acomp.Spec.Type)

// 		if acomp.Spec.Type == "ServiceMesh" {
// 			mesh, ok := acomp.Spec.Settings["variant"].(string)
// 			if !ok {
// 				rw.WriteHeader(http.StatusBadRequest)
// 				fmt.Fprintf(rw, "invalid Pattern, invalid ServiceMesh variant: %s", err)
// 				return false
// 			}
// 			fmt.Println("REQUESTED MESH NAME: ", mesh)

// 			if err := executeAction(r.Context(), prefObj, mesh, acomp.Namespace, mesh, user.UserID, isDel); err != nil {
// 				rw.WriteHeader(http.StatusInternalServerError)
// 				fmt.Fprintf(rw, "failed to execute action: %s", err)
// 				return false
// 			}
// 		}

// 		if acomp.Spec.Type == "ServiceMeshAddon" {
// 			variant, ok := acomp.Spec.Settings["variant"].(string)
// 			fmt.Println("VARIANT:", variant)
// 			if !ok {
// 				rw.WriteHeader(http.StatusBadRequest)
// 				fmt.Fprintf(rw, "invalid Pattern, invalid ServiceMeshAddon variant: %s", err)
// 				return false
// 			}

// 			mesh, ok := acomp.Spec.Settings["dependsOn"].([]string)
// 			if !ok {
// 				rw.WriteHeader(http.StatusBadRequest)
// 				fmt.Fprintf(rw, "invalid Pattern, invalid ServiceMeshAddon dependency: %s", err)
// 				return false
// 			}

// 			if err := executeAction(
// 				r.Context(),
// 				prefObj, variant+"-addon",
// 				acomp.Namespace,
// 				mesh[0],
// 				user.UserID,
// 				isDel,
// 			); err != nil {
// 				rw.WriteHeader(http.StatusInternalServerError)
// 				fmt.Fprintf(rw, "failed to execute action: %s", err)
// 				return false
// 			}
// 		}
// Apply application configuration
// for _, comp := range aConfig.Spec.Components {
// 	// Istio
// 	if comp.ComponentName == "istio" {
// 		for _, trait := range comp.Traits {
// 			// mTLS
// 			if trait.Name == "mTLS" {
// 				policy, ok := trait.Properties["policy"].(string)
// 				if !ok {
// 					// fallback to strict
// 					policy = "strict"
// 				}

// 				// create name of the operation
// 				opName := fmt.Sprintf("%s-mtls-policy-operation", policy)

// 				if err := executeAction(
// 					r.Context(),
// 					prefObj,
// 					opName,
// 					aConfig.Namespace,
// 					comp.ComponentName,
// 					user.UserID,
// 					isDel,
// 				); err != nil {
// 					rw.WriteHeader(http.StatusInternalServerError)
// 					fmt.Fprintf(rw, "failed to execute action: %s", err)
// 				}
// 			}

// 			// automaticSidecarInjection
// 			if trait.Name == "automaticSidecarInjection" {
// 				ns, ok := trait.Properties["namespace"].(string)
// 				if !ok {
// 					// fallback to default
// 					ns = "default"
// 				}

// 				// create name of the operation
// 				opName := "label-namespace"

// 				if err := executeAction(
// 					r.Context(),
// 					prefObj,
// 					opName,
// 					ns,
// 					comp.ComponentName,
// 					user.UserID,
// 					isDel,
// 				); err != nil {
// 					rw.WriteHeader(http.StatusInternalServerError)
// 					fmt.Fprintf(rw, "failed to execute action: %s", err)
// 				}
// 			}
// 		}
// 	}
// }

func executeAction(
	ctx context.Context,
	prefObj *models.Preference,
	adapter,
	userID string,
	delete bool,
	oamComps []string,
	oamConfig string,
) error {
	logrus.Debugf("Adapter to execute operations on: %s", adapter)

	if prefObj.K8SConfig == nil || !prefObj.K8SConfig.InClusterConfig && (prefObj.K8SConfig.Config == nil || len(prefObj.K8SConfig.Config) == 0) {
		return fmt.Errorf("no valid kubernetes config found")
	}

	mClient, err := meshes.CreateClient(ctx, prefObj.K8SConfig.Config, prefObj.K8SConfig.ContextName, adapter)
	if err != nil {
		return fmt.Errorf("error creating a mesh client: %v", err)
	}
	defer func() {
		_ = mClient.Close()
	}()

	_, err = mClient.MClient.ProcessOAM(ctx, &meshes.ProcessOAMRequest{
		Username:  userID,
		DeleteOp:  delete,
		OamComps:  oamComps,
		OamConfig: oamConfig,
	})
	if err != nil {
		return err
	}

	return nil
}
