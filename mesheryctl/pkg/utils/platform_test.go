// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	v1core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// The function are related to download should be test in meshkit package, please do not add test here.

func TestPlatform(t *testing.T) {
	t.Run("GetManifestTreeURL", func(t *testing.T) {
		_, err := GetManifestTreeURL("master")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestListManifests(t *testing.T) {
	t.Run("ListManifests with empty manifest", func(t *testing.T) {
		url := "https://api.github.com/repos/meshery/meshery/git/trees/47c634a49e6d143a54d734437a26ad233146ddf5"

		_, err := ListManifests(url)
		if err != nil {
			t.Errorf("ListManifests failed: %v", err)
		}
	})
}

func TestGetManifestURL(t *testing.T) {
	t.Run("GetManifestURL with empty manifest", func(t *testing.T) {
		manifest := Manifest{
			Typ:  "blob",
			Path: "testdata/manifest.yaml",
		}
		rawManifestsURL := "https://raw.githubusercontent.com/meshery/meshery/master/mesheryctl/pkg/utils/"

		manifestURL := GetManifestURL(manifest, rawManifestsURL)
		expectedURL := "https://raw.githubusercontent.com/meshery/meshery/master/mesheryctl/pkg/utils/testdata/manifest.yaml"

		if manifestURL != expectedURL {
			t.Errorf("GetManifestURL failed: expected %s, but got %s", expectedURL, manifestURL)
		}
	})
}

func TestCanUseCachedOperatorManifests(t *testing.T) {
	createDirectories := func(t *testing.T) {
		err := os.MkdirAll(filepath.Join(MesheryFolder, ManifestsFolder), os.ModePerm)
		if err != nil {
			t.Fatalf("failed to create directories: %v", err)
		}
	}

	createFile := func(t *testing.T, filename string) {
		_, err := os.Create(filepath.Join(MesheryFolder, ManifestsFolder, filename))
		if err != nil {
			t.Fatalf("failed to create %s file: %v", filename, err)
		}
	}

	removeFile := func(t *testing.T, filename string) {
		err := os.Remove(filepath.Join(MesheryFolder, ManifestsFolder, filename))
		if err != nil {
			t.Fatalf("failed to remove %s file: %v", filename, err)
		}
	}

	tests := []struct {
		name          string
		setupMock     func(t *testing.T)
		expectedError string
	}{
		{
			name: "All files exist",
			setupMock: func(t *testing.T) {
				createDirectories(t)
				createFile(t, MesheryOperator)
				createFile(t, MesheryOperatorBroker)
				createFile(t, MesheryOperatorMeshsync)
			},
			expectedError: "",
		},
		{
			name: "Operator file missing",
			setupMock: func(t *testing.T) {
				createDirectories(t)
				removeFile(t, MesheryOperator)
				createFile(t, MesheryOperatorBroker)
				createFile(t, MesheryOperatorMeshsync)
			},
			expectedError: "operator manifest file does not exist",
		},
		{
			name: "Broker file missing",
			setupMock: func(t *testing.T) {
				createDirectories(t)
				createFile(t, MesheryOperator)
				removeFile(t, MesheryOperatorBroker)
				createFile(t, MesheryOperatorMeshsync)
			},
			expectedError: "broker manifest file does not exist",
		},
		{
			name: "Meshsync file missing",
			setupMock: func(t *testing.T) {
				createDirectories(t)
				createFile(t, MesheryOperator)
				createFile(t, MesheryOperatorBroker)
				removeFile(t, MesheryOperatorMeshsync)
			},
			expectedError: "meshsync manifest file does not exist",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock(t)

			err := CanUseCachedOperatorManifests(nil)

			if tc.expectedError == "" && err != nil {
				t.Fatalf("expected no error, but got %v", err)
			}
			if tc.expectedError != "" && err == nil {
				t.Fatalf("expected error %v, but got no error", tc.expectedError)
			}
			if tc.expectedError != "" && err != nil && err.Error() != tc.expectedError {
				t.Fatalf("expected error %v, but got %v", tc.expectedError, err.Error())
			}
		})
	}

	// Cleanup mock environment
	os.RemoveAll(filepath.Join(MesheryFolder, ManifestsFolder))
}

func TestIsAdapterValid(t *testing.T) {
	testCases := []struct {
		name              string
		manifestArr       []Manifest
		componentManifest string
		expected          bool
	}{
		{
			name: "Valid Adapter",
			manifestArr: []Manifest{
				{Path: "component1.yaml"},
				{Path: "component2.yaml"},
			},
			componentManifest: "component1.yaml",
			expected:          true,
		},
		{
			name: "Invalid Adapter",
			manifestArr: []Manifest{
				{Path: "component1.yaml"},
				{Path: "component2.yaml"},
			},
			componentManifest: "component3.yaml",
			expected:          false,
		},
		{
			name:              "Empty Manifest Array",
			manifestArr:       []Manifest{},
			componentManifest: "component1.yaml",
			expected:          false,
		},
		{
			name: "Empty Component Manifest",
			manifestArr: []Manifest{
				{Path: "component1.yaml"},
				{Path: "component2.yaml"},
			},
			componentManifest: "",
			expected:          false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := IsAdapterValid(tc.manifestArr, tc.componentManifest)
			if result != tc.expected {
				t.Errorf("IsAdapterValid failed for %s: expected %v, but got %v", tc.name, tc.expected, result)
			}
		})
	}
}

func TestGetRequiredPods(t *testing.T) {
	testCases := []struct {
		name           string
		specifiedPods  []string
		availablePods  []v1core.Pod
		expectedResult map[string]string
		expectError    bool
	}{
		{
			name:          "All specified pods exist",
			specifiedPods: []string{"pod1", "pod2"},
			availablePods: []v1core.Pod{
				{ObjectMeta: metav1.ObjectMeta{Name: "pod1"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "pod2"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "pod3"}},
			},
			expectedResult: map[string]string{"pod1": "pod1", "pod2": "pod2"},
			expectError:    false,
		},
		{
			name:          "Some specified pods do not exist",
			specifiedPods: []string{"pod1", "pod4"},
			availablePods: []v1core.Pod{
				{ObjectMeta: metav1.ObjectMeta{Name: "pod1"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "pod2"}},
			},
			expectedResult: nil,
			expectError:    true,
		},
		{
			name:          "No specified pods exist",
			specifiedPods: []string{"pod4", "pod5"},
			availablePods: []v1core.Pod{
				{ObjectMeta: metav1.ObjectMeta{Name: "pod1"}},
				{ObjectMeta: metav1.ObjectMeta{Name: "pod2"}},
			},
			expectedResult: nil,
			expectError:    true,
		},
		{
			name:           "Specified pods list is empty",
			specifiedPods:  []string{},
			availablePods:  []v1core.Pod{},
			expectedResult: map[string]string{},
			expectError:    false,
		},
		{
			name:           "Available pods list is empty",
			specifiedPods:  []string{"pod1", "pod2"},
			availablePods:  []v1core.Pod{},
			expectedResult: nil,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GetRequiredPods(tc.specifiedPods, tc.availablePods)
			if tc.expectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error but got: %v", err)
				}
				if !reflect.DeepEqual(result, tc.expectedResult) {
					t.Errorf("expected %v, but got %v", tc.expectedResult, result)
				}
			}
		})
	}
}

func TestGetCleanPodName(t *testing.T) {
	testCases := []struct {
		name     string
		podName  string
		expected string
	}{
		{"Valid pod name", "meshery-istio-123abc", "meshery-istio"},
		{"Pod name with no useful information", "meshery-123abc", "meshery"},
		{"Pod name with special characters", "meshery-istio-123abc-!@#", "meshery-istio"},
		{"Already clean pod name", "meshery-istio", "meshery-istio"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetCleanPodName(tc.podName)
			if result != tc.expected {
				t.Errorf("GetCleanPodName failed: expected %s, but got %s", tc.expected, result)
			}
		})
	}
}
