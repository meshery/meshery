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

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/meshes"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshery/server/models/pattern/patterns"
	"github.com/spf13/viper"

	"github.com/layer5io/meshery/server/models/pattern/patterns/k8s"
	patternutils "github.com/layer5io/meshery/server/models/pattern/utils"

	"github.com/layer5io/meshery/server/models/pattern/stages"
	"github.com/layer5io/meshkit/logger"
	events "github.com/layer5io/meshkit/models/events"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/layer5io/meshkit/utils"
	meshkube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	token, _ := r.Context().Value(models.TokenCtxKey).(string)
	var payload models.MesheryPatternFileDeployPayload
	var patternFileByte []byte

	// Read the PatternFile
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusInternalServerError)

		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "failed to read request body: %s", err)
		return
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusInternalServerError)

		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "failed to unmarshal request body: %s", err)
		return
	}

	patternFileByte = []byte(payload.PatternFile)

	queryParams, _ := extractBoolQueryParams(r, "dryRun", "skipCRD", "verify", "upgrade")
	isDryRun := queryParams["dryRun"]
	skipCRDAndOperator := queryParams["skipCRD"]
	validate := queryParams["verify"]
	upgradeExistingRelease := queryParams["upgrade"]

	isDelete := r.Method == http.MethodDelete
	action := "deploy"
	if isDelete {
		action = "undeploy"
	}
	patternID := payload.PatternID

	isDesignInAlpha2Format, err := patternutils.IsDesignInAlpha2Format(payload.PatternFile)
	if err != nil {
		err = ErrPatternFile(err)
		event := events.NewEvent().ActedUpon(payload.PatternID).FromSystem(*h.SystemID).FromUser(userID).WithCategory("pattern").WithAction("view").WithDescription("Failed to parse design").WithMetadata(map[string]interface{}{"error": err, "id": payload.PatternID}).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if isDesignInAlpha2Format {
		eventBuilder := events.NewEvent().ActedUpon(patternID).FromSystem(*h.SystemID).FromUser(userID).WithCategory("pattern").WithAction("convert")

		_, patternFileStr, err := h.convertV1alpha2ToV1beta1(&models.MesheryPattern{
			ID:          &patternID,
			PatternFile: payload.PatternFile,
		}, eventBuilder)

		event := eventBuilder.Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)

		if err != nil {
			h.log.Error(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		patternFileByte = []byte(patternFileStr)
	}

	patternFile, err := core.NewPatternFile(patternFileByte)
	if err != nil {
		h.log.Error(ErrPatternFile(err))
		http.Error(rw, ErrPatternFile(err).Error(), http.StatusInternalServerError)
		return
	}

	if patternID == uuid.Nil {
		patternID = patternFile.Id
	}

	patternFile.Id = patternID
	// Generate the pattern file object
	description := fmt.Sprintf("%sed design '%s'", action, patternFile.Name)
	if isDryRun {
		action = "verify"
		description = fmt.Sprintf("%sed design '%s'", action, patternFile.Name)
	}

	opts := &core.ProcessPatternOptions{
		Context:                r.Context(),
		Provider:               provider,
		Pattern:                patternFile,
		PrefObj:                prefObj,
		UserID:                 user.ID,
		IsDelete:               isDelete,
		Validate:               validate,
		DryRun:                 isDryRun,
		SkipCRDAndOperator:     skipCRDAndOperator,
		UpgradeExistingRelease: upgradeExistingRelease,
		SkipPrintLogs:          viper.GetBool("DEBUG"),
		Registry:               h.registryManager,
		EventBroadcaster:       h.config.EventBroadcaster,
		Log:                    h.log,
	}
	response, err := _processPattern(opts)

	eventBuilder := events.NewEvent().ActedUpon(patternID).FromUser(userID).FromSystem(*h.SystemID).WithCategory("pattern").WithAction(action)

	if err != nil {
		err := ErrPatternDeploy(err, patternFile.Name)
		metadata := map[string]interface{}{
			"error": err,
		}

		event := eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Failed to %s design '%s'.", action, patternFile.Name)).WithMetadata(metadata).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)

		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	metadata := map[string]interface{}{
		"summary": response,
	}

	serverURL, _ := r.Context().Value(models.MesheryServerURL).(string)

	if action == "deploy" {
		viewLink := fmt.Sprintf("%s/extension/meshmap?mode=visualize&design=%s", serverURL, patternID)
		description = fmt.Sprintf("%s.", description)
		metadata["view_link"] = viewLink
	}

	event := eventBuilder.WithSeverity(events.Success).WithDescription(description).WithMetadata(metadata).Build()
	_ = provider.PersistEvent(event)
	go func() {
		h.config.EventBroadcaster.Publish(userID, event)
		err = provider.PublishEventToProvider(token, *event)
		if err != nil {
			h.log.Warn(ErrPersistEventToRemoteProvider(err))
		}
	}()

	ec := json.NewEncoder(rw)
	_ = ec.Encode(response)
}

func _processPattern(opts *core.ProcessPatternOptions) (map[string]interface{}, error) {
	resp := make(map[string]interface{})

	// Get the token from the context
	token, ok := opts.Context.Value(models.TokenCtxKey).(string)
	if !ok {
		return nil, ErrRetrieveUserToken(fmt.Errorf("token not found in the context"))
	}

	// Get the kubehandler from the context
	k8scontexts, ok := opts.Context.Value(models.KubeClustersKey).([]models.K8sContext)
	if !ok || len(k8scontexts) == 0 {
		return nil, ErrInvalidKubeHandler(fmt.Errorf("Meshery server failed to interact with the Kubernetes cluster due to the cluster not being available."), opts.Pattern.Name)
	}

	var ctxToconfig = make(map[string]string)
	for _, ctx := range k8scontexts {
		cfg, err := ctx.GenerateKubeConfig()
		if err != nil {
			return nil, ErrInvalidKubeConfig(fmt.Errorf("failed to find Kubernetes config"), "_processPattern couldn't find a valid Kubernetes config")
		}
		ctxToconfig[ctx.ID] = string(cfg)
	}

	internal := func(mk8scontext []models.K8sContext) (map[string]interface{}, error) {
		sip := &serviceInfoProvider{
			token:      token,
			provider:   opts.Provider,
			opIsDelete: opts.IsDelete,
		}
		sap := &serviceActionProvider{
			token:                  token,
			log:                    opts.Log,
			provider:               opts.Provider,
			prefObj:                opts.PrefObj,
			userID:                 opts.UserID,
			registry:               opts.Registry,
			skipPrintLogs:          opts.SkipPrintLogs,
			skipCrdAndOperator:     opts.SkipCRDAndOperator,
			upgradeExistingRelease: opts.UpgradeExistingRelease,
			ctxTokubeconfig:        ctxToconfig,
			accumulatedMsgs:        []string{},
			err:                    nil,
			eventsChannel:          opts.EventBroadcaster,
			opIsDelete:             opts.IsDelete,
			patternName:            strings.ToLower(opts.Pattern.Name),
		}
		fmt.Println("line 244 reached")

		chain := stages.CreateChain()
		chain.
			// Add(stages.Import(sip, sap)).
			Add(stages.Format()).
			// Add(stages.ServiceIdentifierAndMutator(sip, sap)).
			Add(stages.Filler(opts.SkipPrintLogs)).
			// Calling this stage `The Validation stage` is a bit deceiving considering
			// that the validation stage also formats the `data` (chain function parameter) that the
			// subsequent stages depend on.
			// We are skipping the `Validation` based on "verify" query paramerter
			Add(stages.Validator(sip, sap, opts.Validate)) // not required as client side RJSF validation is enough, but for mesheryctl client it's required
		if opts.DryRun {
			chain.Add(stages.DryRun(sip, sap))
		}
		if !opts.DryRun {
			chain.
				Add(stages.Provision(sip, sap, sap.log))
			// Removed Persist stage
			// Add(stages.Persist(sip, sap, sap.log))
		}
		chain.
			Add(func(data *stages.Data, err error, next stages.ChainStageNextFunction) {
				data.Lock.Lock()
				for k, v := range data.Other {
					var key string
					if strings.HasSuffix(k, stages.ProvisionSuffixKey) {
						key, _ = strings.CutSuffix(k, stages.ProvisionSuffixKey)
					}
					if k == stages.DryRunResponseKey {
						key = k
					}
					resp[key] = v
				}
				data.Lock.Unlock()
				sap.err = err
			}).
			Process(&stages.Data{
				Pattern:                       &opts.Pattern,
				Other:                         map[string]interface{}{},
				DeclartionToDefinitionMapping: make(map[uuid.UUID]component.ComponentDefinition),
			})
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
	skipCrdAndOperator     bool
	upgradeExistingRelease bool
	skipPrintLogs          bool
	accumulatedMsgs        []string
	err                    error
	eventsChannel          *models.Broadcast
	registry               *meshmodel.RegistryManager
	patternName            string
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
		sap.log.Error(err)
	}
	sap.err = err
}
func (sap *serviceActionProvider) Mutate(p *pattern.PatternFile) {
	//TODO: externalize these mutation rules with policies.
	//1. Enforce the deployment of CRDs before other resources
	for _, component := range p.Components {
		if component.Component.Kind == "CustomResourceDefinition" {
			for _, comp := range p.Components {
				if comp.Component.Kind != "CustomResourceDefinition" {
					dependsOnSlice, err := utils.Cast[[]string](comp.Metadata.AdditionalProperties["dependsOn"])
					if err != nil {
						err = errors.Wrapf(err, "Failed to cast 'dependsOn' to []string for component %s", comp.DisplayName)
						sap.log.Error(err)
						sap.Terminate(err)
					}
					dependsOnSlice = append(dependsOnSlice, comp.Id.String())
					comp.Metadata.AdditionalProperties["dependsOn"] = dependsOnSlice
				}
			}
		}
	}
}

// v1.StatusApplyConfiguration has deprecated, needed to find a different option to do this
// NOTE: Currently tied to kubernetes
// Returns ComponentName->ContextID->Response
func (sap *serviceActionProvider) DryRun(comps []*component.ComponentDefinition) (resp map[string]map[string]core.DryRunResponseWrapper, err error) {
	for _, cmp := range comps {
		for ctxID, kc := range sap.ctxTokubeconfig {
			cl, err := meshkube.New([]byte(kc))
			if err != nil {
				return resp, err
			}
			dResp, err := dryRunComponent(cl, cmp)
			if err != nil {
				return resp, err
			}
			if resp == nil {
				resp = make(map[string]map[string]core.DryRunResponseWrapper)
			}
			if resp[cmp.DisplayName] == nil {
				resp[cmp.DisplayName] = make(map[string]core.DryRunResponseWrapper)
			}
			resp[cmp.DisplayName][ctxID] = dResp
		}
	}
	return
}

func dryRunComponent(cl *meshkube.Client, cmd *component.ComponentDefinition) (core.DryRunResponseWrapper, error) {
	st, ok, err := k8s.DryRunHelper(cl, *cmd)
	dResp := core.DryRunResponseWrapper{Success: ok, Component: cmd}
	if ok {
		dResp.Component.Configuration = filterConfiguration(st)
	} else if err != nil {
		dResp.Error = &core.DryRunResponse{Status: err.Error()}
	} else {
		dResp.Error = parseDryRunFailure(st, cmd.DisplayName)
	}
	return dResp, nil
}

func filterConfiguration(configuration map[string]interface{}) map[string]interface{} {
	filteredConfiguration := make(map[string]interface{})
	for k, v := range configuration {
		if k != "apiVersion" && k != "kind" && k != "metadata" {
			filteredConfiguration[k] = v
		}
	}
	return filteredConfiguration
}

func parseDryRunFailure(settings map[string]interface{}, name string) *core.DryRunResponse {
	byt, err := json.Marshal(settings)
	if err != nil {
		return nil
	}
	var a metav1.Status
	err = json.Unmarshal(byt, &a)
	if err != nil {
		return nil
	}
	dResp := core.DryRunResponse{}
	if a.Status != "" {
		dResp.Status = a.Status
	}
	if a.Details != nil {
		dResp.Causes = make([]core.DryRunFailureCause, 0)
		for _, c := range a.Details.Causes {
			msg := ""
			field := ""
			typ := ""
			if c.Message != "" {
				msg = c.Message
			}
			if c.Field != "" {
				field = name + "." + getComponentFieldPathFromK8sFieldPath(c.Field)
			}
			if c.Type != "" {
				typ = string(c.Type)
			}
			failureCase := core.DryRunFailureCause{Message: msg, FieldPath: field, Type: typ}
			dResp.Causes = append(dResp.Causes, failureCase)
		}
	}
	return &dResp
}

func getComponentFieldPathFromK8sFieldPath(path string) (newpath string) {
	if strings.HasPrefix(path, "metadata.") {
		path = strings.TrimPrefix(path, "metadata.")
		paths := strings.Split(path, ".")
		if len(paths) != 0 {
			if paths[0] == "name" || paths[0] == "namespace" || paths[0] == "labels" || paths[0] == "annotations" {
				return paths[0]
			}
		}
		return
	}
	return fmt.Sprintf("%s.%s", "settings", path)
}

func (sap *serviceActionProvider) Provision(ccp stages.CompConfigPair) ([]patterns.DeploymentMessagePerContext, error) { // Marshal the component
	// jsonComp, err := json.Marshal(ccp.Component)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to serialize the data: %s", err)
	// }

	// Marshal the configuration
	// configuration attribute removed.
	// jsonConfig, err := json.Marshal(ccp.Configuration)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to serialize the data: %s", err)
	// }

	msgs := []patterns.DeploymentMessagePerContext{}
	for _, host := range ccp.Hosts {
		// Hack until adapters fix the concurrent client
		// creation issue: https://github.com/layer5io/meshery-adapter-library/issues/32
		time.Sleep(50 * time.Microsecond)
		sap.log.Debug("Execute operations on: ", host.Kind)

		_hostPort, ok := host.Metadata["port"]
		// Local call
		if !ok {
			//TODO: Accommodate internal calls to use context mapping with kubeconfig
			var kconfigs []string
			for _, v := range sap.ctxTokubeconfig {
				kconfigs = append(kconfigs, v)
			}
			resp, err := patterns.Process(
				kconfigs,
				[]component.ComponentDefinition{ccp.Component},
				sap.opIsDelete,
				sap.patternName,
				sap.eventsChannel,
				sap.userID,
				sap.provider,
				host,
				sap.skipCrdAndOperator,
				sap.upgradeExistingRelease,
			)
			return resp, err
		}

		hostName, err := utils.Cast[string](host.Metadata["host_name"])
		if err != nil {
			return nil, fmt.Errorf("error execute operation on %v: %v", host, err)
		}

		hostPort, err := utils.Cast[int](_hostPort)
		if err != nil {
			return nil, fmt.Errorf("error execute operation on %v: %v", host, err)
		}

		addr := hostName
		if hostPort != 0 {
			addr += ":" + strconv.Itoa(hostPort)
		}
		// Create mesh client
		mClient, err := meshes.CreateClient(
			context.TODO(),
			addr,
		)
		if err != nil {
			return nil, fmt.Errorf("error creating a mesh client: %v", err)
		}
		defer func() {
			_ = mClient.Close()
		}()

		// Else it is an  adapter call
		//TODO: Accommodate gRPC calls to use context mapping with kubeconfig
		var kconfigs []string
		for _, v := range sap.ctxTokubeconfig {
			kconfigs = append(kconfigs, v)
		}
		compStr, err := utils.Marshal(ccp.Component)
		if err != nil {
			err = errors.Wrapf(err, "error marshalling component \"%s\" of type : %s", ccp.Component.DisplayName, ccp.Component.Component.Kind)
			return nil, err
		}
		resp, err := mClient.MClient.Provision(context.TODO(), &meshes.ProvisionRequest{
			Username:     sap.userID,
			DeleteOp:     sap.opIsDelete,
			KubeConfigs:  kconfigs,
			Declarations: []string{compStr},
		})
		sucess := err == nil
		msgs = append(msgs, patterns.DeploymentMessagePerContext{
			SystemName: hostName,
			Location:   fmt.Sprintf("%s:%s", hostName, strconv.Itoa(hostPort)),
			Summary: []patterns.DeploymentMessagePerComp{
				{
					Kind:       ccp.Component.Component.Kind,
					Model:      ccp.Component.Model.Name,
					CompName:   ccp.Component.DisplayName,
					DesignName: sap.patternName,
					Success:    sucess,
					Message:    resp.GetMessage(),
					Error:      err,
				},
			},
		})
	}

	return msgs, nil
}

// func (sap *serviceActionProvider) Persist(name string, svc core.Service, isUpdate bool) error {
// 	if !sap.opIsDelete {
// 		if isUpdate {
// 			// Do nothing
// 			return nil
// 		}

// 		_, err := sap.provider.SaveMesheryPatternResource(
// 			sap.token,
// 			&models.PatternResource{
// 				ID:        svc.ID,
// 				Name:      name,
// 				Namespace: svc.Namespace,
// 				Type:      svc.Type,
// 				OAMType:   "workload",
// 			},
// 		)

// 		return err
// 	}

// 	return sap.provider.DeleteMesheryPatternResource(
// 		sap.token,
// 		svc.ID.String(),
// 	)
// }
