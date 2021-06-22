package constants

var (
	version        = "Not Set"
	commitsha      = "Not Set"
	releasechannel = "Not Set"
)

func GetMesheryctlVersion() string {
	return version
}

func GetMesheryctlCommitsha() string {
	return commitsha
}

func GetMesheryctlReleaseChannel() string {
	return releasechannel
}

var (
	token = "Not Set"
)

func GetAuthenticationToken() string {
	// TODO
	// Check if token is downloaded and present on the location
	// Return token
	return token
}
