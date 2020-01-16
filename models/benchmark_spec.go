package models

import "time"

type Environment struct {
	Kubernetes string `yaml:"kubernetes,omitempty"`
	NodeCount  int    `yaml:"node_count,omitempty"`
}

type MeshConfig struct {
	MeshPolicyEnabled    bool `yaml:"mesh_policy_enabled,omitempty"`
	MeshTelemetryEnabled bool `yaml:"mesh_telemetry_enabled,omitempty"`
	MtlsEnabled          bool `yaml:"mtls_enabled,omitempty"`
	ProxyConcurrency     int  `yaml:"proxy_concurrency,omitempty"`
}

type MeshClientConfig struct {
	Internal    bool         `yaml:"internal,omitempty"`
	Protocol    string       `yaml:"protocol,omitempty"`
	Connections int          `yaml:"connections,omitempty"`
	Rps         float64      `yaml:"rps,omitempty"`
	LatenciesMs *LatenciesMs `yaml:"latencies_ms,omitempty"`
}

type LatenciesMs struct {
	Min     float64 `yaml:"min,omitempty"`
	Average float64 `yaml:"average,omitempty"`
	P50     float64 `yaml:"p50,omitempty"`
	P90     float64 `yaml:"p90,omitempty"`
	P99     float64 `yaml:"p99,omitempty"`
	Max     float64 `yaml:"max,omitempty"`
}

type IngressGateway struct {
	Count     int     `yaml:"count,omitempty"`
	CPUMCores float64 `yaml:"cpu_mCores,omitempty"`
	MemMb     float64 `yaml:"mem_mb,omitempty"`
	Rps       float64 `yaml:"rps,omitempty"`
	Bps       float64 `yaml:"bps,omitempty"`
}

type Sidecars struct {
	Count     int     `yaml:"count,omitempty"`
	CPUMCores float64 `yaml:"cpu_mCores,omitempty"`
	MemMb     float64 `yaml:"mem_mb,omitempty"`
	Rps       float64 `yaml:"rps,omitempty"`
	Bps       float64 `yaml:"bps,omitempty"`
}

type MeshTelemetry struct {
	Count     int     `yaml:"count,omitempty"`
	CPUMCores float64 `yaml:"cpu_mCores,omitempty"`
	MemMb     float64 `yaml:"mem_mb,omitempty"`
	Rps       float64 `yaml:"rps,omitempty"`
}

type MeshPolicy struct {
	Count        int     `yaml:"count,omitempty"`
	CPUMCores    float64 `yaml:"cpu_mCores,omitempty"`
	MemMb        float64 `yaml:"mem_mb,omitempty"`
	Rps          float64 `yaml:"rps,omitempty"`
	CacheHitRate float64 `yaml:"cache_hit_rate,omitempty"`
}

type MeshControlPlane struct {
	Count            int     `yaml:"count,omitempty"`
	CPUMCores        float64 `yaml:"cpu_mCores,omitempty"`
	MemMb            float64 `yaml:"mem_mb,omitempty"`
	Endpoints        int     `yaml:"endpoints,omitempty"`
	Services         int     `yaml:"services,omitempty"`
	Sidecars         int     `yaml:"sidecars,omitempty"`
	VirtualServices  int     `yaml:"virtual_services,omitempty"`
	DestinationRules int     `yaml:"destination_rules,omitempty"`
	LdsLatencyMs     float64 `yaml:"lds_latency_ms,omitempty"`
	CdsLatencyMs     float64 `yaml:"cds_latency_ms,omitempty"`
}

type Workload struct {
	Name      string  `yaml:"name,omitempty"`
	Count     int     `yaml:"count,omitempty"`
	CPUMCores float64 `yaml:"cpu_mCores,omitempty"`
	MemMb     float64 `yaml:"mem_mb,omitempty"`
}

type Metrics struct {
	IngressGateway     *IngressGateway   `yaml:"ingress_gateway,omitempty"`
	Sidecars           *Sidecars         `yaml:"sidecars,omitempty"`
	MeshTelemetry      *MeshTelemetry    `yaml:"mesh_telemetry,omitempty"`
	MeshPolicy         *MeshPolicy       `yaml:"mesh_policy,omitempty"`
	MeshControlPlane   *MeshControlPlane `yaml:"mesh_control_plane,omitempty"`
	IndividualWorkload *Workload         `yaml:"individual_workload_1,omitempty"`
}

type BenchmarkSpec struct {
	StartTime    time.Time         `yaml:"start_time,omitempty"`
	EndTime      time.Time         `yaml:"end_time,omitempty"`
	MeshBuild    string            `yaml:"mesh_build,omitempty"`
	ProxyBuild   string            `yaml:"proxy_build,omitempty"`
	ExpGroupUUID string            `yaml:"exp_group_uuid,omitempty"`
	ExpUUID      string            `yaml:"exp_uuid,omitempty"`
	Profile      string            `yaml:"profile,omitempty"`
	DetailsURI   string            `yaml:"details_uri,omitempty"`
	EndpointURL  string            `yaml:"endpoint_url,omitempty"`
	Env          *Environment      `yaml:"env,omitempty"`
	Config       *MeshConfig       `yaml:"config,omitempty"`
	Client       *MeshClientConfig `yaml:"client,omitempty"`
	Metrics      *Metrics          `yaml:"metrics,omitempty"`
}

// func (b *BenchmarkSpec) ConvertToMesheryResult() *MesheryResult {

// }
