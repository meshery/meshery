package helpers

import (
	"strings"

	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// NOT TO BE UPDATED at runtime
var meshesMeta = map[string][]string{
	"Istio": {
		"docker.io/istio/citadel",
		"docker.io/istio/proxyv2",
		"docker.io/istio/galley",
		"docker.io/istio/pilot",
	},
	"Linkerd": {
		"gcr.io/linkerd-io/controller",
		"gcr.io/linkerd-io/proxy",
	},
	"Consul": {
		"hashicorp/consul-k8s",
	},
	"Network Service Mesh": {
		"docker.io/networkservicemesh",
	},
	"Citrix": {
		"quay.io/citrix/citrix-istio-adaptor",
		"quay.io/citrix/citrix-k8s-cpx-ingress",
	},
	"osm": {
		"openservicemesh/osm-controller",
		"openservicemesh/init",
	},
}

// ScanKubernetes - Runs a quick scan on kubernetes to find out the version of service meshes deployed
func ScanKubernetes(kubeconfig []byte, contextName string) (map[string][]v1.Deployment, error) {
	clientset, err := getK8SClientSet(kubeconfig, contextName)
	if err != nil {
		return nil, err
	}
	namespacelist, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		err = errors.Wrap(err, "unable to get the list of namespaces")
		logrus.Error(err)
		return nil, err
	}
	newResults := map[string][]v1.Deployment{}
	for _, ns := range namespacelist.Items {
		logrus.Debugf("Listing deployments in namespace %q", ns.GetName())

		deploymentsClient := clientset.AppsV1().Deployments(ns.GetName())
		deplist, err := deploymentsClient.List(metav1.ListOptions{})
		if err != nil {
			err = errors.Wrapf(err, "unable to get deployments in the %s namespace", ns.GetName())
			logrus.Error(err)
			return nil, err
		}

		for _, d := range deplist.Items {
			logrus.Debugf(" * %s (%d replicas)", d.Name, *d.Spec.Replicas)
			meshIdentifier := ""
			for _, cont := range d.Spec.Template.Spec.Containers {
				// logrus.Debugf("    - name: %s, image: %s", cont.Name, cont.Image)
				for meshName, imageNames := range meshesMeta {
					for _, imageName := range imageNames {
						if strings.HasPrefix(cont.Image, imageName) || strings.Contains(cont.Image, imageName+":") {
							meshIdentifier = meshName
						}
					}
				}
			}
			if meshIdentifier != "" {
				newResults[meshIdentifier] = append(newResults[meshIdentifier], d)
			}
		}
	}
	return newResults, nil
}

// ScanPromGrafana - Runs a quick scan for Prometheus & Grafanas
func ScanPromGrafana(kubeconfig []byte, contextName string) (map[string][]string, error) {
	imageNames := []string{"prometheus", "grafana"}

	return detectServiceForDeploymentImage(kubeconfig, contextName, imageNames)
}

// ScanPrometheus - Runs a quick scan for Prometheus
func ScanPrometheus(kubeconfig []byte, contextName string) (map[string][]string, error) {
	imageNames := []string{"prometheus"}

	return detectServiceForDeploymentImage(kubeconfig, contextName, imageNames)
}

// ScanGrafana - Runs a quick scan for Grafanas
func ScanGrafana(kubeconfig []byte, contextName string) (map[string][]string, error) {
	imageNames := []string{"grafana"}

	return detectServiceForDeploymentImage(kubeconfig, contextName, imageNames)
}

func detectServiceForDeploymentImage(kubeconfig []byte, contextName string, imageNames []string) (map[string][]string, error) {
	clientset, err := getK8SClientSet(kubeconfig, contextName)
	if err != nil {
		return nil, err
	}
	namespacelist, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		err = errors.Wrap(err, "unable to get the list of namespaces")
		logrus.Error(err)
		return nil, err
	}
	result := map[string][]string{}

	for _, ns := range namespacelist.Items {
		logrus.Debugf("Listing deployments in namespace %q", ns.GetName())

		deploymentsClient := clientset.AppsV1().Deployments(ns.GetName())
		deplist, err := deploymentsClient.List(metav1.ListOptions{})
		if err != nil {
			err = errors.Wrapf(err, "unable to get deployments in the %s namespace", ns.GetName())
			logrus.Error(err)
			return nil, err
		}

		for _, d := range deplist.Items {
			logrus.Debugf(" * %s (%d replicas)", d.Name, *d.Spec.Replicas)
			foundDeployment := false
			for _, cont := range d.Spec.Template.Spec.Containers {
				for _, imageName := range imageNames {
					if strings.HasPrefix(cont.Image, imageName) || strings.Contains(cont.Image, imageName+":") {
						foundDeployment = true
						break
					}
				}
				if foundDeployment {
					break
				}
			}
			if foundDeployment {
				logrus.Debugf("found deployment: %s", d.GetName())
				lbls := d.Spec.Template.ObjectMeta.GetLabels()
				svcClient := clientset.CoreV1().Services(ns.GetName())
				svcList, err := svcClient.List(metav1.ListOptions{
					LabelSelector: labels.SelectorFromSet(lbls).String(),
				})
				if err != nil {
					err = errors.Wrapf(err, "unable to get deployments in the %s namespace", ns.GetName())
					logrus.Error(err)
					return nil, err
				}
				for _, sv := range svcList.Items {
					logrus.Debugf("Service Name: %s", sv.GetName())
					logrus.Debugf("Service type: %s", sv.Spec.Type)
					ports := []string{}
					for _, spr := range sv.Spec.Ports {
						logrus.Debugf("protocol: %s, port: %d", spr.Protocol, spr.Port)
						ports = append(ports, fmt.Sprintf("%d", spr.Port))
					}
					result[sv.GetName()+"."+sv.GetNamespace()] = ports
				}
			}
		}
	}
	logrus.Debugf("Derived tags: %s", result)

	// use that to go thru services with the given tags
	// from there get the ports and service type
	return result, nil
}
