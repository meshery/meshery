package helpers

import (
	"strings"

	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// NOT TO BE UPDATED at runtime
var meshesMeta = map[string][]string{
	"Istio": []string{
		"docker.io/istio/citadel",
		"docker.io/istio/proxyv2",
		"docker.io/istio/galley",
		"docker.io/istio/pilot",
	},
	"Linkerd": []string{
		"gcr.io/linkerd-io/controller",
		"gcr.io/linkerd-io/proxy",
	},
	"Consul": []string{
		"hashicorp/consul-k8s",
	},
	"Network Service Mesh": []string{
		"docker.io/networkservicemesh/nsmdp",
		"docker.io/networkservicemesh/nsmd",
		"docker.io/networkservicemesh/nsmd-k8s",
	},
}

// ScanKubernetes - Runs a quick scan on kubernetes to find out the version of service meshes deployed
func ScanKubernetes(kubeconfig []byte, contextName string) (map[string]string, error) {
	clientset, err := getK8SClientSet(kubeconfig, contextName)
	if err != nil {
		return nil, err
	}
	result := map[string]string{}
	namespacelist, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		err = errors.Wrap(err, "unable to get the list of namespaces")
		logrus.Error(err)
		return nil, err
	}
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
			for _, cont := range d.Spec.Template.Spec.Containers {
				// logrus.Debugf("    - name: %s, image: %s", cont.Name, cont.Image)
				for meshName, imageNames := range meshesMeta {
					for _, imageName := range imageNames {
						if strings.HasPrefix(cont.Image, imageName) || strings.Contains(cont.Image, imageName+":") {
							versionInfo := strings.Split(cont.Image, ":")[1]
							result[meshName] = versionInfo
						}
					}
				}
			}
		}
	}
	logrus.Debugf("Derived mesh versions: %s", result)
	return result, nil
}

// ScanKubernetes - Runs a quick scan on kubernetes to find out the version of service meshes deployed
func ScanPromGrafana(kubeconfig []byte, contextName string) (map[string][]string, error) {
	imageNames := []string{"prometheus", "grafana"}

	return detectServiceForDeploymentImage(kubeconfig, contextName, imageNames)
}

// ScanKubernetes - Runs a quick scan on kubernetes to find out the version of service meshes deployed
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
