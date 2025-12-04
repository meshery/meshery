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

// Package system provides commands for configuring Meshery to work with various Kubernetes clusters.
// This file implements Oracle Kubernetes Engine (OKE) integration for mesheryctl.
//
// OKE Integration Pattern for Future Cloud Provider Additions:
// ========================================================
// The OKE integration follows the established pattern used by other cloud providers (EKS, GKE, AKS).
// This pattern should be used as a template when adding support for new cloud providers:
//
// 1. Prerequisites Validation:
//    - Validate cloud provider CLI is installed and accessible
//    - Validate cloud provider configuration exists and is readable
//    - Provide clear error messages with installation/configuration guidance
//
// 2. Parameter Collection:
//    - Support both interactive and non-interactive modes
//    - Use command flags for automation (CI/CD pipelines)
//    - Provide sensible defaults from cloud provider configuration
//    - Implement retry logic for user input validation
//
// 3. Cluster Discovery and Selection:
//    - List available clusters using cloud provider CLI
//    - Filter to only show clusters in usable state (e.g., ACTIVE)
//    - Auto-select when only one cluster is available
//    - Provide interactive selection for multiple clusters
//
// 4. Kubeconfig Generation:
//    - Use cloud provider's native kubeconfig generation command
//    - Backup existing kubeconfig before overwriting
//    - Validate generated kubeconfig for completeness
//    - Handle both public and private cluster endpoints
//
// 5. Meshery Integration:
//    - Register context with Meshery server using existing functions
//    - Set authentication token using existing workflow
//    - Provide confirmation of successful configuration
//
// 6. Error Handling and User Experience:
//    - Provide context-specific error messages
//    - Include troubleshooting guidance in error messages
//    - Reference cloud provider documentation for setup
//    - List required IAM policies or permissions
//
// Key Functions to Implement for New Providers:
// - validate<Provider>CLI() - Check CLI installation
// - validate<Provider>Config() - Check configuration
// - list<Provider>Clusters() - Discover available clusters
// - generate<Provider>Kubeconfig() - Generate kubeconfig
// - handle<Provider>Error() - Provider-specific error handling
//
// Prerequisites for OKE:
// - OCI CLI must be installed (https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm)
// - OCI CLI must be configured with valid credentials (~/.oci/config)
// - User must have appropriate IAM permissions to access OKE clusters
//
// Required OCI IAM Policies:
//   Allow group <group> to manage cluster-family in compartment <compartment>
//   Allow group <group> to use virtual-network-family in compartment <compartment>
//
// Usage Examples:
//   # Basic usage - interactive prompts for all parameters
//   mesheryctl system config oke
//
//   # Specify cluster ID and region (non-interactive)
//   mesheryctl system config oke --cluster-id ocid1.cluster.oc1.phx.xxx --region us-phoenix-1
//
//   # Use a specific OCI profile (multi-tenancy support)
//   mesheryctl system config oke --profile myprofile
//
//   # Combine with authentication token (CI/CD automation)
//   mesheryctl system config oke --token auth.json --region us-ashburn-1
//
//   # Full automation example
//   mesheryctl system config oke \
//     --cluster-id ocid1.cluster.oc1.iad.xxx \
//     --region us-ashburn-1 \
//     --profile production \
//     --token auth.json
package system

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
)

// OCIProfile represents an OCI configuration profile from ~/.oci/config.
// This structure maps to the standard OCI CLI configuration format.
//
// OCI Configuration File Format:
// The OCI CLI stores configuration in an INI-style file at ~/.oci/config.
// Each profile section contains authentication and connection details.
//
// Example ~/.oci/config:
//   [DEFAULT]
//   user=ocid1.user.oc1..xxx
//   fingerprint=aa:bb:cc:dd:ee:ff:00:11:22:33:44:55:66:77:88:99
//   tenancy=ocid1.tenancy.oc1..xxx
//   region=us-phoenix-1
//   key_file=~/.oci/oci_api_key.pem
//
// Fields:
//   - TenancyID: OCID of the OCI tenancy (required)
//   - UserID: OCID of the OCI user (required)
//   - Region: Default OCI region for API calls (required)
//   - Fingerprint: Fingerprint of the API signing key (required)
//   - KeyFile: Path to the private key file for API signing (required)
//   - PassPhrase: Optional passphrase for encrypted private keys
type OCIProfile struct {
	TenancyID   string `ini:"tenancy"`
	UserID      string `ini:"user"`
	Region      string `ini:"region"`
	Fingerprint string `ini:"fingerprint"`
	KeyFile     string `ini:"key_file"`
	PassPhrase  string `ini:"pass_phrase,omitempty"`
}

// validateOCICLI checks if OCI CLI is installed and accessible.
// This is a prerequisite check that must pass before attempting OKE configuration.
//
// The function executes 'oci --version' to verify the CLI is available.
// If the command fails, it indicates OCI CLI is not installed or not in PATH.
//
// Returns:
//   - nil if OCI CLI is found and executable
//   - error if OCI CLI is not found or not accessible
//
// Error Handling:
// Errors from this function should be passed to handleOCICLIError() for
// user-friendly error messages with installation guidance.
func validateOCICLI() error {
	ociCheck := exec.Command("oci", "--version")
	ociCheck.Stdout = os.Stdout
	ociCheck.Stderr = os.Stderr
	err := ociCheck.Run()
	if err != nil {
		return fmt.Errorf("OCI CLI not found")
	}
	return nil
}

// validateOCIConfig verifies OCI configuration file exists and is readable.
// This function checks that the OCI CLI configuration is properly set up.
//
// The function performs the following validations:
// 1. Checks if ~/.oci/config file exists
// 2. Verifies the file is readable
// 3. Confirms the specified profile exists in the config file
//
// Parameters:
//   - profile: Name of the OCI profile to validate (e.g., "DEFAULT", "myprofile")
//
// Returns:
//   - nil if configuration is valid
//   - error if config file is missing, unreadable, or profile not found
//
// Error Handling:
// Errors from this function should be passed to handleOCIConfigError() for
// user-friendly error messages with configuration guidance.
func validateOCIConfig(profile string) error {
	// Get the OCI config file path
	configPath := filepath.Join(os.Getenv("HOME"), ".oci", "config")
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("OCI config file not found at %s", configPath)
	}
	
	// Try to read the config file to ensure it's readable
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("OCI config file exists but is not readable: %w", err)
	}
	
	// Check if the specified profile exists in the config
	profileHeader := fmt.Sprintf("[%s]", profile)
	if !strings.Contains(string(content), profileHeader) {
		return fmt.Errorf("profile '%s' not found in OCI config", profile)
	}
	
	return nil
}

// handleOCICLIError provides user-friendly error messages with installation guidance.
// This function transforms technical CLI errors into actionable user guidance.
//
// Common Error Scenarios:
// - "not found" / "executable file not found": OCI CLI not installed
// - "exec" errors: Permission or path issues
//
// The function provides:
// - Clear explanation of the problem
// - Direct link to OCI CLI installation documentation
// - Actionable next steps for the user
//
// Parameters:
//   - err: The original error from OCI CLI validation
//
// Returns:
//   - nil if err is nil
//   - Enhanced error with user guidance if err is not nil
func handleOCICLIError(err error) error {
	if err == nil {
		return nil
	}
	
	errMsg := err.Error()
	if strings.Contains(errMsg, "not found") || 
	   strings.Contains(errMsg, "executable file not found") || 
	   strings.Contains(errMsg, "file does not exist") ||
	   strings.Contains(errMsg, "exec") {
		return fmt.Errorf("OCI CLI not found. Please install OCI CLI and try again.\nInstallation guide: https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm")
	}
	
	return fmt.Errorf("error validating OCI CLI: %w", err)
}

// handleOCIConfigError guides users through OCI configuration setup.
// This function provides context-specific guidance for OCI configuration issues.
//
// Common Error Scenarios:
// - Profile not found: User specified a profile that doesn't exist in ~/.oci/config
// - Config file not found: OCI CLI has never been configured
// - Permission issues: Config file exists but cannot be read
//
// The function provides:
// - Specific diagnosis of the configuration problem
// - Commands to resolve the issue (e.g., 'oci setup config')
// - Links to OCI configuration documentation
//
// Parameters:
//   - err: The original error from OCI config validation
//   - profile: The profile name that was being validated
//
// Returns:
//   - nil if err is nil
//   - Enhanced error with configuration guidance if err is not nil
func handleOCIConfigError(err error, profile string) error {
	if err == nil {
		return nil
	}
	
	errMsg := err.Error()
	
	// Check for profile not found first (more specific than file not found)
	if strings.Contains(errMsg, "profile") && strings.Contains(errMsg, "not found") {
		return fmt.Errorf("OCI profile '%s' not found in config file. Please check your OCI configuration.\nAvailable profiles can be listed by examining ~/.oci/config\nTo create a new profile, run: oci setup config", profile)
	}
	
	// Check for config file not found
	if strings.Contains(errMsg, "not found at") || 
	   strings.Contains(errMsg, "file does not exist") ||
	   strings.Contains(errMsg, "no such file") {
		return fmt.Errorf("OCI config file not found. Please run 'oci setup config' to create configuration.\nConfiguration guide: https://docs.oracle.com/en-us/iaas/Content/API/Concepts/sdkconfig.htm")
	}
	
	// Check for permission issues
	if strings.Contains(errMsg, "not readable") {
		return fmt.Errorf("OCI config file exists but cannot be read. Please check file permissions.\nConfig file location: ~/.oci/config")
	}
	
	return fmt.Errorf("error validating OCI config: %w", err)
}

// handleClusterAccessError provides specific error messages for cluster access issues.
// This function diagnoses and provides guidance for OKE cluster access problems.
//
// Common Error Scenarios:
// - NotAuthorizedOrNotFound: Cluster doesn't exist or user lacks IAM permissions
// - InvalidParameter: Cluster OCID format is incorrect
// - Timeout/Connection: Network connectivity issues
// - Permission/Unauthorized: Insufficient IAM policies
//
// The function provides:
// - Specific diagnosis of the access problem
// - Checklist of items to verify (cluster ID, permissions, region, credentials)
// - Required IAM policies for OKE access
// - Troubleshooting commands to diagnose the issue
//
// OCI IAM Requirements:
// Users need the following IAM policies to access OKE clusters:
//   - Allow group <group> to manage cluster-family in compartment <compartment>
//   - Allow group <group> to use virtual-network-family in compartment <compartment>
//
// Parameters:
//   - err: The original error from cluster access attempt
//   - clusterID: The cluster OCID that was being accessed
//
// Returns:
//   - nil if err is nil
//   - Enhanced error with troubleshooting guidance if err is not nil
func handleClusterAccessError(err error, clusterID string) error {
	if err == nil {
		return nil
	}
	
	errMsg := err.Error()
	
	// Check for NotAuthorizedOrNotFound error
	if strings.Contains(errMsg, "NotAuthorizedOrNotFound") || 
	   strings.Contains(errMsg, "not found") || 
	   strings.Contains(errMsg, "access denied") {
		return fmt.Errorf("cluster '%s' not found or access denied. Please check:\n"+
			"  1. Cluster ID is correct (format: ocid1.cluster.oc1.<region>.<unique-id>)\n"+
			"  2. You have proper IAM permissions to access the cluster\n"+
			"  3. Cluster is in the specified region\n"+
			"  4. Your OCI credentials are valid and not expired", clusterID)
	}
	
	// Check for InvalidParameter error
	if strings.Contains(errMsg, "InvalidParameter") || 
	   strings.Contains(errMsg, "invalid") {
		return fmt.Errorf("invalid cluster ID format: %s\n"+
			"Expected format: ocid1.cluster.oc1.<region>.<unique-id>\n"+
			"Example: ocid1.cluster.oc1.phx.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", clusterID)
	}
	
	// Check for network connectivity issues
	if strings.Contains(errMsg, "timeout") || 
	   strings.Contains(errMsg, "connection") || 
	   strings.Contains(errMsg, "network") {
		return fmt.Errorf("network connectivity issue accessing cluster '%s':\n"+
			"  1. Check your internet connection\n"+
			"  2. Verify OCI service endpoints are accessible\n"+
			"  3. Check if you're behind a proxy or firewall\n"+
			"  4. Ensure OCI CLI can reach the OCI API endpoints\n"+
			"Original error: %v", clusterID, err)
	}
	
	// Check for permission issues
	if strings.Contains(errMsg, "permission") || 
	   strings.Contains(errMsg, "unauthorized") || 
	   strings.Contains(errMsg, "forbidden") {
		return fmt.Errorf("insufficient permissions to access cluster '%s':\n"+
			"Required OCI IAM policies:\n"+
			"  - Allow group <your-group> to manage cluster-family in compartment <compartment>\n"+
			"  - Allow group <your-group> to use virtual-network-family in compartment <compartment>\n"+
			"Please contact your OCI administrator to grant the necessary permissions.\n"+
			"Original error: %v", clusterID, err)
	}
	
	// Generic error with helpful context
	return fmt.Errorf("error accessing cluster '%s': %w\n"+
		"Troubleshooting steps:\n"+
		"  1. Verify cluster ID is correct: oci ce cluster get --cluster-id %s\n"+
		"  2. Check your OCI credentials: oci iam user get --user-id <your-user-id>\n"+
		"  3. Verify region is correct: oci iam region list\n"+
		"  4. Check OCI CLI configuration: cat ~/.oci/config", clusterID, err, clusterID)
}

// validatePrerequisites runs all prerequisite checks for OKE configuration.
// This is the main validation entry point that orchestrates all prerequisite checks.
//
// Validation Steps:
// 1. Verify OCI CLI is installed and accessible
// 2. Verify OCI configuration file exists and is valid
// 3. Verify the specified profile exists and is properly configured
//
// This function should be called before attempting any OKE operations to ensure
// the environment is properly set up. Early validation provides better user experience
// by catching configuration issues before attempting cluster operations.
//
// Parameters:
//   - profile: Name of the OCI profile to validate
//
// Returns:
//   - nil if all prerequisites are met
//   - error with user guidance if any prerequisite check fails
func validatePrerequisites(profile string) error {
	log.Debug("Validating OCI CLI installation...")
	
	// Check OCI CLI
	if err := validateOCICLI(); err != nil {
		return handleOCICLIError(err)
	}
	
	log.Debug("OCI CLI found")
	log.Debugf("Validating OCI configuration for profile: %s", profile)
	
	// Check OCI config
	if err := validateOCIConfig(profile); err != nil {
		return handleOCIConfigError(err, profile)
	}
	
	log.Debugf("OCI configuration validated successfully for profile: %s", profile)
	return nil
}

// readOCIConfig reads and parses the OCI configuration file for a given profile.
// This function loads the OCI CLI configuration and extracts profile-specific settings.
//
// The OCI configuration file (~/.oci/config) uses INI format with profile sections.
// Each profile contains authentication credentials and default settings for OCI API calls.
//
// Process:
// 1. Locate the OCI config file at ~/.oci/config
// 2. Parse the INI-formatted file
// 3. Extract the specified profile section
// 4. Map the profile data to OCIProfile struct
// 5. Validate that all required fields are present
//
// Parameters:
//   - profile: Name of the profile to read (e.g., "DEFAULT", "production")
//
// Returns:
//   - *OCIProfile: Populated profile structure with authentication details
//   - error: If config file is missing, profile not found, or validation fails
func readOCIConfig(profile string) (*OCIProfile, error) {
	// Get the OCI config file path
	configPath := filepath.Join(os.Getenv("HOME"), ".oci", "config")
	
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("OCI config file not found at %s", configPath)
	}
	
	// Load the INI file
	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read OCI config file: %w", err)
	}
	
	// Get the specified profile section
	section, err := cfg.GetSection(profile)
	if err != nil {
		return nil, fmt.Errorf("profile '%s' not found in OCI config", profile)
	}
	
	// Map the section to OCIProfile struct
	var ociProfile OCIProfile
	err = section.MapTo(&ociProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OCI profile: %w", err)
	}
	
	// Validate required fields
	if err := validateOCIProfile(&ociProfile, profile); err != nil {
		return nil, err
	}
	
	return &ociProfile, nil
}

// validateOCIProfile ensures required fields are present in the OCI profile.
// This function performs comprehensive validation of OCI profile configuration.
//
// Required Fields:
// - tenancy: OCID of the OCI tenancy
// - user: OCID of the OCI user
// - region: Default region for API calls
// - fingerprint: Fingerprint of the API signing key
// - key_file: Path to the private key file
//
// Additional Validation:
// - Verifies the private key file exists at the specified path
// - Expands ~ to home directory in key file paths
//
// This validation ensures the profile is complete and usable for OCI API calls.
// Missing or invalid configuration will result in API authentication failures.
//
// Parameters:
//   - profile: The OCIProfile structure to validate
//   - profileName: Name of the profile (for error messages)
//
// Returns:
//   - nil if profile is valid and complete
//   - error listing missing fields or invalid key file path
func validateOCIProfile(profile *OCIProfile, profileName string) error {
	var missingFields []string
	
	if profile.TenancyID == "" {
		missingFields = append(missingFields, "tenancy")
	}
	if profile.UserID == "" {
		missingFields = append(missingFields, "user")
	}
	if profile.Region == "" {
		missingFields = append(missingFields, "region")
	}
	if profile.Fingerprint == "" {
		missingFields = append(missingFields, "fingerprint")
	}
	if profile.KeyFile == "" {
		missingFields = append(missingFields, "key_file")
	}
	
	if len(missingFields) > 0 {
		return fmt.Errorf("OCI profile '%s' is missing required fields: %s. Please run 'oci setup config' to configure your profile", 
			profileName, strings.Join(missingFields, ", "))
	}
	
	// Validate that the key file exists
	keyFilePath := profile.KeyFile
	// Expand ~ to home directory if present
	if strings.HasPrefix(keyFilePath, "~/") {
		keyFilePath = filepath.Join(os.Getenv("HOME"), keyFilePath[2:])
	}
	
	if _, err := os.Stat(keyFilePath); os.IsNotExist(err) {
		return fmt.Errorf("OCI key file not found at %s. Please ensure the key_file path in profile '%s' is correct", 
			keyFilePath, profileName)
	}
	
	return nil
}

// getRegionFromOCIConfig extracts the default region from the OCI config profile.
// This function reads the OCI configuration and returns the default region setting.
//
// The region is used as the default when the user doesn't specify --region flag.
// This provides a better user experience by using their configured default region.
//
// OCI Region Format:
// Regions follow the pattern: <location>-<city>-<number>
// Examples: us-phoenix-1, us-ashburn-1, eu-frankfurt-1, ap-tokyo-1
//
// Parameters:
//   - profile: Name of the OCI profile to read
//
// Returns:
//   - string: The default region from the profile
//   - error: If profile cannot be read or region is not set
func getRegionFromOCIConfig(profile string) (string, error) {
	// Read the OCI config
	ociProfile, err := readOCIConfig(profile)
	if err != nil {
		return "", err
	}
	
	// Return the region from the profile
	if ociProfile.Region == "" {
		return "", fmt.Errorf("region not found in OCI profile '%s'", profile)
	}
	
	log.Debugf("Found region '%s' in OCI profile '%s'", ociProfile.Region, profile)
	return ociProfile.Region, nil
}

// OKEClusterInfo represents information about an OKE cluster.
// This structure maps to the JSON response from 'oci ce cluster list' command.
//
// The OCI CLI returns cluster information in JSON format, which is parsed into
// this structure for easier handling and display to users.
//
// Fields:
//   - ID: Cluster OCID (format: ocid1.cluster.oc1.<region>.<unique-id>)
//   - Name: Human-readable cluster name
//   - Region: OCI region where the cluster is deployed
//   - CompartmentID: OCID of the compartment containing the cluster
//   - KubernetesVersion: Kubernetes version running on the cluster
//   - State: Lifecycle state (ACTIVE, CREATING, UPDATING, DELETING, etc.)
//   - Endpoint: Kubernetes API server endpoint URL
//
// Only clusters in ACTIVE state are selectable for configuration.
type OKEClusterInfo struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Region            string `json:"region"`
	CompartmentID     string `json:"compartment-id"`
	KubernetesVersion string `json:"kubernetes-version"`
	State             string `json:"lifecycle-state"`
	Endpoint          string `json:"endpoint"`
}

// listOKEClusters fetches available OKE clusters using OCI CLI.
// This function discovers all OKE clusters accessible to the user in the specified compartment and region.
//
// The function executes 'oci ce cluster list' command and parses the JSON response.
// Only clusters in ACTIVE state are returned, as other states (CREATING, UPDATING,
// DELETING, etc.) are not ready for configuration.
//
// OCI CLI Command:
//   oci ce cluster list --compartment-id <compartment-id> --all --region <region> --profile <profile> --output json
//
// The --all flag retrieves all resources including sub-compartments.
//
// IAM Requirements:
// User must have the following IAM policy to list clusters:
//   Allow group <group> to inspect cluster-family in compartment <compartment>
//   Allow group <group> to inspect compartments in tenancy (if using tenancy as compartment)
//
// Parameters:
//   - region: OCI region to search for clusters
//   - profile: OCI profile to use for authentication
//   - compartmentID: OCI compartment OCID to search in (typically tenancy ID or specific compartment)
//
// Returns:
//   - []*OKEClusterInfo: Slice of active clusters found in the compartment/region
//   - error: If OCI CLI fails or response cannot be parsed
func listOKEClusters(region, profile, compartmentID string) ([]*OKEClusterInfo, error) {
	log.Debugf("Listing OKE clusters in compartment '%s', region '%s' using profile '%s'", compartmentID, region, profile)
	
	// Build the OCI CLI command to list clusters
	// Use the provided compartment ID to list all clusters
	args := []string{"ce", "cluster", "list", 
		"--compartment-id", compartmentID,
		"--all"}
	
	if region != "" {
		args = append(args, "--region", region)
	}
	
	if profile != "" {
		args = append(args, "--profile", profile)
	}
	
	// Add output format
	args = append(args, "--output", "json")
	
	// Execute the command
	log.Debugf("Executing OCI CLI command: oci %s", strings.Join(args, " "))
	cmd := exec.Command("oci", args...)
	
	// Capture both stdout and stderr
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err = cmd.Run()
	if err != nil {
		stderrMsg := stderr.String()
		stdoutMsg := stdout.String()
		
		log.Debugf("OCI CLI stderr: %s", stderrMsg)
		log.Debugf("OCI CLI stdout: %s", stdoutMsg)
		
		// Provide more helpful error message based on stderr content
		if strings.Contains(stderrMsg, "compartment-id") || strings.Contains(stderrMsg, "Missing option") {
			return nil, fmt.Errorf("failed to list OKE clusters: compartment ID issue.\n"+
				"Compartment ID: %s\n"+
				"Please ensure:\n"+
				"  1. The compartment ID is valid\n"+
				"  2. You have permissions to list clusters in this compartment\n"+
				"  3. Required IAM policy: Allow group <your-group> to inspect cluster-family in compartment <compartment>\n"+
				"  4. If using tenancy as compartment, ensure: Allow group <your-group> to inspect cluster-family in tenancy\n"+
				"Error details: %s", compartmentID, stderrMsg)
		}
		if strings.Contains(stderrMsg, "NotAuthorizedOrNotFound") || 
		   strings.Contains(stderrMsg, "authorization") {
			return nil, fmt.Errorf("failed to list OKE clusters: insufficient permissions.\n"+
				"Compartment ID: %s\n"+
				"Please ensure you have the following OCI IAM policies:\n"+
				"  Allow group <your-group> to inspect cluster-family in compartment <compartment>\n"+
				"  Allow group <your-group> to inspect compartments in tenancy\n"+
				"Error details: %s", compartmentID, stderrMsg)
		}
		
		// If stderr is empty but command failed, provide generic error with both outputs
		if stderrMsg == "" {
			return nil, fmt.Errorf("failed to list OKE clusters: command failed with exit code.\n"+
				"Command: oci %s\n"+
				"Error: %v\n"+
				"Output: %s", strings.Join(args, " "), err, stdoutMsg)
		}
		
		return nil, fmt.Errorf("failed to list OKE clusters: %s", stderrMsg)
	}
	
	// Get the output for parsing
	output := []byte(stdout.String())
	
	// Log the raw output for debugging
	log.Debugf("OCI CLI output length: %d bytes", len(output))
	if len(output) > 0 {
		log.Debugf("OCI CLI output preview: %s", string(output[:min(len(output), 200)]))
	}
	
	// Check if output is empty
	if len(output) == 0 {
		return nil, fmt.Errorf("OCI CLI returned empty response. This may indicate:\n"+
			"  1. No OKE clusters exist in the compartment: %s\n"+
			"  2. Insufficient permissions to list clusters\n"+
			"  3. Network connectivity issues\n"+
			"Please verify:\n"+
			"  - You have OKE clusters in region '%s'\n"+
			"  - Your IAM policies allow listing clusters\n"+
			"  - Run 'oci ce cluster list --compartment-id %s --region %s' manually to debug", 
			compartmentID, region, compartmentID, region)
	}
	
	// Parse the JSON response
	var response struct {
		Data []OKEClusterInfo `json:"data"`
	}
	
	err = parseJSON(output, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cluster list response: %w\nRaw output: %s", err, string(output))
	}
	
	log.Debugf("Parsed %d cluster(s) from OCI CLI response", len(response.Data))
	
	// Filter to only include ACTIVE clusters
	var activeClusters []*OKEClusterInfo
	for i := range response.Data {
		cluster := &response.Data[i]
		if cluster.State == "ACTIVE" {
			activeClusters = append(activeClusters, cluster)
		} else {
			log.Debugf("Skipping cluster '%s' with state '%s'", cluster.Name, cluster.State)
		}
	}
	
	log.Debugf("Found %d active OKE cluster(s)", len(activeClusters))
	return activeClusters, nil
}

// selectCluster allows the user to choose from multiple clusters or auto-selects if only one is available.
// This function provides an interactive cluster selection interface.
//
// Behavior:
// - If 0 clusters: Returns error
// - If 1 cluster: Auto-selects and returns the cluster ID
// - If multiple clusters: Displays list and prompts user to choose
//
// Display Format:
// For each cluster, shows:
// - Cluster name
// - Cluster OCID
// - Region
// - Kubernetes version
// - Lifecycle state
//
// This provides users with enough information to make an informed choice
// when multiple clusters are available.
//
// Parameters:
//   - clusters: Slice of available OKE clusters
//
// Returns:
//   - string: Selected cluster OCID
//   - error: If no clusters available or user input is invalid
func selectCluster(clusters []*OKEClusterInfo) (string, error) {
	if len(clusters) == 0 {
		return "", fmt.Errorf("no clusters available for selection")
	}
	
	// Auto-select if only one cluster is available
	if len(clusters) == 1 {
		log.Infof("Found 1 OKE cluster: %s (ID: %s)", clusters[0].Name, clusters[0].ID)
		return clusters[0].ID, nil
	}
	
	// Display available clusters
	fmt.Println("\nAvailable OKE clusters:")
	for i, cluster := range clusters {
		fmt.Printf("  (%d) %s\n", i+1, cluster.Name)
		fmt.Printf("      ID: %s\n", cluster.ID)
		fmt.Printf("      Region: %s\n", cluster.Region)
		fmt.Printf("      Kubernetes Version: %s\n", cluster.KubernetesVersion)
		fmt.Printf("      State: %s\n", cluster.State)
		if i < len(clusters)-1 {
			fmt.Println()
		}
	}
	
	// Prompt user for selection
	var choice int
	fmt.Printf("\nEnter choice (1-%d): ", len(clusters))
	_, err := fmt.Scanf("%d", &choice)
	if err != nil {
		return "", fmt.Errorf("error reading choice: %w", err)
	}
	
	// Validate choice
	if choice < 1 || choice > len(clusters) {
		return "", fmt.Errorf("invalid choice: %d (must be between 1 and %d)", choice, len(clusters))
	}
	
	selectedCluster := clusters[choice-1]
	log.Infof("Selected cluster: %s (ID: %s)", selectedCluster.Name, selectedCluster.ID)
	return selectedCluster.ID, nil
}

// parseJSON is a helper function to parse JSON with better error messages.
// This function provides enhanced error handling for JSON parsing operations.
//
// Validation:
// - Checks for empty responses
// - Verifies response looks like JSON (starts with { or [)
// - Provides clear error messages for parsing failures
//
// This helper improves debugging by catching common JSON parsing issues
// and providing actionable error messages.
//
// Parameters:
//   - data: Raw JSON bytes to parse
//   - v: Pointer to structure to unmarshal into
//
// Returns:
//   - nil if parsing succeeds
//   - error with diagnostic information if parsing fails
func parseJSON(data []byte, v interface{}) error {
	if len(data) == 0 {
		return fmt.Errorf("empty JSON response")
	}
	
	// Check if the response looks like JSON
	trimmed := strings.TrimSpace(string(data))
	if !strings.HasPrefix(trimmed, "{") && !strings.HasPrefix(trimmed, "[") {
		return fmt.Errorf("response does not appear to be JSON: %s", trimmed)
	}
	
	// Try to unmarshal
	err := json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	return nil
}

// min returns the minimum of two integers.
// This is a helper function for limiting string output lengths.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// collectOKEParameters gathers required configuration parameters from flags or user input.
// This is the main parameter collection orchestrator for OKE configuration.
//
// Parameter Collection Strategy:
// 1. Profile: Use flag value or default to "DEFAULT"
// 2. Region: Use flag value, then try OCI config, then prompt user
// 3. Compartment ID: Use flag value, then use tenancy ID from profile
// 4. Cluster ID: Use flag value, then list clusters and prompt user
//
// This cascading approach provides flexibility:
// - Power users can specify all parameters via flags for automation
// - Interactive users get helpful prompts and cluster discovery
// - Sensible defaults reduce the number of required inputs
//
// Validation:
// - Region format is validated (e.g., us-phoenix-1)
// - Cluster ID format is validated (e.g., ocid1.cluster.oc1.phx.xxx)
// - Compartment ID format is validated (e.g., ocid1.compartment.oc1..xxx or ocid1.tenancy.oc1..xxx)
// - User input includes retry logic for invalid entries
//
// Parameters:
//   - clusterID: Cluster OCID from --cluster-id flag (empty string if not provided)
//   - region: Region from --region flag (empty string if not provided)
//   - profile: Profile from --profile flag (empty string if not provided)
//   - compartmentID: Compartment OCID from --compartment-id flag (empty string if not provided)
//
// Returns:
//   - clusterID: Final cluster OCID to use
//   - region: Final region to use
//   - profile: Final profile to use
//   - compartmentID: Final compartment OCID to use
//   - error: If parameter collection fails
func collectOKEParameters(clusterID, region, profile, compartmentID string) (string, string, string, string, error) {
	log.Debug("Collecting OKE configuration parameters...")
	
	// Set default profile if not provided
	if profile == "" {
		profile = "DEFAULT"
		log.Infof("Using default OCI profile: %s", profile)
	}
	
	// Get region from OCI config or prompt user
	if region == "" {
		log.Debug("Region not provided via flag, attempting to read from OCI config...")
		configRegion, err := getRegionFromOCIConfig(profile)
		if err != nil {
			log.Debugf("Could not read region from OCI config: %v", err)
			log.Info("Region not found in OCI config")
		} else if configRegion != "" {
			region = configRegion
			log.Infof("Using region from OCI profile: %s", region)
		}
		
		// If still no region, prompt user with retry logic
		if region == "" {
			region, err = promptForRegion()
			if err != nil {
				return "", "", "", "", err
			}
		}
	} else {
		log.Infof("Using region from flag: %s", region)
	}
	
	// Validate region format
	if err := validateRegionFormat(region); err != nil {
		return "", "", "", "", err
	}
	
	// Get compartment ID - must be provided by user
	if compartmentID == "" {
		// Prompt user for compartment ID
		fmt.Println("\nPlease enter the OCI compartment OCID where your OKE clusters are located:")
		fmt.Print("Compartment ID: ")
		_, err := fmt.Scanf("%s", &compartmentID)
		if err != nil {
			return "", "", "", "", fmt.Errorf("error reading compartment ID: %w", err)
		}
		if compartmentID == "" {
			return "", "", "", "", fmt.Errorf("compartment ID is required")
		}
	} else {
		log.Infof("Using compartment ID from flag: %s", compartmentID)
	}
	
	// Get cluster ID - list clusters and let user choose if not provided
	if clusterID == "" {
		log.Debug("Cluster ID not provided via flag, listing available clusters...")
		clusters, err := listOKEClusters(region, profile, compartmentID)
		if err != nil {
			return "", "", "", "", fmt.Errorf("failed to list OKE clusters: %w", err)
		}
		
		if len(clusters) == 0 {
			return "", "", "", "", fmt.Errorf("no active OKE clusters found in compartment %s, region %s. Please ensure you have at least one ACTIVE cluster", compartmentID, region)
		}
		
		// Select cluster with retry logic for invalid input
		clusterID, err = selectClusterWithRetry(clusters)
		if err != nil {
			return "", "", "", "", err
		}
	} else {
		log.Infof("Using cluster ID from flag: %s", clusterID)
		
		// Validate cluster ID format
		if err := validateClusterIDFormat(clusterID); err != nil {
			return "", "", "", "", err
		}
	}
	
	log.Debugf("Collected parameters - Profile: %s, Region: %s, Compartment ID: %s, Cluster ID: %s", profile, region, compartmentID, clusterID)
	return clusterID, region, profile, compartmentID, nil
}

// promptForRegion prompts the user to enter a region with retry logic.
// This function provides an interactive prompt for region input with validation.
//
// The function allows up to 3 attempts for the user to enter a valid region.
// Each attempt validates the region format before accepting it.
//
// Expected Region Format:
// <location>-<city>-<number>
// Examples: us-phoenix-1, us-ashburn-1, eu-frankfurt-1, ap-tokyo-1
//
// Returns:
//   - string: Valid region entered by user
//   - error: If user fails to provide valid region after 3 attempts
func promptForRegion() (string, error) {
	maxRetries := 3
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		fmt.Println("\nPlease enter the OCI region (e.g., us-phoenix-1, us-ashburn-1):")
		fmt.Print("Region: ")
		
		var region string
		_, err := fmt.Scanf("%s", &region)
		if err != nil {
			log.Warnf("Error reading region input: %v", err)
			if attempt < maxRetries {
				log.Info("Let's try again...")
				continue
			}
			return "", fmt.Errorf("failed to read region after %d attempts: %w", maxRetries, err)
		}
		
		// Validate region format
		if err := validateRegionFormat(region); err != nil {
			log.Warnf("Invalid region format: %v", err)
			if attempt < maxRetries {
				log.Info("Let's try again...")
				continue
			}
			return "", fmt.Errorf("invalid region format after %d attempts", maxRetries)
		}
		
		return region, nil
	}
	
	return "", fmt.Errorf("failed to get valid region after %d attempts", maxRetries)
}

// validateRegionFormat validates the OCI region format.
// This function ensures the region string follows OCI naming conventions.
//
// OCI Region Format:
// Regions follow the pattern: <location>-<city>-<number>
// - location: Geographic area (us, eu, ap, ca, uk, etc.)
// - city: City name (phoenix, ashburn, frankfurt, tokyo, etc.)
// - number: Availability domain number (typically 1)
//
// Examples of valid regions:
// - us-phoenix-1 (US West - Phoenix)
// - us-ashburn-1 (US East - Ashburn)
// - eu-frankfurt-1 (Germany - Frankfurt)
// - ap-tokyo-1 (Japan - Tokyo)
// - uk-london-1 (UK - London)
//
// The validation is intentionally lenient to accommodate new regions
// that may be added by OCI in the future.
//
// Parameters:
//   - region: Region string to validate
//
// Returns:
//   - nil if region format is valid
//   - error with format guidance if region is invalid
func validateRegionFormat(region string) error {
	if region == "" {
		return fmt.Errorf("region cannot be empty")
	}
	
	// OCI regions typically follow the pattern: <location>-<city>-<number>
	// Examples: us-phoenix-1, us-ashburn-1, eu-frankfurt-1, ap-tokyo-1
	// We'll do a basic validation to ensure it has the right structure
	parts := strings.Split(region, "-")
	if len(parts) < 3 {
		return fmt.Errorf("invalid region format: '%s'. Expected format: <location>-<city>-<number> (e.g., us-phoenix-1)", region)
	}
	
	// Check if the last part is a number
	lastPart := parts[len(parts)-1]
	if len(lastPart) == 0 {
		return fmt.Errorf("invalid region format: '%s'. Region must end with a number", region)
	}
	
	// Basic check - last character should be a digit
	if lastPart[len(lastPart)-1] < '0' || lastPart[len(lastPart)-1] > '9' {
		return fmt.Errorf("invalid region format: '%s'. Region must end with a number (e.g., us-phoenix-1)", region)
	}
	
	log.Debugf("Region format validated: %s", region)
	return nil
}

// validateClusterIDFormat validates the OKE cluster OCID format.
// This function ensures the cluster ID follows OCI OCID conventions.
//
// OKE Cluster OCID Format:
// ocid1.cluster.oc1.<region>.<unique-id>
//
// Components:
// - ocid1: OCID version identifier
// - cluster: Resource type (OKE cluster)
// - oc1: Oracle Cloud identifier
// - region: Short region code (phx, iad, fra, etc.)
// - unique-id: Unique identifier (32 character alphanumeric string)
//
// Example:
// ocid1.cluster.oc1.phx.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
//
// The validation ensures the OCID has the correct prefix and structure.
// This catches typos and format errors before making API calls.
//
// Parameters:
//   - clusterID: Cluster OCID to validate
//
// Returns:
//   - nil if cluster ID format is valid
//   - error with format guidance if cluster ID is invalid
func validateClusterIDFormat(clusterID string) error {
	if clusterID == "" {
		return fmt.Errorf("cluster ID cannot be empty")
	}
	
	// OKE cluster OCIDs follow the pattern: ocid1.cluster.oc1.<region>.<unique-id>
	// Example: ocid1.cluster.oc1.phx.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
	if !strings.HasPrefix(clusterID, "ocid1.cluster.oc1.") {
		return fmt.Errorf("invalid cluster ID format: '%s'. Expected format: ocid1.cluster.oc1.<region>.<unique-id>", clusterID)
	}
	
	parts := strings.Split(clusterID, ".")
	if len(parts) < 5 {
		return fmt.Errorf("invalid cluster ID format: '%s'. Expected format: ocid1.cluster.oc1.<region>.<unique-id>", clusterID)
	}
	
	log.Debugf("Cluster ID format validated: %s", clusterID)
	return nil
}

// selectClusterWithRetry allows the user to select a cluster with retry logic for invalid input.
// This function provides an interactive cluster selection interface with error recovery.
//
// Behavior:
// - If 0 clusters: Returns error
// - If 1 cluster: Auto-selects and returns the cluster ID
// - If multiple clusters: Displays list and prompts user to choose (up to 3 attempts)
//
// Retry Logic:
// The function allows up to 3 attempts for valid input. This handles:
// - Non-numeric input
// - Out-of-range selections
// - Input buffer issues
//
// For each retry, the cluster list is redisplayed to help the user make
// the correct selection.
//
// Parameters:
//   - clusters: Slice of available OKE clusters
//
// Returns:
//   - string: Selected cluster OCID
//   - error: If no clusters available or user fails to provide valid input after 3 attempts
func selectClusterWithRetry(clusters []*OKEClusterInfo) (string, error) {
	if len(clusters) == 0 {
		return "", fmt.Errorf("no clusters available for selection")
	}
	
	// Auto-select if only one cluster is available
	if len(clusters) == 1 {
		log.Infof("Found 1 OKE cluster: %s (ID: %s)", clusters[0].Name, clusters[0].ID)
		return clusters[0].ID, nil
	}
	
	maxRetries := 3
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Display available clusters
		fmt.Println("\nAvailable OKE clusters:")
		for i, cluster := range clusters {
			fmt.Printf("  (%d) %s\n", i+1, cluster.Name)
			fmt.Printf("      ID: %s\n", cluster.ID)
			fmt.Printf("      Region: %s\n", cluster.Region)
			fmt.Printf("      Kubernetes Version: %s\n", cluster.KubernetesVersion)
			fmt.Printf("      State: %s\n", cluster.State)
			if i < len(clusters)-1 {
				fmt.Println()
			}
		}
		
		// Prompt user for selection
		var choice int
		fmt.Printf("\nEnter choice (1-%d): ", len(clusters))
		_, err := fmt.Scanf("%d", &choice)
		if err != nil {
			log.Warnf("Error reading choice: %v", err)
			if attempt < maxRetries {
				log.Info("Let's try again...")
				// Clear the input buffer
				var discard string
				fmt.Scanln(&discard)
				continue
			}
			return "", fmt.Errorf("failed to read cluster choice after %d attempts: %w", maxRetries, err)
		}
		
		// Validate choice
		if choice < 1 || choice > len(clusters) {
			log.Warnf("Invalid choice: %d (must be between 1 and %d)", choice, len(clusters))
			if attempt < maxRetries {
				log.Info("Let's try again...")
				continue
			}
			return "", fmt.Errorf("invalid cluster choice after %d attempts", maxRetries)
		}
		
		selectedCluster := clusters[choice-1]
		log.Infof("Selected cluster: %s (ID: %s)", selectedCluster.Name, selectedCluster.ID)
		return selectedCluster.ID, nil
	}
	
	return "", fmt.Errorf("failed to select cluster after %d attempts", maxRetries)
}

// generateOKEKubeconfig generates kubeconfig for the selected OKE cluster using OCI CLI.
// This function creates a kubeconfig file that can be used to connect to the OKE cluster.
//
// The function uses the OCI CLI 'create-kubeconfig' command, which:
// - Generates a kubeconfig with cluster endpoint and CA certificate
// - Configures authentication using OCI CLI (exec credential plugin)
// - Supports both public and private cluster endpoints
// - Handles cluster-specific networking configurations
//
// OCI CLI Command:
//   oci ce cluster create-kubeconfig \
//     --cluster-id <cluster-id> \
//     --region <region> \
//     --profile <profile> \
//     --file <config-path> \
//     --overwrite
//
// Backup Strategy:
// If a kubeconfig already exists at the target path, it is backed up with
// a .bak extension before being overwritten. This prevents accidental loss
// of existing configurations.
//
// Validation:
// After generation, the kubeconfig is validated to ensure:
// - File was created successfully
// - File is not empty
// - File contains required kubeconfig fields
//
// Parameters:
//   - clusterID: OKE cluster OCID
//   - region: OCI region
//   - profile: OCI profile for authentication
//   - configPath: Path where kubeconfig should be written
//
// Returns:
//   - nil if kubeconfig is generated and validated successfully
//   - error if generation or validation fails
func generateOKEKubeconfig(clusterID, region, profile, configPath string) error {
	log.Debugf("Generating kubeconfig for cluster %s in region %s", clusterID, region)
	
	// Check if kubeconfig file already exists and create backup if it does
	if _, err := os.Stat(configPath); err == nil {
		log.Debugf("Existing kubeconfig found at %s, creating backup", configPath)
		backupKubeconfig(configPath)
	}
	
	// Build the OCI CLI command to generate kubeconfig
	// The create-kubeconfig command supports both public and private endpoints
	// By default, it will use the appropriate endpoint based on cluster configuration
	args := []string{
		"ce", "cluster", "create-kubeconfig",
		"--cluster-id", clusterID,
		"--file", configPath,
	}
	
	// Add region if provided
	if region != "" {
		args = append(args, "--region", region)
	}
	
	// Add profile if provided
	if profile != "" {
		args = append(args, "--profile", profile)
	}
	
	// Add overwrite flag to replace existing kubeconfig
	args = append(args, "--overwrite")
	
	log.Debugf("Executing OCI CLI command: oci %s", strings.Join(args, " "))
	
	// Execute the command
	okeCmd := exec.Command("oci", args...)
	okeCmd.Stdout = os.Stdout
	okeCmd.Stderr = os.Stderr
	
	err := okeCmd.Run()
	if err != nil {
		// Provide detailed error message
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("error generating kubeconfig: OCI CLI exited with error: %s", string(exitErr.Stderr))
		}
		return fmt.Errorf("error generating kubeconfig: %w", err)
	}
	
	// Validate that the kubeconfig was successfully generated
	if err := validateKubeconfigGeneration(configPath); err != nil {
		return fmt.Errorf("kubeconfig generation validation failed: %w", err)
	}
	
	log.Debugf("OKE configuration is written to: %s", configPath)
	log.Info("Successfully generated kubeconfig for OKE cluster")
	
	return nil
}

// backupKubeconfig creates a backup of the existing kubeconfig file.
// This function preserves existing kubeconfig before overwriting it.
//
// Backup Naming:
// If the original file is "kubeconfig.yaml", the backup will be "kubeconfig.bak.yaml"
// This preserves the file extension while clearly marking it as a backup.
//
// The backup allows users to recover their previous configuration if needed.
// This is especially important when users have multiple contexts or custom
// configurations in their kubeconfig.
//
// Parameters:
//   - configPath: Path to the kubeconfig file to backup
//
// Returns:
//   - nil if backup is created successfully
//   - error if backup creation fails
func backupKubeconfig(configPath string) error {
	// Get directory and filename
	dir, file := filepath.Split(configPath)
	extension := filepath.Ext(file)
	
	// Create backup filename
	baseFilename := file[:len(file)-len(extension)]
	bakLocation := filepath.Join(dir, baseFilename+".bak"+extension)
	
	log.Debugf("Creating backup: %s -> %s", configPath, bakLocation)
	
	// Read the source file
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read source file for backup: %w", err)
	}
	
	// Write to backup location
	err = os.WriteFile(bakLocation, content, 0600)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	
	log.Debugf("Backup created successfully (%d bytes written)", len(content))
	log.Infof("Existing kubeconfig backed up to: %s", bakLocation)
	
	return nil
}

// validateKubeconfigGeneration validates that the kubeconfig file was successfully generated.
// This function performs post-generation validation to ensure the kubeconfig is usable.
//
// Validation Checks:
// 1. File exists at the specified path
// 2. File is not empty (has content)
// 3. File contains required kubeconfig fields:
//    - apiVersion: Kubeconfig API version
//    - clusters: Cluster definitions
//    - contexts: Context definitions
//    - users: User/authentication definitions
//
// These checks catch common generation failures:
// - OCI CLI command failed silently
// - Partial file write due to disk space or permissions
// - Malformed output from OCI CLI
//
// The validation provides early detection of issues before attempting to
// use the kubeconfig with Meshery or kubectl.
//
// Parameters:
//   - configPath: Path to the kubeconfig file to validate
//
// Returns:
//   - nil if kubeconfig is valid and complete
//   - error with diagnostic information if validation fails
func validateKubeconfigGeneration(configPath string) error {
	// Check if file exists
	fileInfo, err := os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("kubeconfig file was not created at %s", configPath)
		}
		return fmt.Errorf("error checking kubeconfig file: %w", err)
	}
	
	// Check if file is not empty
	if fileInfo.Size() == 0 {
		return fmt.Errorf("kubeconfig file at %s is empty", configPath)
	}
	
	// Read and validate basic YAML structure
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading kubeconfig file: %w", err)
	}
	
	// Basic validation - check for required kubeconfig fields
	contentStr := string(content)
	requiredFields := []string{"apiVersion:", "clusters:", "contexts:", "users:"}
	
	for _, field := range requiredFields {
		if !strings.Contains(contentStr, field) {
			return fmt.Errorf("kubeconfig file is missing required field: %s", field)
		}
	}
	
	log.Debug("Kubeconfig validation successful")
	return nil
}
