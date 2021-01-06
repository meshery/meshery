package cfg

import (
	"github.com/spf13/viper"
)

// PerfCfg stores the perf subcommand configurations
type PerfCfg struct {
	baseURL string
	// AuthTokenURI stores the URI for the token endpoint
	AuthTokenURI string `json:"AuthtokenURI"`
	// LoadTestSmpURI stores the URI for the load test SMP endpoint
	LoadTestSmpURI string `json:"loadTestSmpURI"`
	LoadTestURI    string `json:"loadTestURI"`
}

// Version unmarshals the json response from the server's version api
type Version struct {
	Build          string `json:"build,omitempty"`
	CommitSHA      string `json:"commitsha,omitempty"`
	ReleaseChannel string `json:"release_channel,omitempty"`
}

// MesheryCtl stores the configurations used by mesheryctl CLI
// Default config:
//  baseMesheryURL: "http://localhost:9081/api",
//  perf:
//	  authTokenURI:    "/gettoken",
//	  loadTestSmpURI: "/perf/load-test-smp",
//	  loadTestURI: "/perf/load-test",
type MesheryCtl struct {
	// BaseMesheryURL is the base URL of the meshery server
	BaseMesheryURL string `json:"baseMesheryURL"`
	// Perf stores the perf subcommand configurations
	Perf *PerfCfg `json:"perf"`

	CtlVersion *Version `json:"ctl_version"`
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

// GetLoadTestSmpURL returns a fully qualified URL to the Load Test SMP endpoint
func (pc *PerfCfg) GetLoadTestSmpURL() string {
	return pc.baseURL + pc.LoadTestSmpURI
}

// GetLoadTestURL returns a fully qualified URL to the Load Test endpoint
func (pc *PerfCfg) GetLoadTestURL() string {
	return pc.baseURL + pc.LoadTestURI
}

// GetVersion returns the version details of this binary
func (mc *MesheryCtl) GetVersion() *Version {
	return mc.CtlVersion
}

// GetBuild returns the build number for the binary
func (v *Version) GetBuild() string {
	return v.Build
}

// GetCommitSHA returns the commit sha for the binary
func (v *Version) GetCommitSHA() string {
	return v.CommitSHA
}
