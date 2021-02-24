package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/internal/store"
	"github.com/layer5io/meshery/meshes"
	"github.com/layer5io/meshery/models"
	OAM "github.com/layer5io/meshery/models/oam"
	"github.com/layer5io/meshery/models/oam/core/v1alpha1"
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
	patternFile, err := OAM.NewPatternFile(body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "failed to parse to PatternFile: %s", err)
		return
	}

	// Get execution plan
	plan, err := OAM.CreatePlan(patternFile, policies)
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

	msg, err := createCompConfigPairsAndExecuteAction(
		r.Context(),
		plan,
		patternFile,
		prefObj,
		user,
		isDel,
	)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "Messages:\n%s\nErrors:%s", msg, err)
		return
	}

	fmt.Fprintf(rw, "Messages:\n%s", msg)
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
			_, _ = rw.Write([]byte(err.Error()))
			return
		}
	}
	if method == "GET" {
		GETOAMRegisterHandler(typ, rw)
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

// GETOAMRegisterHandler handles the get requests for the OAM objects
func GETOAMRegisterHandler(typ string, rw http.ResponseWriter) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)

	if typ == "workload" {
		res := OAM.GetWorkloads()

		if err := enc.Encode(res); err != nil {
			logrus.Error("failed to encode workload definitions")
		}
	}

	if typ == "trait" {
		res := OAM.GetTraits()

		enc := json.NewEncoder(rw)
		if err := enc.Encode(res); err != nil {
			logrus.Error("failed to encode trait definitions")
		}
	}

	if typ == "scope" {
		res := OAM.GetScopes()

		enc := json.NewEncoder(rw)
		if err := enc.Encode(res); err != nil {
			logrus.Error("failed to encode scope definitions")
		}
	}
}

func createCompConfigPairsAndExecuteAction(
	ctx context.Context,
	plan *OAM.Plan,
	patternFile OAM.Pattern,
	prefObj *models.Preference,
	user *models.User,
	isDel bool,
) (string, error) {
	var internalErrs []error
	var msgs []string

	aConfig, err := patternFile.GenerateApplicationConfiguration()
	if err != nil {
		return "", err
	}

	_ = plan.Execute(func(svcName string, _ OAM.Service) bool {
		compcon := compConfigPair{
			Hosts: make(map[string]bool),
		}

		// Convert the current service into application component
		comp, err := patternFile.GetApplicationComponent(svcName)
		if err != nil {
			internalErrs = append(internalErrs, err)
			return false
		}

		// Get workload definition corresponding to the type of the workload
		workload, err := getWorkloadDefinition(comp.Spec.Type)
		if err != nil {
			internalErrs = append(internalErrs, err)
			return false
		}

		// Validate the configuration against the workloadDefinition schema
		workloadCap, err := OAM.ValidateWorkload(workload, comp)
		if err != nil {
			internalErrs = append(internalErrs, fmt.Errorf("invalid Pattern: %s", err))
			return false
		}

		compcon.Component = comp
		compcon.Hosts[workloadCap.Host] = true

		// Get component from the configuration file
		configComp, ok := getComponentFromConfiguration(aConfig, comp.Name)
		if !ok {
			msg, err := handleCompConfigPairAction(ctx, compcon, prefObj, user, isDel)
			msgs = append(msgs, msg)
			if err != nil {
				internalErrs = append(internalErrs, err)
				return false
			}

			return true
		}

		for _, tr := range configComp.Traits {
			traitDef, err := getTraitDefinition(tr.Name)
			if err != nil {
				internalErrs = append(internalErrs, err)
				return false
			}

			traitCap, err := OAM.ValidateTrait(traitDef, configComp, patternFile)
			if err != nil {
				internalErrs = append(internalErrs, err)
				return false
			}

			compcon.Hosts[traitCap.Host] = true
		}

		compcon.Configuration = aConfig

		msg, err := handleCompConfigPairAction(ctx, compcon, prefObj, user, isDel)
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
	ccp compConfigPair,
	prefObj *models.Preference,
	user *models.User,
	isDel bool,
) (string, error) {
	var msgs []string

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
		)
		if err != nil {
			msgs = append(msgs, err.Error())
			continue
		}

		msgs = append(msgs, msg)
	}
	return strings.Join(msgs, "\n"), nil
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
) (string, error) {
	logrus.Debugf("Adapter to execute operations on: %s", adapter)

	if prefObj.K8SConfig == nil ||
		!prefObj.K8SConfig.InClusterConfig &&
			(prefObj.K8SConfig.Config == nil || len(prefObj.K8SConfig.Config) == 0) {
		return "", fmt.Errorf("no valid kubernetes config found")
	}

	if callType == noneLocal {
		return "success", nil
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
