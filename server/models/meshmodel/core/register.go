package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	cueJson "cuelang.org/go/encoding/json"
	"github.com/layer5io/meshkit/models/meshmodel"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	oamcore "github.com/layer5io/meshkit/models/oam/core/v1alpha1"
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
			Hostname:  "kubernetes",
			ContextID: contextID,
		}, v1alpha1.ComponentDefinition{
			Schema: schema,
			TypeMeta: v1alpha1.TypeMeta{
				APIVersion: def.Spec.Metadata["k8sKind"],
				Kind:       def.Spec.Metadata["k8sAPIVersion"],
			},
			Model: v1alpha1.Models{
				Name:    "kubernetes",
				Version: version,
			},
		})
	}
}

// move to meshmodel
func GetK8sMeshModelComponents(ctx context.Context, kubeconfig []byte) ([]v1alpha1.ComponentDefinition, error) {
	cli, err := kubernetes.New(kubeconfig)
	if err != nil {
		return nil, ErrGetK8sComponents(err)
	}
	req := cli.KubeClient.RESTClient().Get().RequestURI("/openapi/v2")
	k8version, err := cli.KubeClient.ServerVersion()
	if err != nil {
		return nil, ErrGetK8sComponents(err)
	}
	var customResources = make(map[string]bool)
	crdresult, err := cli.KubeClient.RESTClient().Get().RequestURI("/apis/apiextensions.k8s.io/v1/customresourcedefinitions").Do(context.Background()).Raw()
	if err != nil {
		return nil, ErrGetK8sComponents(err)
	}

	var xcrd crd
	err = json.Unmarshal(crdresult, &xcrd)
	if err != nil {
		return nil, ErrGetK8sComponents(err)
	}
	for _, item := range xcrd.Items {
		customResources[item.Spec.Names.Kind] = true
	}
	res := req.Do(context.Background())
	content, err := res.Raw()
	if err != nil {
		return nil, ErrGetK8sComponents(err)
	}
	apiResources, err := getAPIRes(cli)
	if err != nil {
		return nil, ErrGetK8sComponents(err)
	}

	var arrAPIResources []string
	for res := range apiResources {
		arrAPIResources = append(arrAPIResources, res)
	}
	groups, err := getGroupsFromResource(cli) //change this
	if err != nil {
		return nil, err
	}
	manifest := string(content)
	crds, metadata := getCRDsFromManifest(manifest, arrAPIResources)
	components := make([]v1alpha1.ComponentDefinition, 1)
	for name, crd := range crds {
		m := make(map[string]interface{})
		m[customResourceKey] = customResources[metadata[name].K8sKind]
		c := v1alpha1.ComponentDefinition{
			Format: v1alpha1.JSON,
			Schema: crd,
			TypeMeta: v1alpha1.TypeMeta{
				Kind:       metadata[name].K8sKind,
				APIVersion: string(groups[kind(metadata[name].K8sKind)]),
			},
			Metadata:    m,
			DisplayName: manifests.FormatToReadableString(metadata[name].K8sKind),
			Model: v1alpha1.Models{
				Version: k8version.String(),
				Name:    "kubernetes",
			},
		}
		components = append(components, c)
	}
	return components, nil
}

const customResourceKey = "isCustomResource"

func getResolvedManifest(manifest string) (string, error) {
	cuectx := cuecontext.New()
	cueParsedManExpr, err := cueJson.Extract("", []byte(manifest))
	parsedManifest := cuectx.BuildExpr(cueParsedManExpr)
	definitions := parsedManifest.LookupPath(cue.ParsePath("definitions"))
	if err != nil {
		fmt.Printf("%v", err)
		return "", nil
	}
	resol := manifests.ResolveOpenApiRefs{}
	resolved, err := resol.ResolveReferences([]byte(manifest), definitions)
	if err != nil {
		panic(err)
	}
	manifest = string(resolved)
	return manifest, nil
}
func getCRDsFromManifest(manifest string, arrAPIResources []string) (map[string]string, map[string]k8sMetadata) {
	var err error
	manifest, err = getResolvedManifest(manifest)
	if err != nil {
		fmt.Printf("%v", err)
		return nil, nil
	}
	resourceToCRD := make(map[string]string)
	resourceToMetadata := make(map[string]k8sMetadata, 0)
	cuectx := cuecontext.New()
	cueParsedManExpr, err := cueJson.Extract("", []byte(manifest))
	parsedManifest := cuectx.BuildExpr(cueParsedManExpr)
	definitions := parsedManifest.LookupPath(cue.ParsePath("definitions"))
	if err != nil {
		fmt.Printf("%v", err)
		return nil, nil
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
				resourceToMetadata[resource] = k8sMetadata{
					K8sKind: name,
				}
				resourceToCRD[resource] = string(crd)
				// resourceToName[resource] = manifests.FormatToReadableString(name)
			}
		}
	}
	return resourceToCRD, resourceToMetadata
}

type k8sMetadata struct {
	K8sAPIVersion string `json:"k8sAPIVersion"`
	K8sKind       string `json:"k8sKind"`
	Namespaced    bool   `json:"namespaced"`
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

type helperGV struct {
	group   string
	version string
	isCRD   bool
}
type kind string
type groupversion string

// TODO: To be moved in meshkit
func getGroupsFromResource(cli *kubernetes.Client) (hgv map[kind]groupversion, err error) {
	hgv = make(map[kind]groupversion)
	var gl v1.APIGroupList
	gs, err := cli.KubeClient.RESTClient().Get().RequestURI("/apis").Do(context.Background()).Raw()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(gs, &gl)
	if err != nil {
		return nil, err
	}

	for _, g := range gl.Groups {
		groupName := g.Name
		var apig v1.APIGroup
		apigbytes, err := cli.KubeClient.RESTClient().Get().RequestURI("/apis/" + groupName).Do(context.Background()).Raw()
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(apigbytes, &apig)
		if err != nil {
			return nil, err
		}
		apiRes, err := cli.KubeClient.DiscoveryClient.ServerResourcesForGroupVersion(apig.PreferredVersion.GroupVersion)
		if err != nil {
			return nil, err
		}
		if len(g.Versions) != 0 {
			for _, res := range apiRes.APIResources {
				hgv[kind(res.Kind)] = groupversion(g.Versions[0].GroupVersion)
				break
			}
		}
	}
	return
}
