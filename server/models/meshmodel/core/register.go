package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/google/uuid"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	cueJson "cuelang.org/go/encoding/json"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshkit/utils/manifests"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type capability struct {
	ID string `json:"id,omitempty"`
	// Host is the address of the service registering the capability
	Host string `json:"host,omitempty"`
}

// SetID sets the ID of the capability
func (cap *capability) SetID(id string) {
	cap.ID = id
}

// GetID returns the ID of the capability
func (cap *capability) GetID() string {
	return cap.ID
}

type ComponentCapability struct {
	// gorm.Model
	v1alpha1.Component
	capability
}

type crd struct {
	Items []crdhelper `json:"items"`
}
type crdhelper struct {
	Metadata map[string]interface{} `json:"metadata"`
}

// move to meshmodel
func GetK8sMeshModelComponents(ctx context.Context, kubeconfig []byte) ([]meshmodel.Component, error) {
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
		customResources[item.Metadata["name"].(string)] = true
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
	// groups, err := getGroupsFromResource(cli) //change this
	if err != nil {
		return nil, err
	}
	manifest := string(content)
	fmt.Println("manifest is: ", manifest[0:2000])
	crds := getCRDsFromManifest(manifest, arrAPIResources)
	components := make([]meshmodel.Component, 1)
	for name, crd := range crds {
		c := meshmodel.NewComponent()
		c.Spec = crd
		c.Metadata["k8sVersion"] = k8version.String()
		c.Metadata[customResourceKey] = false
		c.Metadata["name"] = name
		// for cr := range customResources {
		// 	if groups[c.GetMetadataValue("name").(string)][cr] {
		// 		c.SetMetadata(customResourceKey, true)
		// 		break
		// 	}
		// }
		fmt.Println("generated: ", c)
		components = append(components, c)
	}
	return components, nil
}

const customResourceKey = "isCustomResource"

func RegisterComponentCapability(dbHandler *database.Handler, mc []meshmodel.Component, host string) (err error) {
	fmt.Println("here")
	for _, c := range mc {
		var cc meshmodel.ComponentCapability
		cc.Component = c
		cc.Host = host
		var cdb meshmodel.ComponentCapabilityDB
		data, _ := json.Marshal(c)
		cdb = meshmodel.ComponentCapabilityDBFromCC(cc)
		cdb.ID = uuid.NewMD5(uuid.New(), data)
		err = dbHandler.DB.Create(&cdb).Error
		if err != nil {
			fmt.Println("wtf ", err.Error())
		}
	}

	return
}

// RegisterComponent will register a component definition into the database
func RegisterComponent(dbHandler *database.Handler, data []byte, host string) (err error) {
	var component meshmodel.ComponentCapability
	if err = json.Unmarshal(data, &component); err != nil {
		return
	}
	component.Host = host
	cdb := meshmodel.ComponentCapabilityDBFromCC(component)
	cdb.ID = uuid.NewMD5(uuid.New(), data)
	err = dbHandler.DB.Create(&cdb).Error
	fmt.Println("err: ", err.Error())
	return
}

// GetComponents return all of the components
func GetComponents(dbHandler *database.Handler) (caps []meshmodel.ComponentCapability) {
	var capsdb []meshmodel.ComponentCapabilityDB
	err := dbHandler.DB.Find(&capsdb).Error
	if err != nil {
		fmt.Println("error: ", err.Error())
	}
	for _, cdb := range capsdb {
		caps = append(caps,
			meshmodel.ComponentCapabilityFromCCDB(cdb),
		)
	}
	return
}

// GetComponents return all of the components
func GetComponentsByName(dbHandler *database.Handler, name string) (caps []meshmodel.ComponentCapability) {
	var capsdb []meshmodel.ComponentCapabilityDB
	err := dbHandler.DB.Find(&capsdb).Error
	if err != nil {
		fmt.Println("error: ", err.Error())
	}
	for _, cdb := range capsdb {
		m := make(map[string]interface{})
		json.Unmarshal(cdb.Metadata, &m)
		if m["name"] == name {
			caps = append(caps,
				meshmodel.ComponentCapabilityFromCCDB(cdb),
			)
		}
	}
	return
}

type ComponentFile struct {
	Name       string                `json:"name"`
	Components []ComponentCapability `json:"components"`
}

func StreamComponents(f io.Reader) (chan ComponentCapability, chan bool) {
	m := make(chan ComponentCapability, 10)
	sync := make(chan bool)
	go func(m chan ComponentCapability) {
		defer func() {
			close(m)
			sync <- true
		}()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println("bruh", err.Error())
			return
		}
		var lol []ComponentFile
		err = yaml.Unmarshal(b, &lol)
		if err != nil {
			fmt.Println("bruhL: ", err.Error())
			return
		}
		for _, compfile := range lol {
			for _, comp := range compfile.Components {
				m <- comp
			}
		}

	}(m)
	return m, sync
}
func SaveComponent(dbHandler *database.Handler, cc chan ComponentCapability, sync chan bool) {
	for {
		select {
		case <-sync:
			break
		default:
			select {
			case c := <-cc:
				err := dbHandler.DB.Create(&c).Error
				if err != nil {
					fmt.Println("err: ", err.Error())
				}
			}
		}
	}
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
// Return a set of resource.groups(with resourcename appended) for a given kind //Example: {"CSIDriver":{"csidrivers.storage.k8s.io"}}
func getGroupsFromResource(cli *kubernetes.Client) (gr map[string]map[string]bool, err error) {
	gr = make(map[string]map[string]bool)

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
		for _, res := range apiRes.APIResources {
			if gr[res.Kind] == nil {
				gr[res.Kind] = make(map[string]bool)
			}
			gr[res.Kind][res.Name+"."+groupName] = true
		}
	}
	return
}
func getResolvedManifest(manifest string) (string, error) {
	fmt.Println("oyr yaha")
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
func getCRDsFromManifest(manifest string, arrApiResources []string) map[string]string {
	fmt.Println("oyr yaha bhai")
	var err error
	manifest, err = getResolvedManifest(manifest)
	if err != nil {
		fmt.Println("fuck: ", err.Error())
	}
	crds := make(map[string]string, 0)
	cuectx := cuecontext.New()
	cueParsedManExpr, err := cueJson.Extract("", []byte(manifest))
	parsedManifest := cuectx.BuildExpr(cueParsedManExpr)
	definitions := parsedManifest.LookupPath(cue.ParsePath("definitions"))
	if err != nil {
		fmt.Printf("%v", err)
		return nil
	}
	for _, resource := range arrApiResources {
		resource = strings.ToLower(resource)
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
				crds[resource] = string(crd)
			}
		}
	}
	return crds
}
