package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	cueJson "cuelang.org/go/encoding/json"
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/helpers"

	"github.com/layer5io/meshery/server/models"

	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/models/meshmodel/registry"
	regv1beta1 "github.com/layer5io/meshkit/models/meshmodel/registry/v1beta1"

	"github.com/layer5io/meshery/server/helpers/utils"
	mesheryutils "github.com/layer5io/meshery/server/helpers/utils"
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"

	"github.com/meshery/schemas/models/v1beta1"

	_component "github.com/layer5io/meshkit/utils/component"
	"github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshkit/utils/manifests"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type crd struct {
	Items []crdhelper `json:"items"`
}
type crdhelper struct {
	Spec spec `json:"spec"`
}
type spec struct {
	Names names `json:"names"`
}
type names struct {
	Kind string `json:"kind"`
}

func RegisterK8sMeshModelComponents(provider *models.Provider, _ context.Context, config []byte, ctxID string, connectionID string, userID string, mesheryInstanceID uuid.UUID, reg *registry.RegistryManager, ec *models.Broadcast, log logger.Handler, ctxName string) (err error) {
	connectionUUID := uuid.FromStringOrNil(connectionID)
	userUUID := uuid.FromStringOrNil(userID)

	man, err := GetK8sMeshModelComponents(config)
	eventMetadata := make(map[string]interface{}, 0)
	ctx := models.K8sContextsFromKubeconfig(*provider, userID, ec, config, &mesheryInstanceID, eventMetadata, log)
	if err != nil {
		return ErrCreatingKubernetesComponents(err, ctxID)
	}
	if man == nil {
		return ErrCreatingKubernetesComponents(errors.New("generated components are nil"), ctxID)
	}
	k8sContext := map[string]interface{}{}
	if len(ctx) > 0 {
		k8sContext, _ = utils.MarshalAndUnmarshal[models.K8sContext, map[string]interface{}](*ctx[0])
	}
	count := 0
	for _, c := range man {
		var isRegistranError bool
		var isModelError bool
		writeK8sMetadata(&c, reg)
		if models.K8sMeshModelMetadata.Capabilities != nil {
			c.Capabilities = models.K8sMeshModelMetadata.Capabilities
		}
		isRegistranError, isModelError, err = reg.RegisterEntity(connection.Connection{
			Kind:     "kubernetes",
			Type:     "registry",
			Metadata: k8sContext,
		}, &c)
		helpers.HandleError(connection.Connection{
			Kind: "kubernetes"}, &c, err, isModelError, isRegistranError)
		count++
	}
	err = helpers.WriteLogsToFiles()
	if err != nil {
		return err
	}
	event := events.NewEvent().ActedUpon(connectionUUID).WithCategory("kubernetes_components").WithAction("registration").FromSystem(mesheryInstanceID).FromUser(userUUID).WithSeverity(events.Informational).WithDescription(fmt.Sprintf("%d Kubernetes components registered for %s", count, ctxName)).WithMetadata(map[string]interface{}{
		"doc": "https://docs.meshery.io/tasks/lifecycle-management",
	}).Build()
	_, err = helpers.FailedEventCompute("Kubernetes", mesheryInstanceID, provider, userID, ec)
	if err != nil {
		return err
	}
	//if want to log we can use the above function in future to log the error in terminal.

	_ = (*provider).PersistEvent(event)
	ec.Publish(userUUID, event)
	return
}

func writeK8sMetadata(comp *component.ComponentDefinition, reg *registry.RegistryManager) {
	ent, _, _, _ := reg.GetEntities(&regv1beta1.ComponentFilter{
		Name:       comp.Component.Kind,
		APIVersion: comp.Component.Version,
	})
	// If component was not available in the registry, then use the generic model level metadata
	if len(ent) == 0 {
		comp.Styles = &models.K8sMeshModelMetadata.Styles
		mesheryutils.WriteSVGsOnFileSystem(comp)
	} else {
		existingComp, ok := ent[0].(*component.ComponentDefinition)
		if !ok {
			comp.Styles = &models.K8sMeshModelMetadata.Styles
			return
		}
		comp.Metadata = existingComp.Metadata
		comp.Styles = existingComp.Styles
		comp.Model = existingComp.Model
	}
}

type OpenAPIV3Response struct {
	Paths map[string]Entry `json:"paths"`
}

type Entry struct {
	URL string `json:"serverRelativeURL"`
}

func mergeAllAPIResults(content []byte, cli *kubernetes.Client) [][]byte {
	var res OpenAPIV3Response
	_ = json.Unmarshal(content, &res)
	m := make([][]byte, 0)
	for k, path := range res.Paths {
		if !strings.HasPrefix(k, "api") {
			continue
		}
		req := cli.KubeClient.RESTClient().Get().RequestURI(path.URL)
		res := req.Do(context.Background())
		content, err := res.Raw()
		if err != nil {
			return nil
		}
		m = append(m, content)
	}
	return m
}

// move to meshmodel
func GetK8sMeshModelComponents(kubeconfig []byte) ([]component.ComponentDefinition, error) {
	cli, err := kubernetes.New(kubeconfig)
	if err != nil {
		return nil, core.ErrGetK8sComponents(err)
	}
	req := cli.KubeClient.RESTClient().Get().RequestURI("/openapi/v3")
	k8sversion, err := cli.KubeClient.ServerVersion()
	if err != nil {
		return nil, core.ErrGetK8sComponents(err)
	}
	var customResources = make(map[string]bool)
	crdresult, err := cli.KubeClient.RESTClient().Get().RequestURI("/apis/apiextensions.k8s.io/v1/customresourcedefinitions").Do(context.Background()).Raw()
	if err != nil {
		return nil, core.ErrGetK8sComponents(err)
	}

	var xcrd crd
	err = json.Unmarshal(crdresult, &xcrd)
	if err != nil {
		return nil, core.ErrGetK8sComponents(err)
	}
	for _, item := range xcrd.Items {
		customResources[item.Spec.Names.Kind] = true
	}
	res := req.Do(context.Background())
	content, err := res.Raw()
	if err != nil {
		return nil, core.ErrGetK8sComponents(err)
	}
	contents := mergeAllAPIResults(content, cli)
	apiResources, err := getAPIRes(cli)
	if err != nil {
		return nil, core.ErrGetK8sComponents(err)
	}

	var arrAPIResources []string
	kindToNamespace := make(map[string]bool)
	for res, api := range apiResources {
		kindToNamespace[api.Kind] = api.Namespaced
		arrAPIResources = append(arrAPIResources, res)
	}
	var crds []crdResponse
	for _, content := range contents {
		crds = append(crds, getCRDsFromManifest(string(content), arrAPIResources)...)
	}
	components := make([]component.ComponentDefinition, 0)
	for _, crd := range crds {
		m := make(map[string]interface{})
		m[customResourceKey] = customResources[crd.kind]
		m[namespacedKey] = kindToNamespace[crd.kind]
		compMetadata := component.ComponentDefinition_Metadata{
			AdditionalProperties: m,
		}
		apiVersion := crd.apiVersion
		c := component.ComponentDefinition{
			SchemaVersion: v1beta1.ComponentSchemaVersion,
			Version:       "v1.0.0",

			Format: component.JSON,
			Component: component.Component{
				Kind:    crd.kind,
				Version: apiVersion,
				Schema:  crd.schema,
			},
			Metadata:    compMetadata,
			DisplayName: manifests.FormatToReadableString(crd.kind),
			Model: model.ModelDefinition{
				SchemaVersion: v1beta1.ModelSchemaVersion,
				Version:       "v1.0.0",

				Model: model.Model{
					Version: k8sversion.String(),
				},
				Name:        "kubernetes",
				DisplayName: "Kubernetes",
				Category: category.CategoryDefinition{
					Name: "Orchestration & Management",
				},
			},
		}
		components = append(components, c)
	}
	return components, nil
}

const customResourceKey = "isCustomResource"
const namespacedKey = "isNamespaced"

func getResolvedManifest(manifest string) (string, error) {
	cuectx := cuecontext.New()
	cueParsedManExpr, err := cueJson.Extract("", []byte(manifest))
	parsedManifest := cuectx.BuildExpr(cueParsedManExpr)
	definitions := parsedManifest.LookupPath(cue.ParsePath("components.schemas"))
	if err != nil {
		return "", err
	}
	resol := manifests.ResolveOpenApiRefs{}
	cache := make(map[string][]byte)
	resolved, err := resol.ResolveReferences([]byte(manifest), definitions, cache)
	if err != nil {
		return "", err
	}
	manifest = string(resolved)
	return manifest, nil
}

type crdResponse struct {
	name       string
	kind       string
	apiVersion string
	schema     string
}

func getCRDsFromManifest(manifest string, arrAPIResources []string) []crdResponse {
	var err error
	res := make([]crdResponse, 0)
	manifest, err = getResolvedManifest(manifest)
	if err != nil {
		fmt.Printf("%v", err)
		return nil
	}
	cuectx := cuecontext.New()
	cueParsedManExpr, err := cueJson.Extract("", []byte(manifest))
	parsedManifest := cuectx.BuildExpr(cueParsedManExpr)
	definitions := parsedManifest.LookupPath(cue.ParsePath("components.schemas"))
	if err != nil {
		return nil
	}
	for _, name := range arrAPIResources {
		resource := strings.ToLower(name)
		fields, err := definitions.Fields()
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}
		for fields.Next() {
			fieldVal := fields.Value()
			kindCue := fieldVal.LookupPath(cue.ParsePath(`"x-kubernetes-group-version-kind"[0].kind`))
			if kindCue.Err() != nil {
				continue
			}
			kind, err := kindCue.String()
			kind = strings.ToLower(kind)
			if err != nil {
				fmt.Printf("%v", err)
				continue
			}
			if kind == resource {
				crd, err := fieldVal.MarshalJSON()
				if err != nil {
					fmt.Printf("%v", err)
					continue
				}
				versionCue := fieldVal.LookupPath(cue.ParsePath(`"x-kubernetes-group-version-kind"[0].version`))
				groupCue := fieldVal.LookupPath(cue.ParsePath(`"x-kubernetes-group-version-kind"[0].group`))
				apiVersion, _ := versionCue.String()
				if g, _ := groupCue.String(); g != "" {
					apiVersion = g + "/" + apiVersion
				}
				modified := make(map[string]interface{}) //Remove the given fields which is either not required by End user (like status) or is prefilled by system (like apiVersion, kind and metadata)
				err = json.Unmarshal(crd, &modified)
				if err != nil {
					fmt.Printf("%v", err)
					continue
				}

				modifiedProps, err := _component.UpdateProperties(fieldVal, cue.ParsePath("properties.spec"), apiVersion)
				if err == nil {
					modified = modifiedProps
				}

				_component.DeleteFields(modified)
				crd, err = json.Marshal(modified)
				if err != nil {
					fmt.Printf("%v", err)
					continue
				}
				res = append(res, crdResponse{
					name:       resource,
					kind:       name,
					schema:     string(crd),
					apiVersion: apiVersion, //add apiVersion
				})
			}
		}
	}
	return res
}

// TODO: To be moved in meshkit
// getAPIRes gets all the available api resources from kube-api server. It is equivalent to the output of `kubectl api-resources`
// Returns a map of api resources with key as api-resource kind and value as api-resource object
func getAPIRes(cli *kubernetes.Client) (map[string]v1.APIResource, error) {
	var apiRes = make(map[string]v1.APIResource)
	lists, err := cli.KubeClient.DiscoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, err
	}
	for _, list := range lists {
		for _, name := range list.APIResources {
			apiRes[name.Kind] = name
		}
	}
	return apiRes, nil
}

// TODO: To be moved in meshkit
// func getGroupsFromResource(cli *kubernetes.Client) (hgv map[kind][]groupversion, err error) {
// 	hgv = make(map[kind][]groupversion)
// 	var gl v1.APIGroupList
// 	gs, err := cli.KubeClient.RESTClient().Get().RequestURI("/apis").Do(context.Background()).Raw()
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = json.Unmarshal(gs, &gl)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, g := range gl.Groups {
// 		groupName := g.Name
// 		var apig v1.APIGroup
// 		apigbytes, err := cli.KubeClient.RESTClient().Get().RequestURI("/apis/" + groupName).Do(context.Background()).Raw()
// 		if err != nil {
// 			return nil, err
// 		}
// 		err = json.Unmarshal(apigbytes, &apig)
// 		if err != nil {
// 			return nil, err
// 		}
// 		for _, v := range apig.Versions {
// 			apiRes, err := cli.KubeClient.DiscoveryClient.ServerResourcesForGroupVersion(v.GroupVersion)
// 			if err != nil {
// 				return nil, err
// 			}
// 			if err != nil {
// 				return nil, err
// 			}
// 			for _, res := range apiRes.APIResources {
// 				if v.GroupVersion != "" {
// 					hgv[kind(res.Kind)] = append(hgv[kind(res.Kind)], groupversion(v.GroupVersion))
// 				} else {
// 					hgv[kind(res.Kind)] = append(hgv[kind(res.Kind)], groupversion(v.Version))
// 				}
// 			}
// 		}
// 		apiRes, err := cli.KubeClient.DiscoveryClient.ServerResourcesForGroupVersion("v1")
// 		if err != nil {
// 			return nil, err
// 		}
// 		if err != nil {
// 			return nil, err
// 		}
// 		for _, res := range apiRes.APIResources {
// 			hgv[kind(res.Kind)] = append(hgv[kind(res.Kind)], groupversion("v1"))
// 		}
// 	}
// 	return
// }
