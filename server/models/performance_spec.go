package models

import (
	"time"
)

// Environment - represents a kubernetes environment
type Environment struct {
	Kubernetes string `yaml:"kubernetes,omitempty"`
	NodeCount  int    `yaml:"node_count,omitempty"`
}

// MeshConfig - represents a service mesh config
type MeshConfig struct {
	MeshPolicyEnabled    bool `yaml:"mesh_policy_enabled,omitempty"`
	MeshTelemetryEnabled bool `yaml:"mesh_telemetry_enabled,omitempty"`
	MtlsEnabled          bool `yaml:"mtls_enabled,omitempty"`
	ProxyConcurrency     int  `yaml:"proxy_concurrency,omitempty"`
}

// MeshClientConfig - represents a load test client config
type MeshClientConfig struct {
	Internal    bool         `yaml:"internal,omitempty"`
	Protocol    string       `yaml:"protocol,omitempty"`
	Connections int          `yaml:"connections,omitempty"`
	Rps         float64      `yaml:"rps,omitempty"`
	LatenciesMs *LatenciesMs `yaml:"latencies_ms,omitempty"`
}

// LatenciesMs - represents a collection of important latencies
type LatenciesMs struct {
	Min     float64 `yaml:"min,omitempty"`
	Average float64 `yaml:"average,omitempty"`
	P50     float64 `yaml:"p50,omitempty"`
	P90     float64 `yaml:"p90,omitempty"`
	P99     float64 `yaml:"p99,omitempty"`
	Max     float64 `yaml:"max,omitempty"`
}

// IngressGateway - holds ingress gateway info
type IngressGateway struct {
	Count     int     `yaml:"count,omitempty"`
	CPUMCores float64 `yaml:"cpu_mCores,omitempty"`
	MemMb     float64 `yaml:"mem_mb,omitempty"`
	Rps       float64 `yaml:"rps,omitempty"`
	Bps       float64 `yaml:"bps,omitempty"`
}

// Sidecars - holds sidecars info
type Sidecars struct {
	Count     int     `yaml:"count,omitempty"`
	CPUMCores float64 `yaml:"cpu_mCores,omitempty"`
	MemMb     float64 `yaml:"mem_mb,omitempty"`
	Rps       float64 `yaml:"rps,omitempty"`
	Bps       float64 `yaml:"bps,omitempty"`
}

// MeshTelemetry - holds overall Mesh info
type MeshTelemetry struct {
	Count     int     `yaml:"count,omitempty"`
	CPUMCores float64 `yaml:"cpu_mCores,omitempty"`
	MemMb     float64 `yaml:"mem_mb,omitempty"`
	Rps       float64 `yaml:"rps,omitempty"`
}

// MeshPolicy - holds MeshPolicy info
type MeshPolicy struct {
	Count        int     `yaml:"count,omitempty"`
	CPUMCores    float64 `yaml:"cpu_mCores,omitempty"`
	MemMb        float64 `yaml:"mem_mb,omitempty"`
	Rps          float64 `yaml:"rps,omitempty"`
	CacheHitRate float64 `yaml:"cache_hit_rate,omitempty"`
}

// MeshControlPlane - holds control plan info
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

// Workload - holds workload info
type Workload struct {
	Name      string  `yaml:"name,omitempty"`
	Count     int     `yaml:"count,omitempty"`
	CPUMCores float64 `yaml:"cpu_mCores,omitempty"`
	MemMb     float64 `yaml:"mem_mb,omitempty"`
}

// Metrics - holds overall metrics info
type Metrics struct {
	IngressGateway     *IngressGateway   `yaml:"ingress_gateway,omitempty"`
	Sidecars           *Sidecars         `yaml:"sidecars,omitempty"`
	MeshTelemetry      *MeshTelemetry    `yaml:"mesh_telemetry,omitempty"`
	MeshPolicy         *MeshPolicy       `yaml:"mesh_policy,omitempty"`
	MeshControlPlane   *MeshControlPlane `yaml:"mesh_control_plane,omitempty"`
	IndividualWorkload *Workload         `yaml:"individual_workload_1,omitempty"`
}

// PerformanceSpec - represents SMP, see here https://github.com/layer5io/service-mesh-performance-specification
type PerformanceSpec struct {
	SMPVersion   string            `yaml:"smp_version,omitempty"`
	id           string            `yaml:"id,omitempty"`
	labels       map[string]string `yaml:"labels,omitempty"`
	StartTime    time.Time         `yaml:"start_time,omitempty"`
	EndTime      time.Time         `yaml:"end_time,omitempty"`
	Latencies    *LatenciesMs      `yaml:"latencies_ms,omitempty"`
	ActualQPS    float64           `yaml:"actual_qps,omitempty"`
	DetailsURI   string            `yaml:"details_uri,omitempty"`
	TestID       string            `yaml:"test_id,omitempty"`
	MeshConfigID string            `yaml:"mesh_config_id,omitempty"`
	EnvID        string            `yaml:"env_id,omitempty"`
}
