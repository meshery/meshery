package utils

import (
	"context"
	"fmt"
	"testing"

	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var versionCheck = []struct {
	version  string // input
	expected [3]int // expected output
}{
	{"1.12.0", [3]int{1, 12, 0}},
	{"1.12.1", [3]int{1, 12, 1}},
	{"1.12.2", [3]int{1, 12, 2}},
	{"1.12.3", [3]int{1, 12, 3}},
	{"1.12.4", [3]int{1, 12, 4}},
	{"1.12.5", [3]int{1, 12, 5}},
	{"1.12.6", [3]int{1, 12, 6}},
	{"1.12.7", [3]int{1, 12, 7}},
	{"1.12.8", [3]int{1, 12, 8}},
	{"1.12.9", [3]int{1, 12, 9}},
	{"1.12.10", [3]int{1, 12, 10}},
}

func TestGetK8sVersion(t *testing.T) {
	for _, tt := range versionCheck {
		t.Run(tt.version, func(t *testing.T) {
			if got, _ := getK8sVersion(tt.version); got != tt.expected {
				t.Errorf("getK8sVersion() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCheckKubectlVersion(t *testing.T) {
	tests := []struct {
		version  string // input
		expected [3]int // expected output
	}{
		{
			version:  "v1.12.0",
			expected: [3]int{1, 12, 0},
		},
		{
			version:  "v1.12.1",
			expected: [3]int{1, 12, 1},
		},
		{
			version:  "v1.12.2",
			expected: [3]int{1, 12, 2},
		},
		{
			version:  "v1.12.3",
			expected: [3]int{1, 12, 3},
		},
		{
			version:  "v1.12.4",
			expected: [3]int{1, 12, 4},
		},
		{
			version:  "v1.12.5",
			expected: [3]int{1, 12, 5},
		},
		{
			version:  "v1.12.6",
			expected: [3]int{1, 12, 6},
		},
		{
			version:  "v1.12.7",
			expected: [3]int{1, 12, 7},
		},
		{
			version:  "v1.12.8",
			expected: [3]int{1, 12, 8},
		},
		{
			version:  "v1.12.9",
			expected: [3]int{1, 12, 9},
		},
		{
			version:  "v1.12.10",
			expected: [3]int{1, 12, 10},
		},
	}

	_, err := parseKubectlShortVersion("Not really a version")
	if err == nil {
		t.Fatalf("Expected to get an error")
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			if got, _ := parseKubectlShortVersion(tt.version); got != tt.expected {
				t.Errorf("checkKubectlVersion() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCheckK8sVersion(t *testing.T) {
	tests := []struct {
		version string
		expected [3]int
	}{
		{
			version: "v1.12.0",
			expected: [3]int{1, 12, 0},
		},
		{
			version: "v1.12.1",
			expected: [3]int{1, 12, 1},
		},
		{
			version: "v1.12.2",
			expected: [3]int{1, 12, 2},
		},
		{
			version: "v1.12.3",
			expected: [3]int{1, 12, 3},
		},
		{
			version: "v1.12.4",
			expected: [3]int{1, 12, 4},
		},
		{
			version: "v1.12.5",
			expected: [3]int{1, 12, 5},
		},
		{
			version: "v1.12.6",
			expected: [3]int{1, 12, 6},
		},
		{
			version: "v1.12.7",
			expected: [3]int{1, 12, 7},
		},
		{
			version: "v1.12.8",
			expected: [3]int{1, 12, 8},
		},
		{
			version: "v1.12.9",
			expected: [3]int{1, 12, 9},
		},
		{
			version: "v1.12.10",
			expected: [3]int{1, 12, 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			if got, _:= getK8sVersion(tt.version); got != tt.expected {
				t.Errorf("getK8sVersion() = %v, want %v", got, tt.expected)
			}
		})
	}

}


func TestIsPodRunning(t *testing.T) {
	tests := []struct {
		namespace string
		podname string
		expected string
	}{
		{
			namespace: "test0",
			podname: "test-pod0",
			expected: "Running",
		},
		{
			namespace: "test1",
			podname: "test-pod1",
			expected: "Running",
		},
		{
			namespace: "test2",
			podname: "test-pod2",
			expected: "Running",
		},
		{
			namespace: "test3",
			podname: "test-pod3",
			expected: "Running",
		},
		{
			namespace: "test4",
			podname: "test-pod4",
			expected: "Running",
		},
		{
			namespace: "",
			podname: "test-pod",
			expected: "Running",
		},
		{
			namespace: "test",
			podname: "",
			expected: "Running",
		},
		{
			namespace: "",
			podname: "",
			expected: "Running",
		},
		{
			podname: "test-pod",
			expected: "Failed",
		},
		{
			namespace: "test",
			expected: "Failed",
		},
	}

	for _, tt := range tests {
		c, _ := meshkitkube.New([]byte(""))
		got, _ := c.KubeClient.CoreV1().Pods(tt.namespace).Get(context.TODO(), tt.podname, metav1.GetOptions{})
		if got.Status.Phase != v1.PodPhase(tt.expected) {
			fmt.Errorf("Pod is Running  got %v want %v", got, tt.expected)
		}
	}
}
