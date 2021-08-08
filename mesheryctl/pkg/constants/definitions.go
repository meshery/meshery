package constants

import "flag"

const (
	// Meshery Repository Location
	mesheryGitHubOrg  string = "meshery"
	mesheryGitHubRepo string = "meshery"
)

//Update ...Used in testing to update golden files
var Update = flag.Bool("update", false, "update golden files")

// GetMesheryGitHubOrg retrieves the name of the GitHub organization under which the Meshery repository resides.
func GetMesheryGitHubOrg() string {
	return mesheryGitHubOrg
}

// GetMesheryGitHubRepo retrieves the name of the Meshery repository
func GetMesheryGitHubRepo() string {
	return mesheryGitHubRepo
}
