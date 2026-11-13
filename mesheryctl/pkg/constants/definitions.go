package constants

import "time"

const (
	// Meshery Repository Location
	mesheryGitHubOrg  string = "meshery"
	mesheryGitHubRepo string = "meshery"
	CallbackURLENV    string = "MESHERY_SERVER_CALLBACK_URL"
	ProviderENV       string = "PROVIDER"
	ProviderURLsENV   string = "PROVIDER_BASE_URLS"

	// HttpClientTimeout defines the standard timeout for HTTP requests in mesheryctl
	HttpClientTimeout = 10 * time.Second
)

// GetMesheryGitHubOrg retrieves the name of the GitHub organization under which the Meshery repository resides.
func GetMesheryGitHubOrg() string {
	return mesheryGitHubOrg
}

// GetMesheryGitHubRepo retrieves the name of the Meshery repository
func GetMesheryGitHubRepo() string {
	return mesheryGitHubRepo
}
