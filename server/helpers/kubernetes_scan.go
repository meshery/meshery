package helpers

import (
	"context"
	"strings"

	"github.com/meshery/meshkit/logger"
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
}

var TelemetryComps = []string{
	"prometheus", "grafana",
}

// ScanKubernetes scans kubernetes to find the pods for each entity
func ScanKubernetes(kubeconfig []byte, contextName string, log logger.Handler) (map[string][]corev1.Pod, error) {
	clientset, err := getK8SClientSet(kubeconfig, contextName)
	if err != nil {
		return nil, ErrScanKubernetes(err)
	}
	// equivalent to GET request to /api/v1/pods
	podlist, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Debug("[ScanKubernetes] Failed to retrieve Pod List")
		return nil, ErrRetrievePodList(err)
	}

	result := map[string][]corev1.Pod{}

	for _, p := range podlist.Items {
		log.Debug("[ScanKubernetes] Found pod ", p.Name)
		meshIdentifier := ""
		// for _, cont := range p.Name {
		// 	for meshName, imageNames := range meshesMeta {
		// 		for _, imageName := range imageNames {
		// 			if strings.HasPrefix(cont.Image, imageName) || strings.Contains(cont.Image, imageName) {
		// 				meshIdentifier = meshName
		// 			}
		// 		}
		// 	}
		// }
		for meshName, names := range meshesMeta {
			for _, name := range names {
				if strings.Contains(p.Name, name) {
					meshIdentifier = meshName
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

// ScanPromGrafana - Runs a quick scan for Prometheus & Grafanas
func ScanPromGrafana(kubeconfig []byte, contextName string) (map[string][]corev1.Service, error) {
	imageNames := []string{"prometheus", "grafana"}

	return detectServiceWithName(kubeconfig, contextName, imageNames)
}

// ScanPrometheus - Runs a quick scan for Prometheus
func ScanPrometheus(kubeconfig []byte, contextName string) (map[string][]corev1.Service, error) {
	imageNames := []string{"prometheus"}

	return detectServiceWithName(kubeconfig, contextName, imageNames)
}

// ScanGrafana - Runs a quick scan for Grafanas
func ScanGrafana(kubeconfig []byte, contextName string) (map[string][]corev1.Service, error) {
	imageNames := []string{"grafana"}

	return detectServiceWithName(kubeconfig, contextName, imageNames)
}

// detectServiceWithName detects the services in the cluster with the name given in "names" parameter
func detectServiceWithName(kubeconfig []byte, contextName string, names []string) (map[string][]corev1.Service, error) {
	clientset, err := getK8SClientSet(kubeconfig, contextName)
	if err != nil {
		return nil, ErrDetectServiceWithName(err)
	}

	// Get all the running services
	// analogous to GET request to /api/v1/services
	svcList, err := clientset.CoreV1().Services("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, ErrDetectServiceWithName(err)
	}

	result := map[string][]corev1.Service{}

	for _, svc := range svcList.Items {
		svcName := svc.GetName()

		for _, query := range names {
			if strings.Contains(strings.ToLower(svcName), query) {
				result[query] = append(result[query], svc)
			}
		}
	}

	return result, nil
}
