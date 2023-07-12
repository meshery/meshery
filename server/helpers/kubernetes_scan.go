package helpers

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOT TO BE UPDATED at runtime
var meshesMeta = map[string][]string{
	"Istio": {
		"istio",
	},
	"Linkerd": {
		"linkerd",
	},
	"Consul": {
		"consul",
	},
	"Network Service Mesh": {
		"networkservicemesh",
		"nsm",
	},
	"Citrix": {
		"citrix",
	},
	"osm": {
		"openservicemesh",
		"osm",
	},
}

var TelemetryComps = []string{
	"prometheus", "grafana",
}

// ScanKubernetes scans kubernetes to find the pods for each service mesh
func ScanKubernetes(kubeconfig []byte, contextName string) (map[string][]corev1.Pod, error) {
	clientset, err := getK8SClientSet(kubeconfig, contextName)
	if err != nil {
		return nil, ErrScanKubernetes(err)
	}
	// equivalent to GET request to /api/v1/pods
	podlist, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		logrus.Debug("[ScanKubernetes] Failed to retrieve Pod List")
		return nil, ErrRetrievePodList(err)
	}

	result := map[string][]corev1.Pod{}

	for _, p := range podlist.Items {
		logrus.Debugf("[ScanKubernetes] Found pod %s", p.Name)
		meshIdentifier := ""
		for meshName, names := range meshesMeta {
			for _, name := range names {
				if strings.Contains(p.Name, name) {
					meshIdentifier = meshName
					//Terminating the loop after a condition match.
					break
				}
			}
		}

		// Ignoring "kube-system" pods
		if meshIdentifier != "" && p.Namespace != "kube-system" {
			result[meshIdentifier] = append(result[meshIdentifier], p)
		}
	}

	return result, nil
}
