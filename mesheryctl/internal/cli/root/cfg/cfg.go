package cfg

import "github.com/spf13/viper"

// PerfCfg stores the perf subcommand configurations
type PerfCfg struct {
	baseURL string
	// AuthTokenURI stores the URI for the token endpoint
	AuthTokenURI string `json:"AuthtokenURI"`
	// LoadTestSmpsURI stores the URI for the load test SMPS endpoint
	LoadTestSmpsURI string `json:"loadTestSmpsURI"`
}

// MesheryCtl stores the configurations used by mesheryctl CLI
// Default config:
//  baseMesheryURL: "http://localhost:9081/api",
//  perf:
//	  authTokenURI:    "/gettoken",
//	  loadTestSmpsURI: "/load-test-smps",
type MesheryCtl struct {
	// BaseMesheryURL is the base URL of the meshery server
	BaseMesheryURL string `json:"baseMesheryURL"`
	// Perf stores the perf subcommand configurations
	Perf *PerfCfg `json:"perf"`
}

// GetMesheryCtl returns a reference to the mesheryctl configuration object.
func GetMesheryCtl(v *viper.Viper) (*MesheryCtl, error) {
	c := &MesheryCtl{}
	// Load the config data into the object
	err := v.Unmarshal(&c)
	if err != nil {
		return nil, err
	}
	return c, err
}

// GetBaseMesheryURL returns the base meshery server URL
func (mc *MesheryCtl) GetBaseMesheryURL() string {
	return mc.BaseMesheryURL
}

// GetPerf returns a reference to the perf configuration object
func (mc *MesheryCtl) GetPerf() *PerfCfg {
	mc.Perf.baseURL = mc.BaseMesheryURL
	return mc.Perf
}

// GetAuthTokenURL returns a fully qualified URL to the token endpoint
func (pc *PerfCfg) GetAuthTokenURL() string {
	return pc.baseURL + pc.AuthTokenURI
}

// GetLoadTestSmpsURL returns a fully qualified URL to the Load Test SMPS endpoint
func (pc *PerfCfg) GetLoadTestSmpsURL() string {
	return pc.baseURL + pc.LoadTestSmpsURI
}
