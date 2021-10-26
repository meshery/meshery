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
	"github.com/layer5io/meshery/meshes"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshery/models/pattern/core"
	"github.com/layer5io/meshery/models/pattern/patterns"
	"github.com/layer5io/meshery/models/pattern/stages"
	meshkube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/sirupsen/logrus"
)

// patternCallType is custom type for pattern
// based calls on the adapter
type patternCallType string

const (
	rawAdapter patternCallType = "<raw-adapter>"
	noneLocal  patternCallType = "<none-local>"
	oamAdapter patternCallType = ""
)

// swagger:route POST /api/pattern/deploy PatternsAPI idPostDeployPattern
// Handle POST request for Pattern Deploy
//
// Deploy an attached pattern with the request
// responses:
// 	200:

// swagger:route DELETE /api/pattern/deploy PatternsAPI idDeleteDeployPattern
// Handle DELETE request for Pattern Deploy
//
// Delete a deployed pattern with the request
// responses:
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
	patternFile, err := core.NewPatternFile(body)
	if err != nil {
		h.log.Error(ErrPatternFile(err))
		http.Error(rw, ErrPatternFile(err).Error(), http.StatusInternalServerError)
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

	msg, err := _processPattern(
		token,
		provider,
		patternFile,
		prefObj,
		h.config.KubeClient,
		user.UserID,
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
		h.GETOAMRegisterHandler(typ, rw, r.URL.Query().Get("trim") == "true")
	}
}

func (h *Handler) OAMComponentDetailsHandler(rw http.ResponseWriter, r *http.Request) {
	typ := mux.Vars(r)["type"]

	if !(typ == "workload" || typ == "trait" || typ == "scope") {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	name := mux.Vars(r)["name"]
	res := []interface{}{}

	if typ == "workload" {
		data := core.GetWorkload(name)
		for _, d := range data {
			res = append(res, d)
		}
	}

	if typ == "trait" {
		data := core.GetTrait(name)
		for _, d := range data {
			res = append(res, d)
		}
	}

	if typ == "scope" {
		data := core.GetScope(name)
		for _, d := range data {
			res = append(res, d)
		}
	}

	if err := json.NewEncoder(rw).Encode(res); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		h.log.Debug(err)
		_, _ = rw.Write([]byte(err.Error()))
	}
}

func (h *Handler) OAMComponentDetailByIDHandler(rw http.ResponseWriter, r *http.Request) {
	typ := mux.Vars(r)["type"]

	if !(typ == "workload" || typ == "trait" || typ == "scope") {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	name := mux.Vars(r)["name"]
	id := mux.Vars(r)["id"]
	var res interface{}

	if typ == "workload" {
		res = core.GetWorkloadByID(name, id)
	}

	if typ == "trait" {
		res = core.GetTraitByID(name, id)
	}

	if typ == "scope" {
		res = core.GetScopeByID(name, id)
	}

	if res == nil {
		http.Error(rw, "not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(rw).Encode(res); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		h.log.Debug(err)
		_, _ = rw.Write([]byte(err.Error()))
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
		return core.RegisterWorkload(body)
	}
	if typ == "trait" {
		return core.RegisterTrait(body)
	}
	if typ == "scope" {
		return core.RegisterScope(body)
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
func (h *Handler) GETOAMRegisterHandler(typ string, rw http.ResponseWriter, trim bool) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)

	if typ == "workload" {
		res := core.GetWorkloads()

		// If trim is set to true then remove the schema from the response
		if trim {
			for i := range res {
				res[i].OAMRefSchema = ""
			}
		}

		if err := enc.Encode(res); err != nil {
			h.log.Error(ErrWorkloadDefinition(err))
			http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
		}
	}

	if typ == "trait" {
		res := core.GetTraits()

		// If trim is set to true then remove the schema from the response
		if trim {
			for i := range res {
				res[i].OAMRefSchema = ""
			}
		}

		enc := json.NewEncoder(rw)
		if err := enc.Encode(res); err != nil {
			h.log.Error(ErrTraitDefinition(err))
			http.Error(rw, ErrScopeDefinition(err).Error(), http.StatusInternalServerError)
		}
	}

	if typ == "scope" {
		res := core.GetScopes()

		// If trim is set to true then remove the schema from the response
		if trim {
			for i := range res {
				res[i].OAMRefSchema = ""
			}
		}

		enc := json.NewEncoder(rw)
		if err := enc.Encode(res); err != nil {
			h.log.Error(ErrScopeDefinition(err))
			http.Error(rw, ErrScopeDefinition(err).Error(), http.StatusInternalServerError)
		}
	}
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

func _processPattern(
	token string,
	provider models.Provider,
	pattern core.Pattern,
	prefObj *models.Preference,
	kubeClient *meshkube.Client,
	userID string,
	isDelete bool,
) (string, error) {
	sip := &serviceInfoProvider{
		token:      token,
		provider:   provider,
		opIsDelete: isDelete,
	}
	sap := &serviceActionProvider{
		token:      token,
		provider:   provider,
		prefObj:    prefObj,
		kubeClient: kubeClient,
		opIsDelete: isDelete,
		userID:     userID,

		accumulatedMsgs: []string{},
		err:             nil,
	}

	chain := stages.CreateChain()
	chain.
		Add(stages.ServiceIdentifier(sip, sap)).
		Add(stages.Filler).
		Add(stages.Validator(sip, sap)).
		Add(stages.Provision(sip, sap)).
		Add(stages.Persist(sip, sap)).
		Add(func(data *stages.Data, err error, next stages.ChainStageNextFunction) {
			data.Lock.Lock()
			for k, v := range data.Other {
				if strings.HasSuffix(k, stages.ProvisionSuffixKey) {
					msg, ok := v.(string)
					if ok {
						sap.accumulatedMsgs = append(sap.accumulatedMsgs, msg)
					}
				}
			}
			data.Lock.Unlock()

			sap.err = err
		}).
		Process(&stages.Data{
			Pattern: &pattern,
			Other:   map[string]interface{}{},
		})

	return mergeMsgs(sap.accumulatedMsgs), sap.err
}

type serviceInfoProvider struct {
	provider   models.Provider
	token      string
	opIsDelete bool
}

func (sip *serviceInfoProvider) GetMesheryPatternResource(name, namespace, typ, oamType string) (*uuid.UUID, error) {
	const page = "0"
	const pageSize = "1"
	res, err := sip.provider.GetMesheryPatternResources(sip.token, pageSize, page, "", "", name, namespace, typ, oamType)
	if err != nil {
		return nil, err
	}

	if len(res.Resources) > 0 {
		return res.Resources[0].ID, nil
	}

	return nil, fmt.Errorf("resource not found")
}

func (sip *serviceInfoProvider) GetServiceMesh() (string, string) {
	return "", ""
}

func (sip *serviceInfoProvider) GetAPIVersionForKind(string) string {
	return ""
}

func (sip *serviceInfoProvider) IsDelete() bool {
	return sip.opIsDelete
}

type serviceActionProvider struct {
	token      string
	provider   models.Provider
	prefObj    *models.Preference
	kubeClient *meshkube.Client
	opIsDelete bool
	userID     string

	accumulatedMsgs []string
	err             error
}

func (sap *serviceActionProvider) Terminate(err error) {
	logrus.Error(err)
	sap.err = err
}

func (sap *serviceActionProvider) Provision(ccp stages.CompConfigPair) (string, error) {
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

	for adapter := range ccp.Hosts {
		// Hack until adapters fix the concurrent client
		// creation issue: https://github.com/layer5io/meshery-adapter-library/issues/32
		time.Sleep(50 * time.Microsecond)

		logrus.Debugf("Adapter to execute operations on: %s", adapter)

		if sap.prefObj.K8SConfig == nil ||
			!sap.prefObj.K8SConfig.InClusterConfig &&
				(sap.prefObj.K8SConfig.Config == nil || len(sap.prefObj.K8SConfig.Config) == 0) {
			return "", fmt.Errorf("no valid kubernetes config found")
		}

		// Local call
		if strings.HasPrefix(adapter, string(noneLocal)) {
			resp, err := patterns.ProcessOAM(
				sap.kubeClient,
				[]string{string(jsonComp)},
				string(jsonConfig),
				sap.opIsDelete,
			)

			return resp, err
		}

		// Create mesh client
		mClient, err := meshes.CreateClient(
			context.TODO(),
			sap.prefObj.K8SConfig.Config,
			sap.prefObj.K8SConfig.ContextName,
			adapter,
		)
		if err != nil {
			return "", fmt.Errorf("error creating a mesh client: %v", err)
		}
		defer func() {
			_ = mClient.Close()
		}()

		// Execute operation on the adapter with raw data
		if strings.HasPrefix(adapter, string(rawAdapter)) {
			resp, err := mClient.MClient.ApplyOperation(context.TODO(), &meshes.ApplyRuleRequest{
				Username:  sap.userID,
				DeleteOp:  sap.opIsDelete,
				OpName:    "custom",
				Namespace: "",
			})

			return resp.String(), err
		}

		// Else it is an OAM adapter call
		resp, err := mClient.MClient.ProcessOAM(context.TODO(), &meshes.ProcessOAMRequest{
			Username:  sap.userID,
			DeleteOp:  sap.opIsDelete,
			OamComps:  []string{string(jsonComp)},
			OamConfig: string(jsonConfig),
		})

		return resp.GetMessage(), err
	}

	return "", nil
}

func (sap *serviceActionProvider) Persist(name string, svc core.Service, isUpdate bool) error {
	if !sap.opIsDelete {
		if isUpdate {
			// Do nothing
			return nil
		}

		_, err := sap.provider.SaveMesheryPatternResource(
			sap.token,
			&models.PatternResource{
				ID:        svc.ID,
				Name:      name,
				Namespace: svc.Namespace,
				Type:      svc.Type,
				OAMType:   "workload",
			},
		)

		return err
	}

	return sap.provider.DeleteMesheryPatternResource(
		sap.token,
		svc.ID.String(),
	)
}
