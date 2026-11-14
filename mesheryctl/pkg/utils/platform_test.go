// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestGetManifestTreeURL(t *testing.T) {
	t.Run("returns error on non-200 HTTP status code (e.g. 403 Rate Limit)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"message": "API rate limit exceeded"}`))
		}))
		defer server.Close()

		origURL := githubBaseURL
		githubBaseURL = server.URL
		defer func() { githubBaseURL = origURL }()

		_, err := GetManifestTreeURL("v0.6.0")
		expectedURL := githubBaseURL + "/repos/" + "meshery" + "/" + "meshery" + "/git/trees/" + "v0.6.0" + "?recursive=1"
		expectedErr := ErrGitHubAPIResponse(http.StatusForbidden, expectedURL, `{"message": "API rate limit exceeded"}`)
		AssertMeshkitErrorsEqual(t, err, expectedErr)
	})

	t.Run("happy path with valid ManifestList", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			respBody := `{
				"sha": "12345",
				"tree": [
					{
						"path": "install/deployment_yamls/k8s",
						"url": "https://api.github.com/repos/meshery/meshery/git/trees/k8s-tree-sha"
					}
				]
			}`
			_, _ = w.Write([]byte(respBody))
		}))
		defer server.Close()

		origURL := githubBaseURL
		githubBaseURL = server.URL
		defer func() { githubBaseURL = origURL }()

		treeURL, err := GetManifestTreeURL("v0.6.0")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		expected := "https://api.github.com/repos/meshery/meshery/git/trees/k8s-tree-sha"
		if treeURL != expected {
			t.Errorf("expected treeURL %q, got %q", expected, treeURL)
		}
	})
}

func TestListManifests(t *testing.T) {
	t.Run("returns error on non-200 HTTP status code (e.g. 403 Rate Limit)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"message": "API rate limit exceeded"}`))
		}))
		defer server.Close()

		_, err := ListManifests(server.URL)
		expectedErr := ErrGitHubAPIResponse(http.StatusForbidden, server.URL, `{"message": "API rate limit exceeded"}`)
		AssertMeshkitErrorsEqual(t, err, expectedErr)
	})

	t.Run("happy path with valid ManifestList", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			respBody := `{
				"sha": "12345",
				"tree": [
					{
						"path": "install/deployment_yamls/k8s/meshery-deployment.yaml",
						"type": "blob"
					}
				]
			}`
			_, _ = w.Write([]byte(respBody))
		}))
		defer server.Close()

		manifests, err := ListManifests(server.URL)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(manifests) != 1 {
			t.Fatalf("expected 1 manifest, got %d", len(manifests))
		}
		expectedPath := "install/deployment_yamls/k8s/meshery-deployment.yaml"
		if manifests[0].Path != expectedPath {
			t.Errorf("expected manifest path %q, got %q", expectedPath, manifests[0].Path)
		}
	})
}

func TestGetCleanPodName(t *testing.T) {
	tests := []struct {
		name     string
		podName  string
		expected string
	}{
		{
			name:     "standard deployment pod name",
			podName:  "meshery-broker-7d9f8c9b47-x8vqp",
			expected: "meshery-broker",
		},
		{
			name:     "two segments, alphabetic second segment joins both",
			podName:  "brokerpod-abc",
			expected: "brokerpod-abc",
		},
		{
			name:     "statefulset pod name, numeric ordinal falls back to first segment",
			podName:  "broker-0",
			expected: "broker",
		},
		{
			name:     "single word pod name with no hyphen",
			podName:  "meshery",
			expected: "meshery",
		},
		{
			name:     "empty pod name",
			podName:  "",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Names with fewer than two hyphen-separated segments used to panic with
			// "slice bounds out of range"; GetCleanPodName must now return cleanly.
			got := GetCleanPodName(tc.podName)
			if got != tc.expected {
				t.Errorf("GetCleanPodName(%q) = %q, want %q", tc.podName, got, tc.expected)
			}
		})
	}
}

func TestCreateManifestsFolder(t *testing.T) {
	// a manifests folder in the working directory belongs to the user, not mesheryctl
	t.Chdir(t.TempDir())
	if err := os.MkdirAll(filepath.Join(ManifestsFolder, "keep"), os.ModePerm); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	originalMesheryFolder := MesheryFolder
	defer func() { MesheryFolder = originalMesheryFolder }()
	MesheryFolder = t.TempDir()
	stale := filepath.Join(MesheryFolder, ManifestsFolder, "stale")
	if err := os.MkdirAll(stale, os.ModePerm); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if err := CreateManifestsFolder(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if _, err := os.Stat(filepath.Join(ManifestsFolder, "keep")); err != nil {
		t.Errorf("expected ./%s in the working directory to be untouched, got: %v", ManifestsFolder, err)
	}
	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Errorf("expected stale files in %s to be removed, got: %v", filepath.Join(MesheryFolder, ManifestsFolder), err)
	}
	if _, err := os.Stat(filepath.Join(MesheryFolder, ManifestsFolder)); err != nil {
		t.Errorf("expected %s to exist, got: %v", filepath.Join(MesheryFolder, ManifestsFolder), err)
	}
}

type interceptingTransport struct {
	body    *trackingBody
	handler func(req *http.Request) (*http.Response, error)
}

func (t *interceptingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.handler != nil {
		resp, err := t.handler(req)
		if err == nil && resp != nil && resp.Body != nil {
			t.body = &trackingBody{ReadCloser: resp.Body}
			resp.Body = t.body
		}
		return resp, err
	}
	return http.DefaultTransport.RoundTrip(req)
}

func TestInstallprereqDocker_BodyClosed(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("InstallprereqDocker shells out to uname, which is not available on Windows")
	}

	_ = SetupMeshkitLoggerTesting(t, false)

	tmpDir := t.TempDir()
	unameSh := filepath.Join(tmpDir, "uname")
	if err := os.WriteFile(unameSh, []byte("#!/bin/sh\necho test\n"), 0o755); err != nil {
		t.Fatalf("failed to create mock uname: %v", err)
	}
	t.Setenv("PATH", tmpDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	transport := &interceptingTransport{
		handler: func(req *http.Request) (*http.Response, error) {
			rec := httptest.NewRecorder()
			rec.WriteHeader(http.StatusOK)
			_, _ = rec.WriteString("invalid json")
			return rec.Result(), nil
		},
	}
	previousTransport := http.DefaultTransport
	http.DefaultTransport = transport
	t.Cleanup(func() { http.DefaultTransport = previousTransport })

	err := InstallprereqDocker()
	if err == nil {
		t.Fatal("expected error from unmarshal, got nil")
	}
	if transport.body == nil || !transport.body.closed {
		t.Fatal("expected response body to be closed in InstallprereqDocker")
	}
}

func newDeploymentClientset(names ...string) *k8sfake.Clientset {
	objects := make([]k8sruntime.Object, 0, len(names))
	for _, name := range names {
		objects = append(objects, &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: MesheryNamespace},
		})
	}
	return k8sfake.NewSimpleClientset(objects...)
}

func forbiddenDeleteReactor(group, resource string) k8stesting.ReactionFunc {
	return func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
		return true, nil, apierrors.NewForbidden(
			schema.GroupResource{Group: group, Resource: resource},
			action.(k8stesting.DeleteAction).GetName(),
			errors.New(resource+" is forbidden"),
		)
	}
}

func listErrorReactor(resource string) k8stesting.ReactionFunc {
	return func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
		return true, nil, errors.New(resource + " list failed")
	}
}

func TestCleanupCluster(t *testing.T) {
	t.Run("continues after multiple deletion failures and aggregates them", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset(
			&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "meshery-svc", Namespace: MesheryNamespace}},
			&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "meshery-secret", Namespace: MesheryNamespace}},
			&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "meshery-config", Namespace: MesheryNamespace}},
			&rbacv1.ClusterRole{ObjectMeta: metav1.ObjectMeta{Name: "meshery-role"}},
			&rbacv1.ClusterRole{ObjectMeta: metav1.ObjectMeta{Name: "other-role"}},
		)
		client.PrependReactor("delete", "services", forbiddenDeleteReactor("", "services"))
		client.PrependReactor("delete", "secrets", forbiddenDeleteReactor("", "secrets"))

		err := cleanupCluster(client)
		if err == nil {
			t.Fatal("expected the deletion failures to be reported, got nil")
		}
		for _, want := range []string{"service meshery-svc", "secret meshery-secret"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("expected the aggregated error to contain %q, got %q", want, err.Error())
			}
		}

		// Resource kinds after the failures must still have been cleaned up.
		configMaps, listErr := client.CoreV1().ConfigMaps(MesheryNamespace).List(context.TODO(), metav1.ListOptions{})
		if listErr != nil {
			t.Fatalf("failed to list configmaps: %v", listErr)
		}
		if len(configMaps.Items) != 0 {
			t.Errorf("expected the configmap to be deleted after earlier failures, got %d remaining", len(configMaps.Items))
		}

		// Only Meshery cluster-scoped resources are removed.
		clusterRoles, listErr := client.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{})
		if listErr != nil {
			t.Fatalf("failed to list clusterroles: %v", listErr)
		}
		if len(clusterRoles.Items) != 1 || clusterRoles.Items[0].GetName() != "other-role" {
			t.Errorf("expected only non-Meshery clusterroles to remain, got %+v", clusterRoles.Items)
		}
	})

	t.Run("a list failure does not prevent later kinds from being cleaned up", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset(
			&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "meshery-svc", Namespace: MesheryNamespace}},
			&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "meshery-config", Namespace: MesheryNamespace}},
		)
		// Fail only the list performed by cleanupCluster so the assertions below
		// can still inspect the cluster state.
		listAttempts := 0
		client.PrependReactor("list", "services", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
			listAttempts++
			if listAttempts == 1 {
				return true, nil, errors.New("services list failed")
			}
			return false, nil, nil
		})

		err := cleanupCluster(client)
		if err == nil || !strings.Contains(err.Error(), "cannot list services") {
			t.Fatalf("expected the list failure to be reported with its resource kind, got %v", err)
		}

		// A kind that could still be listed must still be cleaned up.
		configMaps, listErr := client.CoreV1().ConfigMaps(MesheryNamespace).List(context.TODO(), metav1.ListOptions{})
		if listErr != nil {
			t.Fatalf("failed to list configmaps: %v", listErr)
		}
		if len(configMaps.Items) != 0 {
			t.Errorf("expected the configmap to be deleted after the service list failure, got %d remaining", len(configMaps.Items))
		}

		services, listErr := client.CoreV1().Services(MesheryNamespace).List(context.TODO(), metav1.ListOptions{})
		if listErr != nil {
			t.Fatalf("failed to list services: %v", listErr)
		}
		if len(services.Items) != 1 {
			t.Errorf("expected the unlistable service to be left untouched, got %d remaining", len(services.Items))
		}
	})

	t.Run("aggregates multiple list and deletion failures", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset(
			&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: "meshery-svc", Namespace: MesheryNamespace}},
			&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "meshery-secret", Namespace: MesheryNamespace}},
			&rbacv1.Role{ObjectMeta: metav1.ObjectMeta{Name: "meshery-role", Namespace: MesheryNamespace}},
			&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "meshery-config", Namespace: MesheryNamespace}},
		)
		client.PrependReactor("list", "services", listErrorReactor("services"))
		client.PrependReactor("delete", "secrets", forbiddenDeleteReactor("", "secrets"))
		client.PrependReactor("list", "roles", listErrorReactor("roles"))

		err := cleanupCluster(client)
		if err == nil {
			t.Fatal("expected the list and deletion failures to be reported, got nil")
		}
		for _, want := range []string{"cannot list services", "secret meshery-secret", "cannot list roles"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("expected the aggregated error to contain %q, got %q", want, err.Error())
			}
		}

		aggregate, ok := err.(interface{ Errors() []error })
		if !ok {
			t.Fatalf("expected an aggregated error, got %T: %v", err, err)
		}
		if len(aggregate.Errors()) != 3 {
			t.Errorf("expected 3 aggregated failures, got %d: %v", len(aggregate.Errors()), err)
		}

		// Kinds whose List and Delete succeed are still cleaned up.
		configMaps, listErr := client.CoreV1().ConfigMaps(MesheryNamespace).List(context.TODO(), metav1.ListOptions{})
		if listErr != nil {
			t.Fatalf("failed to list configmaps: %v", listErr)
		}
		if len(configMaps.Items) != 0 {
			t.Errorf("expected the configmap to be deleted, got %d remaining", len(configMaps.Items))
		}

		secrets, listErr := client.CoreV1().Secrets(MesheryNamespace).List(context.TODO(), metav1.ListOptions{})
		if listErr != nil {
			t.Fatalf("failed to list secrets: %v", listErr)
		}
		if len(secrets.Items) != 1 {
			t.Errorf("expected the forbidden secret to remain, got %d remaining", len(secrets.Items))
		}
	})

	t.Run("ignores resources that are already gone", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset(
			&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: "meshery-config", Namespace: MesheryNamespace}},
		)
		client.PrependReactor("delete", "configmaps", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
			return true, nil, apierrors.NewNotFound(
				schema.GroupResource{Resource: "configmaps"},
				action.(k8stesting.DeleteAction).GetName())
		})

		if err := cleanupCluster(client); err != nil {
			t.Fatalf("expected NotFound to be ignored, got %v", err)
		}
	})
}

func TestDeleteDeployments(t *testing.T) {
	t.Run("wraps list failures with the resource kind", func(t *testing.T) {
		client := newDeploymentClientset("meshery")
		client.PrependReactor("list", "deployments", listErrorReactor("deployments"))

		err := deleteDeployments(context.TODO(), client.AppsV1().Deployments(MesheryNamespace))
		if err == nil {
			t.Fatal("expected a list failure to be reported, got nil")
		}
		if !strings.Contains(err.Error(), "cannot list deployments") {
			t.Errorf("expected the list failure to name the resource kind, got %v", err)
		}
		if !strings.Contains(err.Error(), "deployments list failed") {
			t.Errorf("expected the underlying list error to be preserved, got %v", err)
		}
	})

	t.Run("returns error when a deletion is forbidden", func(t *testing.T) {
		client := newDeploymentClientset("meshery")
		client.PrependReactor("delete", "deployments", forbiddenDeleteReactor("apps", "deployments"))

		err := deleteDeployments(context.TODO(), client.AppsV1().Deployments(MesheryNamespace))
		if err == nil {
			t.Fatal("expected a forbidden deletion to be reported, got nil")
		}
		if !strings.Contains(err.Error(), "forbidden") || !strings.Contains(err.Error(), "meshery") {
			t.Fatalf("expected the forbidden failure to name the deployment, got %v", err)
		}
	})

	t.Run("aggregates multiple deletion failures and tries every deployment", func(t *testing.T) {
		client := newDeploymentClientset("meshery-a", "meshery-b", "meshery-c")
		attempted := 0
		client.PrependReactor("delete", "deployments", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
			attempted++
			name := action.(k8stesting.DeleteAction).GetName()
			if name == "meshery-c" {
				return true, nil, nil
			}
			return true, nil, apierrors.NewForbidden(
				schema.GroupResource{Group: "apps", Resource: "deployments"}, name, errors.New("forbidden"))
		})

		err := deleteDeployments(context.TODO(), client.AppsV1().Deployments(MesheryNamespace))
		if err == nil {
			t.Fatal("expected aggregated deletion failures, got nil")
		}
		if attempted != 3 {
			t.Fatalf("expected all 3 deployments to be attempted, got %d", attempted)
		}

		aggregate, ok := err.(interface{ Errors() []error })
		if !ok {
			t.Fatalf("expected an aggregated error, got %T: %v", err, err)
		}
		if len(aggregate.Errors()) != 2 {
			t.Fatalf("expected 2 aggregated failures, got %d: %v", len(aggregate.Errors()), err)
		}
		for _, want := range []string{"meshery-a", "meshery-b"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("expected aggregated error to name %q, got %q", want, err.Error())
			}
		}
	})

	t.Run("ignores deployments that are already gone", func(t *testing.T) {
		client := newDeploymentClientset("meshery")
		client.PrependReactor("delete", "deployments", func(action k8stesting.Action) (bool, k8sruntime.Object, error) {
			return true, nil, apierrors.NewNotFound(
				schema.GroupResource{Group: "apps", Resource: "deployments"},
				action.(k8stesting.DeleteAction).GetName())
		})

		if err := deleteDeployments(context.TODO(), client.AppsV1().Deployments(MesheryNamespace)); err != nil {
			t.Fatalf("expected NotFound to be ignored, got %v", err)
		}
	})

	t.Run("deletes every deployment on success", func(t *testing.T) {
		client := newDeploymentClientset("meshery-a", "meshery-b")
		deployments := client.AppsV1().Deployments(MesheryNamespace)

		if err := deleteDeployments(context.TODO(), deployments); err != nil {
			t.Fatalf("expected cleanup to succeed, got %v", err)
		}

		remaining, err := deployments.List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			t.Fatalf("failed to list deployments: %v", err)
		}
		if len(remaining.Items) != 0 {
			t.Fatalf("expected no deployments to remain, got %d", len(remaining.Items))
		}
	})
}
