package handlers

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMergeResults(t *testing.T) {
	svc1 := corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "prometheus-1"}}
	svc2 := corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "prometheus-2"}}
	svc3 := corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "grafana-1"}}

	input := []map[string][]corev1.Service{
		{
			"prometheus": {svc1},
			"grafana":    {svc3},
		},
		{
			"prometheus": {svc2},
		},
	}

	result := make(map[string][]corev1.Service)
	for _, res := range input {
		for k, v := range res {
			result[k] = append(result[k], v...)
		}
	}

	if len(result["prometheus"]) != 2 {
		t.Fatalf("expected 2 prometheus services, got %d", len(result["prometheus"]))
	}

	if len(result["grafana"]) != 1 {
		t.Fatalf("expected 1 grafana service, got %d", len(result["grafana"]))
	}

	if result["prometheus"][0].Name != "prometheus-1" || result["prometheus"][1].Name != "prometheus-2" {
		t.Errorf("incorrect order or content of prometheus services")
	}
}
