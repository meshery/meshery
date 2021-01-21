package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/helpers/oam"
	"github.com/layer5io/meshery/helpers/oam/core/v1alpha1"
	"github.com/layer5io/meshery/meshes"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// policies are hardcoded here BUT they should be
// fetched from a "service registry" as soon as we have
// one
var policies = [][2]string{
	{"prometheus", "grafana"},
}

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
	isDel := r.Method == http.MethodDelete

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "failed to read request body: %s", err)
		return
	}

	patternFile, err := oam.NewPatternFile(body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "failed to parse to PatternFile: %s", err)
		return
	}

	// Get the application configuration
	aConfig, err := patternFile.ToApplicationConfiguration()
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "invalid Pattern, application configuration couldn't be parsed: %s", err)
		return
	}

	// Get execution plan
	plan, err := oam.CreatePlan(patternFile, policies)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "failed to create an execution plan: %s", err)
		return
	}

	// Check for feasibility
	if feasible := plan.IsFeasible(); !feasible {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, "invalid Pattern, execution is infeasible")
		return
	}

	// Apply application component
	err = plan.Execute(func(acomp v1alpha1.Component) bool {
		fmt.Println("COMPONENT:", acomp.Spec.Type)

		if acomp.Spec.Type == "ServiceMesh" {
			mesh, ok := acomp.Spec.Settings["variant"].(string)
			if !ok {
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "invalid Pattern, invalid ServiceMesh variant: %s", err)
				return false
			}
			fmt.Println("REQUESTED MESH NAME: ", mesh)

			if err := executeAction(r.Context(), prefObj, mesh, acomp.Namespace, mesh, user.UserID, isDel); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(rw, "failed to execute action: %s", err)
				return false
			}
		}

		if acomp.Spec.Type == "ServiceMeshAddon" {
			variant, ok := acomp.Spec.Settings["variant"].(string)
			fmt.Println("VARIANT:", variant)
			if !ok {
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "invalid Pattern, invalid ServiceMeshAddon variant: %s", err)
				return false
			}

			mesh, ok := acomp.Spec.Settings["dependsOn"].([]string)
			if !ok {
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "invalid Pattern, invalid ServiceMeshAddon dependency: %s", err)
				return false
			}

			if err := executeAction(
				r.Context(),
				prefObj, variant+"-addon",
				acomp.Namespace,
				mesh[0],
				user.UserID,
				isDel,
			); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(rw, "failed to execute action: %s", err)
				return false
			}
		}

		return true
	})

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "invalid Pattern, failed to execute the plan: %s", err)
		return
	}

	// Apply application configuration
	for _, comp := range aConfig.Spec.Components {
		// Istio
		if comp.ComponentName == "istio" {
			for _, trait := range comp.Traits {
				// mTLS
				if trait.Name == "mTLS" {
					policy, ok := trait.Properties["policy"].(string)
					if !ok {
						// fallback to strict
						policy = "strict"
					}

					// create name of the operation
					opName := fmt.Sprintf("%s-mtls-policy-operation", policy)

					if err := executeAction(
						r.Context(),
						prefObj,
						opName,
						aConfig.Namespace,
						comp.ComponentName,
						user.UserID,
						isDel,
					); err != nil {
						rw.WriteHeader(http.StatusInternalServerError)
						fmt.Fprintf(rw, "failed to execute action: %s", err)
					}
				}

				// automaticSidecarInjection
				if trait.Name == "automaticSidecarInjection" {
					ns, ok := trait.Properties["namespace"].(string)
					if !ok {
						// fallback to default
						ns = "default"
					}

					// create name of the operation
					opName := "label-namespace"

					if err := executeAction(
						r.Context(),
						prefObj,
						opName,
						ns,
						comp.ComponentName,
						user.UserID,
						isDel,
					); err != nil {
						rw.WriteHeader(http.StatusInternalServerError)
						fmt.Fprintf(rw, "failed to execute action: %s", err)
					}
				}
			}
		}
	}
}

func executeAction(
	ctx context.Context,
	prefObj *models.Preference,
	opName,
	namespace,
	adapter,
	userID string,
	delete bool,
) error {
	meshAdapters := prefObj.MeshAdapters
	if meshAdapters == nil {
		meshAdapters = []*models.Adapter{}
	}

	logrus.Debugf("Adapter Name to execute operations on: %s.", adapter)

	aID := -1
	for i, ad := range meshAdapters {
		if adapter == ad.Name {
			aID = i
		}
	}
	if aID < 0 {
		return errors.New("unable to find a valid adapter for the given adapter name: " + adapter)
	}

	if namespace == "" {
		namespace = "default"
	}

	if prefObj.K8SConfig == nil || !prefObj.K8SConfig.InClusterConfig && (prefObj.K8SConfig.Config == nil || len(prefObj.K8SConfig.Config) == 0) {
		return fmt.Errorf("no valid kubernetes config found")
	}

	mClient, err := meshes.CreateClient(ctx, prefObj.K8SConfig.Config, prefObj.K8SConfig.ContextName, meshAdapters[aID].Location)
	if err != nil {
		return fmt.Errorf("error creating a mesh client: %v", err)
	}
	defer func() {
		_ = mClient.Close()
	}()

	operationID, err := uuid.NewV4()

	if err != nil {
		return fmt.Errorf("error generating an operation id: %v", err)
	}

	_, err = mClient.MClient.ApplyOperation(ctx, &meshes.ApplyRuleRequest{
		OperationId: operationID.String(),
		OpName:      opName,
		Username:    userID,
		Namespace:   namespace,
		CustomBody:  "",
		DeleteOp:    delete,
	})
	if err != nil {
		return err
	}

	return nil
}
