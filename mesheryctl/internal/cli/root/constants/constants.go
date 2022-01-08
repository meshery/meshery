package constants

var (
	version        = "Not Set"
	commitsha      = "Not Set"
	releasechannel = "Not Set"
	// TokenFlag sets token location passed by user with --token
	TokenFlag = "Not Set"
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
