package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/meshes"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshery/server/models/pattern/patterns"
	"github.com/layer5io/meshery/server/models/pattern/patterns/k8s"
	"github.com/layer5io/meshery/server/models/pattern/stages"
	"github.com/layer5io/meshkit/logger"
	events "github.com/layer5io/meshkit/models/events"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	meshkube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/client-go/applyconfigurations/meta/v1"
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
	userID := uuid.FromStringOrNil(user.ID)

	// Read the PatternFile
	body, err := io.ReadAll(r.Body)
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

	isDel := r.Method == http.MethodDelete
	action := "deploy"
	if isDel {
		action = "undeploy"
	}
	// Generate the pattern file object
	patternFile, err := core.NewPatternFile(body)
	if err != nil {
		h.log.Error(ErrPatternFile(err))
		http.Error(rw, ErrPatternFile(err).Error(), http.StatusInternalServerError)
		return
	}

	response, err := _processPattern(
		r.Context(),
		provider,
		patternFile,
		prefObj,
		user.ID,
		isDel,
		r.URL.Query().Get("verify") == "true",
		r.URL.Query().Get("dryRun") == "true",
		r.URL.Query().Get("skipCRD") == "true",
		false,
		h.registryManager,
		h.config.EventBroadcaster,
		h.log,
	)

	patternID := uuid.FromStringOrNil(patternFile.PatternID)
	eventBuilder := events.NewEvent().ActedUpon(patternID).FromUser(userID).FromSystem(*h.SystemID).WithCategory("pattern").WithAction(action)

	if err != nil {
		err := ErrCompConfigPairs(err)
		metadata := map[string]interface{}{
			"error": err,
		}

		event := eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Error %sing pattern %s", action, patternFile.Name)).WithMetadata(metadata).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)

		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	metadata := map[string]interface{}{
		"summary": response,
	}

	event := eventBuilder.WithSeverity(events.Informational).WithDescription(fmt.Sprintf("Pattern %s deployed", patternFile.Name)).WithMetadata(metadata).Build()
	_ = provider.PersistEvent(event)
	go h.config.EventBroadcaster.Publish(userID, event)

	ec := json.NewEncoder(rw)
	_ = ec.Encode(response)
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
	ctx context.Context,
	provider models.Provider,
	pattern core.Pattern,
	prefObj *models.Preference,
	userID string,
	isDelete bool,
	verify bool,
	dryRun bool,
	skipCrdAndOperator bool,
	skipPrintLogs bool,
	registry *meshmodel.RegistryManager,
	ec *models.Broadcast,
	l logger.Handler,
) (map[string]interface{}, error) {
	resp := make(map[string]interface{})

	// Get the token from the context
	token, ok := ctx.Value(models.TokenCtxKey).(string)
	if !ok {
		return nil, ErrRetrieveUserToken(fmt.Errorf("token not found in the context"))
	}
	// // Get the kubehandler from the context
	k8scontexts, ok := ctx.Value(models.KubeClustersKey).([]models.K8sContext)
	if !ok || len(k8scontexts) == 0 {
		return nil, ErrInvalidKubeHandler(fmt.Errorf("failed to find k8s handler"), "_processPattern couldn't find a valid k8s handler")
	}

	// // Get the kubernetes config from the context
	// kubecfg, ok := ctx.Value(models.KubeConfigKey).([]byte)
	// if !ok || kubecfg == nil {
	// 	return "", ErrInvalidKubeConfig(fmt.Errorf("failed to find k8s config"), "_processPattern couldn't find a valid k8s config")
	// }

	// // Get the kubernetes context from the context
	// mk8scontext, ok := ctx.Value(models.KubeContextKey).(*models.K8sContext)
	// if !ok || mk8scontext == nil {
	// 	return "", ErrInvalidKubeContext(fmt.Errorf("failed to find k8s context"), "_processPattern couldn't find a valid k8s context")
	// }
	var ctxToconfig = make(map[string]string)
	for _, ctx := range k8scontexts {
		cfg, err := ctx.GenerateKubeConfig()
		if err != nil {
			return nil, ErrInvalidKubeConfig(fmt.Errorf("failed to find k8s config"), "_processPattern couldn't find a valid k8s config")
		}
		ctxToconfig[ctx.ID] = string(cfg)
		// configs = append(configs, string(cfg))
	}
	internal := func(mk8scontext []models.K8sContext) (map[string]interface{}, error) {
		sip := &serviceInfoProvider{
			token:      token,
			provider:   provider,
			opIsDelete: isDelete,
		}
		sap := &serviceActionProvider{
			token:    token,
			log:      l,
			provider: provider,
			prefObj:  prefObj,
			// kubeClient:    kubeClient,
			opIsDelete: isDelete,
			userID:     userID,
			registry:   registry,
			// kubeconfig:    kubecfg,
			// kubecontext:   mk8scontext,
			skipPrintLogs:      skipPrintLogs,
			skipCrdAndOperator: skipCrdAndOperator,
			ctxTokubeconfig:    ctxToconfig,
			accumulatedMsgs:    []string{},
			err:                nil,
			eventsChannel:      ec,
			patternName:        strings.ToLower(pattern.Name),
		}
		chain := stages.CreateChain()
		chain.
			Add(stages.Import(sip, sap)).
			Add(stages.ServiceIdentifierAndMutator(sip, sap)).
			Add(stages.Filler(skipPrintLogs)).
			// Calling this stage `The Validation stage` is a bit deceiving considering
			// that the validation stage also formats the `data` (chain function parameter) that the
			// subsequent stages depend on.
			// We are skipping the `Validation` part in case of dryRun
			Add(stages.Validator(sip, sap, dryRun))
		if dryRun {
			chain.Add(stages.DryRun(sip, sap))
		}
		if !verify && !dryRun {
			chain.
				Add(stages.Provision(sip, sap)).
				Add(stages.Persist(sip, sap))
		}
		chain.
			Add(func(data *stages.Data, err error, next stages.ChainStageNextFunction) {
				data.Lock.Lock()
				for k, v := range data.Other {
					if strings.HasSuffix(k, stages.ProvisionSuffixKey) {
						msg, ok := v.(string)
						if ok {
							sap.accumulatedMsgs = append(sap.accumulatedMsgs, msg)
						}
					}
					if k == stages.DryRunResponseKey {
						if v != nil {
							resp["dryRunResponse"] = v
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
		resp["messages"] = mergeMsgs(sap.accumulatedMsgs)
		return resp, sap.err
	}
	return internal(k8scontexts)
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
	token    string
	log      logger.Handler
	provider models.Provider
	prefObj  *models.Preference
	// kubeClient      *meshkube.Client
	ctxTokubeconfig map[string]string
	opIsDelete      bool
	userID          string
	// kubeconfig  []byte
	// kubecontext     *models.K8sContext
	skipCrdAndOperator bool
	skipPrintLogs      bool
	accumulatedMsgs    []string
	err                error
	eventsChannel      *models.Broadcast
	registry           *meshmodel.RegistryManager
	patternName        string
}

func (sap *serviceActionProvider) GetRegistry() *meshmodel.RegistryManager {
	return sap.registry
}

func (sap *serviceActionProvider) Log(msg string) {
	if sap.log != nil {
		sap.log.Info(msg)
	}
}
func (sap *serviceActionProvider) Terminate(err error) {
	if !sap.skipPrintLogs {
		logrus.Error(err)
	}
	sap.err = err
}
func (sap *serviceActionProvider) Mutate(p *core.Pattern) {
	//TODO: externalize these mutation rules with policies.
	//1. Enforce the deployment of CRDs before other resources
	for name, svc := range p.Services {
		if svc.Type == "CustomResourceDefinition.K8s" {
			for _, svc := range p.Services {
				if svc.Type != "CustomResourceDefinition.K8s" {
					svc.DependsOn = append(svc.DependsOn, name)
				}
			}
		}
	}
}

// NOTE: Currently tied to kubernetes
// Returns ComponentName->ContextID->Response
func (sap *serviceActionProvider) DryRun(comps []v1alpha1.Component) (resp map[string]map[string]core.DryRunResponseWrapper, err error) {
	for _, cmp := range comps {
		for ctxID, kc := range sap.ctxTokubeconfig {
			cl, err := meshkube.New([]byte(kc))
			if err != nil {
				return resp, err
			}

			// status represents kubernetes status object
			status, ok, err := k8s.DryRunHelper(cl, cmp)
			dResp := core.DryRunResponseWrapper{Success: ok, Component: &core.Service{
				Name:        cmp.Name,
				Type:        cmp.Spec.Type,
				Namespace:   cmp.Namespace,
				APIVersion:  cmp.Spec.APIVersion,
				Version:     cmp.Spec.Version,
				Model:       cmp.Spec.Model,
				Labels:      cmp.Labels,
				Annotations: cmp.Annotations,
			}}

			// Dry run was success
			if ok {
				dResp.Component.Settings = make(map[string]interface{})
				for k, v := range status {
					if k == "apiVersion" || k == "kind" || k == "metadata" {
						continue
					}
					dResp.Component.Settings[k] = v
				}
			} else if err != nil { //Dry run failed due to some error eg: K8s server could not identify the resource
				dResp.Error = &core.DryRunResponse{
					Status: err.Error(),
				}
			} else { //Dry run failure returned with an error wrapped in kubernetes custom error					
				dResp.Error, err = convertRawDryRunResponse(cmp.Name, status)
				if err != nil {
					return nil, err
				}
			}
			if resp == nil {
				resp = make(map[string]map[string]core.DryRunResponseWrapper)
			}
			if resp[cmp.Name] == nil {
				resp[cmp.Name] = make(map[string]core.DryRunResponseWrapper)
			}
			resp[cmp.Name][ctxID] = dResp
		}
	}
	return
}

func convertRawDryRunResponse(componentName string, status map[string]interface{}) (*core.DryRunResponse, error) {
	response := core.DryRunResponse{}

	byt, err := json.Marshal(status)
	if err != nil {
		return nil, err
	}

	var a v1.StatusApplyConfiguration
	err = json.Unmarshal(byt, &a)
	if err != nil {
		return nil, err
	}
	
	if a.Status != nil {
		response.Status = *a.Status
	}

	response.Causes = make([]core.DryRunFailureCause, 0)
	if a.Details != nil {
		for _, cause := range a.Details.Causes {
			msg := ""
			field := ""
			typ := ""
			if cause.Message != nil {
				msg = *cause.Message
			}
			if cause.Field != nil {
				field = componentName + "." + utils.GetComponentFieldPathFromK8sFieldPath(*cause.Field)
			}
			if cause.Type != nil {
				typ = string(*cause.Type)
			}
			failureCase := core.DryRunFailureCause{Message: msg, FieldPath: field, Type: typ}
			response.Causes = append(response.Causes, failureCase)
		}
	}

	if len(response.Causes) == 0 && a.Message != nil {
		response.Status = *a.Message
	}
	return &response, nil
}

func (sap *serviceActionProvider) Provision(ccp stages.CompConfigPair) (string, error) { // Marshal the component
	jsonComp, err := json.Marshal(ccp.Component)
	if err != nil {
		return "", fmt.Errorf("failed to serialize the data: %s", err)
	}

	// Marshal the configuration
	jsonConfig, err := json.Marshal(ccp.Configuration)
	if err != nil {
		return "", fmt.Errorf("failed to serialize the data: %s", err)
	}

	for host := range ccp.Hosts {
		// Hack until adapters fix the concurrent client
		// creation issue: https://github.com/layer5io/meshery-adapter-library/issues/32
		time.Sleep(50 * time.Microsecond)

		logrus.Debugf("Adapter to execute operations on: %s", host.Hostname)

		// Local call
		if host.Port == 0 {
			//TODO: Accommodate internal calls to use context mapping with kubeconfig
			var kconfigs []string
			for _, v := range sap.ctxTokubeconfig {
				kconfigs = append(kconfigs, v)
			}
			resp, err := patterns.ProcessOAM(
				kconfigs,
				[]string{string(jsonComp)},
				string(jsonConfig),
				sap.opIsDelete,
				sap.patternName,
				sap.eventsChannel,
				sap.userID,
				sap.provider,
				host.IHost,
				sap.skipCrdAndOperator,
			)
			return resp, err
		}
		addr := host.Hostname
		if host.Port != 0 {
			addr += ":" + strconv.Itoa(host.Port)
		}
		// Create mesh client
		mClient, err := meshes.CreateClient(
			context.TODO(),
			addr,
		)
		if err != nil {
			return "", fmt.Errorf("error creating a mesh client: %v", err)
		}
		defer func() {
			_ = mClient.Close()
		}()

		// Execute operation on the adapter with raw data
		// if strings.HasPrefix(adapter, string(rawAdapter)) {
		// 	resp, err := mClient.MClient.ApplyOperation(context.TODO(), &meshes.ApplyRuleRequest{
		// 		Username:    sap.userID,
		// 		DeleteOp:    sap.opIsDelete,
		// 		OpName:      "custom",
		// 		Namespace:   "",
		// 		KubeConfigs: sap.kubeconfigs,
		// 	})

		// 	return resp.String(), err
		// }

		// Else it is an OAM adapter call
		//TODO: Accommodate gRPC calls to use context mapping with kubeconfig
		var kconfigs []string
		for _, v := range sap.ctxTokubeconfig {
			kconfigs = append(kconfigs, v)
		}
		resp, err := mClient.MClient.ProcessOAM(context.TODO(), &meshes.ProcessOAMRequest{
			Username:    sap.userID,
			DeleteOp:    sap.opIsDelete,
			OamComps:    []string{string(jsonComp)},
			OamConfig:   string(jsonConfig),
			KubeConfigs: kconfigs,
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
