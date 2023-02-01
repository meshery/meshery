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
// func ScanPromGrafana(kubeconfig []byte, contextName string) (map[string][]corev1.Service, error) {
// 	imageNames := []string{"prometheus", "grafana"}

// 	return detectServiceWithName(kubeconfig, contextName, imageNames)
// }

// ScanPrometheus - Runs a quick scan for Prometheus
// func ScanPrometheus(kubeconfig []byte, contextName string) (map[string][]corev1.Service, error) {
// 	imageNames := []string{"prometheus"}

// 	return detectServiceWithName(kubeconfig, contextName, imageNames)
// }

// ScanGrafana - Runs a quick scan for Grafanas
// func ScanGrafana(kubeconfig []byte, contextName string) (map[string][]corev1.Service, error) {
// 	imageNames := []string{"grafana"}

// 	return detectServiceWithName(kubeconfig, contextName, imageNames)
// }

// func detectServiceForDeploymentImage(kubeconfig []byte, contextName string, imageNames []string) (map[string][]string, error) {
// 	clientset, err := getK8SClientSet(kubeconfig, contextName)
// 	if err != nil {
// 		return nil, ErrDetectServiceForDeploymentImage(err)
// 	}
// 	namespacelist, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
// 	if err != nil {
// 		return nil, ErrRetrieveNamespacesList(err)
// 	}
// 	result := map[string][]string{}

// 	for _, ns := range namespacelist.Items {
// 		logrus.Debugf("Listing deployments in namespace %q", ns.GetName())

// 		deploymentsClient := clientset.AppsV1().Deployments(ns.GetName())
// 		deplist, err := deploymentsClient.List(context.Background(), metav1.ListOptions{})
// 		if err != nil {
// 			return nil, ErrGetNamespaceDeployments(err, ns.GetName())
// 		}

// 		for _, d := range deplist.Items {
// 			logrus.Debugf(" * %s (%d replicas)", d.Name, *d.Spec.Replicas)
// 			foundDeployment := false
// 			for _, cont := range d.Spec.Template.Spec.Containers {
// 				for _, imageName := range imageNames {
// 					if strings.HasPrefix(cont.Image, imageName) || strings.Contains(cont.Image, imageName+":") {
// 						foundDeployment = true
// 						break
// 					}
// 				}
// 				if foundDeployment {
// 					break
// 				}
// 			}
// 			if foundDeployment {
// 				logrus.Debugf("found deployment: %s", d.GetName())
// 				lbls := d.Spec.Template.ObjectMeta.GetLabels()
// 				svcClient := clientset.CoreV1().Services(ns.GetName())
// 				svcList, err := svcClient.List(context.Background(), metav1.ListOptions{
// 					LabelSelector: labels.SelectorFromSet(lbls).String(),
// 				})
// 				if err != nil {
// 					return nil, ErrGetNamespaceDeployments(err, ns.GetName())
// 				}
// 				for _, sv := range svcList.Items {
// 					logrus.Debugf("Service Name: %s", sv.GetName())
// 					logrus.Debugf("Service type: %s", sv.Spec.Type)
// 					ports := []string{}
// 					for _, spr := range sv.Spec.Ports {
// 						logrus.Debugf("protocol: %s, port: %d", spr.Protocol, spr.Port)
// 						ports = append(ports, fmt.Sprintf("%d", spr.Port))
// 					}
// 					result[sv.GetName()+"."+sv.GetNamespace()] = ports
// 				}
// 			}
// 		}
// 	}
// 	logrus.Debugf("Derived tags: %s", result)

// 	// use that to go thru services with the given tags
// 	// from there get the ports and service type
// 	return result, nil
// }

// detectServiceWithName detects the services in the cluster with the name given in "names" parameter
// func detectServiceWithName(kubeconfig []byte, contextName string, names []string) (map[string][]corev1.Service, error) {
// 	clientset, err := getK8SClientSet(kubeconfig, contextName)
// 	if err != nil {
// 		return nil, ErrDetectServiceWithName(err)
// 	}

// 	// Get all the running services
// 	// analogous to GET request to /api/v1/services
// 	svcList, err := clientset.CoreV1().Services("").List(context.Background(), metav1.ListOptions{})
// 	if err != nil {
// 		return nil, ErrDetectServiceWithName(err)
// 	}

// 	result := map[string][]corev1.Service{}

// 	for _, svc := range svcList.Items {
// 		svcName := svc.GetName()
// 		logrus.Debugf("[DetectServiceWithName] Service Name: %s", svcName)
// 		logrus.Debugf("[DetectServiceWithName] Service type: %s", svc.Spec.Type)

// 		for _, query := range names {
// 			if strings.Contains(strings.ToLower(svcName), query) {
// 				result[query] = append(result[query], svc)
// 			}
// 		}
// 	}

// 	return result, nil
// }
