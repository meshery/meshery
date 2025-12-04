// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandleOCICLIError(t *testing.T) {
	tests := []struct {
		name           string
		inputError     error
		expectedSubstr string
	}{
		{
			name:           "nil error",
			inputError:     nil,
			expectedSubstr: "",
		},
		{
			name:           "command not found error",
			inputError:     &os.PathError{Op: "exec", Path: "oci", Err: os.ErrNotExist},
			expectedSubstr: "OCI CLI not found",
		},
		{
			name:           "executable not found error",
			inputError:     os.ErrNotExist,
			expectedSubstr: "OCI CLI not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleOCICLIError(tt.inputError)
			
			if tt.inputError == nil {
				if err != nil {
					t.Errorf("handleOCICLIError() with nil input should return nil, got: %v", err)
				}
				return
			}
			
			if err == nil {
				t.Errorf("handleOCICLIError() should return error, got nil")
				return
			}
			
			if tt.expectedSubstr != "" && !strings.Contains(err.Error(), tt.expectedSubstr) {
				t.Errorf("handleOCICLIError() error = %v, should contain %v", err.Error(), tt.expectedSubstr)
			}
		})
	}
}

func TestHandleOCIConfigError(t *testing.T) {
	tests := []struct {
		name           string
		inputError     error
		profile        string
		expectedSubstr string
	}{
		{
			name:           "nil error",
			inputError:     nil,
			profile:        "DEFAULT",
			expectedSubstr: "",
		},
		{
			name:           "config file not found at path",
			inputError:     os.ErrNotExist,
			profile:        "DEFAULT",
			expectedSubstr: "OCI config file not found",
		},
		{
			name:           "config file not readable",
			inputError:     os.ErrPermission,
			profile:        "DEFAULT",
			expectedSubstr: "error validating OCI config",
		},
		{
			name:           "profile not found in config",
			inputError:     fmt.Errorf("profile 'NONEXISTENT' not found in OCI config"),
			profile:        "NONEXISTENT",
			expectedSubstr: "profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleOCIConfigError(tt.inputError, tt.profile)
			
			if tt.inputError == nil {
				if err != nil {
					t.Errorf("handleOCIConfigError() with nil input should return nil, got: %v", err)
				}
				return
			}
			
			if err == nil {
				t.Errorf("handleOCIConfigError() should return error, got nil")
				return
			}
			
			if tt.expectedSubstr != "" && !strings.Contains(err.Error(), tt.expectedSubstr) {
				t.Errorf("handleOCIConfigError() error = %v, should contain %v", err.Error(), tt.expectedSubstr)
			}
		})
	}
}

func TestHandleClusterAccessError(t *testing.T) {
	tests := []struct {
		name           string
		inputError     error
		clusterID      string
		expectedSubstr string
	}{
		{
			name:           "nil error",
			inputError:     nil,
			clusterID:      "ocid1.cluster.oc1.phx.test",
			expectedSubstr: "",
		},
		{
			name:           "NotAuthorizedOrNotFound error",
			inputError:     fmt.Errorf("NotAuthorizedOrNotFound: cluster not found"),
			clusterID:      "ocid1.cluster.oc1.phx.test",
			expectedSubstr: "not found or access denied",
		},
		{
			name:           "access denied error",
			inputError:     fmt.Errorf("access denied to cluster"),
			clusterID:      "ocid1.cluster.oc1.phx.test",
			expectedSubstr: "not found or access denied",
		},
		{
			name:           "InvalidParameter error",
			inputError:     fmt.Errorf("InvalidParameter: invalid cluster ID"),
			clusterID:      "invalid-cluster-id",
			expectedSubstr: "invalid cluster ID format",
		},
		{
			name:           "network timeout error",
			inputError:     fmt.Errorf("connection timeout while accessing cluster"),
			clusterID:      "ocid1.cluster.oc1.phx.test",
			expectedSubstr: "network connectivity issue",
		},
		{
			name:           "permission error",
			inputError:     fmt.Errorf("permission denied to access cluster"),
			clusterID:      "ocid1.cluster.oc1.phx.test",
			expectedSubstr: "insufficient permissions",
		},
		{
			name:           "generic error",
			inputError:     fmt.Errorf("some other error occurred"),
			clusterID:      "ocid1.cluster.oc1.phx.test",
			expectedSubstr: "error accessing cluster",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleClusterAccessError(tt.inputError, tt.clusterID)
			
			if tt.inputError == nil {
				if err != nil {
					t.Errorf("handleClusterAccessError() with nil input should return nil, got: %v", err)
				}
				return
			}
			
			if err == nil {
				t.Errorf("handleClusterAccessError() should return error, got nil")
				return
			}
			
			if tt.expectedSubstr != "" && !strings.Contains(err.Error(), tt.expectedSubstr) {
				t.Errorf("handleClusterAccessError() error = %v, should contain %v", err.Error(), tt.expectedSubstr)
			}
		})
	}
}

func TestValidateOCIConfig(t *testing.T) {
	// Create a temporary directory for test config files
	tempDir := t.TempDir()
	
	// Save original HOME and restore after test
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	
	// Set HOME to temp directory
	os.Setenv("HOME", tempDir)
	
	// Create .oci directory
	ociDir := filepath.Join(tempDir, ".oci")
	err := os.MkdirAll(ociDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test .oci directory: %v", err)
	}
	
	tests := []struct {
		name        string
		setupFunc   func() error
		profile     string
		expectError bool
	}{
		{
			name: "valid config with DEFAULT profile",
			setupFunc: func() error {
				configContent := `[DEFAULT]
user=ocid1.user.oc1..test
fingerprint=aa:bb:cc:dd:ee:ff
tenancy=ocid1.tenancy.oc1..test
region=us-phoenix-1
key_file=~/.oci/oci_api_key.pem
`
				return os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
			},
			profile:     "DEFAULT",
			expectError: false,
		},
		{
			name: "valid config with custom profile",
			setupFunc: func() error {
				configContent := `[DEFAULT]
user=ocid1.user.oc1..test
fingerprint=aa:bb:cc:dd:ee:ff
tenancy=ocid1.tenancy.oc1..test
region=us-phoenix-1
key_file=~/.oci/oci_api_key.pem

[CUSTOM]
user=ocid1.user.oc1..custom
fingerprint=11:22:33:44:55:66
tenancy=ocid1.tenancy.oc1..custom
region=us-ashburn-1
key_file=~/.oci/custom_key.pem
`
				return os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
			},
			profile:     "CUSTOM",
			expectError: false,
		},
		{
			name: "config file does not exist",
			setupFunc: func() error {
				// Remove config file if it exists
				configPath := filepath.Join(ociDir, "config")
				os.Remove(configPath)
				return nil
			},
			profile:     "DEFAULT",
			expectError: true,
		},
		{
			name: "profile not found in config",
			setupFunc: func() error {
				configContent := `[DEFAULT]
user=ocid1.user.oc1..test
fingerprint=aa:bb:cc:dd:ee:ff
tenancy=ocid1.tenancy.oc1..test
region=us-phoenix-1
key_file=~/.oci/oci_api_key.pem
`
				return os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
			},
			profile:     "NONEXISTENT",
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			if tt.setupFunc != nil {
				if err := tt.setupFunc(); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}
			
			// Run validation
			err := validateOCIConfig(tt.profile)
			
			// Check result
			if (err != nil) != tt.expectError {
				t.Errorf("validateOCIConfig() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestValidatePrerequisites(t *testing.T) {
	// Create a temporary directory for test config files
	tempDir := t.TempDir()
	
	// Save original HOME and restore after test
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	
	// Set HOME to temp directory
	os.Setenv("HOME", tempDir)
	
	// Create .oci directory with valid config
	ociDir := filepath.Join(tempDir, ".oci")
	err := os.MkdirAll(ociDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test .oci directory: %v", err)
	}
	
	configContent := `[DEFAULT]
user=ocid1.user.oc1..test
fingerprint=aa:bb:cc:dd:ee:ff
tenancy=ocid1.tenancy.oc1..test
region=us-phoenix-1
key_file=~/.oci/oci_api_key.pem
`
	err = os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	
	t.Run("validate prerequisites with valid config", func(t *testing.T) {
		// This test will check if OCI CLI is installed
		// If OCI CLI is not installed, it should return an error with helpful message
		err := validatePrerequisites("DEFAULT")
		
		// We expect this to either succeed (if OCI CLI is installed) or fail with helpful error
		if err != nil {
			// Check that error message contains helpful guidance
			if !strings.Contains(err.Error(), "OCI CLI") {
				t.Errorf("validatePrerequisites() error should mention OCI CLI, got: %v", err)
			}
		}
	})
	
	t.Run("validate prerequisites with invalid profile", func(t *testing.T) {
		err := validatePrerequisites("NONEXISTENT")
		
		// Should fail because profile doesn't exist
		if err == nil {
			t.Error("validatePrerequisites() should fail with nonexistent profile")
		}
		
		// Check that error message is helpful
		if !strings.Contains(err.Error(), "profile") && !strings.Contains(err.Error(), "OCI CLI") {
			t.Errorf("validatePrerequisites() error should be helpful, got: %v", err)
		}
	})
}

func TestReadOCIConfig(t *testing.T) {
	// Create a temporary directory for test config files
	tempDir := t.TempDir()
	
	// Save original HOME and restore after test
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	
	// Set HOME to temp directory
	os.Setenv("HOME", tempDir)
	
	// Create .oci directory
	ociDir := filepath.Join(tempDir, ".oci")
	err := os.MkdirAll(ociDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test .oci directory: %v", err)
	}
	
	// Create a test key file
	keyFilePath := filepath.Join(ociDir, "test_key.pem")
	err = os.WriteFile(keyFilePath, []byte("test key content"), 0600)
	if err != nil {
		t.Fatalf("Failed to create test key file: %v", err)
	}
	
	tests := []struct {
		name          string
		setupFunc     func() error
		profile       string
		expectError   bool
		validateFunc  func(*testing.T, *OCIProfile)
	}{
		{
			name: "valid config with all required fields",
			setupFunc: func() error {
				configContent := `[DEFAULT]
user=ocid1.user.oc1..testuser
fingerprint=aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99
tenancy=ocid1.tenancy.oc1..testtenant
region=us-phoenix-1
key_file=~/.oci/test_key.pem
`
				return os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
			},
			profile:     "DEFAULT",
			expectError: false,
			validateFunc: func(t *testing.T, profile *OCIProfile) {
				if profile.UserID != "ocid1.user.oc1..testuser" {
					t.Errorf("Expected UserID 'ocid1.user.oc1..testuser', got '%s'", profile.UserID)
				}
				if profile.TenancyID != "ocid1.tenancy.oc1..testtenant" {
					t.Errorf("Expected TenancyID 'ocid1.tenancy.oc1..testtenant', got '%s'", profile.TenancyID)
				}
				if profile.Region != "us-phoenix-1" {
					t.Errorf("Expected Region 'us-phoenix-1', got '%s'", profile.Region)
				}
				if profile.Fingerprint != "aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99" {
					t.Errorf("Expected Fingerprint 'aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99', got '%s'", profile.Fingerprint)
				}
				if profile.KeyFile != "~/.oci/test_key.pem" {
					t.Errorf("Expected KeyFile '~/.oci/test_key.pem', got '%s'", profile.KeyFile)
				}
			},
		},
		{
			name: "valid config with passphrase",
			setupFunc: func() error {
				configContent := `[DEFAULT]
user=ocid1.user.oc1..testuser
fingerprint=aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99
tenancy=ocid1.tenancy.oc1..testtenant
region=us-ashburn-1
key_file=~/.oci/test_key.pem
pass_phrase=mypassphrase
`
				return os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
			},
			profile:     "DEFAULT",
			expectError: false,
			validateFunc: func(t *testing.T, profile *OCIProfile) {
				if profile.PassPhrase != "mypassphrase" {
					t.Errorf("Expected PassPhrase 'mypassphrase', got '%s'", profile.PassPhrase)
				}
			},
		},
		{
			name: "config file does not exist",
			setupFunc: func() error {
				configPath := filepath.Join(ociDir, "config")
				os.Remove(configPath)
				return nil
			},
			profile:     "DEFAULT",
			expectError: true,
		},
		{
			name: "profile not found in config",
			setupFunc: func() error {
				configContent := `[DEFAULT]
user=ocid1.user.oc1..testuser
fingerprint=aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99
tenancy=ocid1.tenancy.oc1..testtenant
region=us-phoenix-1
key_file=~/.oci/test_key.pem
`
				return os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
			},
			profile:     "NONEXISTENT",
			expectError: true,
		},
		{
			name: "config missing required field - user",
			setupFunc: func() error {
				configContent := `[DEFAULT]
fingerprint=aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99
tenancy=ocid1.tenancy.oc1..testtenant
region=us-phoenix-1
key_file=~/.oci/test_key.pem
`
				return os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
			},
			profile:     "DEFAULT",
			expectError: true,
		},
		{
			name: "config missing required field - region",
			setupFunc: func() error {
				configContent := `[DEFAULT]
user=ocid1.user.oc1..testuser
fingerprint=aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99
tenancy=ocid1.tenancy.oc1..testtenant
key_file=~/.oci/test_key.pem
`
				return os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
			},
			profile:     "DEFAULT",
			expectError: true,
		},
		{
			name: "config with non-existent key file",
			setupFunc: func() error {
				configContent := `[DEFAULT]
user=ocid1.user.oc1..testuser
fingerprint=aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99
tenancy=ocid1.tenancy.oc1..testtenant
region=us-phoenix-1
key_file=~/.oci/nonexistent_key.pem
`
				return os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
			},
			profile:     "DEFAULT",
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			if tt.setupFunc != nil {
				if err := tt.setupFunc(); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}
			
			// Run readOCIConfig
			profile, err := readOCIConfig(tt.profile)
			
			// Check error expectation
			if (err != nil) != tt.expectError {
				t.Errorf("readOCIConfig() error = %v, expectError %v", err, tt.expectError)
				return
			}
			
			// If no error expected, validate the profile
			if !tt.expectError && tt.validateFunc != nil {
				tt.validateFunc(t, profile)
			}
		})
	}
}

func TestValidateOCIProfile(t *testing.T) {
	// Create a temporary directory for test key files
	tempDir := t.TempDir()
	
	// Save original HOME and restore after test
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	
	// Set HOME to temp directory
	os.Setenv("HOME", tempDir)
	
	// Create .oci directory
	ociDir := filepath.Join(tempDir, ".oci")
	err := os.MkdirAll(ociDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test .oci directory: %v", err)
	}
	
	// Create a test key file
	keyFilePath := filepath.Join(ociDir, "test_key.pem")
	err = os.WriteFile(keyFilePath, []byte("test key content"), 0600)
	if err != nil {
		t.Fatalf("Failed to create test key file: %v", err)
	}
	
	tests := []struct {
		name        string
		profile     *OCIProfile
		profileName string
		expectError bool
		errorSubstr string
	}{
		{
			name: "valid profile with all required fields",
			profile: &OCIProfile{
				UserID:      "ocid1.user.oc1..testuser",
				TenancyID:   "ocid1.tenancy.oc1..testtenant",
				Region:      "us-phoenix-1",
				Fingerprint: "aa:bb:cc:dd:ee:ff",
				KeyFile:     "~/.oci/test_key.pem",
			},
			profileName: "DEFAULT",
			expectError: false,
		},
		{
			name: "missing user field",
			profile: &OCIProfile{
				TenancyID:   "ocid1.tenancy.oc1..testtenant",
				Region:      "us-phoenix-1",
				Fingerprint: "aa:bb:cc:dd:ee:ff",
				KeyFile:     "~/.oci/test_key.pem",
			},
			profileName: "DEFAULT",
			expectError: true,
			errorSubstr: "user",
		},
		{
			name: "missing tenancy field",
			profile: &OCIProfile{
				UserID:      "ocid1.user.oc1..testuser",
				Region:      "us-phoenix-1",
				Fingerprint: "aa:bb:cc:dd:ee:ff",
				KeyFile:     "~/.oci/test_key.pem",
			},
			profileName: "DEFAULT",
			expectError: true,
			errorSubstr: "tenancy",
		},
		{
			name: "missing region field",
			profile: &OCIProfile{
				UserID:      "ocid1.user.oc1..testuser",
				TenancyID:   "ocid1.tenancy.oc1..testtenant",
				Fingerprint: "aa:bb:cc:dd:ee:ff",
				KeyFile:     "~/.oci/test_key.pem",
			},
			profileName: "DEFAULT",
			expectError: true,
			errorSubstr: "region",
		},
		{
			name: "missing fingerprint field",
			profile: &OCIProfile{
				UserID:    "ocid1.user.oc1..testuser",
				TenancyID: "ocid1.tenancy.oc1..testtenant",
				Region:    "us-phoenix-1",
				KeyFile:   "~/.oci/test_key.pem",
			},
			profileName: "DEFAULT",
			expectError: true,
			errorSubstr: "fingerprint",
		},
		{
			name: "missing key_file field",
			profile: &OCIProfile{
				UserID:      "ocid1.user.oc1..testuser",
				TenancyID:   "ocid1.tenancy.oc1..testtenant",
				Region:      "us-phoenix-1",
				Fingerprint: "aa:bb:cc:dd:ee:ff",
			},
			profileName: "DEFAULT",
			expectError: true,
			errorSubstr: "key_file",
		},
		{
			name: "multiple missing fields",
			profile: &OCIProfile{
				Region: "us-phoenix-1",
			},
			profileName: "DEFAULT",
			expectError: true,
			errorSubstr: "missing required fields",
		},
		{
			name: "key file does not exist",
			profile: &OCIProfile{
				UserID:      "ocid1.user.oc1..testuser",
				TenancyID:   "ocid1.tenancy.oc1..testtenant",
				Region:      "us-phoenix-1",
				Fingerprint: "aa:bb:cc:dd:ee:ff",
				KeyFile:     "~/.oci/nonexistent_key.pem",
			},
			profileName: "DEFAULT",
			expectError: true,
			errorSubstr: "key file not found",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateOCIProfile(tt.profile, tt.profileName)
			
			if (err != nil) != tt.expectError {
				t.Errorf("validateOCIProfile() error = %v, expectError %v", err, tt.expectError)
				return
			}
			
			if tt.expectError && tt.errorSubstr != "" {
				if !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tt.errorSubstr)) {
					t.Errorf("validateOCIProfile() error = %v, should contain %v", err.Error(), tt.errorSubstr)
				}
			}
		})
	}
}

func TestGetRegionFromOCIConfig(t *testing.T) {
	// Create a temporary directory for test config files
	tempDir := t.TempDir()
	
	// Save original HOME and restore after test
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	
	// Set HOME to temp directory
	os.Setenv("HOME", tempDir)
	
	// Create .oci directory
	ociDir := filepath.Join(tempDir, ".oci")
	err := os.MkdirAll(ociDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test .oci directory: %v", err)
	}
	
	// Create a test key file
	keyFilePath := filepath.Join(ociDir, "test_key.pem")
	err = os.WriteFile(keyFilePath, []byte("test key content"), 0600)
	if err != nil {
		t.Fatalf("Failed to create test key file: %v", err)
	}
	
	tests := []struct {
		name           string
		setupFunc      func() error
		profile        string
		expectedRegion string
		expectError    bool
	}{
		{
			name: "get region from DEFAULT profile",
			setupFunc: func() error {
				configContent := `[DEFAULT]
user=ocid1.user.oc1..testuser
fingerprint=aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99
tenancy=ocid1.tenancy.oc1..testtenant
region=us-phoenix-1
key_file=~/.oci/test_key.pem
`
				return os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
			},
			profile:        "DEFAULT",
			expectedRegion: "us-phoenix-1",
			expectError:    false,
		},
		{
			name: "get region from custom profile",
			setupFunc: func() error {
				configContent := `[DEFAULT]
user=ocid1.user.oc1..testuser
fingerprint=aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99
tenancy=ocid1.tenancy.oc1..testtenant
region=us-phoenix-1
key_file=~/.oci/test_key.pem

[CUSTOM]
user=ocid1.user.oc1..customuser
fingerprint=11:22:33:44:55:66:77:88:99:aa:bb:cc:dd:ee:ff:00
tenancy=ocid1.tenancy.oc1..customtenant
region=us-ashburn-1
key_file=~/.oci/test_key.pem
`
				return os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
			},
			profile:        "CUSTOM",
			expectedRegion: "us-ashburn-1",
			expectError:    false,
		},
		{
			name: "profile does not exist",
			setupFunc: func() error {
				configContent := `[DEFAULT]
user=ocid1.user.oc1..testuser
fingerprint=aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99
tenancy=ocid1.tenancy.oc1..testtenant
region=us-phoenix-1
key_file=~/.oci/test_key.pem
`
				return os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
			},
			profile:     "NONEXISTENT",
			expectError: true,
		},
		{
			name: "config file does not exist",
			setupFunc: func() error {
				configPath := filepath.Join(ociDir, "config")
				os.Remove(configPath)
				return nil
			},
			profile:     "DEFAULT",
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			if tt.setupFunc != nil {
				if err := tt.setupFunc(); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}
			
			// Run getRegionFromOCIConfig
			region, err := getRegionFromOCIConfig(tt.profile)
			
			// Check error expectation
			if (err != nil) != tt.expectError {
				t.Errorf("getRegionFromOCIConfig() error = %v, expectError %v", err, tt.expectError)
				return
			}
			
			// If no error expected, validate the region
			if !tt.expectError {
				if region != tt.expectedRegion {
					t.Errorf("getRegionFromOCIConfig() region = %v, expected %v", region, tt.expectedRegion)
				}
			}
		})
	}
}

func TestSelectCluster(t *testing.T) {
	tests := []struct {
		name           string
		clusters       []*OKEClusterInfo
		expectError    bool
		expectedID     string
		errorSubstr    string
	}{
		{
			name:        "empty cluster list",
			clusters:    []*OKEClusterInfo{},
			expectError: true,
			errorSubstr: "no clusters available",
		},
		{
			name: "single cluster auto-selection",
			clusters: []*OKEClusterInfo{
				{
					ID:                "ocid1.cluster.oc1.phx.test1",
					Name:              "test-cluster-1",
					Region:            "us-phoenix-1",
					KubernetesVersion: "v1.28.2",
					State:             "ACTIVE",
				},
			},
			expectError: false,
			expectedID:  "ocid1.cluster.oc1.phx.test1",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clusterID, err := selectCluster(tt.clusters)
			
			if (err != nil) != tt.expectError {
				t.Errorf("selectCluster() error = %v, expectError %v", err, tt.expectError)
				return
			}
			
			if tt.expectError && tt.errorSubstr != "" {
				if !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf("selectCluster() error = %v, should contain %v", err.Error(), tt.errorSubstr)
				}
			}
			
			if !tt.expectError && clusterID != tt.expectedID {
				t.Errorf("selectCluster() clusterID = %v, expected %v", clusterID, tt.expectedID)
			}
		})
	}
}

func TestParseJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expectError bool
		errorSubstr string
	}{
		{
			name:        "empty input",
			input:       []byte{},
			expectError: true,
			errorSubstr: "empty JSON response",
		},
		{
			name:        "non-JSON input",
			input:       []byte("this is not json"),
			expectError: true,
			errorSubstr: "does not appear to be JSON",
		},
		{
			name:        "valid JSON object",
			input:       []byte(`{"key": "value"}`),
			expectError: false,
		},
		{
			name:        "valid JSON array",
			input:       []byte(`[{"key": "value"}]`),
			expectError: false,
		},
		{
			name:        "invalid JSON structure",
			input:       []byte(`{"key": "value"`),
			expectError: true,
			errorSubstr: "failed to parse JSON",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			err := parseJSON(tt.input, &result)
			
			if (err != nil) != tt.expectError {
				t.Errorf("parseJSON() error = %v, expectError %v", err, tt.expectError)
				return
			}
			
			if tt.expectError && tt.errorSubstr != "" {
				if !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf("parseJSON() error = %v, should contain %v", err.Error(), tt.errorSubstr)
				}
			}
		})
	}
}

func TestValidateRegionFormat(t *testing.T) {
	tests := []struct {
		name        string
		region      string
		expectError bool
		errorSubstr string
	}{
		{
			name:        "valid region - us-phoenix-1",
			region:      "us-phoenix-1",
			expectError: false,
		},
		{
			name:        "valid region - us-ashburn-1",
			region:      "us-ashburn-1",
			expectError: false,
		},
		{
			name:        "valid region - eu-frankfurt-1",
			region:      "eu-frankfurt-1",
			expectError: false,
		},
		{
			name:        "valid region - ap-tokyo-1",
			region:      "ap-tokyo-1",
			expectError: false,
		},
		{
			name:        "valid region - uk-london-1",
			region:      "uk-london-1",
			expectError: false,
		},
		{
			name:        "empty region",
			region:      "",
			expectError: true,
			errorSubstr: "region cannot be empty",
		},
		{
			name:        "invalid format - missing number",
			region:      "us-phoenix",
			expectError: true,
			errorSubstr: "invalid region format",
		},
		{
			name:        "invalid format - only two parts",
			region:      "us-1",
			expectError: true,
			errorSubstr: "invalid region format",
		},
		{
			name:        "invalid format - no number at end",
			region:      "us-phoenix-a",
			expectError: true,
			errorSubstr: "must end with a number",
		},
		{
			name:        "invalid format - single word",
			region:      "phoenix",
			expectError: true,
			errorSubstr: "invalid region format",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRegionFormat(tt.region)
			
			if (err != nil) != tt.expectError {
				t.Errorf("validateRegionFormat() error = %v, expectError %v", err, tt.expectError)
				return
			}
			
			if tt.expectError && tt.errorSubstr != "" {
				if !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf("validateRegionFormat() error = %v, should contain %v", err.Error(), tt.errorSubstr)
				}
			}
		})
	}
}

func TestValidateClusterIDFormat(t *testing.T) {
	tests := []struct {
		name        string
		clusterID   string
		expectError bool
		errorSubstr string
	}{
		{
			name:        "valid cluster ID",
			clusterID:   "ocid1.cluster.oc1.phx.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			expectError: false,
		},
		{
			name:        "valid cluster ID - different region",
			clusterID:   "ocid1.cluster.oc1.iad.bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			expectError: false,
		},
		{
			name:        "valid cluster ID - short unique id",
			clusterID:   "ocid1.cluster.oc1.fra.abc123",
			expectError: false,
		},
		{
			name:        "empty cluster ID",
			clusterID:   "",
			expectError: true,
			errorSubstr: "cluster ID cannot be empty",
		},
		{
			name:        "invalid prefix - missing ocid1",
			clusterID:   "cluster.oc1.phx.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			expectError: true,
			errorSubstr: "invalid cluster ID format",
		},
		{
			name:        "invalid prefix - wrong resource type",
			clusterID:   "ocid1.instance.oc1.phx.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			expectError: true,
			errorSubstr: "invalid cluster ID format",
		},
		{
			name:        "invalid format - missing parts",
			clusterID:   "ocid1.cluster.oc1",
			expectError: true,
			errorSubstr: "invalid cluster ID format",
		},
		{
			name:        "invalid format - only prefix",
			clusterID:   "ocid1.cluster.oc1.",
			expectError: true,
			errorSubstr: "invalid cluster ID format",
		},
		{
			name:        "invalid format - random string",
			clusterID:   "not-a-valid-ocid",
			expectError: true,
			errorSubstr: "invalid cluster ID format",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateClusterIDFormat(tt.clusterID)
			
			if (err != nil) != tt.expectError {
				t.Errorf("validateClusterIDFormat() error = %v, expectError %v", err, tt.expectError)
				return
			}
			
			if tt.expectError && tt.errorSubstr != "" {
				if !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf("validateClusterIDFormat() error = %v, should contain %v", err.Error(), tt.errorSubstr)
				}
			}
		})
	}
}

func TestSelectClusterWithRetry(t *testing.T) {
	tests := []struct {
		name        string
		clusters    []*OKEClusterInfo
		expectError bool
		expectedID  string
		errorSubstr string
	}{
		{
			name:        "empty cluster list",
			clusters:    []*OKEClusterInfo{},
			expectError: true,
			errorSubstr: "no clusters available",
		},
		{
			name: "single cluster auto-selection",
			clusters: []*OKEClusterInfo{
				{
					ID:                "ocid1.cluster.oc1.phx.test1",
					Name:              "test-cluster-1",
					Region:            "us-phoenix-1",
					KubernetesVersion: "v1.28.2",
					State:             "ACTIVE",
				},
			},
			expectError: false,
			expectedID:  "ocid1.cluster.oc1.phx.test1",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clusterID, err := selectClusterWithRetry(tt.clusters)
			
			if (err != nil) != tt.expectError {
				t.Errorf("selectClusterWithRetry() error = %v, expectError %v", err, tt.expectError)
				return
			}
			
			if tt.expectError && tt.errorSubstr != "" {
				if !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf("selectClusterWithRetry() error = %v, should contain %v", err.Error(), tt.errorSubstr)
				}
			}
			
			if !tt.expectError && clusterID != tt.expectedID {
				t.Errorf("selectClusterWithRetry() clusterID = %v, expected %v", clusterID, tt.expectedID)
			}
		})
	}
}

func TestBackupKubeconfig(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(string) (string, error)
		expectError bool
		errorSubstr string
	}{
		{
			name: "backup existing kubeconfig",
			setupFunc: func(tempDir string) (string, error) {
				// Create a test kubeconfig file
				configPath := filepath.Join(tempDir, "kubeconfig.yaml")
				content := []byte("apiVersion: v1\nkind: Config\nclusters: []\ncontexts: []\nusers: []\n")
				err := os.WriteFile(configPath, content, 0600)
				return configPath, err
			},
			expectError: false,
		},
		{
			name: "backup non-existent file",
			setupFunc: func(tempDir string) (string, error) {
				// Return path to non-existent file
				configPath := filepath.Join(tempDir, "nonexistent.yaml")
				return configPath, nil
			},
			expectError: true,
			errorSubstr: "failed to read source file",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir := t.TempDir()
			
			// Setup test environment
			configPath, err := tt.setupFunc(tempDir)
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			
			// Run backup
			err = backupKubeconfig(configPath)
			
			// Check error expectation
			if (err != nil) != tt.expectError {
				t.Errorf("backupKubeconfig() error = %v, expectError %v", err, tt.expectError)
				return
			}
			
			if tt.expectError && tt.errorSubstr != "" {
				if !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf("backupKubeconfig() error = %v, should contain %v", err.Error(), tt.errorSubstr)
				}
				return
			}
			
			// If no error, verify backup was created
			if !tt.expectError {
				dir, file := filepath.Split(configPath)
				extension := filepath.Ext(file)
				baseFilename := file[:len(file)-len(extension)]
				bakLocation := filepath.Join(dir, baseFilename+".bak"+extension)
				
				if _, err := os.Stat(bakLocation); os.IsNotExist(err) {
					t.Errorf("backupKubeconfig() did not create backup file at %s", bakLocation)
				}
				
				// Verify backup content matches original
				originalContent, _ := os.ReadFile(configPath)
				backupContent, _ := os.ReadFile(bakLocation)
				
				if string(originalContent) != string(backupContent) {
					t.Errorf("backupKubeconfig() backup content does not match original")
				}
			}
		})
	}
}

func TestValidateKubeconfigGeneration(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(string) (string, error)
		expectError bool
		errorSubstr string
	}{
		{
			name: "valid kubeconfig file",
			setupFunc: func(tempDir string) (string, error) {
				configPath := filepath.Join(tempDir, "kubeconfig.yaml")
				content := []byte(`apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://example.com
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-context
users:
- name: test-user
  user:
    token: test-token
`)
				err := os.WriteFile(configPath, content, 0600)
				return configPath, err
			},
			expectError: false,
		},
		{
			name: "kubeconfig file does not exist",
			setupFunc: func(tempDir string) (string, error) {
				configPath := filepath.Join(tempDir, "nonexistent.yaml")
				return configPath, nil
			},
			expectError: true,
			errorSubstr: "was not created",
		},
		{
			name: "empty kubeconfig file",
			setupFunc: func(tempDir string) (string, error) {
				configPath := filepath.Join(tempDir, "kubeconfig.yaml")
				err := os.WriteFile(configPath, []byte(""), 0600)
				return configPath, err
			},
			expectError: true,
			errorSubstr: "is empty",
		},
		{
			name: "kubeconfig missing apiVersion field",
			setupFunc: func(tempDir string) (string, error) {
				configPath := filepath.Join(tempDir, "kubeconfig.yaml")
				content := []byte(`kind: Config
clusters: []
contexts: []
users: []
`)
				err := os.WriteFile(configPath, content, 0600)
				return configPath, err
			},
			expectError: true,
			errorSubstr: "missing required field: apiVersion:",
		},
		{
			name: "kubeconfig missing clusters field",
			setupFunc: func(tempDir string) (string, error) {
				configPath := filepath.Join(tempDir, "kubeconfig.yaml")
				content := []byte(`apiVersion: v1
kind: Config
contexts: []
users: []
`)
				err := os.WriteFile(configPath, content, 0600)
				return configPath, err
			},
			expectError: true,
			errorSubstr: "missing required field: clusters:",
		},
		{
			name: "kubeconfig missing contexts field",
			setupFunc: func(tempDir string) (string, error) {
				configPath := filepath.Join(tempDir, "kubeconfig.yaml")
				content := []byte(`apiVersion: v1
kind: Config
clusters: []
users: []
`)
				err := os.WriteFile(configPath, content, 0600)
				return configPath, err
			},
			expectError: true,
			errorSubstr: "missing required field: contexts:",
		},
		{
			name: "kubeconfig missing users field",
			setupFunc: func(tempDir string) (string, error) {
				configPath := filepath.Join(tempDir, "kubeconfig.yaml")
				content := []byte(`apiVersion: v1
kind: Config
clusters: []
contexts: []
`)
				err := os.WriteFile(configPath, content, 0600)
				return configPath, err
			},
			expectError: true,
			errorSubstr: "missing required field: users:",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir := t.TempDir()
			
			// Setup test environment
			configPath, err := tt.setupFunc(tempDir)
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			
			// Run validation
			err = validateKubeconfigGeneration(configPath)
			
			// Check error expectation
			if (err != nil) != tt.expectError {
				t.Errorf("validateKubeconfigGeneration() error = %v, expectError %v", err, tt.expectError)
				return
			}
			
			if tt.expectError && tt.errorSubstr != "" {
				if !strings.Contains(err.Error(), tt.errorSubstr) {
					t.Errorf("validateKubeconfigGeneration() error = %v, should contain %v", err.Error(), tt.errorSubstr)
				}
			}
		})
	}
}

// Integration Tests

func TestOKEConfigurationWorkflow(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(string) error
		clusterID   string
		region      string
		profile     string
		expectError bool
		errorSubstr string
	}{
		{
			name: "complete workflow with valid configuration",
			setupFunc: func(tempDir string) error {
				// Create OCI config directory and file
				ociDir := filepath.Join(tempDir, ".oci")
				if err := os.MkdirAll(ociDir, 0755); err != nil {
					return err
				}
				
				// Create key file
				keyFile := filepath.Join(ociDir, "test_key.pem")
				if err := os.WriteFile(keyFile, []byte("test key"), 0600); err != nil {
					return err
				}
				
				// Create config file
				configContent := `[DEFAULT]
user=ocid1.user.oc1..testuser
fingerprint=aa:bb:cc:dd:ee:ff
tenancy=ocid1.tenancy.oc1..testtenant
region=us-phoenix-1
key_file=` + keyFile + `
`
				return os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
			},
			clusterID:   "ocid1.cluster.oc1.phx.test123",
			region:      "us-phoenix-1",
			profile:     "DEFAULT",
			expectError: false,
		},
		{
			name: "workflow with missing OCI config",
			setupFunc: func(tempDir string) error {
				// Don't create OCI config
				return nil
			},
			clusterID:   "ocid1.cluster.oc1.phx.test123",
			region:      "us-phoenix-1",
			profile:     "DEFAULT",
			expectError: true,
			errorSubstr: "OCI config file not found",
		},
		{
			name: "workflow with invalid cluster ID",
			setupFunc: func(tempDir string) error {
				// Create valid OCI config
				ociDir := filepath.Join(tempDir, ".oci")
				if err := os.MkdirAll(ociDir, 0755); err != nil {
					return err
				}
				
				keyFile := filepath.Join(ociDir, "test_key.pem")
				if err := os.WriteFile(keyFile, []byte("test key"), 0600); err != nil {
					return err
				}
				
				configContent := `[DEFAULT]
user=ocid1.user.oc1..testuser
fingerprint=aa:bb:cc:dd:ee:ff
tenancy=ocid1.tenancy.oc1..testtenant
region=us-phoenix-1
key_file=` + keyFile + `
`
				return os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
			},
			clusterID:   "invalid-cluster-id",
			region:      "us-phoenix-1",
			profile:     "DEFAULT",
			expectError: true,
			errorSubstr: "invalid cluster ID format",
		},
		{
			name: "workflow with invalid region",
			setupFunc: func(tempDir string) error {
				// Create valid OCI config
				ociDir := filepath.Join(tempDir, ".oci")
				if err := os.MkdirAll(ociDir, 0755); err != nil {
					return err
				}
				
				keyFile := filepath.Join(ociDir, "test_key.pem")
				if err := os.WriteFile(keyFile, []byte("test key"), 0600); err != nil {
					return err
				}
				
				configContent := `[DEFAULT]
user=ocid1.user.oc1..testuser
fingerprint=aa:bb:cc:dd:ee:ff
tenancy=ocid1.tenancy.oc1..testtenant
region=us-phoenix-1
key_file=` + keyFile + `
`
				return os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
			},
			clusterID:   "ocid1.cluster.oc1.phx.test123",
			region:      "invalid",
			profile:     "DEFAULT",
			expectError: true,
			errorSubstr: "invalid region format",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir := t.TempDir()
			
			// Save and restore HOME
			originalHome := os.Getenv("HOME")
			defer os.Setenv("HOME", originalHome)
			os.Setenv("HOME", tempDir)
			
			// Setup test environment
			if err := tt.setupFunc(tempDir); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			
			// Test the workflow components
			
			// 1. Validate prerequisites
			err := validatePrerequisites(tt.profile)
			if tt.expectError && strings.Contains(tt.errorSubstr, "OCI config") {
				if err == nil {
					t.Error("validatePrerequisites() should fail with missing config")
				}
				return
			}
			
			// 2. Validate cluster ID format
			if tt.clusterID != "" {
				err = validateClusterIDFormat(tt.clusterID)
				if tt.expectError && strings.Contains(tt.errorSubstr, "cluster ID") {
					if err == nil {
						t.Error("validateClusterIDFormat() should fail with invalid cluster ID")
					}
					return
				}
			}
			
			// 3. Validate region format
			if tt.region != "" {
				err = validateRegionFormat(tt.region)
				if tt.expectError && strings.Contains(tt.errorSubstr, "region") {
					if err == nil {
						t.Error("validateRegionFormat() should fail with invalid region")
					}
					return
				}
			}
			
			// 4. Test kubeconfig generation (mock)
			if !tt.expectError {
				configPath := filepath.Join(tempDir, "kubeconfig.yaml")
				
				// Create a mock kubeconfig
				kubeconfigContent := []byte(`apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://test.oke.oraclecloud.com
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-context
users:
- name: test-user
  user:
    token: test-token
`)
				err = os.WriteFile(configPath, kubeconfigContent, 0600)
				if err != nil {
					t.Fatalf("Failed to create mock kubeconfig: %v", err)
				}
				
				// Validate the generated kubeconfig
				err = validateKubeconfigGeneration(configPath)
				if err != nil {
					t.Errorf("validateKubeconfigGeneration() failed: %v", err)
				}
			}
		})
	}
}

func TestCollectOKEParametersIntegration(t *testing.T) {
	tests := []struct {
		name           string
		setupFunc      func(string) error
		inputClusterID string
		inputRegion    string
		inputProfile   string
		expectError    bool
	}{
		{
			name: "collect parameters with all flags provided",
			setupFunc: func(tempDir string) error {
				// Create valid OCI config
				ociDir := filepath.Join(tempDir, ".oci")
				if err := os.MkdirAll(ociDir, 0755); err != nil {
					return err
				}
				
				keyFile := filepath.Join(ociDir, "test_key.pem")
				if err := os.WriteFile(keyFile, []byte("test key"), 0600); err != nil {
					return err
				}
				
				configContent := `[DEFAULT]
user=ocid1.user.oc1..testuser
fingerprint=aa:bb:cc:dd:ee:ff
tenancy=ocid1.tenancy.oc1..testtenant
region=us-phoenix-1
key_file=` + keyFile + `
`
				return os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
			},
			inputClusterID: "ocid1.cluster.oc1.phx.test123",
			inputRegion:    "us-phoenix-1",
			inputProfile:   "DEFAULT",
			expectError:    false,
		},
		{
			name: "collect parameters with custom profile",
			setupFunc: func(tempDir string) error {
				// Create valid OCI config with custom profile
				ociDir := filepath.Join(tempDir, ".oci")
				if err := os.MkdirAll(ociDir, 0755); err != nil {
					return err
				}
				
				keyFile := filepath.Join(ociDir, "test_key.pem")
				if err := os.WriteFile(keyFile, []byte("test key"), 0600); err != nil {
					return err
				}
				
				configContent := `[DEFAULT]
user=ocid1.user.oc1..testuser
fingerprint=aa:bb:cc:dd:ee:ff
tenancy=ocid1.tenancy.oc1..testtenant
region=us-phoenix-1
key_file=` + keyFile + `

[CUSTOM]
user=ocid1.user.oc1..customuser
fingerprint=11:22:33:44:55:66
tenancy=ocid1.tenancy.oc1..customtenant
region=us-ashburn-1
key_file=` + keyFile + `
`
				return os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
			},
			inputClusterID: "ocid1.cluster.oc1.iad.test456",
			inputRegion:    "us-ashburn-1",
			inputProfile:   "CUSTOM",
			expectError:    false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir := t.TempDir()
			
			// Save and restore HOME
			originalHome := os.Getenv("HOME")
			defer os.Setenv("HOME", originalHome)
			os.Setenv("HOME", tempDir)
			
			// Setup test environment
			if err := tt.setupFunc(tempDir); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			
			// Test parameter collection (without interactive prompts)
			// Since we're providing all parameters, collectOKEParameters should work
			clusterID, region, profile, err := collectOKEParameters(
				tt.inputClusterID,
				tt.inputRegion,
				tt.inputProfile,
			)
			
			if (err != nil) != tt.expectError {
				t.Errorf("collectOKEParameters() error = %v, expectError %v", err, tt.expectError)
				return
			}
			
			if !tt.expectError {
				if clusterID != tt.inputClusterID {
					t.Errorf("collectOKEParameters() clusterID = %v, want %v", clusterID, tt.inputClusterID)
				}
				if region != tt.inputRegion {
					t.Errorf("collectOKEParameters() region = %v, want %v", region, tt.inputRegion)
				}
				if profile != tt.inputProfile {
					t.Errorf("collectOKEParameters() profile = %v, want %v", profile, tt.inputProfile)
				}
			}
		})
	}
}

// TestOKEWorkflowIntegration tests the complete OKE configuration workflow
func TestOKEWorkflowIntegration(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()
	
	// Save original HOME and restore after test
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	
	// Set HOME to temp directory
	os.Setenv("HOME", tempDir)
	
	// Create .oci directory with valid config
	ociDir := filepath.Join(tempDir, ".oci")
	err := os.MkdirAll(ociDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test .oci directory: %v", err)
	}
	
	// Create a test key file
	keyFilePath := filepath.Join(ociDir, "test_key.pem")
	err = os.WriteFile(keyFilePath, []byte("test key content"), 0600)
	if err != nil {
		t.Fatalf("Failed to create test key file: %v", err)
	}
	
	// Create OCI config file
	configContent := `[DEFAULT]
user=ocid1.user.oc1..testuser
fingerprint=aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99
tenancy=ocid1.tenancy.oc1..testtenant
region=us-phoenix-1
key_file=~/.oci/test_key.pem

[CUSTOM]
user=ocid1.user.oc1..customuser
fingerprint=11:22:33:44:55:66:77:88:99:aa:bb:cc:dd:ee:ff:00
tenancy=ocid1.tenancy.oc1..customtenant
region=us-ashburn-1
key_file=~/.oci/test_key.pem
`
	err = os.WriteFile(filepath.Join(ociDir, "config"), []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	
	t.Run("complete workflow with valid configuration", func(t *testing.T) {
		// Test 1: Validate prerequisites
		err := validatePrerequisites("DEFAULT")
		// This may fail if OCI CLI is not installed, which is expected in CI/CD
		if err != nil {
			if !strings.Contains(err.Error(), "OCI CLI not found") {
				t.Logf("Prerequisites validation: %v (expected in environments without OCI CLI)", err)
			}
		}
		
		// Test 2: Read OCI config
		profile, err := readOCIConfig("DEFAULT")
		if err != nil {
			t.Fatalf("readOCIConfig() failed: %v", err)
		}
		
		if profile.Region != "us-phoenix-1" {
			t.Errorf("Expected region 'us-phoenix-1', got '%s'", profile.Region)
		}
		
		// Test 3: Get region from OCI config
		region, err := getRegionFromOCIConfig("DEFAULT")
		if err != nil {
			t.Fatalf("getRegionFromOCIConfig() failed: %v", err)
		}
		
		if region != "us-phoenix-1" {
			t.Errorf("Expected region 'us-phoenix-1', got '%s'", region)
		}
		
		// Test 4: Validate region format
		err = validateRegionFormat(region)
		if err != nil {
			t.Errorf("validateRegionFormat() failed for valid region: %v", err)
		}
		
		// Test 5: Validate cluster ID format
		testClusterID := "ocid1.cluster.oc1.phx.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		err = validateClusterIDFormat(testClusterID)
		if err != nil {
			t.Errorf("validateClusterIDFormat() failed for valid cluster ID: %v", err)
		}
		
		// Test 6: Test kubeconfig backup and validation workflow
		kubeconfigPath := filepath.Join(tempDir, "kubeconfig.yaml")
		validKubeconfig := []byte(`apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0t
    server: https://test-cluster.us-phoenix-1.clusters.oci.oraclecloud.com:6443
  name: cluster-test
contexts:
- context:
    cluster: cluster-test
    user: user-test
  name: context-test
current-context: context-test
users:
- name: user-test
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: oci
      args:
      - ce
      - cluster
      - generate-token
      - --cluster-id
      - ocid1.cluster.oc1.phx.test
`)
		
		// Create initial kubeconfig
		err = os.WriteFile(kubeconfigPath, validKubeconfig, 0600)
		if err != nil {
			t.Fatalf("Failed to create test kubeconfig: %v", err)
		}
		
		// Test backup mechanism
		err = backupKubeconfig(kubeconfigPath)
		if err != nil {
			t.Errorf("backupKubeconfig() failed: %v", err)
		}
		
		// Verify backup was created
		backupPath := filepath.Join(tempDir, "kubeconfig.bak.yaml")
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			t.Errorf("Backup file was not created at %s", backupPath)
		}
		
		// Test validation
		err = validateKubeconfigGeneration(kubeconfigPath)
		if err != nil {
			t.Errorf("validateKubeconfigGeneration() failed for valid kubeconfig: %v", err)
		}
	})
	
	t.Run("workflow with missing OCI config", func(t *testing.T) {
		// Create a separate temp directory for this test
		testTempDir := t.TempDir()
		
		// Save and restore HOME for this subtest
		originalTestHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalTestHome)
		os.Setenv("HOME", testTempDir)
		
		// Test should fail when trying to read config (no config file exists)
		_, err := readOCIConfig("DEFAULT")
		if err == nil {
			t.Error("readOCIConfig() should fail with missing config file")
		}
		
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("Error should mention config not found, got: %v", err)
		}
	})
	
	t.Run("workflow with invalid cluster ID", func(t *testing.T) {
		invalidClusterIDs := []string{
			"",
			"invalid-id",
			"ocid1.instance.oc1.phx.test",
			"cluster-123",
		}
		
		for _, clusterID := range invalidClusterIDs {
			err := validateClusterIDFormat(clusterID)
			if err == nil {
				t.Errorf("validateClusterIDFormat() should fail for invalid cluster ID: %s", clusterID)
			}
		}
	})
	
	t.Run("workflow with invalid region", func(t *testing.T) {
		invalidRegions := []string{
			"",
			"invalid",
			"us-phoenix",
			"phoenix-1",
		}
		
		for _, region := range invalidRegions {
			err := validateRegionFormat(region)
			if err == nil {
				t.Errorf("validateRegionFormat() should fail for invalid region: %s", region)
			}
		}
	})
	
	t.Run("collect parameters with all flags provided", func(t *testing.T) {
		// Create a separate temp directory for this test
		testTempDir := t.TempDir()
		
		// Save and restore HOME for this subtest
		originalTestHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalTestHome)
		os.Setenv("HOME", testTempDir)
		
		// Create .oci directory with valid config for this test
		testOciDir := filepath.Join(testTempDir, ".oci")
		err := os.MkdirAll(testOciDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test .oci directory: %v", err)
		}
		
		// Create a test key file
		testKeyFilePath := filepath.Join(testOciDir, "test_key.pem")
		err = os.WriteFile(testKeyFilePath, []byte("test key content"), 0600)
		if err != nil {
			t.Fatalf("Failed to create test key file: %v", err)
		}
		
		// Create OCI config file
		testConfigContent := `[DEFAULT]
user=ocid1.user.oc1..testuser
fingerprint=aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99
tenancy=ocid1.tenancy.oc1..testtenant
region=us-phoenix-1
key_file=~/.oci/test_key.pem
`
		err = os.WriteFile(filepath.Join(testOciDir, "config"), []byte(testConfigContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}
		
		clusterID := "ocid1.cluster.oc1.phx.test123"
		region := "us-phoenix-1"
		profile := "DEFAULT"
		
		// When all parameters are provided, collectOKEParameters should validate and return them
		// Note: This will try to list clusters which requires OCI CLI, so we just test the validation parts
		
		// Test region validation
		err = validateRegionFormat(region)
		if err != nil {
			t.Errorf("validateRegionFormat() failed: %v", err)
		}
		
		// Test cluster ID validation
		err = validateClusterIDFormat(clusterID)
		if err != nil {
			t.Errorf("validateClusterIDFormat() failed: %v", err)
		}
		
		// Test profile reading
		ociProfile, err := readOCIConfig(profile)
		if err != nil {
			t.Errorf("readOCIConfig() failed: %v", err)
		}
		
		if ociProfile.Region != region {
			t.Errorf("Expected region %s, got %s", region, ociProfile.Region)
		}
	})
	
	t.Run("collect parameters with custom profile", func(t *testing.T) {
		// Create a separate temp directory for this test
		testTempDir := t.TempDir()
		
		// Save and restore HOME for this subtest
		originalTestHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalTestHome)
		os.Setenv("HOME", testTempDir)
		
		// Create .oci directory with valid config for this test
		testOciDir := filepath.Join(testTempDir, ".oci")
		err := os.MkdirAll(testOciDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create test .oci directory: %v", err)
		}
		
		// Create a test key file
		testKeyFilePath := filepath.Join(testOciDir, "test_key.pem")
		err = os.WriteFile(testKeyFilePath, []byte("test key content"), 0600)
		if err != nil {
			t.Fatalf("Failed to create test key file: %v", err)
		}
		
		// Create OCI config file with custom profile
		testConfigContent := `[DEFAULT]
user=ocid1.user.oc1..testuser
fingerprint=aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99
tenancy=ocid1.tenancy.oc1..testtenant
region=us-phoenix-1
key_file=~/.oci/test_key.pem

[CUSTOM]
user=ocid1.user.oc1..customuser
fingerprint=11:22:33:44:55:66:77:88:99:aa:bb:cc:dd:ee:ff:00
tenancy=ocid1.tenancy.oc1..customtenant
region=us-ashburn-1
key_file=~/.oci/test_key.pem
`
		err = os.WriteFile(filepath.Join(testOciDir, "config"), []byte(testConfigContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config file: %v", err)
		}
		
		profile := "CUSTOM"
		
		// Test reading custom profile
		ociProfile, err := readOCIConfig(profile)
		if err != nil {
			t.Errorf("readOCIConfig() failed for custom profile: %v", err)
		}
		
		if ociProfile.Region != "us-ashburn-1" {
			t.Errorf("Expected region us-ashburn-1 for custom profile, got %s", ociProfile.Region)
		}
		
		// Test getting region from custom profile
		region, err := getRegionFromOCIConfig(profile)
		if err != nil {
			t.Errorf("getRegionFromOCIConfig() failed for custom profile: %v", err)
		}
		
		if region != "us-ashburn-1" {
			t.Errorf("Expected region us-ashburn-1, got %s", region)
		}
	})
}
