package models

type MesheryCtlConfig struct {
	Perf           PerformanceEndpoints `json:"perf"`
	Contexts       map[string]Context   `json:"contexts"`
	CurrentContext string               `json:"currentcontext"`
}

type Context struct {
	BaseMesheryURL string `json:"basemesheryurl"`
}

type PerformanceEndpoints struct {
	AuthTokenURI   string `json:"authtokenuri"`
	LoadTestSMPURI string `json:"loadtestsmpuri"`
	LoadTestURI    string `json:"loadtesturi"`
}
