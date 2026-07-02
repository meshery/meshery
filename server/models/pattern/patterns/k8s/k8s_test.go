package k8s

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	restfake "k8s.io/client-go/rest/fake"
)

func TestDryRunReturnsValidationErrorForInvalidResource(t *testing.T) {
	resource := map[string]interface{}{
		"kind":       "",
		"apiVersion": "",
	}

	status, success, err := dryRun(nil, resource, "", false)

	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if success {
		t.Fatal("expected success to be false")
	}

	if status != nil {
		t.Fatalf("expected nil status, got %#v", status)
	}
}

func TestDryRunCreateSuccess(t *testing.T) {
	resource := map[string]interface{}{
		"kind":       "Deployment",
		"apiVersion": "apps/v1",
		"metadata": map[string]interface{}{
			"name": "nginx",
		},
	}

	client := &restfake.RESTClient{
		GroupVersion: schema.GroupVersion{
			Group:   "apps",
			Version: "v1",
		},
		NegotiatedSerializer: scheme.Codecs.WithoutConversion(),
		Client: restfake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {

			if req.Method != http.MethodPost {
				t.Fatalf("expected POST request, got %s", req.Method)
			}

			if req.URL.Query().Get("dryRun") != "All" {
				t.Fatal("expected dryRun=All query parameter")
			}

			if req.URL.Path != "/apis/apps/v1/namespaces/default/deployments" {
				t.Fatalf("unexpected path: %s", req.URL.Path)
			}

			q := req.URL.Query()

			if q.Get("fieldValidation") != "Strict" {
				t.Fatal("expected fieldValidation=Strict")
			}

			if q.Get("fieldManager") != "meshery" {
				t.Fatal("expected fieldManager=meshery")
			}

			header := http.Header{}
			header.Set("Content-Type", runtime.ContentTypeJSON)

			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     header,
				Body:       io.NopCloser(strings.NewReader(`{"kind":"Deployment"}`)),
			}, nil
		}),
	}
	status, success, err := dryRun(client, resource, "default", false)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !success {
		t.Fatal("expected success to be true")
	}

	if status == nil {
		t.Fatal("expected non-nil status")
	}

	if status["kind"] != "Deployment" {
		t.Fatalf("expected kind Deployment, got %v", status["kind"])
	}

}

func TestDryRunDeleteSuccess(t *testing.T) {
	resource := map[string]interface{}{
		"kind":       "Deployment",
		"apiVersion": "apps/v1",
		"metadata": map[string]interface{}{
			"name": "nginx",
		},
	}

	client := &restfake.RESTClient{
		GroupVersion: schema.GroupVersion{
			Group:   "apps",
			Version: "v1",
		},
		NegotiatedSerializer: scheme.Codecs.WithoutConversion(),
		Client: restfake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {

			if req.Method != http.MethodDelete {
				t.Fatalf("expected DELETE request, got %s", req.Method)
			}

			if req.URL.Query().Get("dryRun") != "All" {
				t.Fatal("expected dryRun=All query parameter")
			}

			if req.URL.Path != "/apis/apps/v1/namespaces/default/deployments/nginx" {
				t.Fatalf("unexpected path: %s", req.URL.Path)
			}

			header := http.Header{}
			header.Set("Content-Type", runtime.ContentTypeJSON)

			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     header,
				Body:       io.NopCloser(strings.NewReader(`{"kind":"Deployment"}`)),
			}, nil
		}),
	}
	status, success, err := dryRun(client, resource, "default", true)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !success {
		t.Fatal("expected success to be true")
	}

	if status == nil {
		t.Fatal("expected non-nil status")
	}

	if status["kind"] != "Deployment" {
		t.Fatalf("expected kind Deployment, got %v", status["kind"])
	}
}

func TestDryRunDeleteReturnsErrorWhenMetadataMissing(t *testing.T) {
	resource := map[string]interface{}{
		"kind":       "Deployment",
		"apiVersion": "apps/v1",
	}

	status, success, err := dryRun(nil, resource, "default", true)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if success {
		t.Fatal("expected success to be false")
	}

	if status != nil {
		t.Fatalf("expected nil status, got %#v", status)
	}
}

func TestDryRunDeleteReturnsErrorWhenResourceNameMissing(t *testing.T) {
	resource := map[string]interface{}{
		"kind":       "Deployment",
		"apiVersion": "apps/v1",
		"metadata":   map[string]interface{}{},
	}

	status, success, err := dryRun(nil, resource, "default", true)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if success {
		t.Fatal("expected success to be false")
	}

	if status != nil {
		t.Fatalf("expected nil status, got %#v", status)
	}
}

func TestDryRunReturnsErrorWhenRESTClientFails(t *testing.T) {
	resource := map[string]interface{}{
		"kind":       "Deployment",
		"apiVersion": "apps/v1",
	}

	client := &restfake.RESTClient{
		Err: errors.New("connection failed"),
	}

	status, success, err := dryRun(client, resource, "default", false)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if success {
		t.Fatal("expected success to be false")
	}

	if status != nil {
		t.Fatalf("expected nil status, got %#v", status)
	}
}

func TestDryRunReturnsFailureStatus(t *testing.T) {
	resource := map[string]interface{}{
		"kind":       "Deployment",
		"apiVersion": "apps/v1",
		"metadata": map[string]interface{}{
			"name": "nginx",
		},
	}

	client := &restfake.RESTClient{
		GroupVersion: schema.GroupVersion{
			Group:   "apps",
			Version: "v1",
		},
		NegotiatedSerializer: scheme.Codecs.WithoutConversion(),
		Client: restfake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {
			header := http.Header{}
			header.Set("Content-Type", runtime.ContentTypeJSON)

			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     header,
				Body: io.NopCloser(strings.NewReader(`{
					"kind":"Status",
					"status":"Failure"
				}`)),
			}, nil
		}),
	}

	status, success, err := dryRun(client, resource, "default", false)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if success {
		t.Fatal("expected success to be false")
	}

	if status == nil {
		t.Fatal("expected non-nil status")
	}

	if status["status"] != "Failure" {
		t.Fatalf("expected Failure status, got %v", status["status"])
	}
}

func TestDryRunCreateCoreResourceSuccess(t *testing.T) {
	resource := map[string]interface{}{
		"kind":       "Pod",
		"apiVersion": "v1",
		"metadata": map[string]interface{}{
			"name": "nginx",
		},
	}

	client := &restfake.RESTClient{
		GroupVersion: schema.GroupVersion{
			Version: "v1",
		},
		NegotiatedSerializer: scheme.Codecs.WithoutConversion(),
		Client: restfake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {

			if req.Method != http.MethodPost {
				t.Fatalf("expected POST request, got %s", req.Method)
			}

			if req.URL.Query().Get("dryRun") != "All" {
				t.Fatal("expected dryRun=All query parameter")
			}

			if req.URL.Path != "/api/v1/namespaces/default/pods" {
				t.Fatalf("unexpected path: %s", req.URL.Path)
			}

			q := req.URL.Query()

			if q.Get("fieldValidation") != "Strict" {
				t.Fatal("expected fieldValidation=Strict")
			}

			if q.Get("fieldManager") != "meshery" {
				t.Fatal("expected fieldManager=meshery")
			}

			header := http.Header{}
			header.Set("Content-Type", runtime.ContentTypeJSON)

			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     header,
				Body:       io.NopCloser(strings.NewReader(`{"kind":"Pod"}`)),
			}, nil
		}),
	}
	status, success, err := dryRun(client, resource, "default", false)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !success {
		t.Fatal("expected success to be true")
	}

	if status == nil {
		t.Fatal("expected non-nil status")
	}

	if status["kind"] != "Pod" {
		t.Fatalf("expected kind Pod, got %v", status["kind"])
	}

}
