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
	mutil "github.com/layer5io/meshery/server/helpers/utils"

	"github.com/layer5io/meshery/server/models"

	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshkit/utils"

	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshkit/models/events"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	oamcore "github.com/layer5io/meshkit/models/oam/core/v1alpha1"

	"github.com/layer5io/meshkit/utils/component"
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

func RegisterK8sMeshModelComponents(provider *models.Provider, _ context.Context, config []byte, ctxID string, connectionID string, userID string, mesheryInstanceID uuid.UUID, reg *meshmodel.RegistryManager, ec *models.Broadcast, ctxName string) (err error) {
	connectionUUID := uuid.FromStringOrNil(connectionID)
	userUUID := uuid.FromStringOrNil(userID)

	man, err := GetK8sMeshModelComponents(config)
	if err != nil {
		return ErrCreatingKubernetesComponents(err, ctxID)
	}
	if man == nil {
		return ErrCreatingKubernetesComponents(errors.New("generated components are nil"), ctxID)
	}
	count := 0
	for _, c := range man {
		writeK8sMetadata(&c, reg)
		err = reg.RegisterEntity(meshmodel.Host{
			Hostname: "kubernetes",
			Metadata: ctxID,
		}, c)
		count++
	}
	event := events.NewEvent().ActedUpon(connectionUUID).WithCategory("kubernetes_components").WithAction("registration").FromSystem(mesheryInstanceID).FromUser(userUUID).WithSeverity(events.Informational).WithDescription(fmt.Sprintf("%d Kubernetes components registered for %s", count, ctxName)).WithMetadata(map[string]interface{}{
		"doc": "https://docs.meshery.io/tasks/lifecycle-management",
	}).Build()

	_ = (*provider).PersistEvent(event)
	ec.Publish(userUUID, event)
	return
}

func writeK8sMetadata(comp *v1alpha1.ComponentDefinition, reg *meshmodel.RegistryManager) {
	ent, _, _ := reg.GetEntities(&v1alpha1.ComponentFilter{
		Name:       comp.Kind,
		APIVersion: comp.APIVersion,
	})
	//If component was not available in the registry, then use the generic model level metadata
	if len(ent) == 0 {
		comp.Metadata = utils.MergeMaps(comp.Metadata, models.K8sMeshModelMetadata)
		mutil.WriteSVGsOnFileSystem(comp)
	} else {
		existingComp, ok := ent[0].(v1alpha1.ComponentDefinition)
		if !ok {
			comp.Metadata = utils.MergeMaps(comp.Metadata, models.K8sMeshModelMetadata)
			return
		}
		comp.Metadata = utils.MergeMaps(comp.Metadata, existingComp.Metadata)
		comp.Model = existingComp.Model
	}
}

func RegisterMeshmodelComponentsForCRDS(reg meshmodel.RegistryManager, k8sYaml []byte, contextID string, version string) {
	//TODO: Replace GenerateComponents in meshkit to natively produce MeshModel components to avoid any interconversion
	comp, err := manifests.GenerateComponents(context.Background(), string(k8sYaml), manifests.K8s, manifests.Config{
		CrdFilter: manifests.NewCueCrdFilter(manifests.ExtractorPaths{
			NamePath:    "spec.names.kind",
			IdPath:      "spec.names.kind",
			VersionPath: "spec.versions[0].name",
			GroupPath:   "spec.group",
			SpecPath:    "spec.versions[0].schema.openAPIV3Schema.properties.spec"}, false),
		ExtractCrds: func(manifest string) []string {
			crds := strings.Split(manifest, "---")
			return crds
		},
	})
	if err != nil {
		fmt.Println("err: ", err.Error())
		return
	}
	for i, schema := range comp.Schemas {
		var def oamcore.WorkloadDefinition
		err := json.Unmarshal([]byte(comp.Definitions[i]), &def)
		if err != nil {
			fmt.Println("err here: ", err.Error())
			return
		}
		_ = reg.RegisterEntity(meshmodel.Host{
			Hostname: meshmodel.Kubernetes{}.String(),
			Metadata: contextID,
		}, v1alpha1.ComponentDefinition{
			Schema: schema,
			TypeMeta: v1alpha1.TypeMeta{
				APIVersion: def.Spec.Metadata["k8sKind"],
				Kind:       def.Spec.Metadata["k8sAPIVersion"],
			},
			Model: v1alpha1.Model{
				Name:        "kubernetes",
				Version:     version,
				DisplayName: "Kubernetes",
				Category: v1alpha1.Category{
					Name: "Orchestration & Management",
				},
			},
		})
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
func GetK8sMeshModelComponents(kubeconfig []byte) ([]v1alpha1.ComponentDefinition, error) {
	cli, err := kubernetes.New(kubeconfig)
	if err != nil {
		return nil, core.ErrGetK8sComponents(err)
	}
	req := cli.KubeClient.RESTClient().Get().RequestURI("/openapi/v3")
	k8version, err := cli.KubeClient.ServerVersion()
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
	components := make([]v1alpha1.ComponentDefinition, 0)
	for _, crd := range crds {
		m := make(map[string]interface{})
		m[customResourceKey] = customResources[crd.kind]
		m[namespacedKey] = kindToNamespace[crd.kind]
		apiVersion := crd.apiVersion
		c := v1alpha1.ComponentDefinition{
			Format: v1alpha1.JSON,
			Schema: crd.schema,
			TypeMeta: v1alpha1.TypeMeta{
				Kind:       crd.kind,
				APIVersion: apiVersion,
			},
			Metadata:    m,
			DisplayName: manifests.FormatToReadableString(crd.kind),
			Model: v1alpha1.Model{
				Version:     k8version.String(),
				Name:        "kubernetes",
				DisplayName: "Kubernetes",
				Category: v1alpha1.Category{
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

				modifiedProps, err := component.UpdateProperties(fieldVal, cue.ParsePath("properties.spec"), apiVersion)
				if err == nil {
					modified = modifiedProps
				}

				component.DeleteFields(modified)
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
