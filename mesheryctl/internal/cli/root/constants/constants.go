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
