package istio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/layer5io/meshery/meshes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
)

// CreateIstioClient is a function to create a Istio mesh client
func CreateIstioClient(ctx context.Context) (meshes.MeshClient, error) {
	// iClient, err := newClient(nil, "")
	// if err != nil {
	// 	err = errors.Wrapf(err, "unable to create a new istio client")
	// 	logrus.Error(err)
	// 	return nil, err
	// }
	// return iClient, nil
	return CreateIstioClientWithK8SConfig(ctx, nil, "")
}

func CreateIstioClientWithK8SConfig(ctx context.Context, k8sConfig []byte, contextName string) (meshes.MeshClient, error) {
	iClient, err := newClient(k8sConfig, contextName)
	if err != nil {
		err = errors.Wrapf(err, "unable to create a new istio client")
		logrus.Error(err)
		return nil, err
	}
	return iClient, nil
}

func (iClient *IstioClient) deleteAllCreatedResources(ctx context.Context, namespace string) {
	resourceNames := []string{"productpage", "ratings", "reviews", "details"}
	for _, rs := range resourceNames {
		iClient.deleteResource(ctx, virtualServices, namespace, rs)
	}
}

func (iClient *IstioClient) deleteResource(ctx context.Context, overallType, namespace, resName string) error {
	newRes := iClient.istioNetworkingApi.Delete().
		Namespace(namespace).
		Resource(overallType).SubResource(resName).Do()
	_, err := newRes.Get()
	if err != nil {
		err = errors.Wrapf(err, "unable to delete the requested resource")
		logrus.Error(err)
		return err
	}
	logrus.Infof("Deleted Resource of type: %s and name: %s", overallType, resName)
	return nil
}

// MeshName just returns the name of the mesh the client is representing
func (iClient *IstioClient) MeshName() string {
	return "Istio"
}

func (iClient *IstioClient) applyRulePayload(ctx context.Context, namespace string, newVSBytes []byte) error {
	vs := &VirtualService{}
	err := yaml.Unmarshal(newVSBytes, vs)
	if err != nil {
		err = errors.Wrapf(err, "unable to unmarshal yaml")
		logrus.Error(err)
		return err
	}
	vs.Kind = virtualservice
	vs.APIVersion = istioNetworkingGroupVersion.String()
	newVSBytesJ, err := json.Marshal(vs)
	if err != nil {
		err = errors.Wrapf(err, "unable to marshal virtual service map")
		logrus.Error(err)
		return err
	}

	newVSRes := iClient.istioNetworkingApi.Post().SetHeader("content-type", runtime.ContentTypeJSON).
		Namespace(namespace).
		Resource(virtualServices).Body(newVSBytesJ).Do()
	_, err = newVSRes.Get()
	if err != nil {
		newVSRes = iClient.istioNetworkingApi.Get().SetHeader("content-type", runtime.ContentTypeJSON).
			Namespace(namespace).Name(vs.ObjectMeta.Name).
			Resource(virtualServices).Do()
		newVSResInst, err := newVSRes.Get()
		if err != nil {
			err = errors.Wrapf(err, "unable to get the virtual service instance")
			logrus.Error(err)
			return err
		}
		vs1, _ := newVSResInst.(*VirtualService)
		vs.ObjectMeta.ResourceVersion = vs1.ObjectMeta.ResourceVersion
		newVSBytesJ, err := json.Marshal(vs)
		if err != nil {
			err = errors.Wrapf(err, "unable to marshal virtual service map")
			logrus.Error(err)
			return err
		}

		_, err = iClient.istioNetworkingApi.Put().SetHeader("content-type", runtime.ContentTypeJSON).
			Namespace(namespace).Name(vs.ObjectMeta.Name).
			Resource(virtualServices).Body(newVSBytesJ).Do().Get()
		if err != nil {
			err = errors.Wrapf(err, "unable to get the virtual service instance from result")
			logrus.Error(err)
			return err
		}
	}
	return nil
}

// ApplyRule is a method invoked to apply a particular operation on the mesh in a namespace
func (iClient *IstioClient) ApplyRule(ctx context.Context, opName, username, namespace string) error {
	yamlFile := ""
	reset := false
	for _, op := range supportedOps {
		if op.key == opName {
			yamlFile = op.templateName
			if op.resetOp == true {
				reset = true
			}
		}
	}
	if yamlFile == "" && !reset {
		return fmt.Errorf("error: %s is not a valid operation name", opName)
	}
	if reset {
		iClient.deleteAllCreatedResources(ctx, namespace)
		return nil
	}
	return iClient.applyConfigChange(ctx, path.Join("..", "meshes", "istio", "config_templates", yamlFile), username, namespace)
}

func (iClient *IstioClient) applyConfigChange(ctx context.Context, yamlFile, username, namespace string) error {
	iClient.deleteAllCreatedResources(ctx, namespace)
	tmpl := template.Must(template.ParseFiles(yamlFile))

	buf := bytes.NewBufferString("")
	err := tmpl.Execute(buf, map[string]string{
		"user_name": username,
		"namespace": namespace,
	})
	if err != nil {
		err = errors.Wrapf(err, "unable to parse template")
		logrus.Error(err)
		return err
	}
	completeYaml := buf.String()
	yamls := strings.Split(completeYaml, "---")

	for _, yml := range yamls {
		if strings.TrimSpace(yml) != "" {
			err = iClient.applyRulePayload(ctx, namespace, []byte(yml))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Operations - returns a list of supported operations on the mesh
func (iClient *IstioClient) Operations(ctx context.Context) (map[string]string, error) {
	result := map[string]string{}
	for _, op := range supportedOps {
		result[op.key] = op.name
	}
	return result, nil
}
