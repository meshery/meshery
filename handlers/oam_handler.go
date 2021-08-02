package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/internal/store"
	"github.com/layer5io/meshery/meshes"
	"github.com/layer5io/meshery/models"
	pCore "github.com/layer5io/meshery/models/pattern/core"
	"github.com/layer5io/meshery/models/pattern/patterns"
	pPlanner "github.com/layer5io/meshery/models/pattern/planner"
	pStages "github.com/layer5io/meshery/models/pattern/stages"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	meshkube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/sirupsen/logrus"
)

type compConfigPair struct {
	Component     v1alpha1.Component
	Configuration v1alpha1.Configuration
	Hosts         map[string]bool
}

// patternCallType is custom type for pattern
// based calls on the adapter
type patternCallType string

const (
	rawAdapter patternCallType = "<raw-adapter>"
	noneLocal  patternCallType = "<none-local>"
	oamAdapter patternCallType = ""
)

// policies are hardcoded here BUT they should be
// fetched from a "service registry" as soon as we have
// one
var policies = [][2]string{}

// swagger:route POST /api/pattern/deploy PatternsAPI idPostDeployPattern
// Handle POST request for Pattern Deploy
//
// Deploy an attached pattern with the request
// response:
// 	200:

// swagger:route DELETE /api/pattern/deploy PatternsAPI idDeleteDeployPattern
// Handle DELETE request for Pattern Deploy
//
// Delete a deployed pattern with the request
// response:
// 	200:

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
		h.log.Error(ErrRequestBody(err))
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusInternalServerError)

		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "failed to read request body: %s", err)
		return
	}

	if r.Header.Get("Content-Type") == "application/json" {
		body, err = yaml.JSONToYAML(body)
		if err != nil {
			h.log.Error(ErrPatternFile(err))
			http.Error(rw, ErrPatternFile(err).Error(), http.StatusInternalServerError)
			return
		}
	}

	// Get the user token
	token, err := provider.GetProviderToken(r)
	if err != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(rw, "failed to read request body: %s", err)
		return
	}

	isDel := r.Method == http.MethodDelete

	// Generate the pattern file object
	patternFile, err := pCore.NewPatternFile(body)
	if err != nil {
		h.log.Error(ErrPatternFile(err))
		http.Error(rw, ErrPatternFile(err).Error(), http.StatusInternalServerError)
		return
	}

	// Get execution plan
	plan, err := pPlanner.CreatePlan(patternFile, policies)
	if err != nil {
		h.log.Error(ErrExecutionPlan(err))
		http.Error(rw, ErrExecutionPlan(err).Error(), http.StatusInternalServerError)
		return
	}

	// Check for feasibility
	if feasible := plan.IsFeasible(); !feasible {
		h.log.Error(ErrInvalidPattern(err))
		http.Error(rw, ErrInvalidPattern(err).Error(), http.StatusInternalServerError)
		return
	}

	// If DynamicKubeClient hasn't been created yet then create one
	if h.config.KubeClient.DynamicKubeClient == nil {
		kc, err := meshkube.New(prefObj.K8SConfig.Config)
		if err != nil {
			h.log.Error(ErrInvalidPattern(err))
			http.Error(rw, ErrInvalidPattern(err).Error(), http.StatusInternalServerError)
			return
		}

		h.config.KubeClient = kc
	}

	msg, err := createCompConfigPairsAndExecuteAction(
		r.Context(),
		h,
		prefObj,
		user,
		provider,
		token,
		plan,
		patternFile,
		isDel,
	)

	if err != nil {
		h.log.Error(ErrCompConfigPairs(err))
		http.Error(rw, ErrCompConfigPairs(err).Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(rw, "%s", msg)
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
		if err := h.POSTOAMRegisterHandler(typ, r); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			h.log.Debug(err)
			_, _ = rw.Write([]byte(err.Error()))
			return
		}
	}
	if method == "GET" {
		h.GETOAMRegisterHandler(typ, rw)
	}
}

// swagger:route POST /api/oam/{type} PatternsAPI idPOSTOAMMesheryPattern
// Handles registering OMA objects
//
// Adding a workload/trait/scope
//
// {type} being of either trait, scope, workload; registration of adapter capabilities.
//
// responses:
// 	200:

// POSTOAMRegisterHandler handles registering OMA objects
func (h *Handler) POSTOAMRegisterHandler(typ string, r *http.Request) error {
	// Get the body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if typ == "workload" {
		return pCore.RegisterWorkload(body)
	}
	if typ == "trait" {
		return pCore.RegisterTrait(body)
	}
	if typ == "scope" {
		return pCore.RegisterScope(body)
	}

	return nil
}

// swagger:route GET /api/oam/{type} PatternsAPI idGETOAMMesheryPattern
// Handles the get requests for the OAM objects
//
// Getting list of workloads/traits/scopes
//
// {type} being of either trait, scope, workload; registration of adapter capabilities.
//
// responses:
// 	200:

// GETOAMRegisterHandler handles the get requests for the OAM objects
func (h *Handler) GETOAMRegisterHandler(typ string, rw http.ResponseWriter) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)

	if typ == "workload" {
		res := pCore.GetWorkloads()

		if err := enc.Encode(res); err != nil {
			h.log.Error(ErrWorkloadDefinition(err))
			http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
		}
	}

	if typ == "trait" {
		res := pCore.GetTraits()

		enc := json.NewEncoder(rw)
		if err := enc.Encode(res); err != nil {
			h.log.Error(ErrTraitDefinition(err))
			http.Error(rw, ErrScopeDefinition(err).Error(), http.StatusInternalServerError)
		}
	}

	if typ == "scope" {
		res := pCore.GetScopes()

		enc := json.NewEncoder(rw)
		if err := enc.Encode(res); err != nil {
			h.log.Error(ErrScopeDefinition(err))
			http.Error(rw, ErrScopeDefinition(err).Error(), http.StatusInternalServerError)
		}
	}
}

func createCompConfigPairsAndExecuteAction(
	ctx context.Context,
	h *Handler,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
	token string,
	plan *pPlanner.Plan,
	patternFile pCore.Pattern,
	isDel bool,
) (string, error) {
	var internalErrs []error
	var msgs []string

	aConfig, err := patternFile.GenerateApplicationConfiguration()
	if err != nil {
		return "", err
	}

	_ = plan.Execute(func(svcName string, _ pCore.Service) bool {
		compcon := compConfigPair{
			Hosts: make(map[string]bool),
		}

		// Create an ID for the resource
		_id, err := uuid.NewV4()
		if err != nil {
			internalErrs = append(internalErrs, err)
			return false
		}
		id := _id.String()

		// Convert the current service into application component
		comp, err := patternFile.GetApplicationComponent(svcName)
		if err != nil {
			internalErrs = append(internalErrs, err)
			return false
		}

		// Check if a similar resource has been previously provisioned
		exists := patternResourceExists(
			h,
			prefObj,
			user,
			provider,
			token,
			comp.Name,
			comp.Namespace,
			comp.Spec.Type,
			"workload",
		)

		if exists {
			resourcePage, err := provider.GetMesheryPatternResources(
				token,
				"0",
				"1",
				"",
				"",
				comp.Name,
				comp.Namespace,
				comp.Spec.Type,
				"workload",
			)
			if err != nil {
				internalErrs = append(internalErrs, err)
				return false
			}

			if len(resourcePage.Resources) < 1 {
				err := fmt.Errorf(
					"resource with name: %s, namespace: %s, type: %s, oam_type: %s not found",
					comp.Name,
					comp.Namespace,
					comp.Spec.Type,
					"workload",
				)

				internalErrs = append(internalErrs, err)
				return false
			}

			resource := resourcePage.Resources[0]
			id = resource.ID.String()

			logrus.Debug("resource exists with id: ", id)
		}

		// Assign ID to the component
		comp.SetLabels(map[string]string{
			"resource.pattern.meshery.io/id": id,
		})

		// Get workload definition corresponding to the type of the workload
		workload, err := getWorkloadDefinition(comp.Spec.Type)
		if err != nil {
			internalErrs = append(internalErrs, err)
			return false
		}

		// Validate the configuration against the workloadDefinition schema
		workloadCap, err := pStages.ValidateWorkload(workload, comp)
		if err != nil {
			internalErrs = append(internalErrs, fmt.Errorf("invalid Pattern: %s", err))
			return false
		}

		compcon.Component = comp
		compcon.Hosts[workloadCap.Host] = true

		// Get component from the configuration file
		configComp, ok := getComponentFromConfiguration(aConfig, comp.Name)
		if !ok { // If no configuration exists corresponding to the component then proceed
			msg, err := handleCompConfigPairAction(
				ctx,
				h,
				compcon,
				prefObj,
				user,
				provider,
				token,
				isDel,
				exists && !isDel,
			)
			msgs = append(msgs, msg)
			if err != nil {
				internalErrs = append(internalErrs, err)
				return false
			}

			return true
		}

		// Configuration exists for the component
		// Processing the traits applied to the component
		for _, tr := range configComp.Traits {
			traitDef, err := getTraitDefinition(tr.Name)
			if err != nil {
				internalErrs = append(internalErrs, err)
				return false
			}

			traitCap, err := pStages.ValidateTrait(traitDef, configComp, patternFile)
			if err != nil {
				internalErrs = append(internalErrs, err)
				return false
			}

			compcon.Hosts[traitCap.Host] = true
		}

		compcon.Configuration = aConfig

		msg, err := handleCompConfigPairAction(
			ctx,
			h,
			compcon,
			prefObj,
			user,
			provider,
			token,
			isDel,
			exists && !isDel,
		)
		msgs = append(msgs, msg)
		if err != nil {
			internalErrs = append(internalErrs, err)
			return false
		}

		return true
	})

	return mergeMsgs(msgs), mergeErrors(internalErrs)
}

func handleCompConfigPairAction(
	ctx context.Context,
	h *Handler,
	ccp compConfigPair,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
	token string,
	isDel bool,
	isUpdate bool,
) (string, error) {
	var msgs []string
	var errs []error

	// Marshal the component
	jsonComp, err := json.Marshal(ccp.Component)
	if err != nil {
		return "", fmt.Errorf("failed to serialize the data: %s", err)
	}

	// Marshal the configuration
	jsonConfig, err := json.Marshal(ccp.Configuration)
	if err != nil {
		return "", fmt.Errorf("failed to serialize the data: %s", err)
	}

	// Execute command on the adapters
	for host := range ccp.Hosts {
		// Hack until adapters fix the concurrent client
		// creation issue: https://github.com/layer5io/meshery-adapter-library/issues/32
		time.Sleep(10 * time.Microsecond)

		callType := oamAdapter
		if strings.HasPrefix(host, string(rawAdapter)) {
			callType = rawAdapter
		}
		if strings.HasPrefix(host, string(noneLocal)) {
			callType = noneLocal
		}

		msg, err := executeAction(
			ctx,
			prefObj,
			host,
			user.UserID,
			isDel,
			callType,
			[]string{string(jsonComp)},
			string(jsonConfig),
			h.config.KubeClient,
		)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		msgs = append(msgs, msg)
	}

	// If no errors occurred add the resource and it's NOT and update operation then
	// create an entry corresponding to this resource in the database
	if len(errs) == 0 && !isUpdate {
		id := uuid.FromStringOrNil(ccp.Component.GetLabels()["resource.pattern.meshery.io/id"])
		if isDel {
			if err := provider.DeleteMesheryResource(token, id.String()); err != nil {
				logrus.Error("failed to delete the pattern resource:", err)
			}
		} else {
			if _, err := provider.SaveMesheryPatternResource(token, &models.PatternResource{
				ID:        &id,
				Name:      ccp.Component.Name,
				Namespace: ccp.Component.Namespace,
				Type:      ccp.Component.Spec.Type,
				OAMType:   "workload",
			}); err != nil {
				logrus.Error("failed to save the pattern resource:", err)
			}
		}
	}

	return mergeMsgs(msgs), mergeErrors(errs)
}

func getWorkloadDefinition(key string) (interface{}, error) {
	// Get the schema for the component
	key2 := fmt.Sprintf(
		"/meshery/registry/definition/%s/%s/%s",
		"core.oam.dev/v1alpha1",
		"WorkloadDefinition",
		key,
	)

	workload, ok := store.Get(key2)
	if !ok {
		return workload, fmt.Errorf("invalid Pattern, service type %s does not exist", key)
	}

	return workload, nil
}

func getTraitDefinition(key string) (interface{}, error) {
	// Get the schema for the component
	key2 := fmt.Sprintf(
		"/meshery/registry/definition/%s/%s/%s",
		"core.oam.dev/v1alpha1",
		"TraitDefinition",
		key,
	)

	trait, ok := store.Get(key2)
	if !ok {
		return trait, fmt.Errorf("invalid Pattern, trait %s does not exist", key)
	}

	return trait, nil
}

func getComponentFromConfiguration(
	aConfig v1alpha1.Configuration,
	compName string,
) (v1alpha1.ConfigurationSpecComponent, bool) {
	for _, comps := range aConfig.Spec.Components {
		if comps.ComponentName == compName {
			return comps, true
		}
	}

	return v1alpha1.ConfigurationSpecComponent{}, false
}

func executeAction(
	ctx context.Context,
	prefObj *models.Preference,
	adapter,
	userID string,
	delete bool,
	callType patternCallType,
	oamComps []string,
	oamConfig string,
	kubeClient *meshkube.Client,
) (string, error) {
	logrus.Debugf("Adapter to execute operations on: %s", adapter)

	if prefObj.K8SConfig == nil ||
		!prefObj.K8SConfig.InClusterConfig &&
			(prefObj.K8SConfig.Config == nil || len(prefObj.K8SConfig.Config) == 0) {
		return "", fmt.Errorf("no valid kubernetes config found")
	}

	if callType == noneLocal {
		resp, err := patterns.ProcessOAM(kubeClient, oamComps, oamConfig, delete)

		return resp, err
	}

	mClient, err := meshes.CreateClient(ctx, prefObj.K8SConfig.Config, prefObj.K8SConfig.ContextName, adapter)
	if err != nil {
		return "", fmt.Errorf("error creating a mesh client: %v", err)
	}
	defer func() {
		_ = mClient.Close()
	}()

	if callType == rawAdapter {
		resp, err := mClient.MClient.ApplyOperation(ctx, &meshes.ApplyRuleRequest{
			Username:  userID,
			DeleteOp:  delete,
			OpName:    "custom",
			Namespace: "",
		})

		return resp.String(), err
	}

	if callType == oamAdapter {
		resp, err := mClient.MClient.ProcessOAM(ctx, &meshes.ProcessOAMRequest{
			Username:  userID,
			DeleteOp:  delete,
			OamComps:  oamComps,
			OamConfig: oamConfig,
		})

		return resp.GetMessage(), err
	}

	return "", fmt.Errorf("invalid")
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

func mergeMsgs(msgs []string) string {
	var finalMsgs []string

	for _, msg := range msgs {
		if msg != "" {
			finalMsgs = append(finalMsgs, msg)
		}
	}

	return strings.Join(finalMsgs, "\n")
}

func patternResourceExists(
	h *Handler,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
	token,
	name,
	namespace,
	typ,
	oamType string,
) bool {
	res, err := provider.GetMesheryPatternResources(token, "0", "1", "", "", name, namespace, typ, oamType)
	if err != nil {
		return false
	}

	logrus.Debugf("&&&&&%+v\n", res)

	return len(res.Resources) > 0
}

func _processPattern() {
	// Process pattern in multiple stages
	// Stages:
	// 1. Authorization
	// 2. Filler
	// 3. Validity check
	// 4. Optional Static Analysis
	// 5. Provider persistence
	// 6. Invoke hosts to process the components
}
