package core

import (
	"testing"

	"github.com/meshery/schemas/models/v1beta2/component"
	"github.com/stretchr/testify/assert"
)

func TestParseCPU(t *testing.T) {
	tests := []struct {
		input       interface{}
		expected    int64
		specified   bool
	}{
		{"500m", 500, true},
		{"1", 1000, true},
		{"1.5", 1500, true},
		{"2 cores", 2000, true},
		{"2500m", 2500, true},
		{0.5, 500, true},
		{2, 2000, true},
		{nil, 0, false},
		{"", 0, false},
	}

	for _, tt := range tests {
		val, specified := ParseCPU(tt.input)
		assert.Equal(t, tt.specified, specified)
		if tt.specified {
			assert.Equal(t, tt.expected, val)
		}
	}
}

func TestParseMemory(t *testing.T) {
	tests := []struct {
		input     interface{}
		expected  int64
		specified bool
	}{
		{"512Mi", 512 * 1024 * 1024, true},
		{"2Gi", 2 * 1024 * 1024 * 1024, true},
		{"1024Ki", 1024 * 1024, true},
		{1048576, 1048576, true},
		{nil, 0, false},
		{"", 0, false},
	}

	for _, tt := range tests {
		val, specified := ParseMemory(tt.input)
		assert.Equal(t, tt.specified, specified)
		if tt.specified {
			assert.Equal(t, tt.expected, val)
		}
	}

	// Human formatted values round trip
	val, ok := ParseMemory("2 GiB")
	assert.True(t, ok)
	assert.Equal(t, int64(2*1024*1024*1024), val)

	val, ok = ParseMemory("512 MiB")
	assert.True(t, ok)
	assert.Equal(t, int64(512*1024*1024), val)

	val, ok = ParseMemory("1000B")
	assert.True(t, ok)
	assert.Equal(t, int64(1000), val)
}

func TestFormatCPU(t *testing.T) {
	assert.Equal(t, "500m", FormatCPU(500))
	assert.Equal(t, "1 core", FormatCPU(1000))
	assert.Equal(t, "1.5 cores", FormatCPU(1500))
	assert.Equal(t, "2 cores", FormatCPU(2000))
	assert.Equal(t, "6 cores", FormatCPU(6000))
	assert.Equal(t, "7.5 cores", FormatCPU(7500))
	assert.Equal(t, "0 cores", FormatCPU(0))
}

func TestFormatMemory(t *testing.T) {
	assert.Equal(t, "512 MiB", FormatMemory(512*1024*1024))
	assert.Equal(t, "2 GiB", FormatMemory(2*1024*1024*1024))
	assert.Equal(t, "4 GiB", FormatMemory(4*1024*1024*1024))
	assert.Equal(t, "12 GiB", FormatMemory(12*1024*1024*1024))
	assert.Equal(t, "16 GiB", FormatMemory(16*1024*1024*1024))
	assert.Equal(t, "0 B", FormatMemory(0))
}

func TestExtractWorkloads(t *testing.T) {
	comps := []*component.ComponentDefinition{
		{
			DisplayName: "frontend",
			Component: component.Component{
				Kind: "Deployment",
			},
			Configuration: map[string]interface{}{
				"spec": map[string]interface{}{
					"replicas": 3,
					"template": map[string]interface{}{
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{
									"name": "web",
									"resources": map[string]interface{}{
										"requests": map[string]interface{}{
											"cpu":    "500m",
											"memory": "512Mi",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			DisplayName: "backend",
			Component: component.Component{
				Kind: "Deployment",
			},
			Configuration: map[string]interface{}{
				"spec": map[string]interface{}{
					"replicas": 3,
					"template": map[string]interface{}{
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{
									"name": "api",
									"resources": map[string]interface{}{
										"requests": map[string]interface{}{
											"cpu":    "1.5",
											"memory": "2Gi",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			DisplayName: "postgres",
			Component: component.Component{
				Kind: "StatefulSet",
			},
			Configuration: map[string]interface{}{
				"spec": map[string]interface{}{
					"replicas": 1,
					"template": map[string]interface{}{
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{
									"name": "db",
									"resources": map[string]interface{}{
										"requests": map[string]interface{}{
											"cpu":    "2",
											"memory": "4Gi",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		// Non-workload component that should be ignored by capacity preview
		{
			DisplayName: "my-service",
			Component: component.Component{
				Kind: "Service",
			},
		},
	}

	workloads, totalCPU, totalMem, totalPods := ExtractWorkloads(comps, 3)

	assert.Equal(t, 3, len(workloads))
	assert.Equal(t, 7, totalPods) // 3 + 3 + 1

	// frontend: 3 * 500m = 1500m, 3 * 512Mi = 1.5GiB
	assert.Equal(t, "frontend", workloads[0].Name)
	assert.Equal(t, 3, workloads[0].Replicas)
	assert.Equal(t, "500m", workloads[0].CPURequest)
	assert.Equal(t, "512 MiB", workloads[0].MemoryRequest)
	assert.Equal(t, int64(1500), workloads[0].CPUMillicores)

	// backend: 3 * 1500m = 4500m, 3 * 2Gi = 6GiB
	assert.Equal(t, "backend", workloads[1].Name)
	assert.Equal(t, 3, workloads[1].Replicas)
	assert.Equal(t, "1.5 cores", workloads[1].CPURequest)
	assert.Equal(t, "2 GiB", workloads[1].MemoryRequest)
	assert.Equal(t, int64(4500), workloads[1].CPUMillicores)

	// postgres: 1 * 2000m = 2000m, 1 * 4Gi = 4GiB
	assert.Equal(t, "postgres", workloads[2].Name)
	assert.Equal(t, 1, workloads[2].Replicas)
	assert.Equal(t, "2 cores", workloads[2].CPURequest)
	assert.Equal(t, "4 GiB", workloads[2].MemoryRequest)
	assert.Equal(t, int64(2000), workloads[2].CPUMillicores)

	// Totals:
	// CPU: 1500 + 4500 + 2000 = 8000m = 8 cores
	// Mem: 1.5GiB + 6GiB + 4GiB = 11.5GiB
	assert.Equal(t, int64(8000), totalCPU)
	assert.Equal(t, int64(11.5*1024*1024*1024), totalMem)
}

func TestCalculateResourceCapacity_OverCapacity(t *testing.T) {
	comps := []*component.ComponentDefinition{
		{
			DisplayName: "backend",
			Component: component.Component{
				Kind: "Deployment",
			},
			Configuration: map[string]interface{}{
				"spec": map[string]interface{}{
					"replicas": 5, // 5 * 1.5 cores = 7.5 cores
					"template": map[string]interface{}{
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{
									"name": "api",
									"resources": map[string]interface{}{
										"requests": map[string]interface{}{
											"cpu":    "1.5",
											"memory": "2Gi", // 5 * 2Gi = 10Gi
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	clusterCap := ClusterNodeCapacity{
		NodeCount:           2,
		TotalAllocatableCPU: 8000,                     // 8 cores allocatable
		TotalAllocatedCPU:   2000,                     // 2 cores already used -> 6 cores available
		TotalAllocatableMem: 20 * 1024 * 1024 * 1024, // 20 GiB allocatable
		TotalAllocatedMem:   4 * 1024 * 1024 * 1024,  // 4 GiB already used -> 16 GiB available
		HasData:             true,
	}

	preview := CalculateResourceCapacity("ctx-1", "dev-cluster", comps, clusterCap)

	assert.Equal(t, "ctx-1", preview.ClusterID)
	assert.Equal(t, "dev-cluster", preview.ClusterName)
	assert.Equal(t, 1, preview.WorkloadCount)
	assert.Equal(t, 5, preview.EstimatedPods)

	// CPU requires 7.5 cores, available is 6 cores -> OverCapacity
	assert.Equal(t, "7.5 cores", preview.CPU.RequiredFormatted)
	assert.Equal(t, "6 cores", preview.CPU.AvailableFormatted)
	assert.Equal(t, StatusOverCapacity, preview.CPU.Status)
	assert.Equal(t, "1.5 cores above estimated availability", preview.CPU.DifferenceFormatted)
	assert.Equal(t, 125, preview.CPU.Percentage)

	// Memory requires 10 GiB, available is 16 GiB -> WithinCapacity
	assert.Equal(t, "10 GiB", preview.Memory.RequiredFormatted)
	assert.Equal(t, "16 GiB", preview.Memory.AvailableFormatted)
	assert.Equal(t, StatusWithinCapacity, preview.Memory.Status)
	assert.Equal(t, "6 GiB estimated headroom", preview.Memory.DifferenceFormatted)

	assert.Equal(t, StatusOverCapacity, preview.Status)
	assert.NotEmpty(t, preview.Warnings)
}

func TestCalculateResourceCapacity_UnknownCluster(t *testing.T) {
	comps := []*component.ComponentDefinition{
		{
			DisplayName: "worker",
			Component: component.Component{
				Kind: "Deployment",
			},
			Configuration: map[string]interface{}{
				"spec": map[string]interface{}{
					"replicas": 2,
					"template": map[string]interface{}{
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{
									"name": "job-runner",
									"resources": map[string]interface{}{
										"requests": map[string]interface{}{
											"cpu":    "1",
											"memory": "1Gi",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	clusterCap := ClusterNodeCapacity{
		HasData: false,
	}

	preview := CalculateResourceCapacity("unknown-ctx", "remote-cluster", comps, clusterCap)

	assert.Equal(t, StatusUnknown, preview.Status)
	assert.True(t, preview.CPU.IsUnknown)
	assert.True(t, preview.Memory.IsUnknown)
	assert.Equal(t, "Unknown", preview.CPU.AvailableFormatted)
	assert.Equal(t, "Unknown", preview.Memory.AvailableFormatted)
	assert.Equal(t, EstimateUnknown, preview.Workloads[0].Estimate)
}

func TestCalculateResourceCapacity_MissingRequests(t *testing.T) {
	comps := []*component.ComponentDefinition{
		{
			DisplayName: "unspecified-workload",
			Component: component.Component{
				Kind: "Deployment",
			},
			Configuration: map[string]interface{}{
				"spec": map[string]interface{}{
					"replicas": 2,
					"template": map[string]interface{}{
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{
									"name": "nginx",
									// No resources specified
								},
							},
						},
					},
				},
			},
		},
	}

	clusterCap := ClusterNodeCapacity{
		NodeCount:           1,
		TotalAllocatableCPU: 4000,
		TotalAllocatableMem: 8 * 1024 * 1024 * 1024,
		HasData:             true,
	}

	preview := CalculateResourceCapacity("ctx-2", "k8s-cluster", comps, clusterCap)

	assert.Equal(t, 1, preview.WorkloadCount)
	assert.True(t, preview.Workloads[0].HasUnspecified)
	assert.Equal(t, "Unspecified", preview.Workloads[0].CPURequest)
	assert.Equal(t, "Unspecified", preview.Workloads[0].MemoryRequest)
	assert.Equal(t, EstimateReview, preview.Workloads[0].Estimate)
}

func TestExtractWorkloads_InitContainers(t *testing.T) {
	comps := []*component.ComponentDefinition{
		{
			DisplayName: "init-workload",
			Component: component.Component{
				Kind: "Deployment",
			},
			Configuration: map[string]interface{}{
				"spec": map[string]interface{}{
					"replicas": 1,
					"template": map[string]interface{}{
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{
									"name": "app",
									"resources": map[string]interface{}{
										"requests": map[string]interface{}{
											"cpu":    "500m",
											"memory": "512Mi",
										},
									},
								},
							},
							"initContainers": []interface{}{
								map[string]interface{}{
									"name": "init-db",
									"resources": map[string]interface{}{
										"requests": map[string]interface{}{
											"cpu":    "2",
											"memory": "2Gi",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	workloads, totalCPU, totalMem, totalPods := ExtractWorkloads(comps, 1)
	assert.Equal(t, 1, len(workloads))
	assert.Equal(t, 1, totalPods)
	// Init container 2 cores / 2Gi exceeds app container 500m / 512Mi, so effective request is 2 cores / 2Gi
	assert.Equal(t, int64(2000), totalCPU)
	assert.Equal(t, int64(2*1024*1024*1024), totalMem)
	assert.Equal(t, "2 cores", workloads[0].CPURequest)
	assert.Equal(t, "2 GiB", workloads[0].MemoryRequest)
}

func TestExtractWorkloads_PodLevelRequests(t *testing.T) {
	comps := []*component.ComponentDefinition{
		{
			DisplayName: "pod-level-workload",
			Component: component.Component{
				Kind: "Pod",
			},
			Configuration: map[string]interface{}{
				"spec": map[string]interface{}{
					"resources": map[string]interface{}{
						"requests": map[string]interface{}{
							"cpu":    "4",
							"memory": "4Gi",
						},
					},
					"containers": []interface{}{
						map[string]interface{}{
							"name": "c1",
							"resources": map[string]interface{}{
								"requests": map[string]interface{}{
									"cpu":    "500m",
									"memory": "512Mi",
								},
							},
						},
					},
				},
			},
		},
	}

	workloads, totalCPU, totalMem, totalPods := ExtractWorkloads(comps, 1)
	assert.Equal(t, 1, len(workloads))
	assert.Equal(t, 1, totalPods)
	// Explicit pod-level request of 4 cores / 4Gi takes precedence
	assert.Equal(t, int64(4000), totalCPU)
	assert.Equal(t, int64(4*1024*1024*1024), totalMem)
	assert.Equal(t, "4 cores", workloads[0].CPURequest)
	assert.Equal(t, "4 GiB", workloads[0].MemoryRequest)
	assert.False(t, workloads[0].HasUnspecified)
}

func TestExtractReplicas(t *testing.T) {
	// Job with explicit parallelism 0 (paused)
	jobPaused := map[string]interface{}{
		"spec": map[string]interface{}{
			"parallelism": 0,
		},
	}
	assert.Equal(t, 0, extractReplicas("Job", jobPaused, 1))

	// Job with suspend: true
	jobSuspended := map[string]interface{}{
		"spec": map[string]interface{}{
			"suspend":     true,
			"parallelism": 4,
		},
	}
	assert.Equal(t, 0, extractReplicas("Job", jobSuspended, 1))

	// Job with default parallelism
	jobDefault := map[string]interface{}{
		"spec": map[string]interface{}{},
	}
	assert.Equal(t, 1, extractReplicas("Job", jobDefault, 1))

	// Deployment with explicit replicas: 0
	deployZero := map[string]interface{}{
		"spec": map[string]interface{}{
			"replicas": 0,
		},
	}
	assert.Equal(t, 0, extractReplicas("Deployment", deployZero, 1))

	// Deployment with replicas: nil (omitted, default to 1)
	deployNil := map[string]interface{}{
		"spec": map[string]interface{}{
			"replicas": nil,
		},
	}
	assert.Equal(t, 1, extractReplicas("Deployment", deployNil, 1))

	// CronJob with parallelism
	cronJob := map[string]interface{}{
		"spec": map[string]interface{}{
			"jobTemplate": map[string]interface{}{
				"spec": map[string]interface{}{
					"parallelism": 3,
				},
			},
		},
	}
	assert.Equal(t, 3, extractReplicas("CronJob", cronJob, 1))

	// CronJob suspended
	cronJobSuspended := map[string]interface{}{
		"spec": map[string]interface{}{
			"suspend": true,
			"jobTemplate": map[string]interface{}{
				"spec": map[string]interface{}{
					"parallelism": 3,
				},
			},
		},
	}
	assert.Equal(t, 0, extractReplicas("CronJob", cronJobSuspended, 1))
}
