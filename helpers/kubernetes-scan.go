package helpers

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOT TO BE UPDATED at runtime
var meshesMeta = map[string][]string{
	"istio": []string{
		"docker.io/istio/citadel",
		"docker.io/istio/proxyv2",
		"docker.io/istio/galley",
		"docker.io/istio/pilot",
	},
	"linkerd": []string{
		"gcr.io/linkerd-io/controller",
		"gcr.io/linkerd-io/proxy",
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
			err = errors.Wrapf(err, "unable to get deployments in the %s namespace", ns)
			logrus.Error(err)
			return nil, err
		}

		for _, d := range deplist.Items {
			logrus.Debugf(" * %s (%d replicas)", d.Name, *d.Spec.Replicas)
			for _, cont := range d.Spec.Template.Spec.Containers {
				// logrus.Debugf("    - name: %s, image: %s", cont.Name, cont.Image)
				for meshName, imageNames := range meshesMeta {
					for _, imageName := range imageNames {
						if strings.HasPrefix(cont.Image, imageName) {
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
