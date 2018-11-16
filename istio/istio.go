package istio

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
)

func CreateIstioClient() (*IstioClient, error) {
	iClient, err := NewClient()
	if err != nil {
		err = errors.Wrapf(err, "unable to create a new istio client")
		logrus.Error(err)
		return nil, err
	}
	return iClient, nil
}

func DeleteAllCreatedResources(iClient *IstioClient, namespace string) {
	resourceNames := []string{"productpage", "ratings", "reviews", "details"}
	for _, rs := range resourceNames {
		DeleteResource(iClient, virtualServices, namespace, rs)
	}
}

func DeleteResource(iClient *IstioClient, overallType, namespace, resName string) error {
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

func applyRulePayload(iClient *IstioClient, namespace string, newVSBytes []byte) error {
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
	newVSResInst, err := newVSRes.Get()
	if err != nil {
		newVSRes = iClient.istioNetworkingApi.Get().SetHeader("content-type", runtime.ContentTypeJSON).
			Namespace(namespace).Name(vs.ObjectMeta.Name).
			Resource(virtualServices).Do()
		newVSResInst, err = newVSRes.Get()
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

		newVSRes = iClient.istioNetworkingApi.Put().SetHeader("content-type", runtime.ContentTypeJSON).
			Namespace(namespace).Name(vs.ObjectMeta.Name).
			Resource(virtualServices).Body(newVSBytesJ).Do()
		newVSResInst, err = newVSRes.Get()
		if err != nil {
			err = errors.Wrapf(err, "unable to get the virtual service instance from result")
			logrus.Error(err)
			return err
		}
	}
	return nil
}

func ApplyYamlChange(iClient *IstioClient, yamlFile, username, namespace string) error {
	DeleteAllCreatedResources(iClient, namespace)
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
			err = applyRulePayload(iClient, namespace, []byte(yml))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
